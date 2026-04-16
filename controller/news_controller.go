package controller

import (
	"be_imoca_golang/model"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
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
			log.Printf("[ERROR] Gagal ambil data berita: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"code": "500", "message": "Gagal memuat berita"})
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
		summary := c.PostForm("summary")
		source := c.PostForm("source")
		badge := c.PostForm("badge") // AMBIL BADGE DARI FORM
		url := c.PostForm("url")

		// Validasi: Badge juga harus divalidasi
		if strings.TrimSpace(title) == "" || strings.TrimSpace(summary) == "" ||
			strings.TrimSpace(source) == "" || strings.TrimSpace(badge) == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    "400",
				"message": "Semua field (Judul, Ringkasan, Sumber, Badge) wajib diisi!",
			})
			return
		}

		file, err := c.FormFile("image")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "400", "message": "Gambar berita wajib diunggah!"})
			return
		}

		// Upload Logic
		extension := filepath.Ext(file.Filename)
		newFileName := fmt.Sprintf("news-%d%s", time.Now().UnixNano(), extension)
		targetPath := filepath.Join("storage", "uploads", "news", newFileName)

		if err := c.SaveUploadedFile(file, targetPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": "500", "message": "Gagal simpan gambar"})
			return
		}

		n := model.News{
			Title:   title,
			Summary: summary,
			Source:  source,
			Badge:   badge, // MASUKKAN KE MODEL
			URL:     url,
			Image:   newFileName,
		}

		if err := model.CreateNews(db, &n); err != nil {
			log.Printf("[DB ERROR]: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"code": "500", "message": "Gagal simpan ke database"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"code":    "201",
			"message": "Berita berhasil diterbitkan!",
			"data":    n,
		})
	}
}

// UpdateNewsHandler
func UpdateNewsHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var uri NewsUri
		if err := c.ShouldBindUri(&uri); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "400", "message": "ID tidak valid"})
			return
		}

		oldNews, err := model.GetNewsByID(db, uri.ID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"code": "404", "message": "Berita tidak ditemukan"})
			return
		}

		title := c.PostForm("title")
		summary := c.PostForm("summary")
		source := c.PostForm("source")
		badge := c.PostForm("badge") // AMBIL BADGE DARI FORM
		url := c.PostForm("url")

		if strings.TrimSpace(title) == "" || strings.TrimSpace(badge) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"code": "400", "message": "Judul dan Badge tidak boleh kosong!"})
			return
		}

		updatedNews := model.News{
			ID:      uri.ID,
			Title:   title,
			Summary: summary,
			Source:  source,
			Badge:   badge, // MASUKKAN KE MODEL
			URL:     url,
			Image:   oldNews.Image,
		}

		// Handle Image Update (Sama seperti sebelumnya)
		file, err := c.FormFile("image")
		if err == nil {
			extension := filepath.Ext(file.Filename)
			newFileName := fmt.Sprintf("news-%d%s", time.Now().UnixNano(), extension)
			targetPath := filepath.Join("storage", "uploads", "news", newFileName)

			if err := c.SaveUploadedFile(file, targetPath); err == nil {
				os.Remove(filepath.Join("storage", "uploads", "news", oldNews.Image))
				updatedNews.Image = newFileName
			}
		}

		if err := model.UpdateNews(db, &updatedNews); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": "500", "message": "Gagal update database"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    "200",
			"message": "Berita berhasil diperbarui!",
			"data":    updatedNews,
		})
	}
}

// Handler lainnya (Delete, GetByID) tetap sama karena sudah menggunakan ID
func DeleteNewsHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		if err := model.DeleteNews(db, id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": "500", "message": "Gagal menghapus berita"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": "200", "message": "Berita berhasil dihapus"})
	}
}

func GetNewsByIDHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var uri NewsUri
		if err := c.ShouldBindUri(&uri); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "400", "message": "ID tidak valid"})
			return
		}
		news, err := model.GetNewsByID(db, uri.ID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"code": "404", "message": "Berita tidak ditemukan"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": "200", "data": news})
	}
}