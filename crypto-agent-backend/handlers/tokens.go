package handlers

import (
	"crypto-agent-backend/cache"
	"crypto-agent-backend/config"
	"crypto-agent-backend/models"
	"crypto-agent-backend/services"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// TokenHandler handles token-related endpoints
type TokenHandler struct {
	aggregator *services.Aggregator
	scorer     *services.EnhancedScorer
	cache      *cache.CacheManager
	config     *config.Config
}

// NewTokenHandler creates a new token handler
func NewTokenHandler(
	aggregator *services.Aggregator,
	scorer *services.EnhancedScorer,
	cacheManager *cache.CacheManager,
	cfg *config.Config,
) *TokenHandler {
	return &TokenHandler{
		aggregator: aggregator,
		scorer:     scorer,
		cache:      cacheManager,
		config:     cfg,
	}
}

// GetTokens handles GET /api/tokens
func (h *TokenHandler) GetTokens(c *gin.Context) {
	startTime := time.Now()

	// Parse query parameters
	params := h.parseFilterParams(c)

	// Create cache key based on params
	cacheKey := h.buildCacheKey(params)

	// Try cache first
	if cached, found := h.cache.GetTokens(cacheKey); found {
		log.Printf("✓ Cache hit for key: %s", cacheKey)

		tokens := cached.([]models.Token)
		filtered := h.filterAndSortTokens(tokens, params)

		c.JSON(http.StatusOK, models.TokensResponse{
			Status:      "success",
			Timestamp:   time.Now(),
			Total:       len(filtered),
			Data:        filtered,
			FetchTimeMs: time.Since(startTime).Milliseconds(),
		})
		return
	}

	log.Printf("✗ Cache miss for key: %s - fetching fresh data", cacheKey)

	// Fetch fresh data
	tokens, err := h.aggregator.FetchAllTokenData(c.Request.Context())
	if err != nil {
		log.Printf("❌ Failed to fetch token data: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Status:    "error",
			Message:   "Failed to fetch token data: " + err.Error(),
			Timestamp: time.Now(),
		})
		return
	}

	// Calculate scores for all tokens
	tokens = h.scorer.CalculateScoresForAll(tokens)

	// Cache the results
	h.cache.SetTokens(cacheKey, tokens)
	log.Printf("✓ Cached %d tokens with key: %s", len(tokens), cacheKey)

	// Filter and sort
	filtered := h.filterAndSortTokens(tokens, params)

	fetchDuration := time.Since(startTime)
	log.Printf("✅ Request completed in %v: %d tokens returned", fetchDuration, len(filtered))

	c.JSON(http.StatusOK, models.TokensResponse{
		Status:      "success",
		Timestamp:   time.Now(),
		Total:       len(filtered),
		Data:        filtered,
		FetchTimeMs: fetchDuration.Milliseconds(),
	})
}

// GetMarketStats handles GET /api/market/stats
func (h *TokenHandler) GetMarketStats(c *gin.Context) {
	stats, err := h.aggregator.FetchGlobalMarketStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}

// GetTokenByID handles GET /api/tokens/:id
func (h *TokenHandler) GetTokenByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	// Try to find the token in the "tokens_all" cache
	if cached, found := h.cache.GetTokens("tokens_all"); found {
		tokens := cached.([]models.Token)
		for _, token := range tokens {
			// Check if ID (coingecko id) or Symbol matches
			// Note: We use lowercase comparison for loose matching
			normalizedId := strings.ToLower(id)
			if strings.ToLower(token.ID) == normalizedId || strings.ToLower(token.Symbol) == normalizedId || strings.ToLower(token.Name) == normalizedId {
				c.JSON(http.StatusOK, token)
				return
			}
		}
	}

	// If not in cache or not found, try a fresh fetch
	tokens, err := h.aggregator.FetchAllTokenData(c.Request.Context())
	if err == nil {
		h.cache.SetTokens("tokens_all", tokens)
		for _, token := range tokens {
			normalizedId := strings.ToLower(id)
			if strings.ToLower(token.ID) == normalizedId || strings.ToLower(token.Symbol) == normalizedId || strings.ToLower(token.Name) == normalizedId {
				c.JSON(http.StatusOK, token)
				return
			}
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "token not found"})
}

// parseFilterParams extracts filter parameters from request
func (h *TokenHandler) parseFilterParams(c *gin.Context) models.FilterParams {
	params := models.FilterParams{
		Limit: 50, // Default limit
	}

	// Min market cap
	if minMcap := c.Query("min_mcap"); minMcap != "" {
		if val, err := strconv.ParseFloat(minMcap, 64); err == nil {
			params.MinMcap = val
		}
	}

	// Max market cap
	if maxMcap := c.Query("max_mcap"); maxMcap != "" {
		if val, err := strconv.ParseFloat(maxMcap, 64); err == nil {
			params.MaxMcap = val
		}
	}

	// Category filter
	if category := c.Query("category"); category != "" {
		params.Category = category
	}

	// Limit
	if limit := c.Query("limit"); limit != "" {
		if val, err := strconv.Atoi(limit); err == nil {
			if val > 0 && val <= 2000 {
				params.Limit = val
			}
		}
	}

	return params
}

// buildCacheKey creates a cache key from filter params
func (h *TokenHandler) buildCacheKey(params models.FilterParams) string {
	return "tokens_all" // Simple key for now; can be enhanced with params
}

// filterAndSortTokens applies filters and sorting
func (h *TokenHandler) filterAndSortTokens(tokens []models.Token, params models.FilterParams) []models.Token {
	filtered := make([]models.Token, 0)

	for _, token := range tokens {
		// Filter out tokens with incomplete data
		if token.Price <= 0 || token.MarketCap <= 0 {
			continue // Skip tokens without basic data
		}

		// Market cap filter
		if params.MinMcap > 0 && token.MarketCap < params.MinMcap {
			continue
		}
		if params.MaxMcap > 0 && token.MarketCap > params.MaxMcap {
			continue
		}

		// Category filter
		if params.Category != "" && token.Category != params.Category {
			continue
		}

		filtered = append(filtered, token)
	}

	// Sort by trust score (descending)
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].TrustScore > filtered[j].TrustScore
	})

	// Apply limit
	if params.Limit > 0 && len(filtered) > params.Limit {
		filtered = filtered[:params.Limit]
	}

	// Add rankings
	for i := range filtered {
		filtered[i].Rank = i + 1
	}

	return filtered
}
