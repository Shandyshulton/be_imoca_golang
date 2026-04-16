package controller

import (
	"be_imoca_golang/model"
	"database/sql"
	"net/http" 

	"github.com/gin-gonic/gin"
)

func GetMissionHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		data, err := model.GetMission(db)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    "404",
				"message": "Data mission belum ada",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": "200", 
			"data": data,
		})
	}
}

func UpdateMissionHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input model.MissionSection
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    "400",
				"message": "Format JSON salah",
			})
			return
		}

		if err := model.UpdateMission(db, &input); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    "500",
				"message": "Gagal update database",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    "200",
			"message": "Mission Section berhasil diperbarui!",
		})
	}
}