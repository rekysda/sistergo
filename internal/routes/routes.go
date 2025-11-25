package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rekysda/sistergo/internal/config"
	"github.com/rekysda/sistergo/internal/controllers"
	"github.com/rekysda/sistergo/internal/middleware"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB, cfg *config.Config) {
	api := r.Group("/api")
	{
		api.POST("/register", controllers.Register(db))
		api.POST("/login", controllers.Login(db, cfg))
		api.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })

		// Protected routes
		protected := api.Group("")
		protected.Use(middleware.JWTAuth(cfg))
		{
			protected.GET("/me", controllers.Me(db, cfg))
			protected.GET("/profile", controllers.GetProfile(db))
			protected.PUT("/profile", controllers.UpdateProfile(db))
			protected.POST("/profile/change-password", controllers.ChangePassword(db))

			students := protected.Group("/students")
			{
				students.GET("", controllers.ListStudents(db))
				students.POST("", controllers.CreateStudent(db))
				students.GET(":id", controllers.GetStudent(db))
				students.PUT(":id", controllers.UpdateStudent(db))
				students.DELETE(":id", controllers.DeleteStudent(db))
			}

			// Admin-like management endpoints
				admin := protected.Group("/admin")
				admin.Use(middleware.AdminOnly(cfg, db))
			{
				admin.GET("/roles", controllers.ListRoles(db))
				admin.POST("/roles", controllers.CreateRole(db))
				admin.GET("/roles/:id", controllers.GetRole(db))
				admin.PUT("/roles/:id", controllers.UpdateRole(db))
				admin.DELETE("/roles/:id", controllers.DeleteRole(db))

				admin.GET("/menus", controllers.ListMenus(db))
				admin.POST("/menus", controllers.CreateMenu(db))
				admin.PUT("/menus/:id", controllers.UpdateMenu(db))
				admin.DELETE("/menus/:id", controllers.DeleteMenu(db))

				admin.POST("/submenus", controllers.CreateSubmenu(db))
				admin.PUT("/submenus/:id", controllers.UpdateSubmenu(db))
				admin.DELETE("/submenus/:id", controllers.DeleteSubmenu(db))

				admin.GET("/settings", controllers.GetSettings(db))
				admin.POST("/settings", controllers.UpdateSetting(db))
			}
		}
	}
}
