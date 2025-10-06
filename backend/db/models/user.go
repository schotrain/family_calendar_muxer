package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	GivenName      string `gorm:"not null;size:100"`
	FamilyName     string `gorm:"not null;size:100"`
	Email          string `gorm:"not null;size:255"`
	AuthProvider   string `gorm:"not null;size:50;index:idx_auth_provider_id;check:auth_provider <> ''"`
	AuthProviderID string `gorm:"not null;size:255;index:idx_auth_provider_id;check:auth_provider_id <> ''"`
}
