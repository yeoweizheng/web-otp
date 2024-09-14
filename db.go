package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
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
	accountStmt, _ := db.Prepare(`CREATE TABLE IF NOT EXISTS accounts ("id" integer PRIMARY KEY AUTOINCREMENT, "name" TEXT, "token" TEXT, "userId" integer REFERENCES "users" ("id"));`)
	userStmt.Exec()
	accountStmt.Exec()
	fmt.Println("Database initialized.")
}

func GetDBFromCtx(c *gin.Context) *sql.DB {
	ctxDb, _ := c.Get("db")
	return ctxDb.(*sql.DB)
}

func GetUserIdFromCtx(c *gin.Context) int {
	userId, _ := c.Get("userId")
	return int(userId.(float64))
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
		rows.Scan(&user.Id, &user.Username, &user.Password)
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

func VerifyAndGetUserId(db *sql.DB, username string, password string) (int, error) {
	stmt, _ := db.Prepare(`SELECT id, password FROM users WHERE username = ?`)
	var id int
	var hash string
	stmt.QueryRow(username).Scan(&id, &hash)
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err == nil {
		return id, nil
	} else {
		return -1, fmt.Errorf("failed to verify user")
	}
}

func GetAccounts(db *sql.DB, userId int) []Account {
	var accounts []Account
	stmt, _ := db.Prepare(`SELECT id, name, token FROM accounts WHERE userId = ?`)
	rows, _ := stmt.Query(userId)
	for rows.Next() {
		var account Account
		rows.Scan(&account.Id, &account.Name, &account.Token)
		accounts = append(accounts, account)
	}
	return accounts
}

func CreateAccount(db *sql.DB, userId int, name string, token string) int64 {
	stmt, _ := db.Prepare(`INSERT INTO accounts (userId, name, token) VALUES (?, ?, ?)`)
	result, _ := stmt.Exec(userId, name, token)
	lastInsertId, _ := result.LastInsertId()
	return lastInsertId
}

func UpdateAccount(db *sql.DB, userId int, accountId int64, name string, token string) int64 {
	stmt, _ := db.Prepare(`UPDATE accounts SET name = ?, token = ? WHERE userId = ? AND id = ?`)
	result, _ := stmt.Exec(name, token, userId, accountId)
	rowCount, _ := result.RowsAffected()
	return rowCount
}

func DeleteAccount(db *sql.DB, userId int, accountId int64) int64 {
	stmt, _ := db.Prepare(`DELETE FROM accounts WHERE userId = ? AND id = ?`)
	result, _ := stmt.Exec(userId, accountId)
	rowCount, _ := result.RowsAffected()
	return rowCount
}
