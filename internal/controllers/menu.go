package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rekysda/sistergo/internal/models"
	"gorm.io/gorm"
)

func ListMenus(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var menus []models.UserMenu
		if err := db.Preload("Submenus").Find(&menus).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"menus": menus})
	}
}

func CreateMenu(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var menu models.UserMenu
		if err := c.ShouldBindJSON(&menu); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := db.Create(&menu).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"menu": menu})
	}
}

func UpdateMenu(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var menu models.UserMenu
		if err := db.First(&menu, c.Param("id")).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "menu not found"})
			return
		}
		var payload models.UserMenu
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		menu.Menu = payload.Menu
		menu.Icon = payload.Icon
		menu.Ordering = payload.Ordering
		if err := db.Save(&menu).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"menu": menu})
	}
}

func DeleteMenu(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := db.Delete(&models.UserMenu{}, c.Param("id")).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"deleted": true})
	}
}

// Submenu
func CreateSubmenu(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var submenu models.UserSubMenu
		if err := c.ShouldBindJSON(&submenu); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := db.Create(&submenu).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"submenu": submenu})
	}
}

func UpdateSubmenu(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var submenu models.UserSubMenu
		if err := db.First(&submenu, c.Param("id")).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "submenu not found"})
			return
		}
		var payload models.UserSubMenu
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		submenu.Title = payload.Title
		submenu.URL = payload.URL
		submenu.IsActive = payload.IsActive
		submenu.Ordering = payload.Ordering
		if err := db.Save(&submenu).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"submenu": submenu})
	}
}

func DeleteSubmenu(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := db.Delete(&models.UserSubMenu{}, c.Param("id")).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"deleted": true})
	}
}
