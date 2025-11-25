package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rekysda/sistergo/internal/models"
	"github.com/rekysda/sistergo/internal/utils"
	"gorm.io/gorm"
)

// UserListResponse for paginated user list
type UserListResponse struct {
	Users       []models.User `json:"users"`
	TotalUsers  int64         `json:"total_users"`
	ActiveUsers int64         `json:"active_users"`
	TotalRoles  int64         `json:"total_roles"`
	Page        int           `json:"page"`
	PerPage     int           `json:"per_page"`
	TotalPages  int           `json:"total_pages"`
}

// CreateUserRequest payload for creating a user
type CreateUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	RoleID   uint   `json:"role_id" binding:"required"`
	IsActive *bool  `json:"is_active"`
}

// UpdateUserRequest payload for updating a user
type UpdateUserRequest struct {
	Name     string  `json:"name" binding:"required"`
	Username string  `json:"username" binding:"required"`
	Email    string  `json:"email" binding:"required,email"`
	Password *string `json:"password"` // optional
	RoleID   uint    `json:"role_id" binding:"required"`
	IsActive *bool   `json:"is_active"`
}

// ListUsers returns paginated list of users with stats
func ListUsers(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		page := 1
		perPage := 10

		if p := c.Query("page"); p != "" {
			if _, err := c.GetQuery("page"); err {
				page = 1
			}
		}
		if pp := c.Query("per_page"); pp != "" {
			if _, err := c.GetQuery("per_page"); err {
				perPage = 10
			}
		}

		var users []models.User
		var totalUsers, activeUsers, totalRoles int64

		db.Model(&models.User{}).Count(&totalUsers)
		db.Model(&models.User{}).Where("is_active = ?", true).Count(&activeUsers)
		db.Model(&models.UserRole{}).Count(&totalRoles)

		offset := (page - 1) * perPage
		db.Preload("Role").Preload("Logs").Offset(offset).Limit(perPage).Find(&users)

		totalPages := int(totalUsers) / perPage
		if int(totalUsers)%perPage > 0 {
			totalPages++
		}

		response := UserListResponse{
			Users:       users,
			TotalUsers:  totalUsers,
			ActiveUsers: activeUsers,
			TotalRoles:  totalRoles,
			Page:        page,
			PerPage:     perPage,
			TotalPages:  totalPages,
		}

		c.JSON(http.StatusOK, response)
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

// CreateUser creates a new user (admin only)
func CreateUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Check if username exists
		var existingUser models.User
		if err := db.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "username already exists"})
			return
		}

		// Check if email exists
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

		isActive := true
		if req.IsActive != nil {
			isActive = *req.IsActive
		}

		now := time.Now()
		user := models.User{
			Name:        req.Name,
			Username:    req.Username,
			Email:       req.Email,
			Password:    hashedPassword,
			RoleID:      req.RoleID,
			IsActive:    isActive,
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

// DeleteUser deletes a user (prevents admin deletion)
func DeleteUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		if err := db.First(&user, c.Param("id")).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		// Prevent deletion of admin users
		if user.RoleID == 1 || user.Username == "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "cannot delete admin user"})
			return
		}

		if err := db.Delete(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"deleted": true, "message": "User deleted successfully"})
	}
}
