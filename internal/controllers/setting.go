package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rekysda/sistergo/internal/models"
	"gorm.io/gorm"
)

func GetSettings(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var settings []models.WebSetting
		if err := db.Find(&settings).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"settings": settings})
	}
}

func UpdateSetting(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload models.WebSetting
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var setting models.WebSetting
		if err := db.Where("name = ?", payload.Name).First(&setting).Error; err != nil {
			// if not exists, create
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&payload).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusOK, gin.H{"setting": payload})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		setting.Value = payload.Value
		setting.IsActive = payload.IsActive
		if err := db.Save(&setting).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"setting": setting})
	}
}
