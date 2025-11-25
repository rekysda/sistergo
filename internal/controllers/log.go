package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rekysda/sistergo/internal/models"
	"gorm.io/gorm"
)

// ListLogs returns a list of all user activity logs with pagination
func ListLogs(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
		if page < 1 {
			page = 1
		}
		if limit < 1 || limit > 100 {
			limit = 50
		}
		offset := (page - 1) * limit

		var logs []models.UserLog
		var total int64

		db.Model(&models.UserLog{}).Count(&total)
		if err := db.Preload("User").Order("created_at desc").Offset(offset).Limit(limit).Find(&logs).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"logs":  logs,
			"total": total,
			"page":  page,
			"limit": limit,
		})
	}
}

// UserActivityStats represents user activity statistics
type UserActivityStats struct {
	TotalUsers           int64 `json:"total_users"`
	TotalActiveUsers     int64 `json:"total_active_users"`
	TotalLoginActivities int64 `json:"total_login_activities"`
}

// UserWithLastLogin represents a user with their last login time
type UserWithLastLogin struct {
	models.User
	LastLogin *string `json:"last_login,omitempty"`
}

// GetUserActivity returns user login activity with last login information
func GetUserActivity(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
		if page < 1 {
			page = 1
		}
		if limit < 1 || limit > 100 {
			limit = 50
		}
		offset := (page - 1) * limit

		// KPI data
		var stats UserActivityStats
		db.Model(&models.User{}).Count(&stats.TotalUsers)
		db.Model(&models.User{}).Where("is_active = ?", true).Count(&stats.TotalActiveUsers)
		db.Model(&models.UserLog{}).Where("activity LIKE ?", "%login%").Count(&stats.TotalLoginActivities)

		// Get active users with their roles
		var users []models.User
		var total int64

		db.Model(&models.User{}).Where("is_active = ?", true).Count(&total)
		if err := db.Preload("Role").Where("is_active = ?", true).Order("name").Offset(offset).Limit(limit).Find(&users).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Get last login for each user
		usersWithLogin := make([]gin.H, len(users))
		for i, user := range users {
			var lastLog models.UserLog
			var lastLogin *string
			if err := db.Where("user_id = ? AND activity LIKE ?", user.ID, "%login%").Order("created_at desc").First(&lastLog).Error; err == nil {
				loginTime := lastLog.CreatedAt.Format("2006-01-02 15:04:05")
				lastLogin = &loginTime
			}
			usersWithLogin[i] = gin.H{
				"id":         user.ID,
				"name":       user.Name,
				"username":   user.Username,
				"email":      user.Email,
				"role":       user.Role,
				"is_active":  user.IsActive,
				"last_login": lastLogin,
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"users": usersWithLogin,
			"stats": stats,
			"total": total,
			"page":  page,
			"limit": limit,
		})
	}
}

// CreateLog creates a new activity log entry
func CreateLog(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var log models.UserLog
		if err := c.ShouldBindJSON(&log); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := db.Create(&log).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"log": log})
	}
}
