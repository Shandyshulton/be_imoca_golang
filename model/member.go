package model

import (
	"database/sql"
	"errors" 
)

type Member struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Position string `json:"position"`
}

func GetAllMembers(db *sql.DB, keyword string) ([]Member, error) {
    var rows *sql.Rows
    var err error

    if keyword != "" {
        query := "SELECT id, name, position FROM members WHERE name LIKE ?"
        rows, err = db.Query(query, "%"+keyword+"%")
    } else {
        query := "SELECT id, name, position FROM members"
        rows, err = db.Query(query)
    }

    if err != nil {
        return nil, err
    }
    defer rows.Close()

    members := []Member{} 

    for rows.Next() {
        var m Member
        if err := rows.Scan(&m.ID, &m.Name, &m.Position); err != nil {
            return nil, err
        }
        members = append(members, m)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }
    return members, nil
}

// Create Member
func CreateMember(db *sql.DB, m *Member) error { 
    query := "INSERT INTO members (name, position) VALUES (?, ?)"
    result, err := db.Exec(query, m.Name, m.Position)
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

// UpdateMember - Hapus parameter string id
func UpdateMember(db *sql.DB, m *Member) error {
    query := `UPDATE members SET name = ?, position = ? WHERE id = ?`
    
    // Langsung ambil m.ID dari struct
    result, err := db.Exec(query, m.Name, m.Position, m.ID)
    if err != nil {
        return err
    }

    rowsAffected, _ := result.RowsAffected()
    if rowsAffected == 0 {
        return errors.New("tidak ada data yang diubah atau ID tidak ditemukan")
    }

    return nil
}

// Delete Member
func DeleteMember(db *sql.DB, id string) error {
	query := "DELETE FROM members WHERE id = ?"
	result, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	// Cek apakah ada baris yang benar-benar terhapus
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("data member tidak ditemukan")
	}

	return nil
}