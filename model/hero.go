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