package controllers

import (
	"be_imoca_golang/model"
	"database/sql"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthRequest struct {
	FullName string `json:"full_name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Register Handler
func RegisterHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req AuthRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    "400",
				"message": "Format data tidak sesuai",
			})
			return
		}

		// Validasi data kosong
		if req.FullName == "" || req.Username == "" || req.Password == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    "400",
				"message": "Nama lengkap, username, dan password wajib diisi!",
			})
			return
		}

		if err := model.CreateUser(db, req.FullName, req.Username, req.Password); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    "500",
				"message": "Gagal buat user: " + err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"code":    "201",
			"message": "User admin berhasil dibuat!",
		})
	}
}

// LoginHandler 
func LoginHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req AuthRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    "400",
				"message": "Format request tidak valid",
			})
			return
		}

		u, err := model.GetUserByUsername(db, req.Username)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    "401",
				"message": "Username atau password salah",
			})
			return
		}

		// Validasi Password
		err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    "401",
				"message": "Username atau password salah",
			})
			return
		}

		// Generate JWT Token
		secret := os.Getenv("JWT_SECRET")
		claims := jwt.MapClaims{
			"user_id":   u.ID,
			"username":  u.Username,
			"full_name": u.FullName, 
			"exp":       time.Now().Add(time.Hour * 24).Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(secret))

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    "500",
				"message": "Gagal generate token login",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    "200",
			"message": "Selamat datang, " + u.FullName + "!",
			"token":   tokenString,
			"user": gin.H{
				"id":        u.ID,
				"full_name": u.FullName,
				"username":  u.Username,
			},
		})
	}
}