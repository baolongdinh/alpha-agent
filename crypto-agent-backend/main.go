package main

import (
	"crypto-agent-backend/cache"
	"crypto-agent-backend/config"
	"crypto-agent-backend/handlers"
	"crypto-agent-backend/services"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	log.Println("🚀 Starting AlphaAgent Backend...")
	log.Printf("📍 Server will run on port: %s", cfg.Port)
	log.Printf("🔧 Gin mode: %s", cfg.GinMode)

	// Set Gin mode
	gin.SetMode(cfg.GinMode)

	// Initialize services
	cacheManager := cache.NewCacheManager(cfg.TokenCacheDuration, cfg.AnalysisCacheDuration)
	log.Printf("✅ Cache manager initialized (Token: %v, Analysis: %v)", cfg.TokenCacheDuration, cfg.AnalysisCacheDuration)

	aggregator := services.NewAggregator(cfg)
	log.Println("✅ Data aggregator initialized")

	scorer := services.NewEnhancedScorer()
	log.Println("✅ Enhanced scorer initialized")

	// Initialize AI service (may fail if API key not set)
	aiService, err := services.NewAIService(cfg)
	if err != nil {
		log.Printf("⚠️  AI service initialization failed: %v", err)
		log.Println("⚠️  /api/analyze endpoint will not work without Gemini API key")
	}

	// Initialize handlers
	tokenHandler := handlers.NewTokenHandler(aggregator, scorer, cacheManager, cfg)
	log.Println("✅ Token handler initialized")

	var analyzeHandler *handlers.AnalyzeHandler
	if aiService != nil {
		analyzeHandler = handlers.NewAnalyzeHandler(aiService, cacheManager, cfg)
		log.Println("✅ Analyze handler initialized")

		// Ensure AI client is closed on shutdown
		defer aiService.Close()
	}

	// Create Gin router
	router := gin.Default()

	// Configure CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // In production, specify exact origins
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "crypto-agent-backend",
			"version": "1.0.0",
			"cache": gin.H{
				"tokens":   cacheManager.GetTokenCacheItemCount(),
				"analysis": cacheManager.GetAnalysisCacheItemCount(),
			},
		})
	})

	// API routes group
	api := router.Group("/api")
	{
		// Token endpoints
		api.GET("/tokens", tokenHandler.GetTokens)
		api.GET("/market/stats", tokenHandler.GetMarketStats)

		api.GET("/tokens/:id", tokenHandler.GetTokenByID)

		// Analysis endpoint (only if AI service is available)
		if analyzeHandler != nil {
			api.POST("/analyze", analyzeHandler.AnalyzeToken)
		} else {
			api.POST("/analyze", func(c *gin.Context) {
				c.JSON(503, gin.H{
					"status":  "error",
					"message": "AI service not available - GEMINI_API_KEY not configured",
				})
			})
		}
	}

	// Print available endpoints
	log.Println("✅ Server configured successfully")
	log.Println("📊 Available endpoints:")
	log.Println("   - GET  /health                 (Health check)")
	log.Println("   - GET  /api/tokens             (List tokens with filtering)")
	log.Println("   - POST /api/analyze            (AI token analysis)")
	log.Println("")
	log.Printf("🌐 Server starting on http://localhost:%s", cfg.Port)

	// Setup graceful shutdown
	go func() {
		if err := router.Run(":" + cfg.Port); err != nil {
			log.Fatalf("❌ Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server...")
	log.Println("✅ Server shutdown complete")
}
