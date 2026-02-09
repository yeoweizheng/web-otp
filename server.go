package main

import (
	"database/sql"
	"fmt"
	"io/fs"
	"net/http"
	"strconv"
	"strings"
	"time"

	"embed"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/pquerna/otp/totp"
)

//go:embed static/*
var staticFS embed.FS

func StartServer(db *sql.DB) {
	r := gin.New()

	// For frontend development - serve from directory
	// r.NoRoute(gin.WrapH(http.FileServer(http.Dir("static"))))

	// For production - serve from embedded FS
	staticSubFS, _ := fs.Sub(staticFS, "static")
	r.NoRoute(gin.WrapH(http.FileServerFS(staticSubFS)))

	api := r.Group("/api")
	api.Use(DbMiddleware(db))
	api.POST("/login/", LoginAPI)
	api.POST("/refresh/", RefreshAPI)
	api.POST("/logout/", LogoutAPI)

	authApi := api.Group("/")
	authApi.Use(AuthRequired())
	authApi.GET("/account_otps/", GetAccountOTPsAPI)
	authApi.POST("/add_account/", AddAccountAPI)
	authApi.POST("/reveal_account_token/:id/", RevealAccountTokenAPI)
	authApi.PATCH("/update_account/:id/", UpdateAccountAPI)
	authApi.DELETE("/delete_account/:id/", DeleteAccountAPI)

	r.Run(HOST_PORT)
}

func createToken(userId int, tokenType string, maxAgeSeconds int) (string, error) {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
		"type":   tokenType,
		"iat":    now.Unix(),
		"exp":    now.Add(time.Duration(maxAgeSeconds) * time.Second).Unix(),
	})
	return token.SignedString([]byte(JWT_SECRET))
}

func setAuthCookies(c *gin.Context, accessToken string, refreshToken string) {
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie(ACCESS_TOKEN_COOKIE_NAME, accessToken, ACCESS_TOKEN_MAX_AGE_SECONDS, "/", "", true, true)
	c.SetCookie(REFRESH_TOKEN_COOKIE_NAME, refreshToken, REFRESH_TOKEN_MAX_AGE_SECONDS, "/", "", true, true)
}

func clearAuthCookies(c *gin.Context) {
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie(ACCESS_TOKEN_COOKIE_NAME, "", -1, "/", "", true, true)
	c.SetCookie(REFRESH_TOKEN_COOKIE_NAME, "", -1, "/", "", true, true)
}

func parseTokenClaims(tokenString string, expectedType string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}
		return []byte(JWT_SECRET), nil
	})
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}
	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != expectedType {
		return nil, fmt.Errorf("invalid token type")
	}
	return claims, nil
}

func LoginAPI(c *gin.Context) {
	db := GetDBFromCtx(c)
	var data struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(400, gin.H{"details": "Invalid request"})
		return
	}
	userId, err := VerifyAndGetUserId(db, data.Username, data.Password)
	if err == nil {
		accessToken, err := createToken(userId, "access", ACCESS_TOKEN_MAX_AGE_SECONDS)
		if err != nil {
			c.JSON(500, gin.H{"details": "Failed to create access token"})
			return
		}
		refreshToken, err := createToken(userId, "refresh", REFRESH_TOKEN_MAX_AGE_SECONDS)
		if err != nil {
			c.JSON(500, gin.H{"details": "Failed to create refresh token"})
			return
		}
		setAuthCookies(c, accessToken, refreshToken)
		c.JSON(200, gin.H{"username": data.Username})
	} else {
		c.JSON(401, gin.H{"details": "Invalid username / password"})
	}
}

