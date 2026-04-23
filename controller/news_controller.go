package controller

import (
	"be_imoca_golang/model"
	"be_imoca_golang/util"
	"database/sql"
	"time"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type NewsUri struct {
	ID int `uri:"id" binding:"required"`
}

// Get News
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

// Create News
func CreateNewsHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		title := c.PostForm("title")
		summary := c.PostForm("summary")
		source := c.PostForm("source")
		badge := c.PostForm("badge")
		url := c.PostForm("url")

		if strings.TrimSpace(title) == "" || strings.TrimSpace(summary) == "" ||
			strings.TrimSpace(source) == "" || strings.TrimSpace(badge) == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    "400",
				"message": "Semua field (Judul, Ringkasan, Sumber, Badge) wajib diisi!",
			})
			return
		}

		newImageName, err := util.HandleFileUpload(c, "image", "news", "")
		if err != nil {
			log.Printf("[ERROR] Gagal upload gambar berita: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    "400",
				"message": err.Error(), // Pesan "ukuran file terlalu besar..." muncul di sini
			})
			return
		}

		if newImageName == "" {
			c.JSON(http.StatusBadRequest, gin.H{"code": "400", "message": "Gambar berita wajib diunggah!"})
			return
		}

		n := model.News{
			Title:   title,
			Summary: summary,
			Source:  source,
			Badge:   badge,
			URL:     url,
			Image:   newImageName,
		}

		// Baca published_at dari form, fallback ke hari ini kalau kosong
		publishedAtStr := c.PostForm("published_at")
		if publishedAtStr != "" {
			if parsed, err := time.Parse("2006-01-02", publishedAtStr); err == nil {
				n.PublishedAt = parsed
			} else {
				n.PublishedAt = time.Now()
			}
		} else {
			n.PublishedAt = time.Now()
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

// Update News
func UpdateNewsHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
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
		summary := c.PostForm("summary")
		source := c.PostForm("source")
		badge := c.PostForm("badge")
		url := c.PostForm("url")

		if strings.TrimSpace(title) == "" || strings.TrimSpace(badge) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"code": "400", "message": "Judul dan Badge tidak boleh kosong!"})
			return
		}

		newImageName, err := util.HandleFileUpload(c, "image", "news", oldNews.Image)
		if err != nil {
			log.Printf("[ERROR] Gagal proses gambar berita: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    "400",
				"message": err.Error(),
			})
			return
		}

		publishedAtStr := c.PostForm("published_at")
		publishedAt := oldNews.PublishedAt // fallback ke date lama kalau kosong
		if publishedAtStr != "" {
			if parsed, err := time.Parse("2006-01-02", publishedAtStr); err == nil {
				publishedAt = parsed
			}
		}

		updatedNews := model.News{
			ID:          uri.ID,
			Title:       title,
			Summary:     summary,
			Source:      source,
			Badge:       badge,
			URL:         url,
			Image:       newImageName,
			PublishedAt: publishedAt, // ← tambahkan ini
		}

		if err := model.UpdateNews(db, &updatedNews); err != nil {
			log.Printf("[DATABASE ERROR] Gagal update berita ID %d: %v", uri.ID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"code": "500", "message": "Gagal update database"})
			return
		}

		log.Printf("[SUCCESS] Berita ID %d berhasil diperbarui", uri.ID)
		c.JSON(http.StatusOK, gin.H{
			"code":    "200",
			"message": "Berita berhasil diperbarui!",
			"data":    updatedNews,
		})
	}
}

// Delete News
func DeleteNewsHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "400", "message": "ID tidak valid"})
			return
		}

		news, err := model.GetNewsByID(db, id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"code": "404", "message": "Berita tidak ditemukan"})
			return
		}

		if err := model.DeleteNews(db, id); err != nil {
			log.Printf("[DB ERROR] Gagal hapus berita ID %d: %v", id, err)
			c.JSON(http.StatusInternalServerError, gin.H{"code": "500", "message": "Gagal menghapus berita dari database"})
			return
		}

		util.DeleteFile("news", news.Image)

		c.JSON(http.StatusOK, gin.H{
			"code":    "200",
			"message": "Berita dan file gambar berhasil dihapus",
		})
	}
}

// Get News ID
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
