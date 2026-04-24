package model

import (
	"time"
)

type Partner struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(255);not null" json:"name" binding:"required"`
	Email       string    `gorm:"type:varchar(255);not null" json:"email" binding:"required"`
	Phone       string    `gorm:"type:varchar(50)" json:"phone"`
	Description string    `gorm:"type:text" json:"description"`
	Link        string    `gorm:"type:varchar(255)" json:"link"`
	Image       string    `gorm:"type:varchar(255)" json:"image"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Partner) TableName() string {
	return "partners"
}