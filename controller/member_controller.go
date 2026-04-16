package controller

import (
	"be_imoca_golang/model"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type MemberUri struct {
	ID int `uri:"id" binding:"required"`
}

// GetMembersHandler 
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

// GetMemberByIDHandler 
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

// CreateMemberHandler 
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

		file, err := c.FormFile("image")
		if err == nil {
			if !isAllowedExtension(file.Filename) {
				c.JSON(http.StatusBadRequest, gin.H{"code": "400", "message": "Format file tidak didukung"})
				return
			}

			extension := filepath.Ext(file.Filename)
			newFileName := fmt.Sprintf("member-%d%s", time.Now().UnixNano(), extension)
			targetPath := filepath.Join("storage", "uploads", "members", newFileName)

			if err := c.SaveUploadedFile(file, targetPath); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"code": "500", "message": "Gagal simpan gambar"})
				return
			}
			m.Image = newFileName
		}

		if err := model.CreateMember(db, &m); err != nil {
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

// UpdateMemberHandler
func UpdateMemberHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var uri MemberUri
		if err := c.ShouldBindUri(&uri); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "400", "message": "ID tidak valid"})
			return
		}

		oldMember, err := model.GetMemberByID(db, uri.ID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"code": "404", "message": "Anggota tidak ditemukan"})
			return
		}

		m := model.Member{
			ID:          uri.ID,
			Name:        c.PostForm("name"),
			Email:       c.PostForm("email"),
			Phone:       c.PostForm("phone"),
			Description: c.PostForm("description"),
			Link:        c.PostForm("link"),
			Image:       oldMember.Image,
		}

		if strings.TrimSpace(m.Name) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"code": "400", "message": "Nama wajib diisi"})
			return
		}

		file, err := c.FormFile("image")
		if err == nil {
			if !isAllowedExtension(file.Filename) {
				c.JSON(http.StatusBadRequest, gin.H{"code": "400", "message": "Format gambar tidak didukung!"})
				return
			}

			extension := filepath.Ext(file.Filename)
			newFileName := fmt.Sprintf("member-%d%s", time.Now().UnixNano(), extension)
			targetPath := filepath.Join("storage", "uploads", "members", newFileName)

			if err := c.SaveUploadedFile(file, targetPath); err == nil {
				if oldMember.Image != "" {
					os.Remove(filepath.Join("storage", "uploads", "members", oldMember.Image))
				}
				m.Image = newFileName
			}
		}

		if err := model.UpdateMember(db, &m); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": "500", "message": "Gagal update database"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"code": "200", "message": "Data berhasil diperbarui", "data": m})
	}
}

// DeleteMemberHandler 
func DeleteMemberHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, _ := strconv.Atoi(idStr)

		m, err := model.GetMemberByID(db, id)
		if err == nil && m.Image != "" {
			os.Remove(filepath.Join("storage", "uploads", "members", m.Image))
		}

		if err := model.DeleteMember(db, id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": "500", "message": "Gagal hapus data"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"code": "200", "message": "Anggota berhasil dihapus"})
	}
}