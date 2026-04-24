package model

import "time"

type MissionSection struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	SectionTitle string    `gorm:"type:varchar(255)" json:"section_title"`
	Description  string    `gorm:"type:text" json:"description"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (MissionSection) TableName() string {
	return "mission_section"
}