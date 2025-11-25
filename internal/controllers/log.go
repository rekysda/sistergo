package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rekysda/sistergo/internal/models"
	"gorm.io/gorm"
)

// LogListResponse for paginated log list
type LogListResponse struct {
	Logs       []models.UserLog `json:"logs"`
	TotalLogs  int64            `json:"total_logs"`
	Page       int              `json:"page"`
	PerPage    int              `json:"per_page"`
	TotalPages int              `json:"total_pages"`
}

// ListLogs returns paginated list of user activity logs
func ListLogs(db *gorm.DB) gin.HandlerFunc {
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

		var logs []models.UserLog
		var totalLogs int64

		db.Model(&models.UserLog{}).Count(&totalLogs)

		offset := (page - 1) * perPage
		db.Preload("User").Order("created_at DESC").Offset(offset).Limit(perPage).Find(&logs)

		totalPages := int(totalLogs) / perPage
		if int(totalLogs)%perPage > 0 {
			totalPages++
		}

		response := LogListResponse{
			Logs:       logs,
			TotalLogs:  totalLogs,
			Page:       page,
			PerPage:    perPage,
			TotalPages: totalPages,
		}

		c.JSON(http.StatusOK, response)
	}
}

// GetLogsByUser returns logs for a specific user
func GetLogsByUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("user_id")

		var logs []models.UserLog
		if err := db.Preload("User").Where("user_id = ?", userID).Order("created_at DESC").Find(&logs).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"logs": logs})
	}
}

// CreateLog creates a new log entry
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
