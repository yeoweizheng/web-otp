package main

import (
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("Usage: web-otp [initdb / manageusers / start]")
		return
	} else {
		db := OpenDB("db.sqlite3")
		defer db.Close()
		if args[0] == "initdb" {
			InitDB(db)
		} else if args[0] == "manageusers" {
			ManageUsers(db)
		}
	}
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
