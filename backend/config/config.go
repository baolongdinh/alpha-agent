package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	// Server settings
	Port    string
	GinMode string

	// API Keys
	GeminiAPIKey string

	// Cache settings
	TokenCacheDuration    time.Duration
	AnalysisCacheDuration time.Duration

	// External API URLs
	DefiLlamaAPIURL   string
	CoinGeckoAPIURL   string
	DexScreenerAPIURL string

	// New Data Sources
	CoinMarketCapAPIURL string
	CoinMarketCapAPIKey string
	MessariAPIURL       string
	MessariAPIKey       string
	TokenTerminalAPIURL string
	TokenTerminalAPIKey string
	LunarCrushAPIURL    string
	LunarCrushAPIKey    string
	GlassnodeAPIURL     string
	GlassnodeAPIKey     string
	CryptoCompareAPIURL string
}

var AppConfig *Config

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	// Load .env file if exists (ignore error in production)
	_ = godotenv.Load()

	config := &Config{
		// Server defaults
		Port:    getEnv("PORT", "8080"),
		GinMode: getEnv("GIN_MODE", "release"),

		// API Keys
		GeminiAPIKey: getEnv("GEMINI_API_KEY", ""),

		// Cache durations
		TokenCacheDuration:    parseDuration(getEnv("TOKEN_CACHE_DURATION", "5m"), 5*time.Minute),
		AnalysisCacheDuration: parseDuration(getEnv("ANALYSIS_CACHE_DURATION", "60m"), 60*time.Minute),

		// External APIs
		DefiLlamaAPIURL:   getEnv("DEFILLAMA_API_URL", "https://api.llama.fi"),
		CoinGeckoAPIURL:   getEnv("COINGECKO_API_URL", "https://api.coingecko.com/api/v3"),
		DexScreenerAPIURL: getEnv("DEXSCREENER_API_URL", "https://api.dexscreener.com/latest"),

		// Enhanced Data Sources
		CoinMarketCapAPIURL: getEnv("CMC_API_URL", "https://pro-api.coinmarketcap.com"),
		CoinMarketCapAPIKey: getEnv("CMC_API_KEY", ""),
		MessariAPIURL:       getEnv("MESSARI_API_URL", "https://data.messari.io/api"),
		MessariAPIKey:       getEnv("MESSARI_API_KEY", ""),
		TokenTerminalAPIURL: getEnv("TOKENTERMINAL_API_URL", "https://api.tokenterminal.com"),
		TokenTerminalAPIKey: getEnv("TOKENTERMINAL_API_KEY", ""),
		LunarCrushAPIURL:    getEnv("LUNARCRUSH_API_URL", "https://api.lunarcrush.com"),
		LunarCrushAPIKey:    getEnv("LUNARCRUSH_API_KEY", ""),
		GlassnodeAPIURL:     getEnv("GLASSNODE_API_URL", "https://api.glassnode.com"),
		GlassnodeAPIKey:     getEnv("GLASSNODE_API_KEY", ""),
		CryptoCompareAPIURL: getEnv("CRYPTOCOMPARE_API_URL", "https://min-api.cryptocompare.com"),
	}

	// Validate required fields
	if config.GeminiAPIKey == "" {
		log.Println("WARNING: GEMINI_API_KEY not set - AI analysis will not work")
	}

	AppConfig = config
	return config
}

// Helper functions
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func parseDuration(value string, defaultDuration time.Duration) time.Duration {
	duration, err := time.ParseDuration(value)
	if err != nil {
		return defaultDuration
	}
	return duration
}
