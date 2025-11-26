package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rekysda/sistergo/internal/models"
	"gorm.io/gorm"
)

func Dashboard(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDVal, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		id, ok := userIDVal.(uint)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		var user models.User
		if err := db.Preload("Role").First(&user, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		var recentLogs []models.UserLog
		if user.RoleID == 1 {
			// Admin sees all logs
			db.Preload("User").Order("created_at desc").Limit(10).Find(&recentLogs)
		} else {
			// Regular users only see their own logs
			db.Preload("User").Where("user_id = ?", user.ID).Order("created_at desc").Limit(10).Find(&recentLogs)
		}

		var recentUsers []models.User
		db.Order("date_created desc").Limit(8).Find(&recentUsers)

		// System info (simplified)
		systemInfo := gin.H{
			"go_version": "1.21",
			"environment": "development",
			"database": "PostgreSQL",
		}

		c.JSON(http.StatusOK, gin.H{
			"recent_logs": recentLogs,
			"recent_users": recentUsers,
			"user": user,
			"system_info": systemInfo,
		})
	}
}

func UserActivity(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var totalUsers int64
		var totalActiveUsers int64
		var totalLoginActivities int64

		db.Model(&models.User{}).Count(&totalUsers)
		db.Model(&models.User{}).Where("is_active = ?", true).Count(&totalActiveUsers)
		db.Model(&models.UserLog{}).Where("activity LIKE ?", "%login%").Count(&totalLoginActivities)

		var users []models.User
		db.Preload("Role").Where("is_active = ?", true).Order("name").Limit(50).Find(&users)

		// Add last login to each user (simplified, in real app use subquery)
		for i := range users {
			var lastLogin models.UserLog
			db.Where("user_id = ? AND activity LIKE ?", users[i].ID, "%login%").Order("created_at desc").First(&lastLogin)
			users[i].DateCreated = &lastLogin.CreatedAt
		}

		c.JSON(http.StatusOK, gin.H{
			"users": users,
			"total_users": totalUsers,
			"total_active_users": totalActiveUsers,
			"total_login_activities": totalLoginActivities,
		})
	}
}

func Logs(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var logs []models.UserLog
		db.Preload("User").Order("created_at desc").Limit(50).Find(&logs)

		c.JSON(http.StatusOK, gin.H{"logs": logs})
	}
}
