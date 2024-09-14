package main

import (
	"database/sql"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func DbMiddleware(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	}
}

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.Contains(c.Request.URL.String(), "login") {
			c.Next()
		} else {
			auth := c.GetHeader("Authorization")
			token, err := jwt.Parse(auth, func(token *jwt.Token) (interface{}, error) {
				return []byte(JWT_SECRET), nil
			})
			if err == nil {
				claims, _ := token.Claims.(jwt.MapClaims)
				c.Set("userId", claims["userId"])
				c.Next()
			} else {
				c.JSON(401, gin.H{"details": "Unauthorized"})
			}
		}
	}
}
