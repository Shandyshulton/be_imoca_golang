package controller

import (
	"be_imoca_golang/model"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// GetHeroHandler
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

// UpdateHeroHandler
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

		// Empty Validation
		if strings.TrimSpace(badge) == "" || strings.TrimSpace(title) == "" ||
			strings.TrimSpace(highlight) == "" || strings.TrimSpace(desc) == "" {
			log.Println("[WARN] Update Hero ditolak: Form tidak lengkap")
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    "400",
				"message": "Badge, Title, Highlight, dan Description wajib diisi!",
			})
			return
		}

		updatedHero := model.Hero{
			ID:             oldHero.ID,
			BadgeText:      badge,
			Title:          title,
			TitleHighlight: highlight,
			Description:    desc,
			Image:          oldHero.Image,
		}

		file, err := c.FormFile("image")
		if err == nil {
			// Validatation file extension
			if !isAllowedExtension(file.Filename) {
				log.Printf("[REJECTED] Format file tidak diizinkan: %s", file.Filename)
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    "400",
					"message": "Format file tidak didukung! Gunakan .jpg, .jpeg, atau .png",
				})
				return
			}

			// Validation Size File (Maks 10MB)
			maxMB := 10
			var maxFileSize int64 = int64(maxMB) * 1024 * 1024
			if file.Size > maxFileSize {
				log.Printf("[REJECTED] File terlalu besar: %d bytes", file.Size)
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    "400",
					"message": fmt.Sprintf("Ukuran file terlalu besar! Maksimal adalah %dMB", maxMB),
				})
				return
			}

			newFileName := fmt.Sprintf("hero-%d%s", time.Now().UnixNano(), filepath.Ext(file.Filename))
			targetPath := filepath.Join("storage", "uploads", "hero", newFileName)

			if err := c.SaveUploadedFile(file, targetPath); err == nil {
				if oldHero.Image != "" {
					os.Remove(filepath.Join("storage", "uploads", "hero", oldHero.Image))
				}
				updatedHero.Image = newFileName
			}
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
