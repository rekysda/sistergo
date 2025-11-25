package models

import (
	"time"

	"gorm.io/gorm"
)

// Student model similar to a sample resource in Laravel app
type Student struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Name  string `json:"name" gorm:"size:255;not null"`
	Email string `json:"email" gorm:"size:255;unique;not null"`
	NIS   string `json:"nis" gorm:"size:50;unique;not null"`
}
