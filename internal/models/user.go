package models

import (
	"time"

	"gorm.io/gorm"
)

// User model. Minimal fields for auth demonstration
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `json:"name" gorm:"size:255;not null"`
	Username  string         `json:"username" gorm:"size:255;unique;not null"`
	Email     string         `json:"email" gorm:"size:255;unique;not null"`
	Image     *string        `json:"image,omitempty" gorm:"size:255"`
	Password  string         `json:"-" gorm:"size:255;not null"`
	RoleID    uint           `json:"role_id"`
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	DateCreated *time.Time   `json:"date_created" gorm:"column:date_created"`
	RememberToken *string    `json:"-" gorm:"size:100;column:remember_token"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Role UserRole `gorm:"foreignKey:RoleID" json:"role,omitempty"`
	Logs []UserLog `gorm:"foreignKey:UserID" json:"logs,omitempty"`
}
