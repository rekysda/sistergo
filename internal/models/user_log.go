package models

import (
	"time"

	"gorm.io/gorm"
)

// UserLog maps to user_log table
type UserLog struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `json:"user_id"`
	Activity  string         `gorm:"type:text" json:"activity"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
