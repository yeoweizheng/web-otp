package main

import (
	"database/sql"
	"fmt"
	"os"

	"golang.org/x/term"
)

func ManageUsers(db *sql.DB) {
	var input string
	for {
		fmt.Println("a> Add user")
		fmt.Println("e> Edit user")
		fmt.Println("q> Quit")
		fmt.Print("Enter selection: ")
		fmt.Scan(&input)
		if input == "a" {
			var username string
			fmt.Print("Username: ")
			fmt.Scan(&username)
			fmt.Print("Password: ")
			password, _ := term.ReadPassword(int(os.Stdin.Fd()))
			fmt.Print(username, string(password))
			CreateUser(db, username, string(password))
		} else if input == "e" {

		} else if input == "q" {
			break
		}
	}
}
