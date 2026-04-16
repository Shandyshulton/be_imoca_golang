package controller

import (
	"be_imoca_golang/model"
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetVisionHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		data, err := model.GetVision(db)
		if err != nil {
			log.Printf("[ERROR] Get Vision Gagal: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    "500", 
				"message": "Data Visi gagal dimuat dari server", // Lebih natural
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code": "200", 
			"data": data,
		})
	}
}

func UpdateVisionHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input model.VisionSectionResponse
		
		if err := c.ShouldBindJSON(&input); err != nil {
			log.Printf("[ERROR] Binding Vision Gagal: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    "400", 
				"message": "Format data tidak sesuai, pastikan semua kolom terisi dengan benar",
			})
			return
		}

		if err := model.UpdateVision(db, &input); err != nil {
			log.Printf("[DB ERROR] Update Vision Gagal: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    "500", 
				"message": "Terjadi kesalahan sistem saat menyimpan perubahan Visi",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    "200", 
			"message": "Poin-poin Visi berhasil diperbarui!", // Jauh lebih enak dibaca
		})
	}
}