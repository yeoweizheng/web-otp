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
			CreateUser(db, username, string(password))
		} else if input == "e" {
			var id int
			var editOp string
			users := GetUsers(db)
			for _, user := range users {
				fmt.Println(user.id, user.username)
			}
			fmt.Print("Select user ID: ")
			fmt.Scan(&id)
			if !UserIDExists(db, id) {
				fmt.Println("Invalid user ID selected.")
			} else {
				fmt.Println("u > Change username")
				fmt.Println("p > Change password")
				fmt.Println("d > Delete user")
				fmt.Println("m > Back to main menu")
				fmt.Print("Enter selection: ")
				fmt.Scan(&editOp)
				if editOp == "u" {
					var newUsername string
					fmt.Print("Enter new username: ")
					fmt.Scan(&newUsername)
					UpdateUsername(db, id, newUsername)
				} else if editOp == "p" {
					fmt.Print("Enter new password: ")
					newPassword, _ := term.ReadPassword(int(os.Stdin.Fd()))
					UpdatePassword(db, id, string(newPassword))
				} else if editOp == "d" {
					DeleteUser(db, id)
				}
			}
		} else if input == "q" {
			break
		}
	}
}
