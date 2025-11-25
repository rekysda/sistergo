package models

import (
	"time"

	"gorm.io/gorm"
)

// UserRole maps to user_role table
type UserRole struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Role      string         `gorm:"size:255;not null" json:"role"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Users              []User              `gorm:"foreignKey:RoleID" json:"users,omitempty"`
	MenuAccessRules    []UserAccessMenu    `gorm:"foreignKey:RoleID" json:"menu_access_rules,omitempty"`
	SubmenuAccessRules []UserAccessSubmenu `gorm:"foreignKey:RoleID" json:"submenu_access_rules,omitempty"`
}
