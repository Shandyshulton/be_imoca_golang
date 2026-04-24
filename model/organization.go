package model

import (
	"time"
)

type OrganizationStructure struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	FullName     string    `gorm:"type:varchar(255);not null" json:"full_name" binding:"required"`
	JobPosition  string    `gorm:"type:varchar(255);not null" json:"job_position" binding:"required"`
	TeamCategory string    `gorm:"type:enum('Advisory Board', 'Team Member');not null" json:"team_category" binding:"required"`
	ProfilePhoto string    `gorm:"type:varchar(255)" json:"profile_photo"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (OrganizationStructure) TableName() string {
	return "organization_structures"
}