package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rekysda/sistergo/internal/models"
	"gorm.io/gorm"
)

// ListLogs returns paginated user activity logs
func ListLogs(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
		offset := (page - 1) * limit

		var logs []models.UserLog
		var total int64

		db.Model(&models.UserLog{}).Count(&total)
		if err := db.Preload("User").Order("created_at DESC").Offset(offset).Limit(limit).Find(&logs).Error; err != nil {
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

// LastLoginResult is used to store last login query results
type LastLoginResult struct {
	UserID    uint
	LastLogin time.Time
}

// UserActivity returns users with their last login activity
func UserActivity(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
		offset := (page - 1) * limit

		// Get KPI data
		var totalUsers, totalActiveUsers, totalLoginActivities int64
		db.Model(&models.User{}).Count(&totalUsers)
		db.Model(&models.User{}).Where("is_active = ?", true).Count(&totalActiveUsers)
		db.Model(&models.UserLog{}).Where("activity LIKE ?", "%login%").Count(&totalLoginActivities)

		var users []models.User
		var total int64
		db.Model(&models.User{}).Where("is_active = ?", true).Count(&total)

		if err := db.Preload("Role").Where("is_active = ?", true).
			Order("name ASC").Offset(offset).Limit(limit).Find(&users).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Get user IDs for batch query
		userIDs := make([]uint, len(users))
		for i, user := range users {
			userIDs[i] = user.ID
		}

		// Batch query for last login times (avoid N+1 query problem)
		var lastLogins []LastLoginResult
		if len(userIDs) > 0 {
			db.Model(&models.UserLog{}).
				Select("user_id, MAX(created_at) as last_login").
				Where("user_id IN ? AND activity LIKE ?", userIDs, "%login%").
				Group("user_id").
				Scan(&lastLogins)
		}

		// Create a map for quick lookup
		lastLoginMap := make(map[uint]string)
		for _, ll := range lastLogins {
			lastLoginMap[ll.UserID] = ll.LastLogin.Format("2006-01-02 15:04:05")
		}

		// Build response with last login info
		var result []map[string]interface{}
		for _, user := range users {
			lastLogin := lastLoginMap[user.ID]
			result = append(result, map[string]interface{}{
				"id":         user.ID,
				"name":       user.Name,
				"username":   user.Username,
				"email":      user.Email,
				"role":       user.Role,
				"is_active":  user.IsActive,
				"last_login": lastLogin,
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"users":                  result,
			"total":                  total,
			"page":                   page,
			"limit":                  limit,
			"total_users":            totalUsers,
			"total_active_users":     totalActiveUsers,
			"total_login_activities": totalLoginActivities,
		})
	}
}
