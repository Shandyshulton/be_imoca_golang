package controllers

import (
	"be_imoca_golang/model"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PartnerController struct {
	DB *gorm.DB
}

// Get All Partners
func (pc *PartnerController) GetAll(c *gin.Context) {
	var partners []model.Partner
	keyword := c.Query("name")

	query := pc.DB.Order("id desc")
	if keyword != "" {
		query = query.Where("name LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	if err := query.Find(&partners).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memuat daftar mitra"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": partners})
}

// Get Partner By ID
func (pc *PartnerController) GetByID(c *gin.Context) {
	id := c.Param("id")
	var partner model.Partner
	if err := pc.DB.First(&partner, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Mitra tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": partner})
}

// Create Partner
func (pc *PartnerController) Create(c *gin.Context) {
	var partner model.Partner
	partner.Name = c.PostForm("name")
	partner.Email = c.PostForm("email")
	partner.Phone = c.PostForm("phone")
	partner.Description = c.PostForm("description")
	partner.Link = c.PostForm("link")

	if strings.TrimSpace(partner.Name) == "" || strings.TrimSpace(partner.Email) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama dan Email wajib diisi"})
		return
	}

	file, err := c.FormFile("image")
	if err == nil {
		// Validasi Ukuran (Max 10MB)
		if file.Size > 10*1024*1024 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Ukuran file maksimal 10MB"})
			return
		}

		// Validasi Ekstensi
		ext := strings.ToLower(filepath.Ext(file.Filename))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Hanya mendukung .jpg, .jpeg, .png"})
			return
		}

		uploadDir := "storage/uploads/partners/"
		if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
			os.MkdirAll(uploadDir, os.ModePerm)
		}

		newFileName := fmt.Sprintf("partner-%s%s", uuid.New().String(), ext)
		if err := c.SaveUploadedFile(file, uploadDir+newFileName); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal simpan gambar"})
			return
		}
		partner.Image = newFileName
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Logo mitra wajib diunggah"})
		return
	}

	if err := pc.DB.Create(&partner).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal simpan ke database"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Mitra berhasil didaftarkan", "data": partner})
}

// Update Partner
func (pc *PartnerController) Update(c *gin.Context) {
	id := c.Param("id")
	var partner model.Partner
	if err := pc.DB.First(&partner, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Mitra tidak ditemukan"})
		return
	}

	oldImage := partner.Image
	partner.Name = c.DefaultPostForm("name", partner.Name)
	partner.Email = c.DefaultPostForm("email", partner.Email)
	partner.Phone = c.DefaultPostForm("phone", partner.Phone)
	partner.Description = c.DefaultPostForm("description", partner.Description)
	partner.Link = c.DefaultPostForm("link", partner.Link)

	file, err := c.FormFile("image")
	if err == nil {
		if file.Size > 10*1024*1024 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "File maksimal 10MB"})
			return
		}

		ext := strings.ToLower(filepath.Ext(file.Filename))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Format file tidak didukung"})
			return
		}

		uploadDir := "storage/uploads/partners/"
		newFileName := fmt.Sprintf("partner-%s%s", uuid.New().String(), ext)
		
		if err := c.SaveUploadedFile(file, uploadDir+newFileName); err == nil {
			if oldImage != "" {
				os.Remove(uploadDir + oldImage) // Hapus logo lama agar tidak menumpuk
			}
			partner.Image = newFileName
		}
	}

	if err := pc.DB.Save(&partner).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui database"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Data mitra berhasil diperbarui", "data": partner})
}

// Delete Partner
func (pc *PartnerController) Delete(c *gin.Context) {
	id := c.Param("id")
	var partner model.Partner
	if err := pc.DB.First(&partner, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Mitra tidak ditemukan"})
		return
	}

	imageToDelete := partner.Image
	if err := pc.DB.Delete(&partner).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus data"})
		return
	}

	if imageToDelete != "" {
		os.Remove("storage/uploads/partners/" + imageToDelete) // Hapus file fisik
	}

	c.JSON(http.StatusOK, gin.H{"message": "Mitra dan logo berhasil dihapus"})
}