package model

import (
	"database/sql"
	"fmt"
)

type ServiceItem struct {
	ID              int    `json:"id"`
	Icon            string `json:"icon"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	FullDescription string `json:"full_description"`
}

type ServicesSectionResponse struct {
	SectionTitle string        `json:"section_title"`
	Items        []ServiceItem `json:"items"`
}

// GetServices 
func GetServices(db *sql.DB) (ServicesSectionResponse, error) {
	var resp ServicesSectionResponse
	resp.Items = []ServiceItem{} 

	err := db.QueryRow("SELECT COALESCE(section_title, '') FROM services_settings WHERE id = 1").Scan(&resp.SectionTitle)
	if err != nil && err != sql.ErrNoRows {
		return resp, err
	}

	rows, err := db.Query("SELECT id, icon, name, description, full_description FROM services_items ORDER BY id ASC")
	if err != nil {
		return resp, err
	}
	defer rows.Close()

	for rows.Next() {
		var item ServiceItem
		err := rows.Scan(&item.ID, &item.Icon, &item.Name, &item.Description, &item.FullDescription)
		if err != nil {
			return resp, err
		}
		resp.Items = append(resp.Items, item)
	}

	return resp, nil
}

// UpdateServices
func UpdateServices(db *sql.DB, s *ServicesSectionResponse) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE services_settings SET section_title = ? WHERE id = 1", s.SectionTitle)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal update settings: %v", err)
	}

	query := `UPDATE services_items SET icon = ?, name = ?, description = ?, full_description = ? WHERE id = ?`
	
	for _, item := range s.Items {
		if item.ID == 0 {
			continue 
		}

		_, err := tx.Exec(query, item.Icon, item.Name, item.Description, item.FullDescription, item.ID)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("gagal update item cluster ID %d: %v", item.ID, err)
		}
	}

	return tx.Commit()
}