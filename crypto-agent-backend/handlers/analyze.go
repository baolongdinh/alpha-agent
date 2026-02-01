package handlers

import (
	"crypto-agent-backend/cache"
	"crypto-agent-backend/config"
	"crypto-agent-backend/models"
	"crypto-agent-backend/services"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// AnalyzeHandler handles AI analysis endpoints
type AnalyzeHandler struct {
	aiService *services.AIService
	cache     *cache.CacheManager
	config    *config.Config
}

// NewAnalyzeHandler creates a new analyze handler
func NewAnalyzeHandler(
	aiService *services.AIService,
	cacheManager *cache.CacheManager,
	cfg *config.Config,
) *AnalyzeHandler {
	return &AnalyzeHandler{
		aiService: aiService,
		cache:     cacheManager,
		config:    cfg,
	}
}

// AnalyzeToken handles POST /api/analyze
func (h *AnalyzeHandler) AnalyzeToken(c *gin.Context) {
	// Parse request body
	var req models.AnalysisRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status:    "error",
			Message:   "Invalid request body: " + err.Error(),
			Timestamp: time.Now(),
		})
		return
	}

	// Validate required fields
	if req.Symbol == "" || req.Name == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status:    "error",
			Message:   "Symbol and Name are required",
			Timestamp: time.Now(),
		})
		return
	}

	log.Printf("📊 Analysis request for %s (%s)", req.Name, req.Symbol)

	// Check cache first
	cacheKey := fmt.Sprintf("analysis_%s", req.Symbol)
	if cached, found := h.cache.GetAnalysis(cacheKey); found {
		log.Printf("✓ Cache hit for analysis: %s", req.Symbol)

		analysis := cached.(string)
		c.JSON(http.StatusOK, models.AnalysisResponse{
			Status:      "success",
			Cached:      true,
			Analysis:    analysis,
			GeneratedAt: time.Now(),
		})
		return
	}

	log.Printf("✗ Cache miss for analysis: %s - generating fresh analysis", req.Symbol)

	// Generate AI analysis
	analysis, err := h.aiService.AnalyzeToken(c.Request.Context(), req)
	if err != nil {
		log.Printf("❌ AI analysis failed for %s: %v", req.Symbol, err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Status:    "error",
			Message:   "Failed to generate analysis: " + err.Error(),
			Timestamp: time.Now(),
		})
		return
	}

	// Cache the result
	h.cache.SetAnalysis(cacheKey, analysis)
	log.Printf("✓ Cached analysis for: %s", req.Symbol)

	c.JSON(http.StatusOK, models.AnalysisResponse{
		Status:      "success",
		Cached:      false,
		Analysis:    analysis,
		GeneratedAt: time.Now(),
	})
}
