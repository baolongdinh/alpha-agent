package cache

import (
	"time"

	gocache "github.com/patrickmn/go-cache"
)

// CacheManager manages in-memory caches for different data types
type CacheManager struct {
	tokenCache    *gocache.Cache
	analysisCache *gocache.Cache
}

// NewCacheManager creates a new cache manager with configured durations
func NewCacheManager(tokenCacheDuration, analysisCacheDuration time.Duration) *CacheManager {
	return &CacheManager{
		// Token data cache: shorter duration for fresh market data
		tokenCache: gocache.New(tokenCacheDuration, tokenCacheDuration*2),

		// AI analysis cache: longer duration to reduce API calls
		analysisCache: gocache.New(analysisCacheDuration, analysisCacheDuration*2),
	}
}

// Token Cache Methods

// GetTokens retrieves cached token data
func (c *CacheManager) GetTokens(key string) (interface{}, bool) {
	return c.tokenCache.Get(key)
}

// SetTokens stores token data in cache
func (c *CacheManager) SetTokens(key string, data interface{}) {
	c.tokenCache.Set(key, data, gocache.DefaultExpiration)
}

// InvalidateTokens clears the token cache
func (c *CacheManager) InvalidateTokens() {
	c.tokenCache.Flush()
}

// Analysis Cache Methods

// GetAnalysis retrieves cached AI analysis
func (c *CacheManager) GetAnalysis(key string) (interface{}, bool) {
	return c.analysisCache.Get(key)
}

// SetAnalysis stores AI analysis in cache
func (c *CacheManager) SetAnalysis(key string, data interface{}) {
	c.analysisCache.Set(key, data, gocache.DefaultExpiration)
}

// InvalidateAnalysis clears the analysis cache
func (c *CacheManager) InvalidateAnalysis() {
	c.analysisCache.Flush()
}

// GetTokenCacheItemCount returns the number of items in token cache
func (c *CacheManager) GetTokenCacheItemCount() int {
	return c.tokenCache.ItemCount()
}

// GetAnalysisCacheItemCount returns the number of items in analysis cache
func (c *CacheManager) GetAnalysisCacheItemCount() int {
	return c.analysisCache.ItemCount()
}

// FlushAll clears all caches
func (c *CacheManager) FlushAll() {
	c.tokenCache.Flush()
	c.analysisCache.Flush()
}
