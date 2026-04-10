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

type PartnerUri struct {
	ID int `uri:"id" binding:"required"`
}

// isAllowedExtension 
func isAllowedExtension(fileName string) bool {
	extensions := []string{".jpg", ".jpeg", ".png"}
	ext := strings.ToLower(filepath.Ext(fileName))
	for _, e := range extensions {
		if e == ext {
			return true
		}
	}
	return false
}

// GetPartnersHandler
func GetPartnersHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		keyword := c.Query("name")
		partners, err := model.GetAllPartners(db, keyword)
		if err != nil {
			log.Printf("[ERROR] Gagal ambil data partner: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    "500",
				"message": "Gagal memuat daftar mitra: " + err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    "200",
			"message": "Daftar mitra berhasil dimuat",
			"data":    partners,
		})
	}
}

func CreatePartnerHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("[INFO] Mencoba mendaftarkan mitra baru...")

		name := c.PostForm("name")
		email := c.PostForm("email")
		phone := c.PostForm("phone")
		description := c.PostForm("description")
		link := c.PostForm("link")

		if strings.TrimSpace(name) == "" || strings.TrimSpace(email) == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    "400",
				"message": "Nama dan Email wajib diisi!",
			})
			return
		}

		file, err := c.FormFile("image")
		var fileName string
		if err == nil {
			// Validation Format File
			if !isAllowedExtension(file.Filename) {
				log.Printf("[REJECTED] Format file tidak diizinkan: %s", file.Filename)
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    "400",
					"message": "Format file tidak didukung! Hanya diperbolehkan: .jpg, .jpeg, dan .png",
				})
				return
			}

			// Validation Size File (Max 10MB)
			maxMB := 10
			var maxFileSize int64 = int64(maxMB) * 1024 * 1024
			if file.Size > maxFileSize {
				log.Printf("[REJECTED] File terlalu besar: %s (%d bytes)", file.Filename, file.Size)
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    "400",
					"message": fmt.Sprintf("Ukuran file terlalu besar! Maksimal diperbolehkan adalah %dMB", maxMB),
				})
				return
			}

			// Save File
			extension := filepath.Ext(file.Filename)
			fileName = fmt.Sprintf("partner-%d%s", time.Now().UnixNano(), extension)
			targetPath := filepath.Join("storage", "uploads", "partners", fileName)

			if err := c.SaveUploadedFile(file, targetPath); err != nil {
				log.Printf("[ERROR] Gagal simpan gambar mitra: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    "500",
					"message": "Gagal menyimpan gambar di server",
				})
				return
			}
		}

		// Save DB
		p := model.Partner{
			Name:        name,
			Email:       email,
			Phone:       phone,
			Description: description,
			Link:        link,
			Image:       fileName,
		}

		if err := model.CreatePartner(db, &p); err != nil {
			log.Printf("[DATABASE ERROR] Gagal create partner: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    "500",
				"message": "Gagal simpan ke database: " + err.Error(),
			})
			return
		}

		log.Printf("[SUCCESS] Mitra '%s' berhasil didaftarkan dengan ID %d", p.Name, p.ID)
		c.JSON(http.StatusCreated, gin.H{
			"code":    "201",
			"message": "Mitra berhasil didaftarkan!",
			"data":    p,
		})
	}
}

// UpdatePartnerHandler 
func UpdatePartnerHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("[INFO] Memproses update data mitra...")

		var uri PartnerUri
		if err := c.ShouldBindUri(&uri); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    "400",
				"message": "ID mitra tidak valid",
			})
			return
		}

		oldPartner, err := model.GetPartnerByID(db, uri.ID)
		if err != nil {
			log.Printf("[WARN] Update gagal, mitra ID %d tidak ditemukan", uri.ID)
			c.JSON(http.StatusNotFound, gin.H{
				"code":    "404",
				"message": "Mitra tidak ditemukan",
			})
			return
		}

		name := c.PostForm("name")
		email := c.PostForm("email")
		phone := c.PostForm("phone")
		description := c.PostForm("description")
		link := c.PostForm("link")

		p := model.Partner{
			ID:          uri.ID,
			Name:        name,
			Email:       email,
			Phone:       phone,
			Description: description,
			Link:        link,
			Image:       oldPartner.Image, 
		}

		file, err := c.FormFile("image")
		if err == nil {
			// Validation Format File
			if !isAllowedExtension(file.Filename) {
				log.Printf("[REJECTED] Update gambar ditolak, format salah: %s", file.Filename)
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    "400",
					"message": "Format file tidak didukung! Hanya diperbolehkan: .jpg, .jpeg, dan .png",
				})
				return
			}

			// Validation Size File (Max 10MB)
			maxMB := 10
			var maxFileSize int64 = int64(maxMB) * 1024 * 1024
			if file.Size > maxFileSize {
				log.Printf("[REJECTED] Update gagal, file %s terlalu besar (%d bytes)", file.Filename, file.Size)
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    "400",
					"message": fmt.Sprintf("Ukuran file terlalu besar! Maksimal diperbolehkan adalah %dMB", maxMB),
				})
				return
			}

			// Save New File
			extension := filepath.Ext(file.Filename)
			newFileName := fmt.Sprintf("partner-%d%s", time.Now().UnixNano(), extension)
			targetPath := filepath.Join("storage", "uploads", "partners", newFileName)

			if err := c.SaveUploadedFile(file, targetPath); err != nil {
				log.Printf("[ERROR] Gagal simpan gambar baru: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    "500",
					"message": "Gagal menyimpan gambar baru di server",
				})
				return
			}

			// Remove File If Successfully Uploaded
			if oldPartner.Image != "" {
				oldFilePath := filepath.Join("storage", "uploads", "partners", oldPartner.Image)
				os.Remove(oldFilePath)
				log.Printf("[INFO] File lama dihapus: %s", oldPartner.Image)
			}
			p.Image = newFileName
		}

		if err := model.UpdatePartner(db, &p); err != nil {
			log.Printf("[DATABASE ERROR] Gagal update partner ID %d: %v", p.ID, err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    "500",
				"message": "Gagal memperbarui data: " + err.Error(),
			})
			return
		}

		log.Printf("[SUCCESS] Data mitra ID %d (%s) berhasil diperbarui", p.ID, p.Name)
		c.JSON(http.StatusOK, gin.H{
			"code":    "200",
			"message": "Data mitra berhasil diperbarui!",
			"data":    p,
		})
	}
}

// DeletePartnerHandler 
func DeletePartnerHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")

		if err := model.DeletePartner(db, idStr); err != nil {
			log.Printf("[ERROR] Gagal hapus mitra ID %s: %v", idStr, err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    "500",
				"message": "Gagal menghapus mitra: " + err.Error(),
			})
			return
		}

		log.Printf("[SUCCESS] Mitra ID %s berhasil dihapus", idStr)
		c.JSON(http.StatusOK, gin.H{
			"code":    "200",
			"message": "Mitra berhasil dihapus",
		})
	}
}

func GetPartnerByIDHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("[INFO] Mengambil detail mitra...")

		var uri PartnerUri
		if err := c.ShouldBindUri(&uri); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    "400",
				"message": "ID mitra tidak valid",
			})
			return
		}

		partner, err := model.GetPartnerByID(db, uri.ID)
		if err != nil {
			log.Printf("[WARN] Mitra ID %d tidak ditemukan", uri.ID)
			c.JSON(http.StatusNotFound, gin.H{
				"code":    "404",
				"message": "Mitra tidak ditemukan",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    "200",
			"message": "Detail mitra berhasil dimuat",
			"data":    partner,
		})
	}
}