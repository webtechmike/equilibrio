package main

import (
	"log"
	"os"
	"time"

	"equilibrio-backend/internal/api"
	"equilibrio-backend/internal/config"
	"equilibrio-backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Load configuration
	cfg := config.Load()

	// Initialize services
	marketDataService := services.NewMarketDataService(cfg)
	indicatorService := services.NewIndicatorService()
	cacheService := services.NewCacheService(cfg)

	// Initialize API handlers
	handlers := api.NewHandlers(marketDataService, indicatorService, cacheService)

	// Setup router
	router := gin.Default()

	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// API routes
	api.SetupRoutes(router, handlers)

	// Lightweight daily cron to refresh snapshot at market close
	// In production you'd typically use a proper scheduler, but this demonstrates the concept.
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			// Run shortly after market close ~ 4:05 PM ET equivalent (we approximate via cacheStrategy)
			status := marketDataService.MarketStatus()
			if status == "Market Closed (After Hours)" {
				// Attempt to refresh once per after-hours period; the snapshot CacheDailySnapshot TTL prevents overwork
				if _, err := marketDataService.RefreshDailySnapshot(nil); err != nil {
					log.Printf("snapshot refresh skipped/failed: %v", err)
				} else {
					log.Printf("snapshot refresh completed")
				}
			}
		}
	}()

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
