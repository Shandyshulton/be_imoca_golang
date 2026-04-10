package main

import (
	"be_imoca_golang/config"
	"be_imoca_golang/route"
	"be_imoca_golang/util"
	"io"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// LOAD ENV
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: file .env tidak ditemukan")
	}

	util.InitLoggers()

	f, _ := os.OpenFile("storage/logs/info.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	// KONEKSI DB
	db := config.ConnectDB()
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Database tidak merespon: ", err)
	}

	// SETUP GIN
	r := gin.Default()

	// CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins: 	  []string{"*"},
		AllowMethods:     []string{"GET", "POST", "DELETE", "PUT", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	r.Static("/uploads", "./storage/uploads")

	r.MaxMultipartMemory = 12 << 20

	// ROUTES
	route.InitRoutes(r, db)

	// RUN
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("IMOCA Backend berjalan di http://localhost:%s\n", port)
	if err := r.Run("0.0.0.0:" + port); err != nil {
		log.Fatal("Gagal menjalankan server: ", err)
	}
}