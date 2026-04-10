package controller

import (
	"be_imoca_golang/model"
	"database/sql"
	"fmt"
	"log"
	"os"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type NewsUri struct {
	ID int `uri:"id" binding:"required"`
}

// GetNewsHandler 
func GetNewsHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		keyword := c.Query("title")
		news, err := model.GetAllNews(db, keyword)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    "500",
				"message": "Gagal memuat daftar berita: " + err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    "200",
			"message": "Data berita berhasil diambil",
			"data":    news,
		})
	}
}

// CreateNewsHandler 
func CreateNewsHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		title := c.PostForm("title")
		content := c.PostForm("content")

		if strings.TrimSpace(title) == "" || strings.TrimSpace(content) == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    "400",
				"message": "Judul dan Konten wajib diisi!",
			})
			return
		}

		file, err := c.FormFile("image")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    "400",
				"message": "Gambar berita wajib diunggah!",
			})
			return
		}

		// Validation Format File
		if !isAllowedExtension(file.Filename) {
			log.Printf("[REJECTED] Format file berita ditolak: %s", file.Filename)
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
			log.Printf("[REJECTED] File berita terlalu besar: %s (%d bytes)", file.Filename, file.Size)
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    "400",
				"message": fmt.Sprintf("Ukuran gambar berita terlalu besar! Maksimal adalah %dMB", maxMB),
			})
			return
		}

		// Save File and Path
		extension := filepath.Ext(file.Filename)
		newFileName := fmt.Sprintf("news-%d%s", time.Now().UnixNano(), extension)
		targetPath := filepath.Join("storage", "uploads", "news", newFileName)

		if err := c.SaveUploadedFile(file, targetPath); err != nil {
			log.Printf("[ERROR] Gagal simpan file berita: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    "500",
				"message": "Gagal menyimpan gambar di server",
			})
			return
		}

		// Save DB
		n := model.News{
			Title:   title,
			Content: content,
			Image:   newFileName, 
		}

		if err := model.CreateNews(db, &n); err != nil {
			log.Printf("[DATABASE ERROR] %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    "500",
				"message": "Gagal simpan ke database: " + err.Error(),
			})
			return
		}

		log.Printf("[SUCCESS] Berita '%s' berhasil diterbitkan", n.Title)
		c.JSON(http.StatusCreated, gin.H{
			"code":    "201",
			"message": "Berita berhasil diterbitkan!",
			"data":    n,
		})
	}
}

func UpdateNewsHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("[INFO] Memproses update data berita...")

		var uri NewsUri
		if err := c.ShouldBindUri(&uri); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "400", "message": "ID tidak valid"})
			return
		}

		oldNews, err := model.GetNewsByID(db, uri.ID)
		if err != nil {
			log.Printf("[WARN] Update gagal, berita ID %d tidak ditemukan", uri.ID)
			c.JSON(http.StatusNotFound, gin.H{"code": "404", "message": "Berita tidak ditemukan"})
			return
		}

		title := c.PostForm("title")
		content := c.PostForm("content")
		
		updatedNews := model.News{
			ID:      uri.ID,
			Title:   title,
			Content: content,
			Image:   oldNews.Image, 
		}

		file, err := c.FormFile("image")
		if err == nil { 
			// A. VALIDASI FORMAT FILE
			if !isAllowedExtension(file.Filename) {
				log.Printf("[REJECTED] Update berita ditolak, format salah: %s", file.Filename)
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
				log.Printf("[REJECTED] Update berita gagal, file terlalu besar: %d bytes", file.Size)
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    "400",
					"message": fmt.Sprintf("Ukuran gambar berita terlalu besar! Maksimal adalah %dMB", maxMB),
				})
				return
			}

			// Save New File
			extension := filepath.Ext(file.Filename)
			newFileName := fmt.Sprintf("news-%d%s", time.Now().UnixNano(), extension)
			targetPath := filepath.Join("storage", "uploads", "news", newFileName)

			if err := c.SaveUploadedFile(file, targetPath); err != nil {
				log.Printf("[ERROR] Gagal simpan gambar berita baru: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"code": "500", "message": "Gagal simpan gambar baru"})
				return
			}

			// Remove File If Successfully Uploaded
			if oldNews.Image != "" {
				oldFilePath := filepath.Join("storage", "uploads", "news", oldNews.Image)
				if err := os.Remove(oldFilePath); err != nil {
					log.Printf("[WARN] Gagal menghapus file lama: %v", err)
				} else {
					log.Printf("[INFO] File lama berhasil dihapus: %s", oldNews.Image)
				}
			}

			updatedNews.Image = newFileName
		}

		// Save DB
		if err := model.UpdateNews(db, &updatedNews); err != nil {
			log.Printf("[DATABASE ERROR] %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"code": "500", "message": "Gagal update database: " + err.Error()})
			return
		}

		log.Printf("[SUCCESS] Berita ID %d berhasil diperbarui", updatedNews.ID)
		c.JSON(http.StatusOK, gin.H{
			"code":    "200",
			"message": "Berita berhasil diperbarui!",
			"data":    updatedNews,
		})
	}
}

// DeleteNewsHandler 
func DeleteNewsHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if err := model.DeleteNews(db, id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    "500",
				"message": "Gagal menghapus berita: " + err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    "200",
			"message": "Berita berhasil dihapus",
		})
	}
}

// GetNewsByIDHandler untuk melihat detail satu berita
func GetNewsByIDHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var uri NewsUri
		if err := c.ShouldBindUri(&uri); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "400", "message": "ID berita tidak valid"})
			return
		}

		news, err := model.GetNewsByID(db, uri.ID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"code": "404", "message": "Berita tidak ditemukan"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    "200",
			"message": "Detail berita berhasil diambil",
			"data":    news,
		})
	}
}