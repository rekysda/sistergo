package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rekysda/sistergo/internal/config"
	"github.com/rekysda/sistergo/internal/models"
	"github.com/rekysda/sistergo/internal/utils"
	"gorm.io/gorm"
)

// Register request payload
type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"omitempty"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// Login request payload
type LoginRequest struct {
	Identifier    string `json:"identifier" binding:"required"` // username or email
	Password string `json:"password" binding:"required"`
}

func Register(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		hash, err := utils.HashPassword(req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to hash password"})
			return
		}

		username := req.Username
		if username == "" {
			username = req.Email
		}
		user := models.User{Name: req.Name, Email: req.Email, Username: username, Password: hash}
		if err := db.Create(&user).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Set date_created if not present
		now := time.Now()
		if user.DateCreated == nil {
			user.DateCreated = &now
			db.Save(&user)
		}
		// Log user creation
		db.Create(&models.UserLog{UserID: user.ID, Activity: "register"})
		c.JSON(http.StatusCreated, gin.H{"user": user})
	}
}

func Login(db *gorm.DB, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var user models.User
		// Try find by username first, then by email
		if err := db.Where("username = ?", req.Identifier).First(&user).Error; err != nil {
			if err := db.Where("email = ?", req.Identifier).First(&user).Error; err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
				return
			}
		}
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		if !utils.CheckPassword(user.Password, req.Password) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		// generate token (72 hours)
		token, err := utils.GenerateJWT(cfg.JWTSecret, user.ID, 72*time.Hour)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create token"})
			return
		}

		// Add login activity into user log
		db.Create(&models.UserLog{UserID: user.ID, Activity: "login"})
		c.JSON(http.StatusOK, gin.H{"token": token, "user": user})
	}
}

func Me(db *gorm.DB, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDVal, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		var user models.User
		id, ok := userIDVal.(uint)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		if err := db.First(&user, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"user": user})
	}
}
