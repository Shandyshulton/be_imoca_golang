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
	// 1. LOAD ENV
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: file .env tidak ditemukan")
	}

	util.InitLoggers()

	// Pastikan folder storage/logs sudah ada
	f, _ := os.OpenFile("storage/logs/info.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	// 2. KONEKSI DB
	db := config.ConnectDB()
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Database tidak merespon: ", err)
	}

	// 3. SETUP GIN
	r := gin.Default()

	// 4. CORS - Mengizinkan React mengakses API
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, 
		AllowMethods:     []string{"GET", "POST", "DELETE", "PUT", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// 5. STATIC FILES
	r.Static("/uploads", "./storage/uploads")

	r.MaxMultipartMemory = 12 << 20

	// 6. ROUTES - Dibungkus dengan grup /api
	api := r.Group("/api")
	{
		route.InitRoutes(api, db) // Melemparkan grup /api ke InitRoutes
	}

	// 7. RUN
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("IMOCA Backend berjalan di http://localhost:%s/api\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Gagal menjalankan server: ", err)
	}
}