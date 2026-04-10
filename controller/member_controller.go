package controller

import (
	"be_imoca_golang/model"
	"database/sql"
	"log"
	"net/http"
	"strings" 

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
            log.Printf("[ERROR] Gagal mengambil data member: %v", err)
            c.JSON(http.StatusInternalServerError, gin.H{
                "code":    "500",
                "message": "Gagal mengambil data member: " + err.Error(),
                "data":    nil,
            })
            return
        }

        message := "Berhasil mengambil semua data member"
        if keyword != "" {
            message = "Hasil pencarian untuk nama: " + keyword
        }

        c.JSON(http.StatusOK, gin.H{
            "code":    "200",
            "message": message,
            "data":    members,
        })
    }
}

// CreateMemberHandler
func CreateMemberHandler(db *sql.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var m model.Member
        
        if err := c.ShouldBindJSON(&m); err != nil {
            log.Printf("[ERROR] Payload member tidak valid: %v", err)
            
            c.JSON(http.StatusBadRequest, gin.H{
                "code":    "400",
                "message": "Format data tidak valid",
                "data":    nil,
            })
            return
        }

        // Cek Data Kosong atau Spasi
        if strings.TrimSpace(m.Name) == "" || strings.TrimSpace(m.Position) == "" {
            log.Printf("[VALIDATION ERROR] Upaya tambah member gagal: Nama atau Posisi hanya berisi spasi/kosong")

            c.JSON(http.StatusBadRequest, gin.H{
                "code":    "400",
                "message": "Nama dan Posisi harus diisi (tidak boleh kosong atau hanya spasi)!",
                "data":    nil,
            })
            return
        }

        if err := model.CreateMember(db, &m); err != nil {
            log.Printf("[DATABASE ERROR] Gagal menambahkan member ke DB: %v", err)
            
            c.JSON(http.StatusInternalServerError, gin.H{
                "code":    "500",
                "message": "Gagal menambahkan member: " + err.Error(),
                "data":    nil,
            })
            return
        }

        log.Printf("[SUCCESS] Member baru ditambahkan: %s (%s) dengan ID: %d", m.Name, m.Position, m.ID)

        c.JSON(http.StatusCreated, gin.H{
            "code":    "201",
            "message": "Member '" + m.Name + "' berhasil ditambahkan!",
            "data":    m,
        })
    }
}

// UpdateMemberHandler (PUT)
func UpdateMemberHandler(db *sql.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var uri MemberUri
        if err := c.ShouldBindUri(&uri); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{
                "code":    "400",
                "message": "ID member tidak valid (harus angka)",
                "data":    nil,
            })
            return
        }

        var m model.Member
        if err := c.ShouldBindJSON(&m); err != nil {
            log.Printf("[ERROR] Payload update member tidak valid: %v", err)
            c.JSON(http.StatusBadRequest, gin.H{
                "code":    "400",
                "message": "Format data tidak valid",
                "data":    nil,
            })
            return
        }

        m.ID = uri.ID

        if strings.TrimSpace(m.Name) == "" || strings.TrimSpace(m.Position) == "" {
            log.Printf("[VALIDATION ERROR] Gagal update: Nama atau Posisi kosong.")
            c.JSON(http.StatusBadRequest, gin.H{
                "code":    "400",
                "message": "Nama dan Posisi tidak boleh kosong!",
                "data":    nil,
            })
            return
        }

        if err := model.UpdateMember(db, &m); err != nil {
            log.Printf("[DATABASE ERROR] Gagal update member ID %d: %v", m.ID, err)
            c.JSON(http.StatusInternalServerError, gin.H{
                "code":    "500",
                "message": "Gagal memperbarui data: " + err.Error(),
                "data":    nil,
            })
            return
        }

        log.Printf("[SUCCESS] Member ID %d berhasil diperbarui menjadi: %s", m.ID, m.Name)

        c.JSON(http.StatusOK, gin.H{
            "code":    "200",
            "message": "Data member berhasil diperbarui!",
            "data":    m, 
        })
    }
}

// DeleteMemberHandler
func DeleteMemberHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		if err := model.DeleteMember(db, id); err != nil {
			// LOGGING ERROR: Mencatat gagal hapus
			log.Printf("[ERROR] Gagal menghapus member ID %s: %v", id, err)

			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    "500",
				"message": "Gagal menghapus member: " + err.Error(),
				"data":    nil,
			})
			return
		}

		// LOGGING SUCCESS
		log.Printf("[SUCCESS] Member dengan ID %s berhasil dihapus", id)

		c.JSON(http.StatusOK, gin.H{
			"code":    "200",
			"message": "Member dengan ID " + id + " telah berhasil dihapus",
			"data":    nil,
		})
	}
}