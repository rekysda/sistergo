package controllers

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rekysda/sistergo/internal/models"
	"github.com/rekysda/sistergo/internal/utils"
	"gorm.io/gorm"
)

func GetProfile(db *gorm.DB) gin.HandlerFunc {
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
		c.JSON(http.StatusOK, gin.H{"user": user})
	}
}

// Update profile: name, email, image (image not implemented for now)
func UpdateProfile(db *gorm.DB) gin.HandlerFunc {
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
		var payload struct {
			Name  string  `json:"name" binding:"required"`
			Email string  `json:"email" binding:"required,email"`
			Image *string `json:"image"`
		}
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// unique email check
		var existing models.User
		if err := db.Where("email = ? AND id != ?", payload.Email, user.ID).First(&existing).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "email already used"})
			return
		}
		user.Name = payload.Name
		user.Email = payload.Email
		user.Image = payload.Image
		user.UpdatedAt = time.Now()
		if err := db.Save(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"user": user})
	}
}

// ChangePassword expects: current_password, new_password, new_password_confirmation
func ChangePassword(db *gorm.DB) gin.HandlerFunc {
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
		var payload struct {
			CurrentPassword string `json:"current_password" binding:"required"`
			NewPassword     string `json:"new_password" binding:"required,min=8"`
			ConfirmPassword string `json:"confirm_password" binding:"required"`
		}
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if strings.Compare(payload.NewPassword, payload.ConfirmPassword) != 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "passwords do not match"})
			return
		}
		if !utils.CheckPassword(user.Password, payload.CurrentPassword) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "current password invalid"})
			return
		}
		hashed, err := utils.HashPassword(payload.NewPassword)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
			return
		}
		user.Password = hashed
		if err := db.Save(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "password changed"})
	}
}
