package model

import (
	"database/sql"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int    `json:"id"`
	FullName  string `json:"full_name"`
	Username  string `json:"username"`
	Password  string `json:"-"` 
	CreatedAt string `json:"created_at"`
}

// CreateUser
func CreateUser(db *sql.DB, fullName, username, password string) error {
	// Hash Password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Insert to DB
	query := "INSERT INTO users (full_name, username, password) VALUES (?, ?, ?)"
	_, err = db.Exec(query, fullName, username, hashedPassword)
	return err
}

// GetUserByUsername 
func GetUserByUsername(db *sql.DB, username string) (*User, error) {
	var u User
	query := "SELECT id, full_name, username, password, created_at FROM users WHERE username = ?"
	
	err := db.QueryRow(query, username).Scan(
		&u.ID, 
		&u.FullName, 
		&u.Username, 
		&u.Password, 
		&u.CreatedAt,
	)
	
	if err != nil {
		return nil, err
	}
	return &u, nil
}