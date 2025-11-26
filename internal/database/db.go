package database

import (
	"fmt"

	"github.com/rekysda/sistergo/internal/config"
	"github.com/rekysda/sistergo/internal/models"
	"github.com/rekysda/sistergo/internal/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Student{},
		&models.UserRole{},
		&models.UserMenu{},
		&models.UserSubMenu{},
		&models.UserAccessMenu{},
		&models.UserAccessSubmenu{},
		&models.UserLog{},
		&models.WebSetting{},
	)
}

// Seed initial data if necessary: roles etc.
func Seed(db *gorm.DB) error {
	// Ensure admin and user roles exist
	var count int64
	db.Model(&models.UserRole{}).Count(&count)
	if count == 0 {
		roles := []models.UserRole{{Role: "admin"}, {Role: "user"}}
		if err := db.Create(&roles).Error; err != nil {
			return err
		}
	}
	// Ensure an admin user exists
	var adminUser models.User
	if err := db.Where("username = ?", "admin").First(&adminUser).Error; err == gorm.ErrRecordNotFound {
		// Create admin user (default password: password)
		hashed, err := utils.HashPassword("password")
		if err != nil {
			return err
		}
		admin := models.User{Name: "Administrator", Username: "admin", Email: "admin@example.com", Password: hashed, RoleID: 1, IsActive: true}
		if err := db.Create(&admin).Error; err != nil {
			return err
		}
	}
	return nil
}
