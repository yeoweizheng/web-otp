package main

import (
	"database/sql"
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
	// staticSubFS, _ := fs.Sub(staticFS, "static")
	// r.NoRoute(gin.WrapH(http.FileServerFS(staticSubFS)))

	// For production - serve from embedded FS
	r.NoRoute(gin.WrapH(http.FileServer(http.Dir("static"))))

	api := r.Group("/api")

	api.Use(DbMiddleware(db))
	api.Use(AuthRequired())

	api.POST("/login/", LoginAPI)
	api.GET("/account_otps/", GetAccountOTPsAPI)
	api.POST("/add_account/", AddAccountAPI)
	api.PATCH("/update_account/:id/", UpdateAccountAPI)
	api.DELETE("/delete_account/:id/", DeleteAccountAPI)

	r.Run(HOST_PORT)
}

func LoginAPI(c *gin.Context) {
	db := GetDBFromCtx(c)
	var data map[string]interface{}
	c.BindJSON(&data)
	userId, err := VerifyAndGetUserId(db, data["username"].(string), data["password"].(string))
	if err == nil {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256,
			jwt.MapClaims{"userId": userId, "timestamp": time.Now().Unix()})
		tokenString, _ := token.SignedString([]byte(JWT_SECRET))
		c.JSON(200, gin.H{"token": tokenString})
	} else {
		c.JSON(401, gin.H{"details": "Invalid username / password"})
	}
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
		accountOTP := AccountOTP{account.Id, account.Name, account.Token, code}
		accountOTPs = append(accountOTPs, accountOTP)
	}
	c.JSON(200, accountOTPs)
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
	rowCount := UpdateAccount(db, userId, accountId, data["name"].(string), data["token"].(string))
	if rowCount == 1 {
		c.JSON(200, Account{accountId, data["name"].(string), data["token"].(string)})
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
