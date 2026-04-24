package model

import (
	"time"
)

type Member struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(255);not null" json:"name" binding:"required"`
	Email       string    `gorm:"type:varchar(255)" json:"email"`
	Phone       string    `gorm:"type:varchar(50)" json:"phone"`
	Description string    `gorm:"type:text" json:"description"`
	Link        string    `gorm:"type:varchar(255)" json:"link"`
	Image       string    `gorm:"type:varchar(255)" json:"image"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Member) TableName() string {
	return "members"
}