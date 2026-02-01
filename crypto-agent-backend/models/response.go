package models

import "time"

// APIResponse is the standard API response wrapper
type APIResponse struct {
	Status      string      `json:"status"`
	Timestamp   time.Time   `json:"timestamp"`
	Total       int         `json:"total,omitempty"`
	Data        interface{} `json:"data,omitempty"`
	Message     string      `json:"message,omitempty"`
	FetchTimeMs int64       `json:"fetch_time_ms,omitempty"`
}

// TokensResponse is the response for /api/tokens endpoint
type TokensResponse struct {
	Status      string    `json:"status"`
	Timestamp   time.Time `json:"timestamp"`
	Total       int       `json:"total"`
	Data        []Token   `json:"data"`
	FetchTimeMs int64     `json:"fetch_time_ms"`
}

// AnalysisRequest is the request body for /api/analyze endpoint
type AnalysisRequest struct {
	Symbol            string  `json:"symbol" binding:"required"`
	Name              string  `json:"name" binding:"required"`
	Price             float64 `json:"price"`
	MarketCap         float64 `json:"market_cap"`
	Volume24h         float64 `json:"volume_24h"`
	TVL               float64 `json:"tvl"`
	TrustScore        float64 `json:"trust_score"`
	Change24h         float64 `json:"change_24h"`
	Change7d          float64 `json:"change_7d"`
	Liquidity         float64 `json:"liquidity"`
	Rank              int     `json:"rank"`
	HolderCount       int     `json:"holder_count"`
	CirculatingSupply float64 `json:"circulating_supply"`
	MaxSupply         float64 `json:"max_supply"`
	TotalSupply       float64 `json:"total_supply"`
	Change30d         float64 `json:"change_30d"`
	Change90d         float64 `json:"change_90d"`
}

// AnalysisResponse is the response for /api/analyze endpoint
type AnalysisResponse struct {
	Status      string    `json:"status"`
	Cached      bool      `json:"cached"`
	Analysis    string    `json:"analysis"`
	GeneratedAt time.Time `json:"generated_at"`
}

// ErrorResponse is the standard error response
type ErrorResponse struct {
	Status    string    `json:"status"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}
