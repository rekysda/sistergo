package controllers

import (
	"net/http"
	"os"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/rekysda/sistergo/internal/models"
	"gorm.io/gorm"
)

// DashboardResponse represents the dashboard data
type DashboardResponse struct {
	RecentLogs  []models.UserLog `json:"recent_logs"`
	RecentUsers []models.User    `json:"recent_users"`
	SystemInfo  SystemInfo       `json:"system_info"`
	TotalUsers  int64            `json:"total_users"`
	ActiveUsers int64            `json:"active_users"`
	TotalRoles  int64            `json:"total_roles"`
	TotalMenus  int64            `json:"total_menus"`
}

// SystemInfo represents system information
type SystemInfo struct {
	GoVersion       string `json:"go_version"`
	Environment     string `json:"environment"`
	DatabaseVersion string `json:"database_version"`
}

// GetDashboard returns dashboard data with recent logs, users, and system info
func GetDashboard(db *gorm.DB) gin.HandlerFunc {
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
		if err := db.First(&user, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		var recentLogs []models.UserLog
		// Admin (role_id = 1) sees all logs, others see only their own
		if user.RoleID == 1 {
			db.Preload("User").Order("created_at DESC").Limit(10).Find(&recentLogs)
		} else {
			db.Preload("User").Where("user_id = ?", user.ID).Order("created_at DESC").Limit(10).Find(&recentLogs)
		}

		// Get recent users
		var recentUsers []models.User
		db.Order("date_created DESC").Limit(8).Find(&recentUsers)

		// Get counts
		var totalUsers, activeUsers, totalRoles, totalMenus int64
		db.Model(&models.User{}).Count(&totalUsers)
		db.Model(&models.User{}).Where("is_active = ?", true).Count(&activeUsers)
		db.Model(&models.UserRole{}).Count(&totalRoles)
		db.Model(&models.UserMenu{}).Count(&totalMenus)

		// Get database version
		dbVersion := getDatabaseVersion(db)

		// Get environment from ENV or default to production
		environment := os.Getenv("APP_ENV")
		if environment == "" {
			environment = "production"
		}

		response := DashboardResponse{
			RecentLogs:  recentLogs,
			RecentUsers: recentUsers,
			SystemInfo: SystemInfo{
				GoVersion:       runtime.Version(),
				Environment:     environment,
				DatabaseVersion: dbVersion,
			},
			TotalUsers:  totalUsers,
			ActiveUsers: activeUsers,
			TotalRoles:  totalRoles,
			TotalMenus:  totalMenus,
		}

		c.JSON(http.StatusOK, gin.H{"dashboard": response})
	}
}

func getDatabaseVersion(db *gorm.DB) string {
	var version string
	row := db.Raw("SELECT VERSION()").Row()
	if err := row.Scan(&version); err != nil {
		return "Unknown"
	}
	return version
}
