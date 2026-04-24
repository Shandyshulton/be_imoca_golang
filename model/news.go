package model

import (
	"time"
)

type News struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `gorm:"type:varchar(255);not null" json:"title" binding:"required"`
	Summary     string    `gorm:"type:text" json:"summary"`
	Source      string    `gorm:"type:varchar(255)" json:"source"`
	Badge       string    `gorm:"type:varchar(100);default:'REGULASI'" json:"badge"`
	URL         string    `gorm:"type:varchar(255)" json:"url"`
	Image       string    `gorm:"type:varchar(255)" json:"image"`
	PublishedAt time.Time `json:"published_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (News) TableName() string {
	return "news"
}