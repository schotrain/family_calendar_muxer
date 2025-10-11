package models

import "gorm.io/gorm"

type CalendarMux struct {
	gorm.Model
	CreatedByID uint   `gorm:"not null;index"`
	CreatedBy   User   `gorm:"foreignKey:CreatedByID;constraint:OnDelete:CASCADE"`
	Name        string `gorm:"not null;size:200"`
	Description string `gorm:"size:1000"`
}
