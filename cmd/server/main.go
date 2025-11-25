package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rekysda/sistergo/internal/config"
	"github.com/rekysda/sistergo/internal/database"
	"github.com/rekysda/sistergo/internal/middleware"
	"github.com/rekysda/sistergo/internal/routes"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := database.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Run auto migrations
	err = database.Migrate(db)
	if err != nil {
		log.Fatalf("failed to migrate db: %v", err)
	}

	// Seed initial data
	err = database.Seed(db)
	if err != nil {
		log.Fatalf("failed to seed db: %v", err)
	}

	r := gin.Default()

	// Enable CORS
	r.Use(middleware.CORS(cfg))

	routes.RegisterRoutes(r, db, cfg)

	addr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("Starting server on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
	_ = os.Stdout
}
