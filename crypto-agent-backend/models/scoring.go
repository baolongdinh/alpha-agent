package models

// Enhanced scoring breakdown with 6 categories
type DetailedScoreBreakdown struct {
	// Main Categories (total 100%)
	LiquidityScore    float64 `json:"liquidity_score"`     // 25%
	VolumeScore       float64 `json:"volume_score"`        // 20%
	TVLScore          float64 `json:"tvl_score"`           // 20%
	TrendScore        float64 `json:"trend_score"`         // 20%
	MarketHealthScore float64 `json:"market_health_score"` // 10%
	SocialScore       float64 `json:"social_score"`        // NEW: 10%
	RiskScore         float64 `json:"risk_score"`          // 5%

	// Overall
	TotalScore float64 `json:"total_score"`
	Grade      string  `json:"grade"`      // S, A, B, C, D, F
	Confidence float64 `json:"confidence"` // 0-100

	// Detailed metrics
	Details ScoreDetails `json:"details"`
}

// ScoreDetails contains granular metrics for each category
type ScoreDetails struct {
	// Liquidity Details
	LiquidityRatio float64 `json:"liquidity_ratio"`
	LiquidityDepth float64 `json:"liquidity_depth"`

	// Volume Details
	VolumeToMcapRatio float64 `json:"volume_to_mcap_ratio"`
	VolumeConsistency float64 `json:"volume_consistency"`
	Volume7dGrowth    float64 `json:"volume_7d_growth"`

	// TVL Details
	TVLToMcapRatio float64 `json:"tvl_to_mcap_ratio"`
	TVL7dGrowth    float64 `json:"tvl_7d_growth"`
	TVL30dGrowth   float64 `json:"tvl_30d_growth"`
	TVLVolatility  float64 `json:"tvl_volatility"`

	// Price & Trend Details
	ShortTermTrend  string  `json:"short_term_trend"`  // 7 days
	MediumTermTrend string  `json:"medium_term_trend"` // 30 days
	LongTermTrend   string  `json:"long_term_trend"`   // 90 days
	Volatility30d   float64 `json:"volatility_30d"`
	SharpeRatio     float64 `json:"sharpe_ratio"`
	RSI             float64 `json:"rsi"`

	// Market Health Details
	Top10HoldersPercent float64 `json:"top10_holders_percent"`
	UniqueHolders       int     `json:"unique_holders"`
	MarketCapRank       int     `json:"market_cap_rank"`

	// Risk Indicators
	RugPullRisk        string `json:"rug_pull_risk"` // Low, Medium, High
	CentralizationRisk string `json:"centralization_risk"`
	SmartContractRisk  string `json:"smart_contract_risk"`
}

// Historical data structures
type PriceHistory struct {
	Last7Days  []float64 `json:"last_7_days"`
	Last30Days []float64 `json:"last_30_days"`
	Last90Days []float64 `json:"last_90_days"`
}

type VolumeHistory struct {
	Last7Days  []float64 `json:"last_7_days"`
	Last30Days []float64 `json:"last_30_days"`
}

type TVLHistory struct {
	Last7Days  []float64 `json:"last_7_days"`
	Last30Days []float64 `json:"last_30_days"`
}
