package controllers

import (
	"be_imoca_golang/model"
	"fmt"
	"net/http"
	"os"
	"strings"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type HeroController struct {
	DB *gorm.DB
}

// Get Hero (Public)
func (hc *HeroController) GetHero(c *gin.Context) {
	var hero model.Hero
	if err := hc.DB.First(&hero).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data Hero belum diatur"})
		return
	}
	c.JSON(http.StatusOK, hero)
}

// Update Hero (Admin)
func (hc *HeroController) UpdateHero(c *gin.Context) {
	var hero model.Hero
	if err := hc.DB.First(&hero).Error; err != nil {
		hero = model.Hero{}
	}

	oldImage := hero.Image

	hero.BadgeText = c.DefaultPostForm("badge_text", hero.BadgeText)
	hero.Title = c.DefaultPostForm("title", hero.Title)
	hero.TitleHighlight = c.DefaultPostForm("title_highlight", hero.TitleHighlight)
	hero.Description = c.DefaultPostForm("description", hero.Description)

	// Handle Image Upload
	file, err := c.FormFile("image")
	if err == nil {
		uploadDir := "storage/uploads/hero/"
		if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
			os.MkdirAll(uploadDir, os.ModePerm)
		}

		newFileName := fmt.Sprintf("hero-%s%s", uuid.New().String(), filepath.Ext(file.Filename))
		if err := c.SaveUploadedFile(file, uploadDir+newFileName); err == nil {
			// Hapus gambar lama
			if oldImage != "" {
				os.Remove(uploadDir + oldImage)
			}
			hero.Image = newFileName
		}
	}

	if err := hc.DB.Save(&hero).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update hero"})
		return
	}
	c.JSON(http.StatusOK, hero)
}

// Get All Hero Images
func (hc *HeroController) GetHeroImages(c *gin.Context) {
	var images []model.HeroImage
	hc.DB.Order("sort_order asc").Find(&images)
	c.JSON(http.StatusOK, images)
}

// Delete Hero Image
func (hc *HeroController) DeleteHeroImage(c *gin.Context) {
	id := c.Param("id")
	var img model.HeroImage

	if err := hc.DB.First(&img, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Gambar tidak ditemukan"})
		return
	}

	if err := hc.DB.Delete(&img).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal hapus data"})
		return
	}

	// Hapus file fisik
	os.Remove("storage/uploads/hero/" + img.Image)

	c.JSON(http.StatusOK, gin.H{"message": "Gambar berhasil dihapus"})
}

// Add Hero Image (Admin)
func (hc *HeroController) AddHeroImage(c *gin.Context) {
	var hero model.Hero
	var heroImage model.HeroImage

	if err := hc.DB.First(&hero).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data Hero utama tidak ditemukan. Harap isi data Hero utama terlebih dahulu."})
		return
	}

	heroImage.HeroID = hero.ID

	sortOrder := c.DefaultPostForm("sort_order", "0")
	fmt.Sscanf(sortOrder, "%d", &heroImage.SortOrder)

	// Handle File Upload
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File gambar wajib diunggah"})
		return
	}

	// Validasi Ukuran (Max 10MB)
	if file.Size > 10*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ukuran file terlalu besar (Maksimal 10MB)"})
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Hanya mendukung format .jpg, .jpeg, dan .png"})
		return
	}

	uploadDir := "storage/uploads/hero/"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.MkdirAll(uploadDir, os.ModePerm)
	}

	newFileName := fmt.Sprintf("hero-gallery-%s%s", uuid.New().String(), ext)
	dst := uploadDir + newFileName

	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal simpan gambar ke server"})
		return
	}

	heroImage.Image = newFileName

	if err := hc.DB.Create(&heroImage).Error; err != nil {
		os.Remove(dst)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal simpan ke database: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, heroImage)
}