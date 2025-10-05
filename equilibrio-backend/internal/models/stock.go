package models

import (
	"time"
)

// StockData represents the complete stock information
type StockData struct {
	Symbol                 string    `json:"symbol"`
	Name                   string    `json:"name"`
	Price                  float64   `json:"price"`
	Change                 float64   `json:"change"`
	ChangePercent          float64   `json:"changePercent"`
	Volume                 int64     `json:"volume"`
	Sector                 string    `json:"sector"`
	Industry               string    `json:"industry"`
	MarketCap              float64   `json:"marketCap"`
	RSI                    float64   `json:"rsi"`
	StochRSI               float64   `json:"stochRsi"`
	HistoricRSIAvg         float64   `json:"historicRsiAvg"`
	SMA50                  float64   `json:"sma50"`
	SMA200                 float64   `json:"sma200"`
	EMA20                  float64   `json:"ema20"`
	MACD                   float64   `json:"macd"`
	MACDSignal             float64   `json:"macdSignal"`
	MACDHistogram          float64   `json:"macdHistogram"`
	EquilibriumLevel       float64   `json:"equilibriumLevel"`
	PriceToEquilibrium     float64   `json:"priceToEquilibrium"`
	Trend                  string    `json:"trend"`         // "bullish", "bearish", "neutral"
	Signal                 string    `json:"signal"`        // "buy", "sell", "hold"
	VolumeProfile          string    `json:"volumeProfile"` // "high", "medium", "low"
	DistanceFrom52WeekHigh float64   `json:"distanceFrom52WeekHigh"`
	DistanceFrom52WeekLow  float64   `json:"distanceFrom52WeekLow"`
	LastUpdated            time.Time `json:"lastUpdated"`
}

// StockFilter represents filtering criteria
type StockFilter struct {
	SearchTerm      string   `json:"searchTerm"`
	Sectors         []string `json:"sectors"`
	RSIMin          float64  `json:"rsiMin"`
	RSIMax          float64  `json:"rsiMax"`
	PriceMin        float64  `json:"priceMin"`
	PriceMax        float64  `json:"priceMax"`
	VolumeProfile   []string `json:"volumeProfile"`
	Signals         []string `json:"signals"`
	Trend           []string `json:"trend"`
	EquilibriumZone []string `json:"equilibriumZone"`
}

// StockListRequest represents the request for stock data
type StockListRequest struct {
	Filter    StockFilter `json:"filter"`
	SortField string      `json:"sortField"`
	SortOrder string      `json:"sortOrder"` // "asc" or "desc"
	Page      int         `json:"page"`
	PageSize  int         `json:"pageSize"`
}

// StockListResponse represents the response for stock data
type StockListResponse struct {
	Stocks     []StockData `json:"stocks"`
	Total      int         `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"pageSize"`
	TotalPages int         `json:"totalPages"`
}

// PriceData represents historical price data for calculations
type PriceData struct {
	Date   time.Time `json:"date"`
	Open   float64   `json:"open"`
	High   float64   `json:"high"`
	Low    float64   `json:"low"`
	Close  float64   `json:"close"`
	Volume int64     `json:"volume"`
}

// TechnicalIndicators represents calculated technical indicators
type TechnicalIndicators struct {
	RSI            float64 `json:"rsi"`
	StochRSI       float64 `json:"stochRsi"`
	HistoricRSIAvg float64 `json:"historicRsiAvg"`
	SMA50          float64 `json:"sma50"`
	SMA200         float64 `json:"sma200"`
	EMA20          float64 `json:"ema20"`
	MACD           float64 `json:"macd"`
	MACDSignal     float64 `json:"macdSignal"`
	MACDHistogram  float64 `json:"macdHistogram"`
}
