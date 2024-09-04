package main

import (
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const JWT_SECRET = "secret"

func StartServer(db *sql.DB) {
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	api := r.Group("/api")
	api.POST("/login", LoginEndpoint)

	r.Run(":9000")
}

func LoginEndpoint(c *gin.Context) {
	db := GetDBFromCtx(c)
	var data map[string]interface{}
	c.BindJSON(&data)
	if VerifyUser(db, data["username"].(string), data["password"].(string)) {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256,
			jwt.MapClaims{"username": data["username"]})
		tokenString, err := token.SignedString([]byte(JWT_SECRET))
		if err != nil {
			fmt.Println(err.Error())
		}
		c.JSON(200, gin.H{"token": tokenString})
	} else {
		c.JSON(401, gin.H{"details": "Invalid username / password"})
	}
}
