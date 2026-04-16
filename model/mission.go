package model

import "database/sql"

type MissionSection struct {
	ID           int    `json:"id"`
	SectionTitle string `json:"section_title"`
	Description  string `json:"description"`
}

func GetMission(db *sql.DB) (MissionSection, error) {
	var m MissionSection
	query := `SELECT id, COALESCE(section_title, ''), COALESCE(description, '') 
			  FROM mission_section LIMIT 1`
	err := db.QueryRow(query).Scan(&m.ID, &m.SectionTitle, &m.Description)
	return m, err
}

func UpdateMission(db *sql.DB, m *MissionSection) error {
	query := `UPDATE mission_section SET section_title = ?, description = ? WHERE id = ?`
	_, err := db.Exec(query, m.SectionTitle, m.Description, 1) // Paksa update ID 1
	return err
}