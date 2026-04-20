package controller

import (
	"be_imoca_golang/model"
	"be_imoca_golang/util"
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Get Hero
func GetHeroHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		hero, err := model.GetHero(db)
		if err != nil {
			log.Printf("[ERROR] Gagal ambil hero: %v", err)
			c.JSON(http.StatusNotFound, gin.H{"code": "404", "message": "Data Hero belum diatur"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": "200", "data": hero})
	}
}

// Update Hero
func UpdateHeroHandler(db *sql.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var input model.Hero
        
        if err := c.ShouldBind(&input); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"code": "400", "message": "Gagal membaca form"})
            return
        }

        oldHero, _ := model.GetHero(db)

        newImageName, err := util.HandleFileUpload(c, "image", "hero", oldHero.Image)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"code": "400", "message": err.Error()})
            return
        }

        input.ID = oldHero.ID
        input.Image = newImageName

        if err := model.UpdateHero(db, &input); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"code": "500", "message": "DB Error"})
            return
        }

        c.JSON(http.StatusOK, gin.H{"code": "200", "data": input})
    }
}