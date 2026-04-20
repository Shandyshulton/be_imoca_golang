package controller

import (
	"be_imoca_golang/model"
	"be_imoca_golang/util"
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type PartnerUri struct {
	ID int `uri:"id" binding:"required"`
}

// Get Partners
func GetPartnersHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		keyword := c.Query("name")
		partners, err := model.GetAllPartners(db, keyword)
		if err != nil {
			log.Printf("[ERROR] Gagal ambil data partner: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    "500",
				"message": "Gagal memuat daftar mitra",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    "200",
			"message": "Daftar mitra berhasil dimuat",
			"data":    partners,
		})
	}
}

// Create Partner
func CreatePartnerHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.PostForm("name")
		email := c.PostForm("email")
		phone := c.PostForm("phone")
		description := c.PostForm("description")
		link := c.PostForm("link")

		if strings.TrimSpace(name) == "" || strings.TrimSpace(email) == "" ||
			strings.TrimSpace(phone) == "" || strings.TrimSpace(description) == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    "400",
				"message": "Nama, Email, Telepon, dan Deskripsi wajib diisi!",
			})
			return
		}

		newImageName, err := util.HandleFileUpload(c, "image", "partners", "")
		if err != nil {
			log.Printf("[WARN] Upload ditolak: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    "400",
				"message": err.Error(), // Pesan 10MB akan muncul di sini
			})
			return
		}

		if newImageName == "" {
			c.JSON(http.StatusBadRequest, gin.H{"code": "400", "message": "Logo mitra wajib diunggah!"})
			return
		}

		p := model.Partner{
			Name:        name,
			Email:       email,
			Phone:       phone,
			Description: description,
			Link:        link,
			Image:       newImageName,
		}

		if err := model.CreatePartner(db, &p); err != nil {
			log.Printf("[DATABASE ERROR] %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"code": "500", "message": "Gagal simpan ke database"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"code":    "201",
			"message": "Mitra berhasil didaftarkan!",
			"data":    p,
		})
	}
}

// Update Partner
func UpdatePartnerHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var uri PartnerUri
		if err := c.ShouldBindUri(&uri); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "400", "message": "ID mitra tidak valid"})
			return
		}

		oldPartner, err := model.GetPartnerByID(db, uri.ID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"code": "404", "message": "Mitra tidak ditemukan"})
			return
		}

		newImageName, err := util.HandleFileUpload(c, "image", "partners", oldPartner.Image)
		if err != nil {
			log.Printf("[ERROR] Update gambar ditolak: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    "400",
				"message": err.Error(),
			})
			return
		}

		p := model.Partner{
			ID:          uri.ID,
			Name:        c.PostForm("name"),
			Email:       c.PostForm("email"),
			Phone:       c.PostForm("phone"),
			Description: c.PostForm("description"),
			Link:        c.PostForm("link"),
			Image:       newImageName,
		}

		if strings.TrimSpace(p.Name) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"code": "400", "message": "Nama wajib diisi!"})
			return
		}

		if err := model.UpdatePartner(db, &p); err != nil {
			log.Printf("[DATABASE ERROR] Gagal update: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"code": "500", "message": "Gagal memperbarui database"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    "200",
			"message": "Data mitra berhasil diperbarui!",
			"data":    p,
		})
	}
}

// Delete Partner
func DeletePartnerHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "400", "message": "ID harus angka"})
			return
		}

		partner, err := model.GetPartnerByID(db, id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"code": "404", "message": "Mitra tidak ditemukan"})
			return
		}

		if err := model.DeletePartner(db, id); err != nil {
			log.Printf("[ERROR] %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"code": "500", "message": "Gagal menghapus mitra"})
			return
		}

		util.DeleteFile("partners", partner.Image)

		c.JSON(http.StatusOK, gin.H{
			"code":    "200",
			"message": "Mitra dan file logo berhasil dihapus",
		})
	}
}

// Get Partner ID
func GetPartnerByIDHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var uri PartnerUri
		if err := c.ShouldBindUri(&uri); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "400", "message": "ID mitra tidak valid"})
			return
		}

		partner, err := model.GetPartnerByID(db, uri.ID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"code": "404", "message": "Mitra tidak ditemukan"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    "200",
			"message": "Detail mitra berhasil dimuat",
			"data":    partner,
		})
	}
}