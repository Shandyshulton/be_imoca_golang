package model

import "time"

// Tabel services_settings (Singleton)
type ServiceSetting struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	SectionTitle string    `gorm:"type:varchar(255)" json:"section_title"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Tabel services_items
type ServiceItem struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	Icon            string    `gorm:"type:varchar(100)" json:"icon"`
	Name            string    `gorm:"type:varchar(255)" json:"name"`
	Description     string    `gorm:"type:text" json:"description"`
	FullDescription string    `gorm:"type:text" json:"full_description"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// Struct untuk Response/Request Gabungan
type ServicesResponse struct {
	SectionTitle string        `json:"section_title"`
	Items        []ServiceItem `json:"items"`
}

func (ServiceSetting) TableName() string {
	return "services_settings"
}

func (ServiceItem) TableName() string {
	return "services_items"
}