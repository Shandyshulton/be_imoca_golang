package model

import (
	"time"
)

type Hero struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	BadgeText      string    `gorm:"type:varchar(255)" json:"badge_text"`
	Title          string    `gorm:"type:varchar(255)" json:"title"`
	TitleHighlight string    `gorm:"type:varchar(255)" json:"title_highlight"`
	Description    string    `gorm:"type:text" json:"description"`
	Image          string    `gorm:"type:varchar(255)" json:"image"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type HeroImage struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	HeroID    uint      `json:"hero_id"`
	Image     string    `gorm:"type:varchar(255)" json:"image"`
	SortOrder int       `gorm:"default:0" json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
}

func (Hero) TableName() string {
	return "hero_section"
}