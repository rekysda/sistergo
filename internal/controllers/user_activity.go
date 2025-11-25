package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rekysda/sistergo/internal/models"
	"gorm.io/gorm"
)

// UserActivityResponse for user activity data
type UserActivityResponse struct {
	Users                []UserWithLastLogin `json:"users"`
	TotalUsers           int64               `json:"total_users"`
	TotalActiveUsers     int64               `json:"total_active_users"`
	TotalLoginActivities int64               `json:"total_login_activities"`
	Page                 int                 `json:"page"`
	PerPage              int                 `json:"per_page"`
	TotalPages           int                 `json:"total_pages"`
}

// UserWithLastLogin represents a user with their last login time
type UserWithLastLogin struct {
	models.User
	LastLogin *string `json:"last_login"`
}

// GetUserActivity returns user login activity with last login information
func GetUserActivity(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		page := 1
		perPage := 50

		if p := c.Query("page"); p != "" {
			if val, err := strconv.Atoi(p); err == nil && val > 0 {
				page = val
			}
		}
		if pp := c.Query("per_page"); pp != "" {
			if val, err := strconv.Atoi(pp); err == nil && val > 0 && val <= 100 {
				perPage = val
			}
		}

		// KPI data
		var totalUsers, activeUsers, loginActivities int64
		db.Model(&models.User{}).Count(&totalUsers)
		db.Model(&models.User{}).Where("is_active = ?", true).Count(&activeUsers)
		db.Model(&models.UserLog{}).Where("activity LIKE ?", "%login%").Count(&loginActivities)

		// Get active users
		var users []models.User
		offset := (page - 1) * perPage
		db.Preload("Role").Where("is_active = ?", true).Order("name").Offset(offset).Limit(perPage).Find(&users)

		// Build response with last login info
		usersWithLogin := make([]UserWithLastLogin, len(users))
		for i, user := range users {
			usersWithLogin[i] = UserWithLastLogin{
				User: user,
			}

			// Get last login for this user
			var lastLog models.UserLog
			if err := db.Where("user_id = ? AND activity LIKE ?", user.ID, "%login%").
				Order("created_at DESC").
				First(&lastLog).Error; err == nil {
				lastLoginStr := lastLog.CreatedAt.Format("2006-01-02 15:04:05")
				usersWithLogin[i].LastLogin = &lastLoginStr
			}
		}

		totalPages := int(activeUsers) / perPage
		if int(activeUsers)%perPage > 0 {
			totalPages++
		}

		response := UserActivityResponse{
			Users:                usersWithLogin,
			TotalUsers:           totalUsers,
			TotalActiveUsers:     activeUsers,
			TotalLoginActivities: loginActivities,
			Page:                 page,
			PerPage:              perPage,
			TotalPages:           totalPages,
		}

		c.JSON(http.StatusOK, response)
	}
}
