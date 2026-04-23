package model

import (
	"database/sql"
	"errors"
)

type Hero struct {
    ID             int    `json:"id" form:"id"`
    BadgeText      string `json:"badge_text" form:"badge_text"`
    Title          string `json:"title" form:"title"`
    TitleHighlight string `json:"title_highlight" form:"title_highlight"`
    Description    string `json:"description" form:"description"`
    Image          string `json:"image"`
}

// Get Hero
func GetHero(db *sql.DB) (Hero, error) {
	var h Hero
	query := `SELECT id, badge_text, title, title_highlight, description, image FROM hero_section LIMIT 1`
	err := db.QueryRow(query).Scan(&h.ID, &h.BadgeText, &h.Title, &h.TitleHighlight, &h.Description, &h.Image)
	
	if err == sql.ErrNoRows {
		return h, errors.New("data hero belum ada di database")
	}
	return h, err
}

// Update Hero
func UpdateHero(db *sql.DB, h *Hero) error {
	query := `UPDATE hero_section SET badge_text = ?, title = ?, title_highlight = ?, description = ?, image = ? WHERE id = ?`
	_, err := db.Exec(query, h.BadgeText, h.Title, h.TitleHighlight, h.Description, h.Image, h.ID)
	return err
}

// ── Hero Images ──────────────────────────────────────────────────────────────

type HeroImage struct {
	ID     int    `json:"id"`
	HeroID int    `json:"hero_id"`
	Image  string `json:"image"`
	Order  int    `json:"order"`
}

// Get all images for a hero
func GetHeroImages(db *sql.DB, heroID int) ([]HeroImage, error) {
	query := `SELECT id, hero_id, image, sort_order FROM hero_images WHERE hero_id = ? ORDER BY sort_order ASC`
	rows, err := db.Query(query, heroID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []HeroImage
	for rows.Next() {
		var img HeroImage
		if err := rows.Scan(&img.ID, &img.HeroID, &img.Image, &img.Order); err != nil {
			return nil, err
		}
		images = append(images, img)
	}
	return images, nil
}

// Add a new image
func AddHeroImage(db *sql.DB, img *HeroImage) error {
	query := `INSERT INTO hero_images (hero_id, image, sort_order) VALUES (?, ?, ?)`
	result, err := db.Exec(query, img.HeroID, img.Image, img.Order)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	img.ID = int(id)
	return nil
}

// Delete an image by ID
func DeleteHeroImage(db *sql.DB, imageID int) (string, error) {
	// ambil nama file dulu sebelum dihapus (untuk hapus file fisik)
	var filename string
	err := db.QueryRow(`SELECT image FROM hero_images WHERE id = ?`, imageID).Scan(&filename)
	if err == sql.ErrNoRows {
		return "", errors.New("gambar tidak ditemukan")
	}
	if err != nil {
		return "", err
	}

	_, err = db.Exec(`DELETE FROM hero_images WHERE id = ?`, imageID)
	return filename, err
}