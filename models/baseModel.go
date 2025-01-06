package models

import (
	"time"
)

// AuditFields struct that references User for audit information
type AuditFields struct {
	CreatedAt time.Time `gorm:"autoCreateTime"` // Automatically set when record is created
	CreatedBy *uint     // Foreign Key to the User table (nullable)
	Creator   *User     `gorm:"foreignKey:CreatedBy"` // Reference to User for Creator

	UpdatedAt time.Time `gorm:"autoUpdateTime"` // Automatically set when record is updated
	UpdatedBy *uint     // Foreign Key to the User table (nullable)
	Updator   *User     `gorm:"foreignKey:UpdatedBy"` // Reference to User for Updator

	IsActive bool `gorm:"default:true"` // Active flag
}
