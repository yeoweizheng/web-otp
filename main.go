package main

import (
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	Test()
	// createUserSQL := `CREATE TABLE users ("id" integer PRIMARY KEY AUTOINCREMENT, "username" TEXT, "password" TEXT);`
	// statement, _ := db.Prepare(createUserSQL)
	// statement.Exec()
	// defer db.Close()

	// r := gin.Default()
	// r.GET("/ping", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"message": "pong",
	// 	})
	// })
	// r.Run(":9000")
}
