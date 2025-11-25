package models

import (
	"time"

	"gorm.io/gorm"
)

// UserSubMenu maps to user_sub_menu table
type UserSubMenu struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	MenuID    uint           `json:"menu_id"`
	Title     string         `gorm:"size:255;not null" json:"title"`
	URL       string         `gorm:"size:1024;not null" json:"url"`
	IsActive  bool           `json:"is_active"`
	Ordering  int            `json:"ordering"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Menu        UserMenu            `gorm:"foreignKey:MenuID" json:"menu,omitempty"`
	AccessRules []UserAccessSubmenu `gorm:"foreignKey:SubmenuID" json:"access_rules,omitempty"`
}
