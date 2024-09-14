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

const JWT_SECRET = "secret"

//go:embed static/*
var staticFS embed.FS

func StartServer(db *sql.DB) {
	r := gin.Default()

	// For frontend development - serve from directory
	// staticSubFS, _ := fs.Sub(staticFS, "static")
	// r.NoRoute(gin.WrapH(http.FileServerFS(staticSubFS)))

	// For production - serve from embedded FS
	r.NoRoute(gin.WrapH(http.FileServer(http.Dir("static"))))

	api := r.Group("/api")

	api.Use(DbMiddleware(db))
	api.Use(AuthRequired())

	api.POST("/login", LoginEndpoint)
	api.GET("/account_otps", GetAccountOTPs)

	r.Run(":9000")
}

func LoginEndpoint(c *gin.Context) {
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

func GetAccountOTPs(c *gin.Context) {
	db := GetDBFromCtx(c)
	var timeStruct time.Time
	tsParam := c.Query("timestamp")
	timestamp, err := strconv.ParseInt(tsParam, 10, 64)
	if err == nil {
		timeStruct = time.Unix(timestamp, 0)
	} else {
		timeStruct = time.Now()
	}
	userId, _ := c.Get("userId")
	accounts := GetAccounts(db, int(userId.(float64)))
	var accountOTPs []AccountOTP
	for _, account := range accounts {
		code, _ := totp.GenerateCode(strings.ReplaceAll(account.Token, " ", ""), timeStruct)
		accountOTP := AccountOTP{account.Id, account.Name, code}
		accountOTPs = append(accountOTPs, accountOTP)
	}
	c.JSON(200, accountOTPs)
}
