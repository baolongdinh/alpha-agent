package services

import (
	"crypto-agent-backend/models"
	"crypto-agent-backend/utils"
	"math"
)

// EnhancedScorer implements comprehensive 6-category scoring system
type EnhancedScorer struct {
	// Weights for each category (total = 100%)
	liquidityWeight    float64 // 25%
	volumeWeight       float64 // 20%
	tvlWeight          float64 // 20%
	trendWeight        float64 // 20%
	marketHealthWeight float64 // 10%
	riskWeight         float64 // 5%
}

// NewEnhancedScorer creates a new comprehensive scorer
func NewEnhancedScorer() *EnhancedScorer {
	return &EnhancedScorer{
		liquidityWeight:    25.0,
		volumeWeight:       20.0,
		tvlWeight:          20.0,
		trendWeight:        20.0,
		marketHealthWeight: 10.0,
		riskWeight:         5.0,
	}
}

// CalculateComprehensiveScore computes all scoring components with dynamic weighting
func (s *EnhancedScorer) CalculateComprehensiveScore(token *models.Token) models.DetailedScoreBreakdown {
	breakdown := models.DetailedScoreBreakdown{
		Details: models.ScoreDetails{},
	}

	// 1. Calculate raw scores (0-100 scale for each)
	lScore := s.calculateLiquidityScore(token, &breakdown.Details)          // raw 0-100
	vScore := s.calculateVolumeScore(token, &breakdown.Details)             // raw 0-100
	tvlScore := s.calculateTVLScore(token, &breakdown.Details)              // raw 0-100
	trendScore := s.calculateTrendScore(token, &breakdown.Details)          // raw 0-100
	mHealthScore := s.calculateMarketHealthScore(token, &breakdown.Details) // raw 0-100
	riskScore := s.calculateRiskScore(token, &breakdown.Details)            // raw 0-100

	// 2. Dynamic Weighting Logic
	// If a token lacks TVL or Liquidity data, we redistribute those weights
	weights := map[string]float64{
		"liquidity": s.liquidityWeight,
		"volume":    s.volumeWeight,
		"tvl":       s.tvlWeight,
		"trend":     s.trendWeight,
		"market":    s.marketHealthWeight,
		"risk":      s.riskWeight,
	}

	// Identify missing core data
	if token.TVL == 0 {
		w := weights["tvl"]
		weights["tvl"] = 0
		// Redistribute TVL weight to Volume, Trend, and Market Health
		weights["volume"] += w * 0.4
		weights["trend"] += w * 0.4
		weights["market"] += w * 0.2
	}

	if token.Liquidity == 0 {
		w := weights["liquidity"]
		weights["liquidity"] = 0
		// Redistribute Liquidity weight to Volume and Risk (as safety check)
		weights["volume"] += w * 0.6
		weights["risk"] += w * 0.4
	}

	// 3. Final Weighted Sum & Category Assignment (Category scores stored as raw 0-100 for radar compatibility)
	breakdown.LiquidityScore = lScore
	breakdown.VolumeScore = vScore
	breakdown.TVLScore = tvlScore
	breakdown.TrendScore = trendScore
	breakdown.MarketHealthScore = mHealthScore
	breakdown.RiskScore = riskScore

	// Calculate weighted total
	weightedTotal := (lScore/100.0)*weights["liquidity"] +
		(vScore/100.0)*weights["volume"] +
		(tvlScore/100.0)*weights["tvl"] +
		(trendScore/100.0)*weights["trend"] +
		(mHealthScore/100.0)*weights["market"] +
		(riskScore/100.0)*weights["risk"]

	breakdown.TotalScore = weightedTotal

	// Boost for Top-Tier tokens (Rank 1-50)
	if token.Rank > 0 && token.Rank <= 50 {
		breakdown.TotalScore += 5.0
	} else if token.Rank > 0 && token.Rank <= 100 {
		breakdown.TotalScore += 2.5
	}

	// Clamp to 100
	if breakdown.TotalScore > 100 {
		breakdown.TotalScore = 100
	}

	// Assign grade and confidence
	breakdown.Grade = s.assignGrade(breakdown.TotalScore)
	breakdown.Confidence = s.calculateConfidence(token)

	return breakdown
}

