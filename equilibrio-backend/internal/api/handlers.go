package api

import (
	"net/http"
	"strconv"
	"strings"

	"equilibrio-backend/internal/models"
	"equilibrio-backend/internal/services"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	marketDataService *services.MarketDataService
	indicatorService  *services.IndicatorService
	cacheService      *services.CacheService
}

func NewHandlers(
	marketDataService *services.MarketDataService,
	indicatorService *services.IndicatorService,
	cacheService *services.CacheService,
) *Handlers {
	return &Handlers{
		marketDataService: marketDataService,
		indicatorService:  indicatorService,
		cacheService:      cacheService,
	}
}

// GetStocks handles GET /api/stocks
func (h *Handlers) GetStocks(c *gin.Context) {
	var req models.StockListRequest

	// Parse query parameters
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Manually parse array parameters that Gin doesn't handle well
	// Split comma-separated values
	if sectorsParam := c.Query("sectors"); sectorsParam != "" {
		req.Sectors = strings.Split(sectorsParam, ",")
	}
	if volumeProfileParam := c.Query("volumeProfile"); volumeProfileParam != "" {
		req.VolumeProfile = strings.Split(volumeProfileParam, ",")
	}
	if signalsParam := c.Query("signals"); signalsParam != "" {
		req.Signals = strings.Split(signalsParam, ",")
	}
	if trendParam := c.Query("trend"); trendParam != "" {
		req.Trend = strings.Split(trendParam, ",")
	}
	if equilibriumZoneParam := c.Query("equilibriumZone"); equilibriumZoneParam != "" {
		req.EquilibriumZone = strings.Split(equilibriumZoneParam, ",")
	}

	// Set defaults
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 50
	}
	if req.SortField == "" {
		req.SortField = "symbol"
	}
	if req.SortOrder == "" {
		req.SortOrder = "asc"
	}

	// Set default filter values
	if req.RSIMin == 0 && req.RSIMax == 0 {
		req.RSIMin = 0
		req.RSIMax = 100
	}
	if req.PriceMin == 0 && req.PriceMax == 0 {
		req.PriceMin = 0
		req.PriceMax = 10000
	}

	// Get stocks from service
	stocks, total, err := h.marketDataService.GetStocks(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch stocks"})
		return
	}

	// Calculate total pages
	totalPages := (total + req.PageSize - 1) / req.PageSize

	response := models.StockListResponse{
		Stocks:     stocks,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}

	c.JSON(http.StatusOK, response)
}

// GetStock handles GET /api/stocks/:symbol
func (h *Handlers) GetStock(c *gin.Context) {
	symbol := c.Param("symbol")
	if symbol == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Symbol is required"})
		return
	}

	stock, err := h.marketDataService.GetStock(symbol)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Stock not found"})
		return
	}

	c.JSON(http.StatusOK, stock)
}

// GetStockChart handles GET /api/stocks/:symbol/chart
func (h *Handlers) GetStockChart(c *gin.Context) {
	symbol := c.Param("symbol")
	if symbol == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Symbol is required"})
		return
	}

	// Get optional days parameter (default to 90)
	daysStr := c.DefaultQuery("days", "90")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days < 1 || days > 365 {
		days = 90 // Default to 90 days if invalid
	}

	// Get chart data from market data service
	chartData, err := h.marketDataService.GetStockChartWithDays(symbol, days)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chart data not found"})
		return
	}

	c.JSON(http.StatusOK, chartData)
}

// GetSectors handles GET /api/sectors
func (h *Handlers) GetSectors(c *gin.Context) {
	sectors, err := h.marketDataService.GetSectors()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch sectors"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"sectors": sectors})
}

// CalculateIndicators handles POST /api/indicators
func (h *Handlers) CalculateIndicators(c *gin.Context) {
	var req struct {
		Symbol string `json:"symbol" binding:"required"`
		Period int    `json:"period"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Period <= 0 {
		req.Period = 200 // Default period
	}

	indicators, err := h.indicatorService.CalculateIndicators(req.Symbol, req.Period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate indicators"})
		return
	}

	c.JSON(http.StatusOK, indicators)
}

// RefreshData handles POST /api/refresh
func (h *Handlers) RefreshData(c *gin.Context) {
	err := h.marketDataService.RefreshAllData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to refresh data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data refreshed successfully"})
}

// ExportStocks handles GET /api/export
func (h *Handlers) ExportStocks(c *gin.Context) {
	var req models.StockListRequest

	// Parse query parameters
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set defaults for export
	req.Page = 1
	req.PageSize = 10000 // Large number for export

	stocks, _, err := h.marketDataService.GetStocks(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to export stocks"})
		return
	}

	// Generate CSV
	csvData := h.generateCSV(stocks)

	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename=stocks.csv")
	c.String(http.StatusOK, csvData)
}

// generateCSV creates CSV data from stocks
func (h *Handlers) generateCSV(stocks []models.StockData) string {
	csv := "Symbol,Name,Price,Change%,RSI,Trend,Signal,Equilibrium,Sector\n"

	for _, stock := range stocks {
		csv += stock.Symbol + "," +
			stock.Name + "," +
			strconv.FormatFloat(stock.Price, 'f', 2, 64) + "," +
			strconv.FormatFloat(stock.ChangePercent, 'f', 2, 64) + "," +
			strconv.FormatFloat(stock.RSI, 'f', 1, 64) + "," +
			stock.Trend + "," +
			stock.Signal + "," +
			strconv.FormatFloat(stock.PriceToEquilibrium, 'f', 1, 64) + "%," +
			stock.Sector + "\n"
	}

	return csv
}

// HealthCheck handles GET /health
func (h *Handlers) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "equilibrio-backend",
	})
}
