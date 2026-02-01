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
	return fmt.Sprintf(`Bạn là AlphaAgent - Hệ thống AI phân tích thị trường Crypto chuyên sâu. Nhiệm vụ của bạn là đóng vai một chuyên gia giao dịch (Trader/Analyst) kỳ cựu để phân tích token sau và đưa ra chiến lược giao dịch cụ thể.

Dựa trên dữ liệu thị trường được cung cấp, hãy phân tích và trả về kết quả dưới dạng JSON (Tuyệt đối không kèm text dẫn nhập):

{
  "summary": "Nhận định tổng quan sắc bén về trạng thái token (dưới 30 từ)",
  "growth_potential": {
    "score": (số 0-100),
    "reason": "Lý do cốt lõi cho điểm số này"
  },
  "technical_analysis": {
    "trend": "Xu hướng chính (Uptrend/Downtrend/Accumulation/Distribution)",
    "strength": "Độ mạnh xu hướng (Very Strong/Strong/Weak/Neutral)",
    "key_levels": "Hỗ trợ quan trọng và Kháng cự gần nhất"
  },
  "risk_analysis": {
    "level": "Rủi ro (Low/Medium/High/Extreme)",
    "concerns": ["Rủi ro 1 (ngắn gọn)", "Rủi ro 2"]
  },
  "recommendation": {
    "action": "ACTION (MUA NGAY / CANH MUA / HOLD / BÁN / QUAN SÁT)",
    "entry_zone": "Vùng giá mua tối ưu (cụ thể)",
    "target": "Mục tiêu giá chính"
  },
  "trading_plan": {
    "buy_strategy": "Chiến lược mua chi tiết (VD: DCA tại vùng A và B, hoặc Breakout C)",
    "sell_targets": ["TP1: $Giá (Mô tả nhẹ)", "TP2: $Giá", "TP3: $Giá (Moonbag)"],
    "stop_loss": "Giá cắt lỗ (hoặc điều kiện invalid)",
    "time_horizon": "Khung thời gian (Ngắn hạn/Trung hạn/Dài hạn)"
  },
  "insights": [
    "Insight 1: Phân tích về Liquidity/Volume so với Mcap (Velocity)",
    "Insight 2: Biến động giá 30d/90d nói lên điều gì về dòng tiền",
    "Insight 3: So sánh tương quan với thị trường chung (Beta)"
  ]
}

**DỮ LIỆU ĐẦU VÀO:**
- Token: %s (%s) | Rank: #%d
- Giá hiện tại: $%.6f
- Biến động: 24h: %.2f%% | 7d: %.2f%%
- Xu hướng trung hạn: 30d: %.2f%% | 90d: %.2f%%
- Vốn hóa (Mcap): $%.2f | Định giá pha loãng (FDV): $%.2f
- Cung: Circulating %.1f%% / Max Supply
- Volume 24h: $%.2f (Tỷ lệ Vol/Mcap: %.4f)
- Liquidity: $%.2f | TVL: $%.2f
- Alpha Trust Score: %.1f/100

**LƯU Ý QUAN TRỌNG:**
1. Nếu Liquidity/Mcap thấp (<1%%), cảnh báo rủi ro thanh khoản.
2. Nếu FDV >> Mcap, cảnh báo lạm phát token.
3. Đưa ra các mốc giá (TP/SL) phải dựa trên biến động giá (Change 7d/30d) và mức giá hiện tại, hãy ước lượng hỗ trợ/kháng cự một cách hợp lý.
4. Trả lời hoàn toàn bằng tiếng Việt chuyên ngành Crypto.`,
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