func RefreshAPI(c *gin.Context) {
	db := GetDBFromCtx(c)
	refreshToken, err := c.Cookie(REFRESH_TOKEN_COOKIE_NAME)
	if err != nil {
		c.JSON(401, gin.H{"details": "Unauthorized"})
		return
	}

	claims, err := parseTokenClaims(refreshToken, "refresh")
	if err != nil {
		clearAuthCookies(c)
		c.JSON(401, gin.H{"details": "Unauthorized"})
		return
	}

	userId, ok := claims["userId"].(float64)
	if !ok || !UserIDExists(db, int(userId)) {
		clearAuthCookies(c)
		c.JSON(401, gin.H{"details": "Unauthorized"})
		return
	}

	accessToken, err := createToken(int(userId), "access", ACCESS_TOKEN_MAX_AGE_SECONDS)
	if err != nil {
		c.JSON(500, gin.H{"details": "Failed to create access token"})
		return
	}
	rotatedRefreshToken, err := createToken(int(userId), "refresh", REFRESH_TOKEN_MAX_AGE_SECONDS)
	if err != nil {
		c.JSON(500, gin.H{"details": "Failed to create refresh token"})
		return
	}
	setAuthCookies(c, accessToken, rotatedRefreshToken)
	c.Status(204)
}

func LogoutAPI(c *gin.Context) {
	clearAuthCookies(c)
	c.Status(204)
}

func GetAccountOTPsAPI(c *gin.Context) {
	db := GetDBFromCtx(c)
	userId := GetUserIdFromCtx((c))
	var timeStruct time.Time
	tsParam := c.Query("timestamp")
	timestamp, err := strconv.ParseInt(tsParam, 10, 64)
	if err == nil {
		timeStruct = time.Unix(timestamp, 0)
	} else {
		timeStruct = time.Now()
	}
	accounts := GetAccounts(db, userId)
	var accountOTPs []AccountOTP
	for _, account := range accounts {
		code, _ := totp.GenerateCode(strings.ReplaceAll(account.Token, " ", ""), timeStruct)
		accountOTP := AccountOTP{account.Id, account.Name, code}
		accountOTPs = append(accountOTPs, accountOTP)
	}
	c.JSON(200, accountOTPs)
}

func RevealAccountTokenAPI(c *gin.Context) {
	db := GetDBFromCtx(c)
	userId := GetUserIdFromCtx(c)
	accountId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"details": "Invalid account ID"})
		return
	}

	var data map[string]interface{}
	c.BindJSON(&data)
	password, ok := data["password"].(string)
	if !ok || password == "" {
		c.JSON(400, gin.H{"details": "Password required"})
		return
	}

	if err := VerifyUserPasswordByID(db, userId, password); err != nil {
		c.JSON(403, gin.H{"details": "Invalid password"})
		return
	}

	token, err := GetAccountToken(db, userId, accountId)
	if err != nil {
		c.JSON(404, gin.H{"details": "Account not found"})
		return
	}
	c.JSON(200, gin.H{"token": token})
}

func AddAccountAPI(c *gin.Context) {
	db := GetDBFromCtx(c)
	userId := GetUserIdFromCtx((c))
	var data map[string]interface{}
	c.BindJSON(&data)
	lastInsertId := CreateAccount(db, userId, data["name"].(string), data["token"].(string))
	c.JSON(201, Account{lastInsertId, data["name"].(string), data["token"].(string)})
}

func UpdateAccountAPI(c *gin.Context) {
	db := GetDBFromCtx(c)
	userId := GetUserIdFromCtx((c))
	var data map[string]interface{}
	c.BindJSON(&data)
	accountId, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	name, hasName := data["name"].(string)
	if !hasName || strings.TrimSpace(name) == "" {
		c.JSON(400, gin.H{"details": "Account name required"})
		return
	}

	token, hasToken := data["token"].(string)
	hasReplacementToken := hasToken && strings.TrimSpace(token) != ""
	var rowCount int64
	if hasReplacementToken {
		rowCount = UpdateAccount(db, userId, accountId, name, token)
	} else {
		rowCount = UpdateAccountName(db, userId, accountId, name)
	}

	if rowCount == 1 {
		response := gin.H{"id": accountId, "name": name}
		if hasReplacementToken {
			response["token"] = token
		}
		c.JSON(200, response)
	} else {
		c.JSON(400, gin.H{"details": "Failed to update account"})
	}
}

func DeleteAccountAPI(c *gin.Context) {
	db := GetDBFromCtx(c)
	userId := GetUserIdFromCtx(c)
	accountId, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	rowCount := DeleteAccount(db, userId, accountId)
	if rowCount == 1 {
		c.Status(204)
	} else {
		c.JSON(400, gin.H{"details": "Failed to delete account"})
	}
}
