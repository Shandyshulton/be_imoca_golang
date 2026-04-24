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

type MemberController struct {
	DB *gorm.DB
}

// Get All Members (with Search)
func (mc *MemberController) GetAll(c *gin.Context) {
	var members []model.Member
	keyword := c.Query("name")

	query := mc.DB.Order("id desc")
	if keyword != "" {
		query = query.Where("name LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	if err := query.Find(&members).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data member"})
		return
	}
	c.JSON(http.StatusOK, members)
}

// Get Member By ID
func (mc *MemberController) GetByID(c *gin.Context) {
	id := c.Param("id")
	var member model.Member

	if err := mc.DB.First(&member, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Member tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, member)
}


// Create Member
func (mc *MemberController) Create(c *gin.Context) {
	var member model.Member

	member.Name = c.PostForm("name")
	member.Email = c.PostForm("email")
	member.Phone = c.PostForm("phone")
	member.Description = c.PostForm("description")
	member.Link = c.PostForm("link")

	file, err := c.FormFile("image")
	if err == nil {
		// Validasi Ukuran File (Max 10MB)
		if file.Size > 10*1024*1024 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Ukuran file maksimal 10MB"})
			return
		}

		// Validasi Ekstensi
		ext := strings.ToLower(filepath.Ext(file.Filename))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Hanya mendukung format .jpg, .jpeg, dan .png"})
			return
		}

		uploadDir := "storage/uploads/members/"
		if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
			os.MkdirAll(uploadDir, os.ModePerm)
		}

		newFileName := fmt.Sprintf("member-%s%s", uuid.New().String(), ext)
		dst := uploadDir + newFileName

		if err := c.SaveUploadedFile(file, dst); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal simpan file: " + err.Error()})
			return
		}

		member.Image = newFileName
	}

	if err := mc.DB.Create(&member).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal simpan ke database"})
		return
	}

	c.JSON(http.StatusCreated, member)
}

// Update Member
func (mc *MemberController) Update(c *gin.Context) {
	id := c.Param("id")
	var member model.Member

	// Cari data lama di DB
	if err := mc.DB.First(&member, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Member tidak ditemukan"})
		return
	}

	// Simpan nama file lama
	oldImage := member.Image
	uploadDir := "storage/uploads/members/"

	member.Name = c.DefaultPostForm("name", member.Name)
	member.Email = c.DefaultPostForm("email", member.Email)
	member.Phone = c.DefaultPostForm("phone", member.Phone)
	member.Description = c.DefaultPostForm("description", member.Description)
	member.Link = c.DefaultPostForm("link", member.Link)

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
			c.JSON(http.StatusBadRequest, gin.H{"error": "Hanya mendukung format .jpg, .jpeg, dan .png"})
			return
		}

		if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
			os.MkdirAll(uploadDir, os.ModePerm)
		}

		newFileName := fmt.Sprintf("member-%s%s", uuid.New().String(), ext)
		dst := uploadDir + newFileName
		
		if err := c.SaveUploadedFile(file, dst); err == nil {
			// Hapus file lama
			if oldImage != "" {
				os.Remove(uploadDir + oldImage)
			}
			member.Image = newFileName
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal simpan file baru"})
			return
		}
	}

	if err := mc.DB.Save(&member).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data berhasil diperbarui", "data": member})
}

// Delete Member
func (mc *MemberController) Delete(c *gin.Context) {
	id := c.Param("id")
	var member model.Member

	if err := mc.DB.First(&member, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Member tidak ditemukan"})
		return
	}

	imageToDelete := member.Image

	if err := mc.DB.Delete(&member).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus data di database"})
		return
	}

	if imageToDelete != "" {
		path := "storage/uploads/members/" + imageToDelete
		err := os.Remove(path)
		if err != nil {
			fmt.Println("Log: File fisik tidak ditemukan atau gagal dihapus:", err)
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Member dan foto berhasil dihapus permanen"})
}