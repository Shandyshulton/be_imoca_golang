package controller

import (
	"be_imoca_golang/model"
	"be_imoca_golang/util"
	"database/sql"
	"fmt"
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

// ── Hero Images ──────────────────────────────────────────────────────────────

// Get Hero Images
func GetHeroImagesHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		hero, err := model.GetHero(db)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"code": "404", "message": "Data hero tidak ditemukan"})
			return
		}

		images, err := model.GetHeroImages(db, hero.ID)
		if err != nil {
			log.Printf("[ERROR] Gagal ambil hero images: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"code": "500", "message": "Gagal mengambil data gambar"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"code": "200", "data": images})
	}
}

// Add Hero Image
func AddHeroImageHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		hero, err := model.GetHero(db)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"code": "404", "message": "Data hero tidak ditemukan"})
			return
		}

		// upload file baru (tidak ada old image karena ini insert)
		filename, err := util.HandleFileUpload(c, "image", "hero", "")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "400", "message": err.Error()})
			return
		}

		img := model.HeroImage{
			HeroID: hero.ID,
			Image:  filename,
			Order:  0, // default, bisa diatur dari frontend
		}

		if err := model.AddHeroImage(db, &img); err != nil {
			log.Printf("[ERROR] Gagal tambah hero image: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"code": "500", "message": "Gagal menyimpan gambar"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"code": "200", "data": img})
	}
}

// Delete Hero Image
func DeleteHeroImageHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		var imageID int
		if _, err := fmt.Sscanf(idStr, "%d", &imageID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "400", "message": "ID tidak valid"})
			return
		}

		filename, err := model.DeleteHeroImage(db, imageID)
		if err != nil {
			log.Printf("[ERROR] Gagal hapus hero image: %v", err)
			c.JSON(http.StatusNotFound, gin.H{"code": "404", "message": err.Error()})
			return
		}

		// hapus file fisik dari disk
		if filename != "" {
			util.DeleteFile("hero", filename)
		}

		c.JSON(http.StatusOK, gin.H{"code": "200", "message": "Gambar berhasil dihapus"})
	}
}