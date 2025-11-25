package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rekysda/sistergo/internal/models"
	"gorm.io/gorm"
)

func ListRoles(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var roles []models.UserRole
		if err := db.Find(&roles).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"roles": roles})
	}
}

func CreateRole(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var role models.UserRole
		if err := c.ShouldBindJSON(&role); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := db.Create(&role).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"role": role})
	}
}

func GetRole(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var role models.UserRole
		if err := db.Preload("MenuAccessRules").Preload("SubmenuAccessRules").First(&role, c.Param("id")).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "role not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"role": role})
	}
}

func UpdateRole(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var role models.UserRole
		if err := db.First(&role, c.Param("id")).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "role not found"})
			return
		}
		var payload models.UserRole
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		role.Role = payload.Role
		if err := db.Save(&role).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"role": role})
	}
}

func DeleteRole(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := db.Delete(&models.UserRole{}, c.Param("id")).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"deleted": true})
	}
}
