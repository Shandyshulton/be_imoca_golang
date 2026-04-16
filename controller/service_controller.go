package controller

import (
	"be_imoca_golang/model" 
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UpdateServicesHandler
func UpdateServicesHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input model.ServicesSectionResponse
		
		if err := c.ShouldBindJSON(&input); err != nil {
			log.Printf("[ERROR] Binding Services Gagal: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    "400", 
				"message": "Format data tidak sesuai, pastikan semua bagian terisi dengan benar",
			})
			return
		}

		log.Printf("[DEBUG] Update Services - Title: %s, Items: %d", 
			input.SectionTitle, len(input.Items))

		if err := model.UpdateServices(db, &input); err != nil {
			log.Printf("[DB ERROR] Update Data Services Gagal: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    "500", 
				"message": "Terjadi kesalahan sistem saat menyimpan perubahan layanan",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    "200", 
			"message": "Data layanan dan detail informasi berhasil diperbarui!",
		})
	}
}

// GetServicesHandler 
func GetServicesHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		data, err := model.GetServices(db)
		if err != nil {
			log.Printf("[ERROR] Ambil Data Services Gagal: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    "500", 
				"message": "Gagal mengambil data layanan dari server",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code": "200", 
			"data": data,
		})
	}
}