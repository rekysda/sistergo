package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rekysda/sistergo/internal/config"
	"github.com/rekysda/sistergo/internal/models"
	"gorm.io/gorm"
)

// AdminOnly verifies that the authenticated user has role_id == 1 (super admin)
func AdminOnly(cfg *config.Config, db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDVal, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}
		id, ok := userIDVal.(uint)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}
		var user models.User
		if err := db.First(&user, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			c.Abort()
			return
		}
		// Role 1 is 'admin' in Laravel app convention
		if user.RoleID != 1 {
			c.JSON(http.StatusForbidden, gin.H{"error": "admin only"})
			c.Abort()
			return
		}
		c.Next()
	}
}
