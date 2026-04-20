package model

import (
	"database/sql"
	"fmt"
)

type VisionItem struct {
	ID          int    `json:"id"`
	Icon        string `json:"icon"`        // Emoji
	Title       string `json:"title"`       
	Description string `json:"description"` 
}

type VisionSectionResponse struct {
	SectionTitle string       `json:"section_title"`
	Items        []VisionItem `json:"items"`
}

// GetVision 
func GetVision(db *sql.DB) (VisionSectionResponse, error) {
	var v VisionSectionResponse
	v.Items = []VisionItem{} 

	err := db.QueryRow("SELECT COALESCE(section_title, '') FROM vision_settings WHERE id = 1").Scan(&v.SectionTitle)
	if err != nil && err != sql.ErrNoRows {
		return v, err
	}

	rows, err := db.Query("SELECT id, icon, title, description FROM vision_items ORDER BY id ASC")
	if err != nil {
		return v, err
	}
	defer rows.Close()

	for rows.Next() {
		var item VisionItem
		err := rows.Scan(&item.ID, &item.Icon, &item.Title, &item.Description)
		if err != nil {
			return v, err
		}
		v.Items = append(v.Items, item)
	}

	return v, nil
}

// UpdateVision 
func UpdateVision(db *sql.DB, v *VisionSectionResponse) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE vision_settings SET section_title = ? WHERE id = 1", v.SectionTitle)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal update vision_settings: %v", err)
	}

	query := `UPDATE vision_items SET icon = ?, title = ?, description = ? WHERE id = ?`
	
	for _, item := range v.Items {
		if item.ID == 0 {
			continue 
		}

		_, err := tx.Exec(query, item.Icon, item.Title, item.Description, item.ID)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("gagal update vision item ID %d: %v", item.ID, err)
		}
	}

	return tx.Commit()
}