// 1. Liquidity Scoring (25%)
func (s *EnhancedScorer) calculateLiquidityScore(token *models.Token, details *models.ScoreDetails) float64 {
	var totalRaw float64

	// A. Liquidity to Market Cap Ratio (60%)
	if token.MarketCap > 0 && token.Liquidity > 0 {
		liquidityRatio := token.Liquidity / token.MarketCap
		details.LiquidityRatio = liquidityRatio

		var ratioScore float64
		switch {
		case liquidityRatio >= 0.05: // Excellent (5%+)
			ratioScore = 1.0
		case liquidityRatio >= 0.02: // Good
			ratioScore = 0.9
		case liquidityRatio >= 0.005: // Acceptable
			ratioScore = 0.7
		default: // Thin
			ratioScore = 0.4
		}
		totalRaw += ratioScore * 60.0
	} else if token.Rank > 0 && token.Rank <= 500 {
		// Tier 1 tokens often have deep CEX liquidity not captured by DexScreener
		totalRaw += 50.0 // Default baseline for major tokens
	}

	// B. Liquidity Depth (40%) - Normalized to orderbook presence
	depthScore := s.calculateDepthScore(token)
	details.LiquidityDepth = depthScore
	totalRaw += depthScore * 40.0

	return totalRaw
}

func (s *EnhancedScorer) calculateDepthScore(token *models.Token) float64 {
	if token.MarketCap == 0 {
		return 0
	}

	// Calculate depth as percentage of market cap
	depth1Pct := token.OrderbookDepth1 / token.MarketCap * 100
	depth5Pct := token.OrderbookDepth5 / token.MarketCap * 100

	// Weighted average (1% depth is more important)
	avgDepth := (depth1Pct * 0.7) + (depth5Pct * 0.3)

	switch {
	case avgDepth >= 5.0: // Excellent depth
		return 1.0
	case avgDepth >= 3.0: // Good depth
		return 0.85
	case avgDepth >= 1.0: // Acceptable depth
		return 0.7
	case avgDepth >= 0.5: // Thin depth
		return 0.5
	default: // Very thin
		return 0.3
	}
}

// 2. Volume Scoring (20%)
func (s *EnhancedScorer) calculateVolumeScore(token *models.Token, details *models.ScoreDetails) float64 {
	var totalRaw float64

	// A. Volume to Market Cap Ratio (50%)
	if token.MarketCap > 0 && token.Volume24h > 0 {
		volumeRatio := token.Volume24h / token.MarketCap
		details.VolumeToMcapRatio = volumeRatio

		var ratioScore float64
		switch {
		case volumeRatio >= 0.20: // Healthy volume (20%+)
			ratioScore = 1.0
		case volumeRatio >= 0.10: // Good
			ratioScore = 0.9
		case volumeRatio >= 0.02: // Acceptable
			ratioScore = 0.7
		default: // Low
			ratioScore = 0.4
		}
		totalRaw += ratioScore * 50.0
	}

	// B. Volume Consistency (25%)
	consistencyScore := s.calculateVolumeConsistency(token.VolumeHistory)
	details.VolumeConsistency = consistencyScore
	totalRaw += consistencyScore * 25.0

	// C. Volume Growth (25%)
	volumeGrowth := s.calculateVolumeGrowth(token.VolumeHistory)
	details.Volume7dGrowth = volumeGrowth

	var growthScore float64
	switch {
	case volumeGrowth >= 30: // +30% growth
		growthScore = 1.0
	case volumeGrowth >= 10: // +10% growth
		growthScore = 0.9
	case volumeGrowth >= -5: // Stable
		growthScore = 0.75
	default: // Declining
		growthScore = 0.4
	}
	totalRaw += growthScore * 25.0

	return totalRaw
}

func (s *EnhancedScorer) calculateVolumeConsistency(history models.VolumeHistory) float64 {
	if len(history.Last7Days) < 7 {
		return 0.5 // Insufficient data
	}

	mean := utils.CalculateMean(history.Last7Days)
	if mean == 0 {
		return 0
	}

	stdDev := utils.CalculateStdDev(history.Last7Days, mean)
	cv := stdDev / mean

	// Lower CV = more consistent
	switch {
	case cv <= 0.2:
		return 1.0
	case cv <= 0.4:
		return 0.85
	case cv <= 0.6:
		return 0.7
	case cv <= 0.8:
		return 0.5
	default:
		return 0.3
	}
}

