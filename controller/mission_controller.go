package controllers

import (
	"be_imoca_golang/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MissionController struct {
	DB *gorm.DB
}

// Get Mission
func (mc *MissionController) GetMission(c *gin.Context) {
	var mission model.MissionSection
	
	// Mengambil data pertama (LIMIT 1)
	if err := mc.DB.First(&mission).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    "404",
			"message": "Data mission belum tersedia",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": "200",
		"data": mission,
	})
}

// Update Mission
func (mc *MissionController) UpdateMission(c *gin.Context) {
	var mission model.MissionSection
	
	// Cek apakah data sudah ada
	if err := mc.DB.First(&mission).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    "404",
			"message": "Data mission tidak ditemukan untuk diupdate",
		})
		return
	}

	// Binding input JSON
	var input model.MissionSection
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "400",
			"message": "Format JSON tidak valid",
		})
		return
	}

	// Gunakan Map untuk memaksa update ke baris ID yang ditemukan
	updateData := map[string]interface{}{
		"section_title": input.SectionTitle,
		"description":   input.Description,
	}

	if err := mc.DB.Model(&mission).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "500",
			"message": "Gagal memperbarui database: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    "200",
		"message": "Mission Section berhasil diperbarui!",
		"data":    mission,
	})
}