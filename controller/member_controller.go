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

type MemberUri struct {
	ID int `uri:"id" binding:"required"`
}

// Get Members
func GetMembersHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		keyword := c.Query("name")
		members, err := model.GetAllMembers(db, keyword)
		if err != nil {
			log.Printf("[ERROR] Gagal ambil data member: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    "500",
				"message": "Gagal memuat daftar anggota",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    "200",
			"message": "Daftar anggota berhasil dimuat",
			"data":    members,
		})
	}
}

// Get Member ID
func GetMemberByIDHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var uri MemberUri
		if err := c.ShouldBindUri(&uri); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "400", "message": "ID tidak valid"})
			return
		}

		member, err := model.GetMemberByID(db, uri.ID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"code": "404", "message": "Anggota tidak ditemukan"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    "200",
			"message": "Detail anggota berhasil dimuat",
			"data":    member,
		})
	}
}

// Create Member
func CreateMemberHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var m model.Member

		m.Name = c.PostForm("name")
		m.Email = c.PostForm("email")
		m.Phone = c.PostForm("phone")
		m.Description = c.PostForm("description")
		m.Link = c.PostForm("link")

		if strings.TrimSpace(m.Name) == "" || strings.TrimSpace(m.Email) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"code": "400", "message": "Nama dan Email wajib diisi"})
			return
		}

		newImageName, err := util.HandleFileUpload(c, "image", "members", "")
		if err != nil {
			log.Printf("[ERROR] Gagal upload gambar member: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    "400",
				"message": err.Error(), 
			})
			return
		}
		m.Image = newImageName

		if err := model.CreateMember(db, &m); err != nil {
			log.Printf("[DB ERROR] %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"code": "500", "message": "Database error: " + err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"code":    "201",
			"message": "Anggota berhasil didaftarkan!",
			"data":    m,
		})
	}
}

// Update Member
func UpdateMemberHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var uri MemberUri
		if err := c.ShouldBindUri(&uri); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "400", "message": "ID tidak valid"})
			return
		}

		oldMember, err := model.GetMemberByID(db, uri.ID)
		if err != nil {
			log.Printf("[WARN] Update member ID %d gagal: Tidak ditemukan", uri.ID)
			c.JSON(http.StatusNotFound, gin.H{"code": "404", "message": "Anggota tidak ditemukan"})
			return
		}

		newImageName, err := util.HandleFileUpload(c, "image", "members", oldMember.Image)
		if err != nil {
			log.Printf("[ERROR] Gagal proses gambar member: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    "400",
				"message": err.Error(),
			})
			return
		}

		m := model.Member{
			ID:          uri.ID,
			Name:        c.PostForm("name"),
			Email:       c.PostForm("email"),
			Phone:       c.PostForm("phone"),
			Description: c.PostForm("description"),
			Link:        c.PostForm("link"),
			Image:       newImageName,
		}

		if strings.TrimSpace(m.Name) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"code": "400", "message": "Nama wajib diisi"})
			return
		}

		if err := model.UpdateMember(db, &m); err != nil {
			log.Printf("[DB ERROR] Gagal update member ID %d: %v", uri.ID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"code": "500", "message": "Gagal update database"})
			return
		}

		log.Printf("[SUCCESS] Member ID %d berhasil diperbarui", uri.ID)
		c.JSON(http.StatusOK, gin.H{"code": "200", "message": "Data berhasil diperbarui", "data": m})
	}
}

// Delete Member
func DeleteMemberHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "400", "message": "ID harus berupa angka"})
			return
		}

		member, err := model.GetMemberByID(db, id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"code": "404", "message": "Anggota tidak ditemukan"})
			return
		}

		if err := model.DeleteMember(db, id); err != nil {
			log.Printf("[DB ERROR] Gagal hapus member ID %d: %v", id, err)
			c.JSON(http.StatusInternalServerError, gin.H{"code": "500", "message": "Gagal hapus data dari database"})
			return
		}

		util.DeleteFile("members", member.Image)

		log.Printf("[SUCCESS] Member ID %d dan filenya berhasil dihapus", id)
		c.JSON(http.StatusOK, gin.H{"code": "200", "message": "Anggota berhasil dihapus"})
	}
}