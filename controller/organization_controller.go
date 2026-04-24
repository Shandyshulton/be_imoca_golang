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

type OrganizationController struct {
	DB *gorm.DB
}

// Get All Members
func (oc *OrganizationController) GetAll(c *gin.Context) {
	var members []model.OrganizationStructure
	category := c.Query("category")

	query := oc.DB.Order("id asc")
	if category != "" {
		query = query.Where("team_category = ?", category)
	}

	if err := query.Find(&members).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data"})
		return
	}
	c.JSON(http.StatusOK, members)
}

// Create Member
func (oc *OrganizationController) Create(c *gin.Context) {
	var member model.OrganizationStructure
	member.FullName = c.PostForm("full_name")
	member.JobPosition = c.PostForm("job_position")
	member.TeamCategory = c.PostForm("team_category")

	file, err := c.FormFile("profile_photo")
	if err == nil {
		// Validasi Ukuran (10MB)
		if file.Size > 10*1024*1024 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Ukuran foto maksimal 10MB"})
			return
		}

		// Validasi Ekstensi (jpg, jpeg, png)
		ext := strings.ToLower(filepath.Ext(file.Filename))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Format file tidak didukung (.jpg, .jpeg, .png)"})
			return
		}

		uploadDir := "storage/uploads/teams/"
		if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
			os.MkdirAll(uploadDir, os.ModePerm)
		}

		prefix := strings.ReplaceAll(member.TeamCategory, " ", "-")
		newFileName := fmt.Sprintf("%s-%s%s", prefix, uuid.New().String(), ext)

		if err := c.SaveUploadedFile(file, uploadDir+newFileName); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal simpan gambar"})
			return
		}
		member.ProfilePhoto = newFileName
	}

	if err := oc.DB.Create(&member).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal simpan ke database"})
		return
	}
	c.JSON(http.StatusCreated, member)
}

// Update Member
func (oc *OrganizationController) Update(c *gin.Context) {
	id := c.Param("id")
	var member model.OrganizationStructure

	if err := oc.DB.First(&member, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data tidak ditemukan"})
		return
	}

	oldPhoto := member.ProfilePhoto
	member.FullName = c.DefaultPostForm("full_name", member.FullName)
	member.JobPosition = c.DefaultPostForm("job_position", member.JobPosition)
	member.TeamCategory = c.DefaultPostForm("team_category", member.TeamCategory)

	file, err := c.FormFile("profile_photo")
	if err == nil {
		// Validasi Ukuran
		if file.Size > 10*1024*1024 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Ukuran foto maksimal 10MB"})
			return
		}

		// Validasi Ekstensi
		ext := strings.ToLower(filepath.Ext(file.Filename))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Format ditolak (.jpg, .jpeg, .png)"})
			return
		}

		uploadDir := "storage/uploads/teams/"
		prefix := strings.ReplaceAll(member.TeamCategory, " ", "-")
		newFileName := fmt.Sprintf("%s-%s%s", prefix, uuid.New().String(), ext)

		if err := c.SaveUploadedFile(file, uploadDir+newFileName); err == nil {
			if oldPhoto != "" {
				os.Remove(uploadDir + oldPhoto) // Hapus foto lama
			}
			member.ProfilePhoto = newFileName
		}
	}

	if err := oc.DB.Model(&member).Omit("created_at").Save(&member).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui database"})
		return
	}
	c.JSON(http.StatusOK, member)
}

// Delete Member
func (oc *OrganizationController) Delete(c *gin.Context) {
	id := c.Param("id")
	var member model.OrganizationStructure

	if err := oc.DB.First(&member, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data tidak ditemukan"})
		return
	}

	photoToDelete := member.ProfilePhoto
	if err := oc.DB.Delete(&member).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus data"})
		return
	}

	if photoToDelete != "" {
		os.Remove("storage/uploads/teams/" + photoToDelete) // Hapus fisik
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data berhasil dihapus"})
}
