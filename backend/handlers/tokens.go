package handlers

import (
	"backend/cache"
	"backend/config"
	"backend/models"
	"backend/services"
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
		filtered, total, hasMore := h.filterAndSortTokens(tokens, params)

		c.JSON(http.StatusOK, models.TokensResponse{
			Status:      "success",
			Timestamp:   time.Now(),
			Total:       total,
			Data:        filtered,
			HasMore:     hasMore,
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
	filtered, total, hasMore := h.filterAndSortTokens(tokens, params)

	fetchDuration := time.Since(startTime)
	log.Printf("✅ Request completed in %v: %d tokens returned (total match: %d)", fetchDuration, len(filtered), total)

	c.JSON(http.StatusOK, models.TokensResponse{
		Status:      "success",
		Timestamp:   time.Now(),
		Total:       total,
		Data:        filtered,
		HasMore:     hasMore,
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
		Limit: 5000, // Default limit
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

	// Search
	if search := c.Query("search"); search != "" {
		params.Search = strings.ToLower(search)
	}

	// Price Range
	if minPrice := c.Query("min_price"); minPrice != "" {
		if val, err := strconv.ParseFloat(minPrice, 64); err == nil {
			params.MinPrice = val
		}
	}
	if maxPrice := c.Query("max_price"); maxPrice != "" {
		if val, err := strconv.ParseFloat(maxPrice, 64); err == nil {
			params.MaxPrice = val
		}
	}

	// Score Range
	if minScore := c.Query("min_score"); minScore != "" {
		if val, err := strconv.ParseFloat(minScore, 64); err == nil {
			params.MinScore = val
		}
	}
	if maxScore := c.Query("max_score"); maxScore != "" {
		if val, err := strconv.ParseFloat(maxScore, 64); err == nil {
			params.MaxScore = val
		}
	}

	// Change Range
	if minChange := c.Query("min_change"); minChange != "" {
		if val, err := strconv.ParseFloat(minChange, 64); err == nil {
			params.MinChange = val
		}
	}
	if maxChange := c.Query("max_change"); maxChange != "" {
		if val, err := strconv.ParseFloat(maxChange, 64); err == nil {
			params.MaxChange = val
		}
	}

	// Limit
	if limit := c.Query("limit"); limit != "" {
		if val, err := strconv.Atoi(limit); err == nil {
			if val > 0 && val <= 5000 {
				params.Limit = val
			}
		}
	}

	// Offset
	if offset := c.Query("offset"); offset != "" {
		if val, err := strconv.Atoi(offset); err == nil {
			if val >= 0 {
				params.Offset = val
			}
		}
	}

	return params
}

// buildCacheKey creates a cache key from filter params
func (h *TokenHandler) buildCacheKey(params models.FilterParams) string {
	return "tokens_all" // Simple key for now; can be enhanced with params
}

// filterAndSortTokens applies filters and sorting, then returns the subset, total count and hasMore info
func (h *TokenHandler) filterAndSortTokens(tokens []models.Token, params models.FilterParams) ([]models.Token, int, bool) {
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

		// Search filter (Symbol or Name)
		if params.Search != "" {
			if !strings.Contains(strings.ToLower(token.Symbol), params.Search) &&
				!strings.Contains(strings.ToLower(token.Name), params.Search) {
				continue
			}
		}

		// Price filter
		if params.MinPrice > 0 && token.Price < params.MinPrice {
			continue
		}
		if params.MaxPrice > 0 && token.Price > params.MaxPrice {
			continue
		}

		// Score filter
		if params.MinScore > 0 && token.TrustScore < params.MinScore {
			continue
		}
		if params.MaxScore > 0 && token.TrustScore > params.MaxScore {
			continue
		}

		// Change filter
		if params.MinChange != 0 && token.Change24h < params.MinChange {
			continue
		}
		if params.MaxChange != 0 && token.Change24h > params.MaxChange {
			continue
		}

		filtered = append(filtered, token)
	}

	// Sort by trust score (descending)
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].TrustScore > filtered[j].TrustScore
	})

	// Add rankings before slicing to maintain consistent rank across pages
	for i := range filtered {
		filtered[i].Rank = i + 1
	}

	totalFiltered := len(filtered)
	hasMore := false

	// Apply Offset
	if params.Offset >= totalFiltered {
		return []models.Token{}, 0, false
	}

	// Apply Limit
	start := params.Offset
	end := start + params.Limit
	if end > totalFiltered {
		end = totalFiltered
	} else if end < totalFiltered {
		hasMore = true
	}

	return filtered[start:end], totalFiltered, hasMore
}
