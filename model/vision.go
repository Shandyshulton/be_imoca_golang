package model

import "time"

// Untuk tabel vision_settings (Singleton)
type VisionSetting struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	SectionTitle string    `gorm:"type:varchar(255)" json:"section_title"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Untuk tabel vision_items
type VisionItem struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Icon        string    `gorm:"type:varchar(100)" json:"icon"`
	Title       string    `gorm:"type:varchar(255)" json:"title"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Struct untuk Response Gabungan
type VisionResponse struct {
	SectionTitle string       `json:"section_title"`
	Items        []VisionItem `json:"items"`
}