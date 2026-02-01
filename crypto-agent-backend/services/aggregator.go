package services

import (
	"context"
	"crypto-agent-backend/config"
	"crypto-agent-backend/models"
	"crypto-agent-backend/utils"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
)

// Aggregator (Enhanced) handles multi-source data aggregation
type Aggregator struct {
	config *config.Config
	cache  *cache.Cache
}

// NewAggregator creates a new enhanced aggregator
func NewAggregator(cfg *config.Config) *Aggregator {
	return &Aggregator{
		config: cfg,
		cache:  cache.New(1*time.Minute, 5*time.Minute),
	}
}

// FetchAllTokenData fetches data from ALL sources concurrently
func (a *Aggregator) FetchAllTokenData(ctx context.Context) ([]models.Token, error) {
	if cached, found := a.cache.Get("all_tokens"); found {
		log.Println("⚡ Serving data from cache")
		return cached.([]models.Token), nil
	}

	startTime := time.Now()

	// Data collection structures
	type dataResult struct {
		name string
		data interface{}
		err  error
	}
	resultChan := make(chan dataResult, 10)
	var wg sync.WaitGroup

	// Task list for concurrent execution
	tasks := []struct {
		name string
		fn   func(context.Context) (interface{}, error)
	}{
		{"CoinMarketCap", a.fetchCoinMarketCap},
		{"CoinGecko", a.fetchCoinGeckoAsync},
		{"DeFiLlama", a.fetchDefiLlamaAsync},
		{"DexScreener", a.fetchDexScreenerAsync},
		{"Messari", a.fetchMessari},
	}

	// Launch all fetchers
	for _, task := range tasks {
		wg.Add(1)
		go func(taskName string, fetchFn func(context.Context) (interface{}, error)) {
			defer wg.Done()
			data, err := fetchFn(ctx)
			resultChan <- dataResult{name: taskName, data: data, err: err}
		}(task.name, task.fn)
	}

	// Wait and close
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collect results
	results := make(map[string]interface{})
	for res := range resultChan {
		if res.err != nil {
			log.Printf("✗ %s fetch error: %v", res.name, res.err)
		} else if res.data != nil {
			results[res.name] = res.data
			log.Printf("✓ %s data received", res.name)
		}
	}

	// Merge with user-requested priority (CMC first)
	tokens := utils.MergeEnhancedData(results)

	if len(tokens) > 0 {
		a.cache.Set("all_tokens", tokens, cache.DefaultExpiration)
	}

	log.Printf("✓ Data aggregation completed in %v: %d total tokens", time.Since(startTime), len(tokens))
	return tokens, nil
}

// fetchCoinMarketCap fetches primary data from CMC
func (a *Aggregator) fetchCoinMarketCap(ctx context.Context) (interface{}, error) {
	if a.config.CoinMarketCapAPIKey == "" {
		return nil, nil // Silently skip if no key
	}
	// Removed aux=sparkline to restore 2000 token limit compatibility on free tier
	url := fmt.Sprintf("%s/v1/cryptocurrency/listings/latest?limit=2000", a.config.CoinMarketCapAPIURL)
	headers := map[string]string{"X-CMC_PRO_API_KEY": a.config.CoinMarketCapAPIKey}

	data, err := utils.FetchJSONWithHeaders(url, headers)
	if err != nil {
		return nil, err
	}
	var resp models.CoinMarketCapResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// fetchCoinGeckoAsync wrapper - Fetches multiple pages for wider coverage
func (a *Aggregator) fetchCoinGeckoAsync(ctx context.Context) (interface{}, error) {
	var allMarkets []models.CoinGeckoMarket
	// Fetch 4 pages of 250 (max allowed per page for free tier) = 1000 tokens
	for page := 1; page <= 4; page++ {
		url := fmt.Sprintf("%s/coins/markets?vs_currency=usd&order=market_cap_desc&per_page=250&page=%d&sparkline=true&price_change_percentage=7d", a.config.CoinGeckoAPIURL, page)
		data, err := utils.FetchJSON(url)
		if err != nil {
			log.Printf("CoinGecko page %d fetch error: %v", page, err)
			continue
		}
		var markets []models.CoinGeckoMarket
		if err := json.Unmarshal(data, &markets); err == nil {
			allMarkets = append(allMarkets, markets...)
		}
		// Small delay to respect rate limits if needed (though FetchJSON should handle retries)
		time.Sleep(100 * time.Millisecond)
	}
	return allMarkets, nil
}

// fetchDefiLlamaAsync wrapper
func (a *Aggregator) fetchDefiLlamaAsync(ctx context.Context) (interface{}, error) {
	url := fmt.Sprintf("%s/protocols", a.config.DefiLlamaAPIURL)
	data, err := utils.FetchJSON(url)
	if err != nil {
		return nil, err
	}
	var protocols []models.DefiLlamaProtocol
	if err := json.Unmarshal(data, &protocols); err != nil {
		return nil, err
	}
	return protocols, nil
}

// fetchDexScreenerAsync - Enhanced with dynamic token list
func (a *Aggregator) fetchDexScreenerAsync(ctx context.Context) (interface{}, error) {
	// Solana top tokens (Dynamic selection could be added here)
	tokens := []string{
		"So11111111111111111111111111111111111111112",  // SOL
		"DezXAZ8z7PnrnRJjz3wXBoRgixCa6xjnB7YaB1pPB263", // BONK
		"EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v", // USDC
	}

	dexData := make(map[string]models.DexScreenerResponse)
	for _, addr := range tokens {
		url := fmt.Sprintf("%s/dex/tokens/%s", a.config.DexScreenerAPIURL, addr)
		data, err := utils.FetchJSON(url)
		if err == nil {
			var resp models.DexScreenerResponse
			if json.Unmarshal(data, &resp) == nil && len(resp.Pairs) > 0 {
				symbol := strings.ToUpper(resp.Pairs[0].BaseToken.Symbol)
				dexData[symbol] = resp
			}
		}
	}
	return dexData, nil
}

// fetchMessari fundamentales
func (a *Aggregator) fetchMessari(ctx context.Context) (interface{}, error) {
	if a.config.MessariAPIKey == "" {
		return nil, nil
	}
	// Use Messari v2 endpoint on api.messari.io as requested
	url := "https://api.messari.io/metrics/v2/assets?limit=500"
	headers := map[string]string{"X-Messari-API-Key": a.config.MessariAPIKey}
	data, err := utils.FetchJSONWithHeaders(url, headers)
	if err != nil {
		return nil, err
	}
	var resp models.MessariResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

func (a *Aggregator) FetchGlobalMarketStats(ctx context.Context) (*models.MarketStats, error) {
	if cached, found := a.cache.Get("global_stats"); found {
		return cached.(*models.MarketStats), nil
	}
	url := fmt.Sprintf("%s/global", a.config.CoinGeckoAPIURL)
	data, err := utils.FetchJSON(url)
	if err != nil {
		return nil, err
	}
	var response models.CoinGeckoGlobalResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, err
	}
	stats := &models.MarketStats{
		TotalMarketCap:      response.Data.TotalMarketCap["usd"],
		TotalVolume:         response.Data.TotalVolume["usd"],
		MarketCapPercentage: response.Data.MarketCapPercentage,
		MarketCapChange24h:  response.Data.MarketCapChangePercent,
		BTCDominance:        response.Data.MarketCapPercentage["btc"],
		ETHDominance:        response.Data.MarketCapPercentage["eth"],
	}
	a.cache.Set("global_stats", stats, 10*time.Minute)
	return stats, nil
}
