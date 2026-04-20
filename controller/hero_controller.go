package controller

import (
	"be_imoca_golang/model"
	"be_imoca_golang/util" 
	"database/sql"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Get Hero
func GetHeroHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		hero, err := model.GetHero(db)
		if err != nil {
			log.Printf("[ERROR] Gagal ambil hero: %v", err)
			c.JSON(http.StatusNotFound, gin.H{"code": "404", "message": "Data Hero belum diatur"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": "200", "data": hero})
	}
}

// Update Hero
func UpdateHeroHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("[INFO] Mencoba memperbarui Hero Section...")

		oldHero, err := model.GetHero(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": "500", "message": "Inisialisasi data hero gagal"})
			return
		}

		badge := c.PostForm("badge_text")
		title := c.PostForm("title")
		highlight := c.PostForm("title_highlight")
		desc := c.PostForm("description")

		if strings.TrimSpace(badge) == "" || strings.TrimSpace(title) == "" ||
			strings.TrimSpace(highlight) == "" || strings.TrimSpace(desc) == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    "400",
				"message": "Badge, Title, Highlight, dan Description wajib diisi!",
			})
			return
		}

		newImageName, err := util.HandleFileUpload(c, "image", "hero", oldHero.Image)
		if err != nil {
			log.Printf("[ERROR] Gagal proses gambar hero: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    "400",
				"message": err.Error(), // Pesan "10MB" atau "Format salah" akan muncul di sini
			})
			return
		}

		updatedHero := model.Hero{
			ID:             oldHero.ID,
			BadgeText:      badge,
			Title:          title,
			TitleHighlight: highlight,
			Description:    desc,
			Image:          newImageName, // Nama file hasil olahan helper
		}

		if err := model.UpdateHero(db, &updatedHero); err != nil {
			log.Printf("[DATABASE ERROR] %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"code": "500", "message": "Gagal update database"})
			return
		}

		log.Printf("[SUCCESS] Hero Section berhasil diperbarui")
		c.JSON(http.StatusOK, gin.H{
			"code":    "200",
			"message": "Hero Section berhasil diperbarui!",
			"data":    updatedHero,
		})
	}
}