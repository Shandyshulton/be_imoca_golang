package model

import (
	"database/sql"
	"errors"
)

type Member struct {
	ID          int    `json:"id" form:"id"`
	Name        string `json:"name" form:"name"`
	Email       string `json:"email" form:"email"`
	Phone       string `json:"phone" form:"phone"`
	Description string `json:"description" form:"description"`
	Link        string `json:"link" form:"link"`
	Image       string `json:"image" form:"image"`
}

// Get All Members 
func GetAllMembers(db *sql.DB, keyword string) ([]Member, error) {
	var rows *sql.Rows
	var err error

	baseQuery := `SELECT id, name, email, 
                  COALESCE(phone, ''), 
                  COALESCE(description, ''), 
                  COALESCE(link, ''), 
                  COALESCE(image, '') 
                  FROM members`

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

	members := []Member{}
	for rows.Next() {
		var m Member
		if err := rows.Scan(&m.ID, &m.Name, &m.Email, &m.Phone, &m.Description, &m.Link, &m.Image); err != nil {
			return nil, err
		}
		members = append(members, m)
	}

	return members, nil
}

// Get Member ID 
func GetMemberByID(db *sql.DB, id int) (Member, error) {
	var m Member
	query := `SELECT id, name, email, 
              COALESCE(phone, ''), 
              COALESCE(description, ''), 
              COALESCE(link, ''), 
              COALESCE(image, '') 
              FROM members WHERE id = ?`

	err := db.QueryRow(query, id).Scan(&m.ID, &m.Name, &m.Email, &m.Phone, &m.Description, &m.Link, &m.Image)
	return m, err
}

// Create Member 
func CreateMember(db *sql.DB, m *Member) error {
	query := "INSERT INTO members (name, email, phone, description, link, image) VALUES (?, ?, ?, ?, ?, ?)"

	result, err := db.Exec(query, m.Name, m.Email, m.Phone, m.Description, m.Link, m.Image)
	if err != nil {
		return err
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	m.ID = int(lastID)
	return nil
}

// Update Member 
func UpdateMember(db *sql.DB, m *Member) error {
	var exists int
	err := db.QueryRow("SELECT COUNT(*) FROM members WHERE id = ?", m.ID).Scan(&exists)
	if err != nil {
		return err
	}
	if exists == 0 {
		return errors.New("member tidak ditemukan")
	}

	query := `UPDATE members SET name = ?, email = ?, phone = ?, description = ?, link = ?, image = ? WHERE id = ?`
	_, err = db.Exec(query, m.Name, m.Email, m.Phone, m.Description, m.Link, m.Image, m.ID)

	return err
}

// Delete Member 
func DeleteMember(db *sql.DB, id int) error {
	query := "DELETE FROM members WHERE id = ?"
	result, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("data member tidak ditemukan")
	}

	return nil
}