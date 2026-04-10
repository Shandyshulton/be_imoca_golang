package config

import (
	"database/sql"
	"fmt"
	"os" // Tambahkan ini untuk membaca environment variable

	_ "github.com/go-sql-driver/mysql"
)

func ConnectDB() *sql.DB {
	// Ambil data dari .env menggunakan os.Getenv
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	if dbHost == "" { dbHost = "127.0.0.1" }
	if dbPort == "" { dbPort = "3306" }

	// Susun DSN menggunakan variabel dari .env
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", 
		dbUser, dbPass, dbHost, dbPort, dbName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(fmt.Sprintf("Database %s tidak ditemukan atau MySQL belum jalan: %v", dbName, err))
	}

	return db
}