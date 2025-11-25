package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rekysda/sistergo/internal/models"
	"gorm.io/gorm"
)

// Dashboard returns summary statistics and recent logs
func Dashboard(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDVal, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		userID, ok := userIDVal.(uint)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		// Get current user
		var currentUser models.User
		if err := db.Preload("Role").First(&currentUser, userID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		// Get recent logs based on role
		var recentLogs []models.UserLog
		if currentUser.RoleID == 1 {
			// Admin sees all logs
			db.Preload("User").Order("created_at DESC").Limit(10).Find(&recentLogs)
		} else {
			// Regular users only see their own logs
			db.Preload("User").Where("user_id = ?", userID).Order("created_at DESC").Limit(10).Find(&recentLogs)
		}

		// Get recent users for admin dashboard
		var recentUsers []models.User
		db.Order("date_created DESC").Limit(8).Find(&recentUsers)

		// Get counts
		var totalUsers, totalActiveUsers, totalRoles, totalMenus int64
		db.Model(&models.User{}).Count(&totalUsers)
		db.Model(&models.User{}).Where("is_active = ?", true).Count(&totalActiveUsers)
		db.Model(&models.UserRole{}).Count(&totalRoles)
		db.Model(&models.UserMenu{}).Count(&totalMenus)

		c.JSON(http.StatusOK, gin.H{
			"user":               currentUser,
			"recent_logs":        recentLogs,
			"recent_users":       recentUsers,
			"total_users":        totalUsers,
			"total_active_users": totalActiveUsers,
			"total_roles":        totalRoles,
			"total_menus":        totalMenus,
		})
	}
}
