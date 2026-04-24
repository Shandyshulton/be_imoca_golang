package controllers

import (
	"be_imoca_golang/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type VisionController struct {
	DB *gorm.DB
}

// Get Vision
func (vc *VisionController) GetVision(c *gin.Context) {
	var setting model.VisionSetting
	var items []model.VisionItem

	vc.DB.First(&setting)

	vc.DB.Order("id asc").Find(&items)

	c.JSON(http.StatusOK, gin.H{
		"code": "200",
		"data": model.VisionResponse{
			SectionTitle: setting.SectionTitle,
			Items:        items,
		},
	})
}

// Update Vision
func (vc *VisionController) UpdateVision(c *gin.Context) {
	var input model.VisionResponse
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "400", "message": "Format data salah"})
		return
	}

	// Gunakan Transaction untuk keamanan data
	err := vc.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.VisionSetting{}).Where("id = ?", 1).
			Update("section_title", input.SectionTitle).Error; err != nil {
			return err
		}

		if err := tx.Exec("DELETE FROM vision_items").Error; err != nil {
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
		c.JSON(http.StatusInternalServerError, gin.H{"code": "500", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    "200",
		"message": "Visi berhasil diperbarui!",
	})
}