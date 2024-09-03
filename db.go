package main

import (
	"database/sql"
	"fmt"
	"log"
)

func OpenDB(filename string) (db *sql.DB) {
	db, err := sql.Open("sqlite3", "db.sqlite3")
	if err != nil {
		log.Fatal(err.Error())
	}
	return
}

func InitDB(db *sql.DB) {
	userStmt, _ := db.Prepare(`CREATE TABLE IF NOT EXISTS users ("id" integer PRIMARY KEY AUTOINCREMENT, "username" TEXT, "password" TEXT);`)
	accountStmt, _ := db.Prepare(`CREATE TABLE IF NOT EXISTS accounts ("id" integer PRIMARY KEY AUTOINCREMENT, "name" TEXT, "token" TEXT, "userId" integer REFERENCES "user" ("id"));`)
	userStmt.Exec()
	accountStmt.Exec()
	fmt.Println("Database initialized.")
}

func UserExists(db *sql.DB, username string) {

}

func CreateUser(db *sql.DB, username string, password string) {
	stmt, _ := db.Prepare(`INSERT INTO users (username, password) VALUES (?, ?)`)
	stmt.Exec(username, password)
	fmt.Println("User created.")
}
