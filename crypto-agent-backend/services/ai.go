package services

import (
	"context"
	"crypto-agent-backend/config"
	"crypto-agent-backend/models"
	"fmt"
	"log"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// AIService handles Gemini AI integration
type AIService struct {
	client *genai.Client
	model  *genai.GenerativeModel
	config *config.Config
}

// NewAIService creates a new AI service
func NewAIService(cfg *config.Config) (*AIService, error) {
	if cfg.GeminiAPIKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY not configured")
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(cfg.GeminiAPIKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	// Configure model
	model := client.GenerativeModel("gemini-3-flash-preview")
	model.SetTemperature(0.7)
	model.SetTopP(0.9)
	model.SetMaxOutputTokens(5000)

	log.Println("✅ Gemini AI service initialized successfully")

	return &AIService{
		client: client,
		model:  model,
		config: cfg,
	}, nil
}

// AnalyzeToken generates AI analysis for a token
func (s *AIService) AnalyzeToken(ctx context.Context, req models.AnalysisRequest) (string, error) {
	prompt := s.buildAnalysisPrompt(req)

	log.Printf("🤖 Generating AI analysis for %s (%s)", req.Name, req.Symbol)

	resp, err := s.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("AI generation failed: %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("empty response from AI")
	}

	analysis := fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0])

	// Clean up markdown code blocks if present
	analysis = strings.TrimPrefix(analysis, "```json")
	analysis = strings.TrimPrefix(analysis, "```")
	analysis = strings.TrimSuffix(analysis, "```")
	analysis = strings.TrimSpace(analysis)

	log.Printf("✅ AI analysis generated for %s", req.Symbol)
	return analysis, nil
}

// buildAnalysisPrompt creates a structured prompt for Gemini to return JSON
func (s *AIService) buildAnalysisPrompt(req models.AnalysisRequest) string {
	return fmt.Sprintf(`You are AlphaAgent - An advanced Crypto Market Analysis AI. Your role is to act as a veteran Trader/Analyst to analyze the following token and provide a specific trading strategy.

Based on the provided market data, analyze and return the result strictly in JSON format (Do NOT allow introductory text):

{
  "summary": "Sharp overview of token status (under 30 words)",
  "growth_potential": {
    "score": (number 0-100),
    "reason": "Core reason for this score"
  },
  "technical_analysis": {
    "trend": "Main Trend (Uptrend/Downtrend/Accumulation/Distribution)",
    "strength": "Trend Strength (Very Strong/Strong/Weak/Neutral)",
    "key_levels": "Key Support and Nearest Resistance"
  },
  "risk_analysis": {
    "level": "Risk Level (Low/Medium/High/Extreme)",
    "concerns": ["Risk 1 (concise)", "Risk 2"]
  },
  "fundamental_analysis": {
    "sector": "Primary Sector/Field (e.g., Layer 1, DeFi, AI, Real World Assets)",
    "tokenomics": "Tokenomics Assessment (e.g., Deflationary, High Inflation, Fair Launch, VC Heavy)",
    "economic_moat": "Competitive Advantage/Moat (e.g., Network Effect, Tech Lead, Community)"
  },
  "recommendation": {
    "action": "ACTION (BUY NOW / BUY ZONE / HOLD / SELL / WATCH)",
    "entry_zone": "Optimal Entry Zone (specific)",
    "target": "Primary Price Target"
  },
  "trading_plan": {
    "buy_strategy": "Detailed Buy Strategy (e.g., DCA at zone A and B, or Breakout C)",
    "sell_targets": ["TP1: $Price (Soft Target)", "TP2: $Price", "TP3: $Price (Moonbag)"],
    "stop_loss": "Stop Loss Price (or Invalidation Condition)",
    "time_horizon": "Time Horizon (Short-term/Mid-term/Long-term)"
  },
  "insights": [
    "Insight 1: Analysis of Liquidity/Volume vs Mcap (Velocity)",
    "Insight 2: What 30d/90d price action says about money flow",
    "Insight 3: Correlation with general market (Beta)"
  ]
}

**INPUT DATA:**
- Token: %s (%s) | Rank: #%d
- Current Price: $%.6f
- 24h Change: %.2f%% | 7d Change: %.2f%%
- Mid-term Trend: 30d: %.2f%% | 90d: %.2f%%
- Market Cap: $%.2f | Fully Diluted Valuation (FDV): $%.2f
- Supply: Circulating %.1f%% / Max Supply
- Volume 24h: $%.2f (Vol/Mcap Ratio: %.4f)
- Liquidity: $%.2f | TVL: $%.2f
- Alpha Trust Score: %.1f/100

**IMPORTANT NOTES:**
1. If Liquidity/Mcap is low (<1%%), warn about high liquidity risk.
2. If FDV >> Mcap, warn about token inflation/unlocks.
3. Price targets (TP/SL) must be based on price volatility (Change 7d/30d) and current price, estimate support/resistance reasonably.
4. Respond entirely in professional Crypto English.`,
		req.Name, req.Symbol, req.Rank,
		req.Price,
		req.Change24h, req.Change7d,
		req.Change30d, req.Change90d,
		req.MarketCap, req.Price*req.TotalSupply,
		(req.CirculatingSupply/req.MaxSupply)*100,
		req.Volume24h, req.Volume24h/req.MarketCap,
		req.Liquidity, req.TVL,
		req.TrustScore,
	)
}

// Close closes the AI client
func (s *AIService) Close() {
	if s.client != nil {
		s.client.Close()
		log.Println("✅ Gemini AI client closed")
	}
}
