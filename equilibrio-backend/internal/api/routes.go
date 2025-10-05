package api

import (
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, handlers *Handlers) {
	// Health check
	router.GET("/health", handlers.HealthCheck)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Stock routes
		v1.GET("/stocks", handlers.GetStocks)
		v1.GET("/stocks/:symbol", handlers.GetStock)
		v1.GET("/sectors", handlers.GetSectors)
		v1.GET("/export", handlers.ExportStocks)

		// Data management
		v1.POST("/refresh", handlers.RefreshData)

		// Technical indicators
		v1.POST("/indicators", handlers.CalculateIndicators)
	}

	// Legacy API routes for backward compatibility
	api := router.Group("/api")
	{
		api.GET("/stocks", handlers.GetStocks)
		api.GET("/stocks/:symbol", handlers.GetStock)
		api.GET("/stocks/:symbol/chart", handlers.GetStockChart)
		api.GET("/sectors", handlers.GetSectors)
		api.GET("/export", handlers.ExportStocks)
		api.POST("/refresh", handlers.RefreshData)
		api.POST("/indicators", handlers.CalculateIndicators)
	}
}
