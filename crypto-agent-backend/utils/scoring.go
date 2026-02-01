package utils

import (
	"crypto-agent-backend/models"
	"math"
)

// CalculateTrustScore computes the comprehensive trust score (0-100)
func CalculateTrustScore(token *models.Token) {
	// 1. Weights
	const (
		WLiq    = 25.0
		WVol    = 20.0
		WTVL    = 20.0
		WTrend  = 20.0
		WMarket = 10.0
		WRisk   = 5.0
	)

	// 2. Calculate Sub-Scores (Raw 0-100)
	lScore := calculateLiquidityScore(token)
	vScore := calculateVolumeScore(token)
	tScore := calculateTVLScore(token)
	trendScore := calculateTrendScore(token)
	mScore := calculateMarketHealthScore(token)
	rScore := calculateRiskScore(token)

	// 3. Dynamic Weighting (Handle missing TVL/Liq)
	wLiq, wVol, wTVL, wTrend, wMarket, wRisk := WLiq, WVol, WTVL, WTrend, WMarket, WRisk

	if token.TVL == 0 {
		wTVL = 0
		wVol += WTVL * 0.4
		wTrend += WTVL * 0.4
		wMarket += WTVL * 0.2
	}
	if token.Liquidity == 0 {
		wLiq = 0
		wVol += WLiq * 0.6
		wRisk += WLiq * 0.4
	}

	// 4. Final Weighted Total
	totalScore := (lScore/100.0)*wLiq +
		(vScore/100.0)*wVol +
		(tScore/100.0)*wTVL +
		(trendScore/100.0)*wTrend +
		(mScore/100.0)*wMarket +
		(rScore/100.0)*wRisk

	// 5. Rank Boost
	if token.Rank > 0 && token.Rank <= 50 {
		totalScore += 5.0
	} else if token.Rank > 0 && token.Rank <= 100 {
		totalScore += 2.5
	}

	// Bulletproof Clamp
	if totalScore > 100 {
		totalScore = 100
	}

	// 6. Assign to Token
	token.TrustScore = totalScore
	token.ScoreBreakdown = models.DetailedScoreBreakdown{
		TotalScore:        totalScore,
		LiquidityScore:    lScore,
		VolumeScore:       vScore,
		TVLScore:          tScore,
		TrendScore:        trendScore,
		MarketHealthScore: mScore,
		RiskScore:         rScore,
		Grade:             AssignGrade(totalScore),
	}
}

// AssignGrade converts numerical score to Letter Grade
func AssignGrade(score float64) string {
	switch {
	case score >= 90:
		return "SSS"
	case score >= 80:
		return "AA"
	case score >= 70:
		return "A"
	case score >= 60:
		return "BBB"
	case score >= 50:
		return "BB"
	case score >= 40:
		return "B"
	default:
		return "D"
	}
}

func calculateLiquidityScore(token *models.Token) float64 {
	if token.MarketCap == 0 {
		return 50.0
	}
	ratio := token.Liquidity / token.MarketCap
	if ratio >= 0.05 {
		return 100
	}
	if ratio >= 0.02 {
		return 80
	}
	if ratio >= 0.005 {
		return 60
	}
	return 40
}

func calculateVolumeScore(token *models.Token) float64 {
	if token.MarketCap == 0 || token.Volume24h == 0 {
		return 40
	}
	ratio := token.Volume24h / token.MarketCap
	if ratio >= 0.20 {
		return 100
	}
	if ratio >= 0.10 {
		return 80
	}
	if ratio >= 0.02 {
		return 60
	}
	return 40
}

func calculateTVLScore(token *models.Token) float64 {
	if token.TVL == 0 {
		if token.Rank > 0 && token.Rank <= 100 {
			return 40
		}
		return 0
	}
	if token.MarketCap == 0 {
		return 50
	}
	ratio := token.TVL / token.MarketCap
	if ratio >= 0.5 {
		return 100
	}
	if ratio >= 0.2 {
		return 80
	}
	return 60
}

func calculateTrendScore(token *models.Token) float64 {
	score := 60.0
	if token.Change24h > 5 {
		score += 20
	}
	if token.Change7d > 15 {
		score += 20
	}
	if token.Change24h < -10 {
		score -= 20
	}
	return math.Min(math.Max(score, 0), 100)
}

func calculateMarketHealthScore(token *models.Token) float64 {
	score := 70.0
	if token.Rank > 0 && token.Rank <= 100 {
		score += 20
	}
	if token.Rank > 0 && token.Rank <= 500 {
		score += 10
	}
	if token.Top10HoldersRatio > 0.8 {
		score -= 30
	}
	return math.Min(math.Max(score, 0), 100)
}

func calculateRiskScore(token *models.Token) float64 {
	score := 80.0
	if !token.IsVerified {
		score -= 40
	}
	if token.Liquidity < 100000 {
		score -= 20
	}
	return math.Min(math.Max(score, 0), 100)
}
