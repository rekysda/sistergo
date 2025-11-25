package models

import (
	"time"

	"gorm.io/gorm"
)

// WebSetting maps to web_setting

type WebSetting struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:255;not null;unique" json:"name"`
	Value     string         `gorm:"type:text" json:"value"`
	IsActive  bool           `json:"is_active"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (ws *WebSetting) GetValue(db *gorm.DB, name string) (string, error) {
	var setting WebSetting
	if err := db.Where("name = ?", name).First(&setting).Error; err != nil {
		return "", err
	}
	return setting.Value, nil
}
