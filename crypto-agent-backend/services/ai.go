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
	return fmt.Sprintf(`Bạn là hệ thống AI phân tích AlphaAgent, chuyên gia phân tích on-chain và thị trường cryptocurrency.

Hãy phân tích token sau dựa trên dữ liệu hiện có và trả về một đối tượng JSON duy nhất theo cấu trúc sau:
{
  "summary": "Tóm tắt ngắn gọn trạng thái hiện tại (1 câu)",
  "growth_potential": {
    "score": (số từ 0-100),
    "reason": "Lý do cho điểm số này"
  },
  "technical_analysis": {
    "trend": "Trạng thái xu hướng (Bullish/Neutral/Bearish)",
    "strength": "Sức mạnh xu hướng (Strong/Weak)"
  },
  "risk_analysis": {
    "level": "Mức độ rủi ro (Low/Medium/High)",
    "concerns": ["Điểm đáng lo ngại 1", "Điểm đáng lo ngại 2"]
  },
  "recommendation": {
    "action": "Hành động khuyến nghị (Buy/Hold/Sell/Watch)",
    "entry_zone": "Vùng giá entry đề xuất",
    "target": "Mục tiêu kỳ vọng"
  },
  "insights": ["Insight đặc biệt về holder/liquidity", "Insight về biến động thị trường"]
}

**Thông tin token:**
- Tên: %s (%s) | Rank: #%d
- Giá: $%.2f | Biến động: 24h: %.2f%%, 7d: %.2f%%
- Market Cap: $%.2f | FDV: (Tỷ lệ Circulating/Max supply: %.1f%% / %.1f%%)
- Volume 24h: $%.2f | Liquidity: $%.2f
- TVL: $%.2f | Holder Count: %d
- Alpha Score: %.1f/100

**Yêu cầu:**
- Trả lời bằng tiếng Việt.
- CHỈ TRẢ VỀ JSON, không thêm văn bản giải thích nào khác.
- Đảm bảo JSON hợp lệ.`,
		req.Name, req.Symbol, req.Rank,
		req.Price, req.Change24h, req.Change7d,
		req.MarketCap, req.CirculatingSupply, req.MaxSupply,
		req.Volume24h, req.Liquidity,
		req.TVL, req.HolderCount,
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
