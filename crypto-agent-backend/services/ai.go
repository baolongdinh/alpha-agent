package services

import (
	"context"
	"crypto-agent-backend/config"
	"crypto-agent-backend/models"
	"fmt"
	"log"

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

	// Log response details for debugging
	log.Printf("📊 Response candidates: %d", len(resp.Candidates))
	log.Printf("📊 Response parts: %d", len(resp.Candidates[0].Content.Parts))

	analysis := fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0])

	log.Printf("✅ AI analysis generated for %s (%d chars)", req.Symbol, len(analysis))
	log.Printf("📝 Full response: %s", analysis) // Log full response

	return analysis, nil
}

// buildAnalysisPrompt creates a detailed prompt for Gemini
func (s *AIService) buildAnalysisPrompt(req models.AnalysisRequest) string {
	return fmt.Sprintf(`Bạn là hệ thống AI phân tích của dự án AlphaAgent, một chuyên gia phân tích cryptocurrency chuyên sâu.

Hãy phân tích token sau và đưa ra nhận định ngắn gọn (tối đa 200 từ) về:
1. **Tiềm năng tăng trưởng**: Phân tích dựa trên các chỉ số
2. **Rủi ro cần lưu ý**: Những điểm đáng lo ngại
3. **Khuyến nghị**: Lời khuyên cho nhà đầu tư

**Thông tin token:**
- Tên: %s (%s)
- Giá hiện tại: $%.2f
- Market Cap: $%.2f
- Volume 24h: $%.2f
- TVL: $%.2f
- Trust Score: %.1f/100
- Biến động 24h: %.2f%%
- Biến động 7d: %.2f%%

**Yêu cầu:**
- Trả lời bằng tiếng Việt
- Sử dụng markdown formatting với headers, bold, bullet points
- Ngắn gọn, súc tích, dễ hiểu
- Tập trung vào insights thực tế`,
		req.Name,
		req.Symbol,
		req.Price,
		req.MarketCap,
		req.Volume24h,
		req.TVL,
		req.TrustScore,
		req.Change24h,
		req.Change7d,
	)
}

// Close closes the AI client
func (s *AIService) Close() {
	if s.client != nil {
		s.client.Close()
		log.Println("✅ Gemini AI client closed")
	}
}