func (s *EnhancedScorer) calculateVolumeGrowth(history models.VolumeHistory) float64 {
	if len(history.Last7Days) < 7 {
		return 0
	}

	early := utils.CalculateMean(history.Last7Days[:3])
	recent := utils.CalculateMean(history.Last7Days[4:])

	if early == 0 {
		return 0
	}

	return ((recent - early) / early) * 100
}

// 3. TVL Scoring (20%)
func (s *EnhancedScorer) calculateTVLScore(token *models.Token, details *models.ScoreDetails) float64 {
	var totalRaw float64

	// A. TVL to Market Cap Ratio (50%)
	if token.MarketCap > 0 && token.TVL > 0 {
		tvlRatio := token.TVL / token.MarketCap
		details.TVLToMcapRatio = tvlRatio

		var ratioScore float64
		switch {
		case tvlRatio >= 0.5: // Healthy (50%+)
			ratioScore = 1.0
		case tvlRatio >= 0.2: // Good
			ratioScore = 0.9
		case tvlRatio >= 0.05: // Acceptable
			ratioScore = 0.7
		default: // Low
			ratioScore = 0.4
		}
		totalRaw += ratioScore * 50.0
	} else if token.Rank > 0 && token.Rank <= 100 {
		// Top 100 tokens often have secondary utility not captured by TVL
		totalRaw += 40.0
	}

	// B. TVL Growth (25%)
	tvl7dGrowth := token.TVL7dChange
	tvl30dGrowth := token.TVL30dChange
	details.TVL7dGrowth = tvl7dGrowth
	details.TVL30dGrowth = tvl30dGrowth

	avgGrowth := (tvl7dGrowth * 0.7) + (tvl30dGrowth * 0.3)
	var growthScore float64
	switch {
	case avgGrowth >= 10:
		growthScore = 1.0
	case avgGrowth >= 0:
		growthScore = 0.8
	case avgGrowth >= -10:
		growthScore = 0.5
	default:
		growthScore = 0.3
	}
	totalRaw += growthScore * 25.0

	// C. TVL Stability (25%)
	tvlVolatility := s.calculateTVLVolatility(token.TVLHistory)
	details.TVLVolatility = tvlVolatility
	var stabilityScore float64
	switch {
	case tvlVolatility <= 10:
		stabilityScore = 1.0
	case tvlVolatility <= 25:
		stabilityScore = 0.8
	default:
		stabilityScore = 0.5
	}
	totalRaw += stabilityScore * 25.0

	return totalRaw
}

func (s *EnhancedScorer) calculateTVLVolatility(history models.TVLHistory) float64 {
	if len(history.Last30Days) < 30 {
		return 50.0 // High uncertainty
	}

	return utils.CalculateVolatility(history.Last30Days)
}

// 4. Trend & Momentum Scoring (20%)
func (s *EnhancedScorer) calculateTrendScore(token *models.Token, details *models.ScoreDetails) float64 {
	var totalRaw float64

	// A. Multi-Timeframe Trend (50%)
	trendRaw := s.analyzeMultiTimeframeTrend(token, details)
	totalRaw += trendRaw * 50.0

	// B. Volatility (25%)
	volRaw := s.calculateVolatilityScore(token, details)
	totalRaw += volRaw * 25.0

	// C. Momentum (25%)
	momRaw := s.calculateMomentumScore(token, details)
	totalRaw += momRaw * 25.0

	return totalRaw
}

