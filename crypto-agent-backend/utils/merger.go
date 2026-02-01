package utils

import (
	"crypto-agent-backend/models"
	"fmt"
	"strings"
)

// MergeEnhancedData combines data from all new sources with priority weighting
func MergeEnhancedData(results map[string]interface{}) []models.Token {
	tokenMap := make(map[string]*models.Token)

	// 1. CoinMarketCap - PRIMARY SOURCE (as requested by user)
	if cmcData, ok := results["CoinMarketCap"].([]models.CoinMarketCapCoin); ok {
		for _, coin := range cmcData {
			symbol := NormalizeSymbol(coin.Symbol)

			token := &models.Token{
				ID:     coin.Slug,
				Symbol: symbol,
				Name:   coin.Name,
				Image:  fmt.Sprintf("https://s2.coinmarketcap.com/static/img/coins/64x64/%d.png", coin.ID),
				Category: func() string {
					if len(coin.Tags) > 0 {
						return strings.Title(coin.Tags[0])
					}
					return ""
				}(),
				Rank:            coin.CMCRank,
				Price:           coin.Quote.USD.Price,
				MarketCap:       coin.Quote.USD.MarketCap,
				Volume24h:       coin.Quote.USD.Volume24h,
				Change1h:        coin.Quote.USD.PercentChange1h,
				Change24h:       coin.Quote.USD.PercentChange24h,
				Change7d:        coin.Quote.USD.PercentChange7d,
				Change30d:       coin.Quote.USD.PercentChange30d,
				Change90d:       coin.Quote.USD.PercentChange90d,
				VolumeChange24h: coin.Quote.USD.VolumeChange24h,
				MaxSupply:       coin.MaxSupply,
				CirculatingSupply: func() float64 {
					if coin.CirculatingSupply > 0 {
						return coin.CirculatingSupply
					}
					if coin.Quote.USD.Price > 0 {
						return coin.Quote.USD.MarketCap / coin.Quote.USD.Price
					}
					return 0
				}(),
				TotalSupply:       coin.TotalSupply,
				FullyDilutedValue: coin.Quote.USD.FullyDilutedMarketCap,
				Sparkline:         coin.Quote.USD.Sparkline,
				MarketCapDom:      coin.Quote.USD.MarketCapDominance,
			}
			tokenMap[symbol] = token
		}
	}

	// 2. CoinGecko - Disabled due to rate limits
	/*
		if cgData, ok := results["CoinGecko"].([]models.CoinGeckoMarket); ok {
			for _, cg := range cgData {
				// ... logic removed ...
			}
		}
	*/

	// 3. DeFiLlama - TVL and Categories
	if defiData, ok := results["DeFiLlama"].([]models.DefiLlamaProtocol); ok {
		for _, dl := range defiData {
			symbol := NormalizeSymbol(dl.Symbol)
			if token, exists := tokenMap[symbol]; exists {
				token.TVL = dl.TVL
				if token.Category == "" {
					token.Category = dl.Category
				}
			}
		}
	}

	// 4. DexScreener - Liquidity
	if dexData, ok := results["DexScreener"].(map[string]models.DexScreenerResponse); ok {
		for symbol, dex := range dexData {
			upperSymbol := NormalizeSymbol(symbol)
			if token, exists := tokenMap[upperSymbol]; exists {
				var totalLiquidity float64
				for _, pair := range dex.Pairs {
					totalLiquidity += pair.Liquidity.USD
				}
				token.Liquidity = totalLiquidity
			}
		}
	}

	// 7. Messari - Fundamental health (NEW)
	if messariData, ok := results["Messari"].([]models.MessariAsset); ok {
		for _, ma := range messariData {
			symbol := NormalizeSymbol(ma.Symbol)
			if token, exists := tokenMap[symbol]; exists {
				// Fallback for missing metrics
				if token.Price == 0 {
					token.Price = ma.MarketData.PriceUSD
				}
				if token.MarketCap == 0 {
					token.MarketCap = ma.MarketData.MarketCap
				}
				if token.Volume24h == 0 {
					token.Volume24h = ma.MarketData.Volume24h
				}
			}
		}
	}

	// 7. DexScreener - PRICE FALLBACK for missing tokens
	if dexData, ok := results["DexScreener"].(map[string]models.DexScreenerResponse); ok {
		for symbol, dex := range dexData {
			upperSymbol := NormalizeSymbol(symbol)
			if token, exists := tokenMap[upperSymbol]; exists {
				// If we still don't have a price, take it from DexScreener
				if token.Price == 0 && len(dex.Pairs) > 0 {
					var p float64
					fmt.Sscanf(dex.Pairs[0].PriceUsd, "%f", &p)
					token.Price = p
				}
			}
		}
	}

	// Final conversion and calculations
	tokens := make([]models.Token, 0, len(tokenMap))
	for _, token := range tokenMap {
		if token.Price > 0 {
			// Derive Volatility from Sparkline if possible
			if len(token.Sparkline) > 0 {
				token.Volatility7d = CalculateVolatility(token.Sparkline)
			}

			// Derivative Metrics
			if token.MarketCap > 0 && token.FullyDilutedValue == 0 {
				// Estimate FDV if missing
				if token.CirculatingSupply > 0 && token.TotalSupply > 0 {
					token.FullyDilutedValue = token.Price * token.TotalSupply
				} else {
					token.FullyDilutedValue = token.MarketCap
				}
			}

			CalculateTrustScore(token)
			tokens = append(tokens, *token)
		}
	}

	return tokens
}

// Keep MergeTokenData for backward compatibility if needed, or redirect it
func MergeTokenData(
	defillama []models.DefiLlamaProtocol,
	coingecko []models.CoinGeckoMarket,
	dexscreener map[string]models.DexScreenerResponse,
) []models.Token {
	results := map[string]interface{}{
		"DeFiLlama":   defillama,
		"CoinGecko":   coingecko,
		"DexScreener": dexscreener,
	}
	return MergeEnhancedData(results)
}

// NormalizeSymbol standardizes token symbols
func NormalizeSymbol(symbol string) string {
	return strings.ToUpper(strings.TrimSpace(symbol))
}
