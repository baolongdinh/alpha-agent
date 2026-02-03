package utils

import "math"

// Mathematical utility functions for scoring calculations

// CalculateMean computes the average of a slice of float64 values
func CalculateMean(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

// CalculateStdDev computes standard deviation given values and their mean
func CalculateStdDev(values []float64, mean float64) float64 {
	if len(values) == 0 {
		return 0
	}

	sumSquares := 0.0
	for _, v := range values {
		diff := v - mean
		sumSquares += diff * diff
	}

	variance := sumSquares / float64(len(values))
	return math.Sqrt(variance)
}

// CalculateRSI computes the Relative Strength Index
func CalculateRSI(prices []float64, period int) float64 {
	if len(prices) < period+1 {
		return 50.0 // Neutral
	}

	gains := 0.0
	losses := 0.0

	for i := 1; i <= period; i++ {
		change := prices[i] - prices[i-1]
		if change > 0 {
			gains += change
		} else {
			losses -= change
		}
	}

	avgGain := gains / float64(period)
	avgLoss := losses / float64(period)

	if avgLoss == 0 {
		return 100.0
	}

	rs := avgGain / avgLoss
	rsi := 100 - (100 / (1 + rs))

	return rsi
}

// CalculateSMA computes Simple Moving Average
func CalculateSMA(prices []float64, period int) float64 {
	if len(prices) < period {
		return 0
	}

	sum := 0.0
	for i := len(prices) - period; i < len(prices); i++ {
		sum += prices[i]
	}

	return sum / float64(period)
}

// CalculateEMA computes Exponential Moving Average
func CalculateEMA(prices []float64, period int) float64 {
	if len(prices) < period {
		return 0
	}

	multiplier := 2.0 / float64(period+1)

	// Start with SMA
	ema := CalculateSMA(prices[:period], period)

	// Calculate EMA for remaining prices
	for i := period; i < len(prices); i++ {
		ema = (prices[i]-ema)*multiplier + ema
	}

	return ema
}

// CalculateMACD computes Moving Average Convergence Divergence
func CalculateMACD(prices []float64) float64 {
	if len(prices) < 26 {
		return 0
	}

	ema12 := CalculateEMA(prices, 12)
	ema26 := CalculateEMA(prices, 26)

	return ema12 - ema26
}

// CalculateSharpeRatio computes risk-adjusted returns
func CalculateSharpeRatio(returns []float64, riskFreeRate float64) float64 {
	if len(returns) == 0 {
		return 0
	}

	meanReturn := CalculateMean(returns)
	stdDev := CalculateStdDev(returns, meanReturn)

	if stdDev == 0 {
		return 0
	}

	excessReturn := meanReturn - riskFreeRate
	return excessReturn / stdDev
}

// CalculateVolatility computes coefficient of variation as percentage
func CalculateVolatility(prices []float64) float64 {
	if len(prices) == 0 {
		return 0
	}

	mean := CalculateMean(prices)
	if mean == 0 {
		return 100.0
	}

	stdDev := CalculateStdDev(prices, mean)
	return (stdDev / mean) * 100
}
