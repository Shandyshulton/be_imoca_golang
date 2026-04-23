package model

import (
	"database/sql"
	"errors"
	"time"
)

type News struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Summary     string    `json:"summary"`
	Source      string    `json:"source"`
	Badge       string    `json:"badge"` 
	URL         string    `json:"url"`
	Image       string    `json:"image"`
	PublishedAt time.Time `json:"published_at"`
}

// Get News
func GetAllNews(db *sql.DB, keyword string) ([]News, error) {
	var rows *sql.Rows
	var err error

	baseQuery := `SELECT id, title, COALESCE(summary, ''), COALESCE(source, ''), COALESCE(badge, 'REGULASI'), url, image, published_at FROM news`

	if keyword != "" {
		query := baseQuery + " WHERE title LIKE ? OR summary LIKE ? ORDER BY published_at DESC"
		searchPattern := "%" + keyword + "%"
		rows, err = db.Query(query, searchPattern, searchPattern)
	} else {
		query := baseQuery + " ORDER BY published_at DESC"
		rows, err = db.Query(query)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	newsList := []News{}
	for rows.Next() {
		var n News
		err := rows.Scan(&n.ID, &n.Title, &n.Summary, &n.Source, &n.Badge, &n.URL, &n.Image, &n.PublishedAt)
		if err != nil {
			return nil, err
		}
		newsList = append(newsList, n)
	}
	return newsList, nil
}

// Create News
func CreateNews(db *sql.DB, n *News) error {
	if n.Badge == "" {
		n.Badge = "REGULASI"
	}

	// Gunakan PublishedAt dari input, fallback ke sekarang kalau zero
	if n.PublishedAt.IsZero() {
		n.PublishedAt = time.Now()
	}

	query := "INSERT INTO news (title, summary, source, badge, url, image, published_at) VALUES (?, ?, ?, ?, ?, ?, ?)"
	result, err := db.Exec(query, n.Title, n.Summary, n.Source, n.Badge, n.URL, n.Image, n.PublishedAt)
	if err != nil {
		return err
	}

	lastID, _ := result.LastInsertId()
	n.ID = int(lastID)

	return nil
}

// Update News
func UpdateNews(db *sql.DB, n *News) error {
	var exists int
	checkQuery := "SELECT COUNT(*) FROM news WHERE id = ?"
	err := db.QueryRow(checkQuery, n.ID).Scan(&exists)
	if err != nil {
		return err
	}

	if exists == 0 {
		return errors.New("data berita tidak ditemukan di database")
	}

	query := `UPDATE news SET title = ?, summary = ?, source = ?, badge = ?, url = ?, image = ?, published_at = ? WHERE id = ?`
	_, err = db.Exec(query, n.Title, n.Summary, n.Source, n.Badge, n.URL, n.Image, n.PublishedAt, n.ID)
	if err != nil {
		return err
	}

	return nil
}

// Get News ID
func GetNewsByID(db *sql.DB, id int) (News, error) {
	var n News
	query := "SELECT id, title, COALESCE(summary, ''), COALESCE(source, ''), COALESCE(badge, 'REGULASI'), url, image, published_at FROM news WHERE id = ?"
	err := db.QueryRow(query, id).Scan(&n.ID, &n.Title, &n.Summary, &n.Source, &n.Badge, &n.URL, &n.Image, &n.PublishedAt)
	return n, err
}

// Delete News 
func DeleteNews(db *sql.DB, id int) error {
	query := "DELETE FROM news WHERE id = ?"
	result, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("data berita tidak ditemukan")
	}

	return nil
}