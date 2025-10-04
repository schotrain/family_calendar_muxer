package db_models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name  string `gorm:"not null;size:100"`
	Email string `gorm:"not null;uniqueIndex;size:255"`
	Age   int    `gorm:"not null"`
}
