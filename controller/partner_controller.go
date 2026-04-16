package controller

import (
	"be_imoca_golang/model"
	"database/sql"
	"strconv"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type PartnerUri struct {
	ID int `uri:"id" binding:"required"`
}

// isAllowedExtension 
func isAllowedExtension(fileName string) bool {
	extensions := []string{".jpg", ".jpeg", ".png"}
	ext := strings.ToLower(filepath.Ext(fileName))
	for _, e := range extensions {
		if e == ext {
			return true
		}
	}
	return false
}

// GetPartnersHandler
func GetPartnersHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		keyword := c.Query("name")
		partners, err := model.GetAllPartners(db, keyword)
		if err != nil {
			log.Printf("[ERROR] Gagal ambil data partner: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    "500",
				"message": "Gagal memuat daftar mitra: " + err.Error(),
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

func CreatePartnerHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("[INFO] Mencoba mendaftarkan mitra baru...")

		name := c.PostForm("name")
		email := c.PostForm("email")
		phone := c.PostForm("phone")
		description := c.PostForm("description")
		link := c.PostForm("link")

		// Empty Validation
		if strings.TrimSpace(name) == "" || strings.TrimSpace(email) == "" || 
		   strings.TrimSpace(phone) == "" || strings.TrimSpace(description) == "" {
			log.Println("[WARN] Pendaftaran mitra ditolak: Data tidak lengkap")
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    "400",
				"message": "Nama, Email, Telepon, dan Deskripsi wajib diisi!",
			})
			return
		}

		file, err := c.FormFile("image")
		if err != nil {
			log.Println("[WARN] Pendaftaran mitra ditolak: Gambar tidak ditemukan")
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    "400",
				"message": "Logo atau Gambar mitra wajib diunggah!",
			})
			return
		}

		// Validation Format File
		if !isAllowedExtension(file.Filename) {
			log.Printf("[REJECTED] Format file tidak diizinkan: %s", file.Filename)
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    "400",
				"message": "Format file tidak didukung! Gunakan .jpg, .jpeg, atau .png",
			})
			return
		}

		// Validation Size File (Maks 10MB)
		maxMB := 10
		var maxFileSize int64 = int64(maxMB) * 1024 * 1024
		if file.Size > maxFileSize {
			log.Printf("[REJECTED] File terlalu besar: %d bytes", file.Size)
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    "400",
				"message": fmt.Sprintf("Ukuran file terlalu besar! Maksimal adalah %dMB", maxMB),
			})
			return
		}

		// Save File to Storage
		extension := filepath.Ext(file.Filename)
		newFileName := fmt.Sprintf("partner-%d%s", time.Now().UnixNano(), extension)
		targetPath := filepath.Join("storage", "uploads", "partners", newFileName)

		if err := c.SaveUploadedFile(file, targetPath); err != nil {
			log.Printf("[ERROR] Gagal simpan gambar mitra: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"code": "500", "message": "Gagal menyimpan gambar di server"})
			return
		}

		p := model.Partner{
			Name:        name,
			Email:       email,
			Phone:       phone,
			Description: description,
			Link:        link,
			Image:       newFileName,
		}

		if err := model.CreatePartner(db, &p); err != nil {
			log.Printf("[DATABASE ERROR] Gagal create partner: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"code": "500", "message": "Gagal simpan ke database"})
			return
		}

		log.Printf("[SUCCESS] Mitra '%s' berhasil didaftarkan dengan ID %d", p.Name, p.ID)
		c.JSON(http.StatusCreated, gin.H{
			"code":    "201",
			"message": "Mitra berhasil didaftarkan!",
			"data":    p,
		})
	}
}

