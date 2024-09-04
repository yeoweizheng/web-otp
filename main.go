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
		} else if args[0] == "start" {
			StartServer(db)
		}
	}
}
