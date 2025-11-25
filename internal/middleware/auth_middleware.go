package middleware

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rekysda/sistergo/internal/config"
	"github.com/rekysda/sistergo/internal/utils"
)

func JWTAuth(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing Authorization header"})
			c.Abort()
			return
		}

		// Expect header: "Bearer <token>"
		const prefix = "Bearer "
		if len(auth) <= len(prefix) || auth[:len(prefix)] != prefix {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid Authorization header format"})
			c.Abort()
			return
		}
		tokenStr := auth[len(prefix):]

		claims, err := utils.ParseJWT(cfg.JWTSecret, tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		// user_id may be float64 when parsed from json
		var userID uint
		if v, ok := claims["user_id"]; ok {
			switch val := v.(type) {
			case float64:
				userID = uint(val)
			case string:
				p, _ := strconv.ParseUint(val, 10, 32)
				userID = uint(p)
			default:
				userID = 0
			}
		}
		if userID == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token payload"})
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}
