package main

import (
	"database/sql"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func OpenDB(filename string) *sql.DB {
	db, err := sql.Open("sqlite3", "db.sqlite3")
	if err != nil {
		log.Fatal(err.Error())
	}
	return db
}

func InitDB(db *sql.DB) {
	userStmt, _ := db.Prepare(`CREATE TABLE IF NOT EXISTS users ("id" integer PRIMARY KEY AUTOINCREMENT, "username" TEXT, "password" TEXT);`)
	accountStmt, _ := db.Prepare(`CREATE TABLE IF NOT EXISTS accounts ("id" integer PRIMARY KEY AUTOINCREMENT, "name" TEXT, "token" TEXT, "userId" integer REFERENCES "user" ("id"));`)
	userStmt.Exec()
	accountStmt.Exec()
	fmt.Println("Database initialized.")
}

func UsernameExists(db *sql.DB, username string) bool {
	stmt, _ := db.Prepare(`SELECT COUNT(id) FROM users WHERE username = ?`)
	var count int
	stmt.QueryRow(username).Scan(&count)
	if count > 0 {
		return true
	} else {
		return false
	}
}

func UserIDExists(db *sql.DB, id int) bool {
	stmt, _ := db.Prepare(`SELECT COUNT(id) FROM users WHERE id = ?`)
	var count int
	stmt.QueryRow(id).Scan(&count)
	if count > 0 {
		return true
	} else {
		return false
	}
}

func CreateUser(db *sql.DB, username string, password string) {
	if UsernameExists(db, username) {
		fmt.Println("Username taken. Please choose another username.")
	} else {
		passwordBytes := []byte(password)
		passwordBytes, _ = bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
		stmt, _ := db.Prepare(`INSERT INTO users (username, password) VALUES (?, ?)`)
		stmt.Exec(username, string(passwordBytes))
		fmt.Println("\nUser created.")
	}
}

func GetUsers(db *sql.DB) []User {
	var users []User
	stmt, _ := db.Prepare(`SELECT id, username, password FROM users`)
	rows, _ := stmt.Query()
	for rows.Next() {
		var user User
		rows.Scan(&user.id, &user.username, &user.password)
		users = append(users, user)
	}
	return users
}

func UpdateUsername(db *sql.DB, id int, username string) {
	stmt, _ := db.Prepare(`UPDATE users SET username = ? WHERE id = ?`)
	stmt.Exec(username, id)
	fmt.Println("Username updated.")
}

func UpdatePassword(db *sql.DB, id int, password string) {
	passwordBytes := []byte(password)
	passwordBytes, _ = bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
	stmt, _ := db.Prepare(`UPDATE users SET password = ? WHERE id = ?`)
	stmt.Exec(string(passwordBytes), id)
	fmt.Println("\nPassword updated.")
}

func DeleteUser(db *sql.DB, id int) {
	stmt, _ := db.Prepare(`DELETE FROM users WHERE id = ?`)
	stmt.Exec(id)
	fmt.Println("User deleted.")
}
