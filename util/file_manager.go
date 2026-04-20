package util

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func HandleFileUpload(c *gin.Context, fieldName string, category string, oldFileName string) (string, error) {
	file, err := c.FormFile(fieldName)
	if err != nil {
		return oldFileName, nil
	}

	// Validasi Ekstensi
	extensions := []string{".jpg", ".jpeg", ".png"}
	ext := strings.ToLower(filepath.Ext(file.Filename))
	// TAMBAHKAN LOG INI:
    log.Printf("[DEBUG] File received: %s | Ext: %s", file.Filename, ext)
	isAllowed := false
	for _, e := range extensions {
		if e == ext {
			isAllowed = true
			break
		}
	}
	if !isAllowed {
		log.Printf("[UPLOAD-REJECT] Format tidak didukung: %s", ext)
		return oldFileName, errors.New("format file tidak didukung! Gunakan .jpg, .jpeg, atau .png")
	}

	// Validasi Ukuran (10MB)
	var maxFileSize int64 = 10 * 1024 * 1024
	if file.Size > maxFileSize {
		log.Printf("[UPLOAD-REJECT] File terlalu besar: %d bytes", file.Size)
		return oldFileName, errors.New("Ukuran file terlalu besar! Maksimal adalah 10MB")
	}

	log.Printf("[UPLOAD] Memproses file baru untuk kategori: %s", category)

	newFileName := fmt.Sprintf("%s-%d%s", category, time.Now().UnixNano(), ext)
	targetPath := filepath.Join("storage", "uploads", category, newFileName)

	if err := c.SaveUploadedFile(file, targetPath); err != nil {
		log.Printf("[UPLOAD-ERROR] Gagal simpan file: %v", err)
		return oldFileName, err
	}
	
	log.Printf("[UPLOAD-SUCCESS] File disimpan: %s", newFileName)

	if oldFileName != "" {
		DeleteFile(category, oldFileName)
	}

	return newFileName, nil
}

// Delete File
func DeleteFile(category string, fileName string) {
	if fileName == "" {
		return
	}
	path := filepath.Join("storage", "uploads", category, fileName)
	if _, err := os.Stat(path); err == nil {
		if err := os.Remove(path); err != nil {
			log.Printf("[Cleanup-ERROR] Gagal hapus file %s: %v", path, err)
		} else {
			log.Printf("[Cleanup-SUCCESS] File berhasil dihapus: %s", path)
		}
	} else {
		log.Printf("[Cleanup-WARN] File tidak ditemukan di storage, skip: %s", path)
	}
}