func (s *EnhancedScorer) analyzeMultiTimeframeTrend(token *models.Token, details *models.ScoreDetails) float64 {
	shortTrend := s.determineTrend(token.Change7d, "short")
	mediumTrend := s.determineTrend(token.Change30d, "medium")
	longTrend := s.determineTrend(token.Change90d, "long")

	details.ShortTermTrend = shortTrend
	details.MediumTermTrend = mediumTrend
	details.LongTermTrend = longTrend

	// Score based on trend alignment
	if shortTrend == "strong_uptrend" && mediumTrend == "strong_uptrend" {
		return 1.0
	} else if shortTrend == "uptrend" && mediumTrend == "uptrend" {
		return 0.9
	} else if shortTrend == "uptrend" || mediumTrend == "uptrend" {
		return 0.75
	} else if shortTrend == "sideways" && mediumTrend == "sideways" {
		return 0.6
	} else if shortTrend == "downtrend" && mediumTrend == "downtrend" {
		return 0.3
	}
	return 0.5 // Mixed signals
}

func (s *EnhancedScorer) determineTrend(changePercent float64, timeframe string) string {
	var strongUp, up, strongDown, down float64

	switch timeframe {
	case "short": // 7 days
		strongUp, up, down, strongDown = 15, 5, -5, -15
	case "medium": // 30 days
		strongUp, up, down, strongDown = 30, 10, -10, -30
	case "long": // 90 days
		strongUp, up, down, strongDown = 50, 20, -20, -50
	}

	switch {
	case changePercent >= strongUp:
		return "strong_uptrend"
	case changePercent >= up:
		return "uptrend"
	case changePercent <= strongDown:
		return "strong_downtrend"
	case changePercent <= down:
		return "downtrend"
	default:
		return "sideways"
	}
}

func (s *EnhancedScorer) calculateVolatilityScore(token *models.Token, details *models.ScoreDetails) float64 {
	vol30d := token.Volatility30d
	details.Volatility30d = vol30d

	switch {
	case vol30d <= 20:
		return 1.0
	case vol30d <= 40:
		return 0.8
	case vol30d <= 60:
		return 0.6
	case vol30d <= 80:
		return 0.4
	default:
		return 0.2
	}
}

func (s *EnhancedScorer) calculateMomentumScore(token *models.Token, details *models.ScoreDetails) float64 {
	var score float64

	// RSI Analysis (40%)
	rsi := token.RSI14
	details.RSI = rsi

	var rsiScore float64
	switch {
	case rsi >= 50 && rsi <= 70:
		rsiScore = 1.0
	case rsi >= 40 && rsi < 50:
		rsiScore = 0.7
	case rsi > 70:
		rsiScore = 0.5
	case rsi >= 30 && rsi < 40:
		rsiScore = 0.6
	default:
		rsiScore = 0.4
	}
	score += rsiScore * 0.4

	// Moving Average Convergence (30%)
	priceAboveSMA7 := token.Price > token.SMA7
	priceAboveSMA30 := token.Price > token.SMA30

	if priceAboveSMA7 && priceAboveSMA30 {
		score += 1.0 * 0.3
	} else if priceAboveSMA7 {
		score += 0.7 * 0.3
	} else if priceAboveSMA30 {
		score += 0.5 * 0.3
	} else {
		score += 0.3 * 0.3
	}

	// MACD Signal (30%)
	macd := token.MACD
	var macdScore float64
	if macd > 0 {
		macdScore = 1.0
	} else if macd > -2 {
		macdScore = 0.7
	} else {
		macdScore = 0.4
	}
	score += macdScore * 0.3

	return score
}

// 5. Market Health Scoring (10%)
func (s *EnhancedScorer) calculateMarketHealthScore(token *models.Token, details *models.ScoreDetails) float64 {
	var totalRaw float64

	// A. Holder Distribution (50%)
	top10Ratio := token.Top10HoldersRatio
	details.Top10HoldersPercent = top10Ratio
	var distributionScore float64
	switch {
	case top10Ratio <= 20: // Highly decentralized
		distributionScore = 1.0
	case top10Ratio <= 50: // Good
		distributionScore = 0.85
	case top10Ratio <= 80: // Centralized
		distributionScore = 0.5
	default:
		distributionScore = 0.2
	}
	totalRaw += distributionScore * 50.0

	// B. Market Dominance (50%)
	details.UniqueHolders = token.HolderCount
	details.MarketCapRank = token.Rank
	var rankScore float64
	switch {
	case token.Rank <= 100:
		rankScore = 1.0
	case token.Rank <= 500:
		rankScore = 0.85
	default:
		rankScore = 0.6
	}
	var holderScore float64
	switch {
	case token.HolderCount >= 50000:
		holderScore = 1.0
	case token.HolderCount >= 10000:
		holderScore = 0.8
	default:
		holderScore = 0.5
	}
	dominanceScore := (rankScore * 0.7) + (holderScore * 0.3)
	totalRaw += dominanceScore * 50.0

	return totalRaw
}

