package controllers

import (
	"be_imoca_golang/model"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NewsController struct {
	DB *gorm.DB
}

// GetAll News
func (nc *NewsController) GetAll(c *gin.Context) {
	var newsList []model.News
	keyword := c.Query("title")

	query := nc.DB.Order("published_at desc")
	if keyword != "" {
		query = query.Where("title LIKE ? OR summary LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	if err := query.Find(&newsList).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memuat berita"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": newsList})
}

// GetByID News
func (nc *NewsController) GetByID(c *gin.Context) {
	id := c.Param("id")
	var news model.News

	if err := nc.DB.First(&news, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Berita tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": news})
}

// Create News
func (nc *NewsController) Create(c *gin.Context) {
	file, err := c.FormFile("image")

	if err == nil {
		if file.Size > 10*1024*1024 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Ukuran file terlalu besar! Maksimal adalah 10MB"})
			return
		}

		// Validasi Ekstensi
		ext := strings.ToLower(filepath.Ext(file.Filename))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Format file dilarang! Hanya boleh .jpg, .jpeg, .png"})
			return
		}
	}

	var news model.News
	news.Title = c.PostForm("title")
	news.Summary = c.PostForm("summary")
	news.Source = c.PostForm("source")
	news.Badge = c.DefaultPostForm("badge", "REGULASI")
	news.URL = c.PostForm("url")

	// Parsing Tanggal
	pubDate := c.PostForm("published_at")
	if pubDate != "" {
		t, _ := time.Parse("2006-01-02", pubDate)
		news.PublishedAt = t
	} else {
		news.PublishedAt = time.Now()
	}

	if file != nil {
		uploadDir := "storage/uploads/news/"
		os.MkdirAll(uploadDir, os.ModePerm)

		newFileName := fmt.Sprintf("news-%s%s", uuid.New().String(), strings.ToLower(filepath.Ext(file.Filename)))
		dst := filepath.Join(uploadDir, newFileName)

		if err := c.SaveUploadedFile(file, dst); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan file ke server"})
			return
		}
		news.Image = newFileName
	}

	if err := nc.DB.Create(&news).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal simpan ke database"})
		return
	}

	c.JSON(http.StatusCreated, news)
}

// Update News
func (nc *NewsController) Update(c *gin.Context) {
	id := c.Param("id")
	var news model.News

	// 1. Cari data lama di database
	if err := nc.DB.First(&news, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Berita tidak ditemukan"})
		return
	}

	// 2. Validasi file baru di awal (Early Validation)
	file, err := c.FormFile("image")
	if err == nil {
		if file.Size > 10*1024*1024 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "File terlalu besar! Maksimal 10MB"})
			return
		}

		ext := strings.ToLower(filepath.Ext(file.Filename))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Format dilarang (Gunakan jpg/png)"})
			return
		}
	}

	oldImage := news.Image
	uploadDir := "storage/uploads/news/"

	// 3. Tangkap data dari form-data ke variabel penampung
	title := c.DefaultPostForm("title", news.Title)
	summary := c.DefaultPostForm("summary", news.Summary)
	source := c.DefaultPostForm("source", news.Source)
	badge := c.DefaultPostForm("badge", news.Badge)
	url := c.DefaultPostForm("url", news.URL)
	imageName := news.Image // Default tetap pakai yang lama

	// Parsing PublishedAt jika dikirim
	publishedAt := news.PublishedAt
	pubDate := c.PostForm("published_at")
	if pubDate != "" {
		if t, err := time.Parse("2006-01-02", pubDate); err == nil {
			publishedAt = t
		}
	}

	// 4. Proses Simpan File Baru
	if file != nil {
		if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
			os.MkdirAll(uploadDir, os.ModePerm)
		}

		newFileName := fmt.Sprintf("news-%s%s", uuid.New().String(), strings.ToLower(filepath.Ext(file.Filename)))
		dst := filepath.Join(uploadDir, newFileName)

		if err := c.SaveUploadedFile(file, dst); err == nil {
			if oldImage != "" {
				os.Remove(filepath.Join(uploadDir, oldImage)) // Hapus file lama
			}
			imageName = newFileName
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal simpan gambar baru"})
			return
		}
	}

	// 5. UPDATE MENGGUNAKAN MAP (Solusi agar pasti tersimpan)
	updateData := map[string]interface{}{
		"title":        title,
		"summary":      summary,
		"source":       source,
		"badge":        badge,
		"url":          url,
		"image":        imageName,
		"published_at": publishedAt,
	}

	// Updates(updateData) akan memaksa update kolom yang ada di map saja
	if err := nc.DB.Model(&news).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update database: " + err.Error()})
		return
	}

	// Ambil ulang data terbaru setelah update untuk dikembalikan di respons
	nc.DB.First(&news, id)
	c.JSON(http.StatusOK, news)
}

// Delete News
func (nc *NewsController) Delete(c *gin.Context) {
	id := c.Param("id")
	var news model.News
	if err := nc.DB.First(&news, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Berita tidak ditemukan"})
		return
	}

	imageToDelete := news.Image
	if err := nc.DB.Delete(&news).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal hapus data"})
		return
	}

	if imageToDelete != "" {
		path := filepath.Join("storage/uploads/news/", imageToDelete)
		err := os.Remove(path)
		if err != nil {
			fmt.Printf("Log: Gagal menghapus file fisik %s: %v\n", path, err)
		} else {
			fmt.Printf("Log: File %s berhasil dihapus dari storage\n", imageToDelete)
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Berita berhasil dihapus"})
}
