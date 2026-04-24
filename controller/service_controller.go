package controllers

import (
	"be_imoca_golang/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ServiceController struct {
	DB *gorm.DB
}

// Get Services (Public)
func (sc *ServiceController) GetServices(c *gin.Context) {
	var setting model.ServiceSetting
	var items []model.ServiceItem

	// Ambil judul section
	sc.DB.First(&setting)

	// Ambil semua item layanan
	sc.DB.Order("id asc").Find(&items)

	c.JSON(http.StatusOK, gin.H{
		"code": "200",
		"data": model.ServicesResponse{
			SectionTitle: setting.SectionTitle,
			Items:        items,
		},
	})
}

// Update Services (Admin)
func (sc *ServiceController) UpdateServices(c *gin.Context) {
	var input model.ServicesResponse
	
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "400",
			"message": "Format data tidak sesuai",
		})
		return
	}

	err := sc.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Update Judul Section (ID 1)
		if err := tx.Model(&model.ServiceSetting{}).Where("id = ?", 1).
			Update("section_title", input.SectionTitle).Error; err != nil {
			return err
		}

		if err := tx.Exec("DELETE FROM services_items").Error; err != nil {
			return err
		}

		if len(input.Items) > 0 {
			if err := tx.Create(&input.Items).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "500",
			"message": "Gagal memperbarui data layanan: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    "200",
		"message": "Data layanan berhasil diperbarui!",
	})
}