// 6. Risk Scoring (5%)
func (s *EnhancedScorer) calculateRiskScore(token *models.Token, details *models.ScoreDetails) float64 {
	// Raw risk score (0-100)
	rugPullRisk := s.assessRugPullRisk(token)
	centralizationRisk := s.assessCentralizationRisk(token)
	contractRisk := s.assessSmartContractRisk(token)

	details.RugPullRisk = rugPullRisk
	details.CentralizationRisk = centralizationRisk
	details.SmartContractRisk = contractRisk

	riskScores := []float64{
		s.riskToScore(rugPullRisk),
		s.riskToScore(centralizationRisk),
		s.riskToScore(contractRisk),
	}

	avgRiskScore := utils.CalculateMean(riskScores)
	return avgRiskScore * 100.0 // Return as raw 0-100
}

func (s *EnhancedScorer) riskToScore(risk string) float64 {
	switch risk {
	case "Low":
		return 1.0
	case "Medium":
		return 0.6
	case "High":
		return 0.2
	default:
		return 0.5
	}
}

func (s *EnhancedScorer) assessRugPullRisk(token *models.Token) string {
	riskPoints := 0

	if token.Liquidity < token.MarketCap*0.01 {
		riskPoints += 2
	}
	if token.Top10HoldersRatio > 70 {
		riskPoints += 2
	}
	if token.ContractAge < 90 {
		riskPoints += 1
	}
	if !token.IsVerified {
		riskPoints += 2
	}

	if riskPoints >= 5 {
		return "High"
	} else if riskPoints >= 3 {
		return "Medium"
	}
	return "Low"
}

func (s *EnhancedScorer) assessCentralizationRisk(token *models.Token) string {
	if token.Top10HoldersRatio > 80 {
		return "High"
	} else if token.Top10HoldersRatio > 60 {
		return "Medium"
	}
	return "Low"
}

func (s *EnhancedScorer) assessSmartContractRisk(token *models.Token) string {
	riskPoints := 0

	if !token.IsVerified {
		riskPoints += 2
	}
	if token.AuditStatus != "Passed" {
		riskPoints += 2
	}
	if token.ContractAge < 180 {
		riskPoints += 1
	}

	if riskPoints >= 4 {
		return "High"
	} else if riskPoints >= 2 {
		return "Medium"
	}
	return "Low"
}

// Helper functions
func (s *EnhancedScorer) assignGrade(score float64) string {
	switch {
	case score >= 90:
		return "S" // Exceptional
	case score >= 80:
		return "A" // Excellent
	case score >= 70:
		return "B" // Good
	case score >= 60:
		return "C" // Average
	case score >= 50:
		return "D" // Below Average
	default:
		return "F" // Poor
	}
}

func (s *EnhancedScorer) calculateConfidence(token *models.Token) float64 {
	confidence := 100.0

	// Deduct for missing data
	if token.TVL == 0 {
		confidence -= 15
	}
	if token.Liquidity == 0 {
		confidence -= 15
	}
	if token.HolderCount == 0 {
		confidence -= 10
	}
	if len(token.PriceHistory.Last30Days) < 30 {
		confidence -= 10
	}
	if token.Volume24h < token.MarketCap*0.001 {
		confidence -= 10
	}
	if token.ContractAge < 30 {
		confidence -= 15
	} else if token.ContractAge < 90 {
		confidence -= 10
	}

	return math.Max(confidence, 0)
}

// CalculateScoresForAll calculates scores for all tokens using enhanced scorer
func (s *EnhancedScorer) CalculateScoresForAll(tokens []models.Token) []models.Token {
	for i := range tokens {
		breakdown := s.CalculateComprehensiveScore(&tokens[i])
		tokens[i].TrustScore = breakdown.TotalScore
		tokens[i].ScoreBreakdown = breakdown
	}
	return tokens
}
