package models

import (
	"time"

	"gorm.io/gorm"
)

// UserAccessSubmenu maps to user_access_submenu table
type UserAccessSubmenu struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	RoleID    uint           `json:"role_id"`
	SubmenuID uint           `json:"submenu_id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Role    UserRole    `gorm:"foreignKey:RoleID" json:"role,omitempty"`
	Submenu UserSubMenu `gorm:"foreignKey:SubmenuID" json:"submenu,omitempty"`
}
