package controllers

import (
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/rekysda/sistergo/internal/config"
	"github.com/rekysda/sistergo/internal/models"
	"gorm.io/gorm"
)

// DashboardStats represents dashboard statistics
type DashboardStats struct {
	TotalUsers       int64 `json:"total_users"`
	TotalActiveUsers int64 `json:"total_active_users"`
	TotalRoles       int64 `json:"total_roles"`
	TotalMenus       int64 `json:"total_menus"`
	TotalSubmenus    int64 `json:"total_submenus"`
}

// SystemInfo represents system information
type SystemInfo struct {
	GoVersion    string `json:"go_version"`
	Environment  string `json:"environment"`
	NumCPU       int    `json:"num_cpu"`
	NumGoroutine int    `json:"num_goroutine"`
}

// GetDashboard returns dashboard statistics
func GetDashboard(db *gorm.DB, cfg *config.Config) gin.HandlerFunc {
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

		// Gather statistics
		var stats DashboardStats
		db.Model(&models.User{}).Count(&stats.TotalUsers)
		db.Model(&models.User{}).Where("is_active = ?", true).Count(&stats.TotalActiveUsers)
		db.Model(&models.UserRole{}).Count(&stats.TotalRoles)
		db.Model(&models.UserMenu{}).Count(&stats.TotalMenus)
		db.Model(&models.UserSubMenu{}).Count(&stats.TotalSubmenus)

		// Recent logs - filter based on user role
		var recentLogs []models.UserLog
		if user.RoleID == config.AdminRoleID {
			// Admin sees all logs
			db.Preload("User").Order("created_at desc").Limit(10).Find(&recentLogs)
		} else {
			// Regular users see only their own logs
			db.Preload("User").Where("user_id = ?", user.ID).Order("created_at desc").Limit(10).Find(&recentLogs)
		}

		// Recent users (for admin)
		var recentUsers []models.User
		if user.RoleID == config.AdminRoleID {
			db.Preload("Role").Order("created_at desc").Limit(8).Find(&recentUsers)
		}

		// System info
		sysInfo := SystemInfo{
			GoVersion:    runtime.Version(),
			Environment:  cfg.AppEnv,
			NumCPU:       runtime.NumCPU(),
			NumGoroutine: runtime.NumGoroutine(),
		}

		c.JSON(http.StatusOK, gin.H{
			"stats":        stats,
			"recent_logs":  recentLogs,
			"recent_users": recentUsers,
			"user":         user,
			"system_info":  sysInfo,
		})
	}
}
