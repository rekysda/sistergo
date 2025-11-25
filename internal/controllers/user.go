package controllers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rekysda/sistergo/internal/config"
	"github.com/rekysda/sistergo/internal/models"
	"github.com/rekysda/sistergo/internal/utils"
	"gorm.io/gorm"
)

// CreateUserRequest represents the request payload for creating a user
type CreateUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	RoleID   uint   `json:"role_id" binding:"required"`
	IsActive bool   `json:"is_active"`
}

// UpdateUserRequest represents the request payload for updating a user
type UpdateUserRequest struct {
	Name     string  `json:"name" binding:"required"`
	Username string  `json:"username" binding:"required"`
	Email    string  `json:"email" binding:"required,email"`
	Password *string `json:"password,omitempty"`
	RoleID   uint    `json:"role_id" binding:"required"`
	IsActive *bool   `json:"is_active,omitempty"`
}

// ListUsers returns a list of all users with pagination
func ListUsers(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
		if page < 1 {
			page = 1
		}
		if limit < 1 || limit > 100 {
			limit = 10
		}
		offset := (page - 1) * limit

		var users []models.User
		var total int64

		db.Model(&models.User{}).Count(&total)
		if err := db.Preload("Role").Preload("Logs").Offset(offset).Limit(limit).Find(&users).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// KPI data
		var totalActive int64
		var totalRoles int64
		db.Model(&models.User{}).Where("is_active = ?", true).Count(&totalActive)
		db.Model(&models.UserRole{}).Count(&totalRoles)

		c.JSON(http.StatusOK, gin.H{
			"users":              users,
			"total":              total,
			"page":               page,
			"limit":              limit,
			"total_users":        total,
			"total_active_users": totalActive,
			"total_roles":        totalRoles,
		})
	}
}

// GetUser returns a single user by ID
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

// CreateUser creates a new user
func CreateUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Check if username already exists
		var existingUser models.User
		if err := db.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "username already exists"})
			return
		}

		// Check if email already exists
		if err := db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "email already exists"})
			return
		}

		// Hash password
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
			return
		}

		now := time.Now()
		user := models.User{
			Name:        req.Name,
			Username:    req.Username,
			Email:       req.Email,
			Password:    hashedPassword,
			RoleID:      req.RoleID,
			IsActive:    req.IsActive,
			DateCreated: &now,
		}

		if err := db.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Log user creation
		db.Create(&models.UserLog{UserID: user.ID, Activity: "user_created"})

		c.JSON(http.StatusCreated, gin.H{"user": user, "message": "User created successfully"})
	}
}

// UpdateUser updates an existing user
func UpdateUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		if err := db.First(&user, c.Param("id")).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		var req UpdateUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Check if username is taken by another user
		var existingUser models.User
		if err := db.Where("username = ? AND id != ?", req.Username, user.ID).First(&existingUser).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "username already exists"})
			return
		}

		// Check if email is taken by another user
		if err := db.Where("email = ? AND id != ?", req.Email, user.ID).First(&existingUser).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "email already exists"})
			return
		}

		user.Name = req.Name
		user.Username = req.Username
		user.Email = req.Email
		user.RoleID = req.RoleID

		if req.IsActive != nil {
			user.IsActive = *req.IsActive
		}

		// Update password if provided
		if req.Password != nil && *req.Password != "" {
			hashedPassword, err := utils.HashPassword(*req.Password)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
				return
			}
			user.Password = hashedPassword
		}

		if err := db.Save(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"user": user, "message": "User updated successfully"})
	}
}

// DeleteUser deletes a user
func DeleteUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		if err := db.First(&user, c.Param("id")).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		// Prevent deletion of admin users
		if user.RoleID == config.AdminRoleID || strings.ToLower(user.Username) == config.DefaultAdminUsername {
			c.JSON(http.StatusForbidden, gin.H{"error": "Cannot delete admin user"})
			return
		}

		if err := db.Delete(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"deleted": true, "message": "User deleted successfully"})
	}
}
