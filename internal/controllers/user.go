package controllers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rekysda/sistergo/internal/models"
	"github.com/rekysda/sistergo/internal/utils"
	"gorm.io/gorm"
)

// AdminRoleID is the role ID for administrators
const AdminRoleID = 1

// AdminUsername is the default admin username that cannot be deleted
const AdminUsername = "admin"

// ListUsers returns paginated list of users with role and logs
func ListUsers(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
		offset := (page - 1) * limit

		var users []models.User
		var total int64

		db.Model(&models.User{}).Count(&total)
		if err := db.Preload("Role").Preload("Logs").Offset(offset).Limit(limit).Find(&users).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// KPI data
		var totalUsers, totalActiveUsers, totalRoles int64
		db.Model(&models.User{}).Count(&totalUsers)
		db.Model(&models.User{}).Where("is_active = ?", true).Count(&totalActiveUsers)
		db.Model(&models.UserRole{}).Count(&totalRoles)

		c.JSON(http.StatusOK, gin.H{
			"users":              users,
			"total":              total,
			"page":               page,
			"limit":              limit,
			"total_users":        totalUsers,
			"total_active_users": totalActiveUsers,
			"total_roles":        totalRoles,
		})
	}
}

// CreateUser creates a new user (admin only)
func CreateUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload struct {
			Name     string `json:"name" binding:"required"`
			Username string `json:"username" binding:"required"`
			Email    string `json:"email" binding:"required,email"`
			Password string `json:"password" binding:"required,min=8"`
			RoleID   uint   `json:"role_id" binding:"required"`
		}
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Check username uniqueness
		var existing models.User
		if err := db.Where("username = ?", payload.Username).First(&existing).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "username already exists"})
			return
		}

		// Check email uniqueness
		if err := db.Where("email = ?", payload.Email).First(&existing).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "email already exists"})
			return
		}

		// Check if role exists
		var role models.UserRole
		if err := db.First(&role, payload.RoleID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "role does not exist"})
			return
		}

		hash, err := utils.HashPassword(payload.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
			return
		}

		now := time.Now()
		defaultImage := "default.jpg"
		user := models.User{
			Name:        payload.Name,
			Username:    payload.Username,
			Email:       payload.Email,
			Password:    hash,
			RoleID:      payload.RoleID,
			IsActive:    true,
			Image:       &defaultImage,
			DateCreated: &now,
		}

		if err := db.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"user": user})
	}
}

// GetUser returns a user by ID
func GetUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		if err := db.Preload("Role").Preload("Logs").First(&user, c.Param("id")).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"user": user})
	}
}

// UpdateUser updates a user
func UpdateUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		if err := db.First(&user, c.Param("id")).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		var payload struct {
			Name     string  `json:"name" binding:"required"`
			Username string  `json:"username" binding:"required"`
			Email    string  `json:"email" binding:"required,email"`
			RoleID   uint    `json:"role_id" binding:"required"`
			Password *string `json:"password,omitempty"`
		}
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Check username uniqueness (exclude current user)
		var existing models.User
		if err := db.Where("username = ? AND id != ?", payload.Username, user.ID).First(&existing).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "username already exists"})
			return
		}

		// Check email uniqueness (exclude current user)
		if err := db.Where("email = ? AND id != ?", payload.Email, user.ID).First(&existing).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "email already exists"})
			return
		}

		// Check if role exists
		var role models.UserRole
		if err := db.First(&role, payload.RoleID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "role does not exist"})
			return
		}

		user.Name = payload.Name
		user.Username = payload.Username
		user.Email = payload.Email
		user.RoleID = payload.RoleID

		// Update password if provided (with validation)
		if payload.Password != nil && *payload.Password != "" {
			if len(*payload.Password) < 8 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "password must be at least 8 characters"})
				return
			}
			hash, err := utils.HashPassword(*payload.Password)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
				return
			}
			user.Password = hash
		}

		if err := db.Save(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"user": user})
	}
}

// DeleteUser deletes a user (cannot delete admin)
func DeleteUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		if err := db.First(&user, c.Param("id")).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		// Prevent deletion of admin users
		if user.RoleID == AdminRoleID || strings.ToLower(user.Username) == AdminUsername {
			c.JSON(http.StatusForbidden, gin.H{"error": "cannot delete admin user"})
			return
		}

		if err := db.Delete(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"deleted": true})
	}
}
