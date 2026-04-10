package model

import (
	"database/sql"
	"errors" 
	"time"
)

type News struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	Image       string    `json:"image"`
	PublishedAt time.Time `json:"published_at"`
}

// Get All News
func GetAllNews(db *sql.DB, keyword string) ([]News, error) {
	var rows *sql.Rows
	var err error

	if keyword != "" {
		query := "SELECT id, title, content, image, published_at FROM news WHERE title LIKE ? OR content LIKE ? ORDER BY published_at DESC"
		searchPattern := "%" + keyword + "%"
		rows, err = db.Query(query, searchPattern, searchPattern)
	} else {
		query := "SELECT id, title, content, image, published_at FROM news ORDER BY published_at DESC"
		rows, err = db.Query(query)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	newsList := []News{} 
	for rows.Next() {
		var n News
		err := rows.Scan(&n.ID, &n.Title, &n.Content, &n.Image, &n.PublishedAt)
		if err != nil {
			return nil, err
		}
		newsList = append(newsList, n)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return newsList, nil
}

// Create News
func CreateNews(db *sql.DB, n *News) error {
    now := time.Now() 
    query := "INSERT INTO news (title, content, image, published_at) VALUES (?, ?, ?, ?)"
    
    result, err := db.Exec(query, n.Title, n.Content, n.Image, now)
    if err != nil {
        return err
    }

    lastID, err := result.LastInsertId()
    if err != nil {
        return err
    }

    n.ID = int(lastID)
    n.PublishedAt = now 

    return nil
}

// Update News
func UpdateNews(db *sql.DB, n *News) error {
	query := `UPDATE news SET title = ?, content = ?, image = ? WHERE id = ?`
	
	result, err := db.Exec(query, n.Title, n.Content, n.Image, n.ID)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("berita tidak ditemukan atau tidak ada perubahan data")
	}

	return nil
}

func GetNewsByID(db *sql.DB, id int) (News, error) {
	var n News
	query := "SELECT id, title, content, image FROM news WHERE id = ?"
	err := db.QueryRow(query, id).Scan(&n.ID, &n.Title, &n.Content, &n.Image)
	return n, err
}

// Delete News
func DeleteNews(db *sql.DB, id string) error {
	query := "DELETE FROM news WHERE id = ?"
	result, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	// Cek apakah ada baris yang benar-benar terhapus
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("data berita tidak ditemukan")
	}

	return nil
}