// UpdatePartnerHandler 
func UpdatePartnerHandler(db *sql.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var uri PartnerUri
        if err := c.ShouldBindUri(&uri); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"code": "400", "message": "ID mitra tidak valid"})
            return
        }

        log.Printf("[INFO] Memproses update data mitra ID: %d", uri.ID)

        oldPartner, err := model.GetPartnerByID(db, uri.ID)
        if err != nil {
            log.Printf("[WARN] Update gagal, mitra ID %d tidak ditemukan", uri.ID)
            c.JSON(http.StatusNotFound, gin.H{"code": "404", "message": "Mitra tidak ditemukan"})
            return
        }

        name := c.PostForm("name")
        email := c.PostForm("email")
        phone := c.PostForm("phone")
        description := c.PostForm("description")
        link := c.PostForm("link")

        // Empty Validation
        if strings.TrimSpace(name) == "" || strings.TrimSpace(email) == "" || strings.TrimSpace(phone) == "" || strings.TrimSpace(description) == "" {
            log.Printf("[WARN] Update ID %d ditolak: Data tidak lengkap", uri.ID)
            c.JSON(http.StatusBadRequest, gin.H{
                "code":    "400",
                "message": "Nama, Email, Telepon, dan Deskripsi wajib diisi!",
            })
            return
        }

        p := model.Partner{
            ID:          uri.ID,
            Name:        name,
            Email:       email,
            Phone:       phone,
            Description: description,
            Link:        link,
            Image:       oldPartner.Image, 
        }

        file, err := c.FormFile("image")
        if err == nil {
            if !isAllowedExtension(file.Filename) {
                c.JSON(http.StatusBadRequest, gin.H{"code": "400", "message": "Format gambar tidak didukung!"})
                return
            }

            extension := filepath.Ext(file.Filename)
            newFileName := fmt.Sprintf("partner-%d%s", time.Now().UnixNano(), extension)
            targetPath := filepath.Join("storage", "uploads", "partners", newFileName)

            if err := c.SaveUploadedFile(file, targetPath); err == nil {
                if oldPartner.Image != "" {
                    oldFilePath := filepath.Join("storage", "uploads", "partners", oldPartner.Image)
                    os.Remove(oldFilePath)
                }
                p.Image = newFileName
            }
        }

        // Update Database
        if err := model.UpdatePartner(db, &p); err != nil {
            log.Printf("[DATABASE ERROR] Gagal update partner ID %d: %v", p.ID, err)
            c.JSON(http.StatusInternalServerError, gin.H{"code": "500", "message": "Gagal memperbarui database"})
            return
        }

        log.Printf("[SUCCESS] Data mitra ID %d (%s) berhasil diperbarui", p.ID, p.Name)
        c.JSON(http.StatusOK, gin.H{
            "code":    "200",
            "message": "Data mitra berhasil diperbarui!",
            "data":    p,
        })
    }
}

// DeletePartnerHandler 
func DeletePartnerHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		log.Printf("[INFO] Menerima request hapus mitra ID: %s", idStr)

		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Printf("[ERROR] ID tidak valid: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    "400",
				"message": "ID mitra harus berupa angka!",
			})
			return
		}

		if err := model.DeletePartner(db, id); err != nil {
			log.Printf("[ERROR] Gagal hapus mitra ID %d: %v", id, err)
			
			if err.Error() == "data partner tidak ditemukan" {
				c.JSON(http.StatusNotFound, gin.H{
					"code":    "404",
					"message": "Mitra tidak ditemukan",
				})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    "500",
				"message": "Gagal menghapus mitra",
			})
			return
		}

		log.Printf("[SUCCESS] Mitra ID %d berhasil dihapus dari sistem", id)
		c.JSON(http.StatusOK, gin.H{
			"code":    "200",
			"message": "Mitra berhasil dihapus",
		})
	}
}

func GetPartnerByIDHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("[INFO] Mengambil detail mitra...")

		var uri PartnerUri
		if err := c.ShouldBindUri(&uri); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    "400",
				"message": "ID mitra tidak valid",
			})
			return
		}

		partner, err := model.GetPartnerByID(db, uri.ID)
		if err != nil {
			log.Printf("[WARN] Mitra ID %d tidak ditemukan", uri.ID)
			c.JSON(http.StatusNotFound, gin.H{
				"code":    "404",
				"message": "Mitra tidak ditemukan",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    "200",
			"message": "Detail mitra berhasil dimuat",
			"data":    partner,
		})
	}
}