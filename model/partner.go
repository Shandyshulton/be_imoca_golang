package model

import (
	"database/sql"
	"errors"
)

type Partner struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Description string `json:"description"`
	Link        string `json:"link"`
	Image       string `json:"image"`
}

// Get All Partners 
func GetAllPartners(db *sql.DB, keyword string) ([]Partner, error) {
	var rows *sql.Rows
	var err error

	baseQuery := `SELECT id, name, email, 
				  COALESCE(phone, ''), 
				  COALESCE(description, ''), 
				  COALESCE(link, ''), 
				  COALESCE(image, '') 
				  FROM partners`

	if keyword != "" {
		query := baseQuery + " WHERE name LIKE ? OR description LIKE ?"
		searchPattern := "%" + keyword + "%"
		rows, err = db.Query(query, searchPattern, searchPattern)
	} else {
		rows, err = db.Query(baseQuery)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	partners := []Partner{} 
	for rows.Next() {
		var p Partner
		if err := rows.Scan(&p.ID, &p.Name, &p.Email, &p.Phone, &p.Description, &p.Link, &p.Image); err != nil {
			return nil, err
		}
		partners = append(partners, p)
	}

	return partners, nil
}

// GetPartnerByID 
func GetPartnerByID(db *sql.DB, id int) (Partner, error) {
	var p Partner
	query := `SELECT id, name, email, 
			  COALESCE(phone, ''), 
			  COALESCE(description, ''), 
			  COALESCE(link, ''), 
			  COALESCE(image, '') 
			  FROM partners WHERE id = ?`
	
	err := db.QueryRow(query, id).Scan(&p.ID, &p.Name, &p.Email, &p.Phone, &p.Description, &p.Link, &p.Image)
	return p, err
}

// Create Partner 
func CreatePartner(db *sql.DB, p *Partner) error { 
	query := "INSERT INTO partners (name, email, phone, description, link, image) VALUES (?, ?, ?, ?, ?, ?)"
	
	result, err := db.Exec(query, p.Name, p.Email, p.Phone, p.Description, p.Link, p.Image)
	if err != nil {
		return err
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	p.ID = int(lastID)
	return nil
}

// UpdatePartner 
func UpdatePartner(db *sql.DB, p *Partner) error {
	var exists int
	err := db.QueryRow("SELECT COUNT(*) FROM partners WHERE id = ?", p.ID).Scan(&exists)
	if err != nil {
		return err
	}
	if exists == 0 {
		return errors.New("partner tidak ditemukan")
	}

	query := `UPDATE partners SET name = ?, email = ?, phone = ?, description = ?, link = ?, image = ? WHERE id = ?`
	_, err = db.Exec(query, p.Name, p.Email, p.Phone, p.Description, p.Link, p.Image, p.ID)
	
	return err
}

// DeletePartner 
func DeletePartner(db *sql.DB, id int) error {
	query := "DELETE FROM partners WHERE id = ?"
	result, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("data partner tidak ditemukan")
	}

	return nil
}