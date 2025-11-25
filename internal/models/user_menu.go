package models

import (
	"time"

	"gorm.io/gorm"
)

// UserMenu maps to user_menu table
type UserMenu struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Menu      string         `gorm:"size:255;not null" json:"menu"`
	Icon      *string        `gorm:"size:255" json:"icon,omitempty"`
	Ordering  int            `json:"ordering"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Submenus    []UserSubMenu    `gorm:"foreignKey:MenuID" json:"submenus,omitempty"`
	AccessRules []UserAccessMenu `gorm:"foreignKey:MenuID" json:"access_rules,omitempty"`
}
