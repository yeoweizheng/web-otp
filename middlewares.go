package main

import (
	"database/sql"
	"fmt"

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
		accessToken, err := c.Cookie(ACCESS_TOKEN_COOKIE_NAME)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"details": "Unauthorized"})
			return
		}

		token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("invalid signing method")
			}
			return []byte(JWT_SECRET), nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(401, gin.H{"details": "Unauthorized"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(401, gin.H{"details": "Unauthorized"})
			return
		}

		tokenType, ok := claims["type"].(string)
		if !ok || tokenType != "access" {
			c.AbortWithStatusJSON(401, gin.H{"details": "Unauthorized"})
			return
		}

		userId, ok := claims["userId"].(float64)
		if !ok {
			c.AbortWithStatusJSON(401, gin.H{"details": "Unauthorized"})
			return
		}

		c.Set("userId", userId)
		c.Next()
	}
}
