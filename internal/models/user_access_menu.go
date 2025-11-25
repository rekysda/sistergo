package models

import (
	"time"

	"gorm.io/gorm"
)

// UserAccessMenu maps to user_access_menu table
type UserAccessMenu struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	RoleID    uint           `json:"role_id"`
	MenuID    uint           `json:"menu_id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Role UserRole `gorm:"foreignKey:RoleID" json:"role,omitempty"`
	Menu UserMenu `gorm:"foreignKey:MenuID" json:"menu,omitempty"`
}
