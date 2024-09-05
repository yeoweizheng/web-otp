package main

import (
	"database/sql"
	"io/fs"
	"net/http"

	"embed"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const JWT_SECRET = "secret"

//go:embed static/*
var staticFS embed.FS

func StartServer(db *sql.DB) {
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	r.POST("/api/login", LoginEndpoint)

	staticSubFS, _ := fs.Sub(staticFS, "static")
	r.StaticFS("/", http.FS(staticSubFS))

	r.Run(":9000")
}

func LoginEndpoint(c *gin.Context) {
	db := GetDBFromCtx(c)
	var data map[string]interface{}
	c.BindJSON(&data)
	if VerifyUser(db, data["username"].(string), data["password"].(string)) {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256,
			jwt.MapClaims{"username": data["username"]})
		tokenString, _ := token.SignedString([]byte(JWT_SECRET))
		c.JSON(200, gin.H{"token": tokenString})
	} else {
		c.JSON(401, gin.H{"details": "Invalid username / password"})
	}
}
