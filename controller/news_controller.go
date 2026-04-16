package controller

import (
	"be_imoca_golang/model"
	"database/sql"
	"strconv"
	"fmt"
	"log"
	"os"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type NewsUri struct {
	ID int `uri:"id" binding:"required"`
}

// GetNewsHandler 
func GetNewsHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		keyword := c.Query("title")
		log.Printf("[INFO] Memproses pengambilan data berita. Keyword pencarian: '%s'", keyword)

		news, err := model.GetAllNews(db, keyword)
		if err != nil {
			log.Printf("[ERROR] Gagal mengambil data berita dari database: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    "500",
				"message": "Gagal memuat daftar berita",
			})
			return
		}
		
		log.Printf("[SUCCESS] Berhasil mengambil %d data berita", len(news))
		c.JSON(http.StatusOK, gin.H{
			"code":    "200",
			"message": "Data berita berhasil diambil",
			"data":    news,
		})
	}
}

// CreateNewsHandler 
func CreateNewsHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("[INFO] Mencoba menambahkan berita baru...")

		title := c.PostForm("title")
		summary := c.PostForm("summary")
		source := c.PostForm("source") 
		url := c.PostForm("url")

		// Empty Validation
		if strings.TrimSpace(title) == "" || strings.TrimSpace(summary) == "" || 
		   strings.TrimSpace(source) == "" || strings.TrimSpace(url) == "" {
			log.Println("[WARN] Penambahan berita ditolak: Data tidak lengkap")
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    "400",
				"message": "Semua field (Judul, Ringkasan, Sumber, URL) wajib diisi!",
			})
			return
		}

		file, err := c.FormFile("image")
		if err != nil {
			log.Println("[WARN] Penambahan berita ditolak: Gambar tidak ditemukan")
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    "400",
				"message": "Gambar berita wajib diunggah!",
			})
			return
		}

		if !isAllowedExtension(file.Filename) {
			log.Printf("[REJECTED] Format file tidak didukung: %s", file.Filename)
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    "400",
				"message": "Format file tidak didukung! Gunakan .jpg, .jpeg, atau .png",
			})
			return
		}

		extension := filepath.Ext(file.Filename)
		newFileName := fmt.Sprintf("news-%d%s", time.Now().UnixNano(), extension)
		targetPath := filepath.Join("storage", "uploads", "news", newFileName)

		if err := c.SaveUploadedFile(file, targetPath); err != nil {
			log.Printf("[ERROR] Gagal menyimpan file di server: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"code": "500", "message": "Gagal menyimpan gambar"})
			return
		}

		n := model.News{
			Title:   title,
			Summary: summary,
			Source:  source, 
			URL:     url,
			Image:   newFileName, 
		}

		if err := model.CreateNews(db, &n); err != nil {
			log.Printf("[DATABASE ERROR] Gagal simpan news: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"code": "500", "message": "Gagal simpan ke database"})
			return
		}

		log.Printf("[SUCCESS] Berita '%s' (ID: %d) berhasil diterbitkan", n.Title, n.ID)
		c.JSON(http.StatusCreated, gin.H{
			"code":    "201",
			"message": "Berita berhasil diterbitkan!",
			"data":    n,
		})
	}
}

func UpdateNewsHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var uri NewsUri
		if err := c.ShouldBindUri(&uri); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "400", "message": "ID tidak valid"})
			return
		}

		log.Printf("[INFO] Memproses pembaruan berita ID: %d", uri.ID)

		oldNews, err := model.GetNewsByID(db, uri.ID)
		if err != nil {
			log.Printf("[WARN] Update gagal: Berita ID %d tidak ditemukan", uri.ID)
			c.JSON(http.StatusNotFound, gin.H{"code": "404", "message": "Berita tidak ditemukan"})
			return
		}

		title := c.PostForm("title")
		summary := c.PostForm("summary")
		source := c.PostForm("source") 
		url := c.PostForm("url")

		// Empty Validation 
		if strings.TrimSpace(title) == "" || strings.TrimSpace(summary) == "" || 
		   strings.TrimSpace(source) == "" || strings.TrimSpace(url) == "" {
			log.Printf("[WARN] Update ID %d ditolak: Data tidak lengkap", uri.ID)
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    "400",
				"message": "Semua field wajib diisi!",
			})
			return
		}

		updatedNews := model.News{
			ID:      uri.ID,
			Title:   title,
			Summary: summary,
			Source:  source,
			URL:     url,
			Image:   oldNews.Image, 
		}

		file, err := c.FormFile("image")
		if err == nil {
			if !isAllowedExtension(file.Filename) {
				c.JSON(http.StatusBadRequest, gin.H{"code": "400", "message": "Format gambar tidak didukung!"})
				return
			}

			extension := filepath.Ext(file.Filename)
			newFileName := fmt.Sprintf("news-%d%s", time.Now().UnixNano(), extension)
			targetPath := filepath.Join("storage", "uploads", "news", newFileName)

			if err := c.SaveUploadedFile(file, targetPath); err == nil {
				oldFilePath := filepath.Join("storage", "uploads", "news", oldNews.Image)
				os.Remove(oldFilePath) 
				updatedNews.Image = newFileName
			}
		}

		if err := model.UpdateNews(db, &updatedNews); err != nil {
			log.Printf("[DATABASE ERROR] Gagal update news ID %d: %v", uri.ID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"code": "500", "message": "Gagal update database"})
			return
		}

		log.Printf("[SUCCESS] Berita ID %d berhasil diperbarui", updatedNews.ID)
		c.JSON(http.StatusOK, gin.H{
			"code":    "200",
			"message": "Berita berhasil diperbarui!",
			"data":    updatedNews,
		})
	}
}

func DeleteNewsHandler(db *sql.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        idStr := c.Param("id")
        log.Printf("[INFO] Menerima request hapus berita ID: %s", idStr)

        // Konversi string ke int
        id, err := strconv.Atoi(idStr)
        if err != nil {
            log.Printf("[ERROR] ID tidak valid: %v", err)
            c.JSON(http.StatusBadRequest, gin.H{"code": "400", "message": "ID berita harus angka"})
            return
        }

        if err := model.DeleteNews(db, id); err != nil {
            log.Printf("[ERROR] Gagal menghapus berita ID %d: %v", id, err)
            c.JSON(http.StatusInternalServerError, gin.H{
                "code":    "500",
                "message": "Gagal menghapus berita",
            })
            return
        }

        log.Printf("[SUCCESS] Berita ID %d berhasil dihapus", id)
        c.JSON(http.StatusOK, gin.H{"code": "200", "message": "Berita berhasil dihapus"})
    }
}

func GetNewsByIDHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var uri NewsUri
		if err := c.ShouldBindUri(&uri); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "400", "message": "ID tidak valid"})
			return
		}

		log.Printf("[INFO] Mengambil detail berita ID: %d", uri.ID)

		news, err := model.GetNewsByID(db, uri.ID)
		if err != nil {
			log.Printf("[WARN] Berita ID %d tidak ditemukan", uri.ID)
			c.JSON(http.StatusNotFound, gin.H{"code": "404", "message": "Berita tidak ditemukan"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    "200",
			"message": "Detail berita berhasil diambil",
			"data":    news,
		})
	}
}

