package config

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql" // Tambahkan ini
	"gorm.io/gorm"         // Tambahkan ini
)

// ConnectDB tetap untuk modul lama (SQL Manual)
func ConnectDB() *sql.DB {
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	if dbHost == "" { dbHost = "127.0.0.1" }
	if dbPort == "" { dbPort = "3306" }

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

// ConnectGORM untuk modul baru (Organization)
func ConnectGORM() *gorm.DB {
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	if dbHost == "" { dbHost = "127.0.0.1" }
	if dbPort == "" { dbPort = "3306" }

	// DSN GORM membutuhkan parseTime=True agar CreatedAt/UpdatedAt berfungsi
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", 
		dbUser, dbPass, dbHost, dbPort, dbName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("Gagal koneksi GORM ke database %s: %v", dbName, err))
	}

	return db
}