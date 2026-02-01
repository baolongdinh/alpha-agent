package models

// Token represents aggregated cryptocurrency data from multiple sources
type Token struct {
	// Basic Info
	ID       string `json:"id"`
	Rank     int    `json:"rank"`
	Symbol   string `json:"symbol"`
	Name     string `json:"name"`
	Image    string `json:"image"`
	Category string `json:"category"`

	// Current Metrics
	Price     float64 `json:"price"`
	MarketCap float64 `json:"market_cap"`
	Volume24h float64 `json:"volume_24h"`
	TVL       float64 `json:"tvl"`
	Liquidity float64 `json:"liquidity"`

	// Price Changes
	Change1h        float64 `json:"change_1h"`
	Change24h       float64 `json:"change_24h"`
	Change7d        float64 `json:"change_7d"`
	Change30d       float64 `json:"change_30d"`
	Change90d       float64 `json:"change_90d"`
	VolumeChange24h float64 `json:"volume_change_24h"`
	MarketCapDom    float64 `json:"market_cap_dominance"`

	// ATH & FDV Metrics
	Ath               float64 `json:"ath"`
	AthChange         float64 `json:"ath_change"`
	FullyDilutedValue float64 `json:"fdv"`
	CirculatingSupply float64 `json:"circulating_supply"`
	TotalSupply       float64 `json:"total_supply"`
	MaxSupply         float64 `json:"max_supply"`

	// Social Metrics (NEW)
	SocialScore  float64 `json:"social_score,omitempty"`
	SocialVolume int     `json:"social_volume,omitempty"`
	Sentiment    float64 `json:"sentiment,omitempty"` // -1 to 1

	// On-Chain Metrics (NEW)
	ActiveAddresses              int     `json:"active_addresses,omitempty"`
	TransactionCount             int     `json:"transaction_count,omitempty"`
	WhaleConcentrationNormalized float64 `json:"whale_concentration_normalized,omitempty"`

	// Fundamental Metrics (NEW)
	Revenue30d float64 `json:"revenue_30d,omitempty"`
	PSRatio    float64 `json:"ps_ratio,omitempty"`
	PERatio    float64 `json:"pe_ratio,omitempty"`

	// Sparkline Data
	Sparkline []float64 `json:"sparkline"`

	// Advanced Metrics
	Volume7dAvg  float64 `json:"volume_7d_avg"`
	Volume30dAvg float64 `json:"volume_30d_avg"`
	TVL7dChange  float64 `json:"tvl_7d_change"`
	TVL30dChange float64 `json:"tvl_30d_change"`

	// Technical Indicators
	RSI14 float64 `json:"rsi_14"`
	SMA7  float64 `json:"sma_7"`
	SMA30 float64 `json:"sma_30"`
	SMA90 float64 `json:"sma_90"`
	EMA12 float64 `json:"ema_12"`
	EMA26 float64 `json:"ema_26"`
	MACD  float64 `json:"macd"`

	// Volatility Metrics
	Volatility7d  float64 `json:"volatility_7d"`
	Volatility30d float64 `json:"volatility_30d"`
	BetaToMarket  float64 `json:"beta_to_market"`

	// Holder & Distribution
	HolderCount        int     `json:"holder_count"`
	Top10HoldersRatio  float64 `json:"top10_holders_ratio"`
	WhaleConcentration float64 `json:"whale_concentration"`

	// Liquidity Depth
	BidAskSpread    float64 `json:"bid_ask_spread"`
	OrderbookDepth1 float64 `json:"orderbook_depth_1"` // 1% depth
	OrderbookDepth5 float64 `json:"orderbook_depth_5"` // 5% depth

	// Smart Contract Info
	IsVerified  bool   `json:"is_verified"`
	AuditStatus string `json:"audit_status"`
	ContractAge int    `json:"contract_age_days"`

	// Historical Data (for calculations)
	PriceHistory  PriceHistory  `json:"price_history,omitempty"`
	VolumeHistory VolumeHistory `json:"volume_history,omitempty"`
	TVLHistory    TVLHistory    `json:"tvl_history,omitempty"`

	// Scoring
	TrustScore     float64                `json:"trust_score"`
	ScoreBreakdown DetailedScoreBreakdown `json:"score_breakdown"`
}

// External API response structures

// DefiLlamaProtocol represents data from DeFiLlama API
type DefiLlamaProtocol struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Symbol   string  `json:"symbol"`
	TVL      float64 `json:"tvl"`
	Category string  `json:"category"`
	Change1h float64 `json:"change_1h,omitempty"`
	Change1d float64 `json:"change_1d,omitempty"`
	Change7d float64 `json:"change_7d,omitempty"`
}

// CoinGeckoMarket represents data from CoinGecko API
type CoinGeckoMarket struct {
	ID                       string  `json:"id"`
	Symbol                   string  `json:"symbol"`
	Name                     string  `json:"name"`
	Image                    string  `json:"image"`
	CurrentPrice             float64 `json:"current_price"`
	MarketCap                float64 `json:"market_cap"`
	MarketCapRank            int     `json:"market_cap_rank"`
	FullyDilutedValuation    float64 `json:"fully_diluted_valuation"`
	TotalVolume              float64 `json:"total_volume"`
	Ath                      float64 `json:"ath"`
	AthChangePercentage      float64 `json:"ath_change_percentage"`
	PriceChangePercentage1h  float64 `json:"price_change_percentage_1h_in_currency"`
	PriceChangePercentage24h float64 `json:"price_change_percentage_24h"`
	PriceChangePercentage7d  float64 `json:"price_change_percentage_7d_in_currency"`
	SparklineIn7d            struct {
		Price []float64 `json:"price"`
	} `json:"sparkline_in_7d"`
	CirculatingSupply float64 `json:"circulating_supply"`
	TotalSupply       float64 `json:"total_supply"`
	MaxSupply         float64 `json:"max_supply"`
}

// DexScreenerResponse represents data from DexScreener API
type DexScreenerResponse struct {
	Pairs []DexPair `json:"pairs"`
}

type DexPair struct {
	BaseToken struct {
		Symbol string `json:"symbol"`
	} `json:"baseToken"`
	Liquidity struct {
		USD float64 `json:"usd"`
	} `json:"liquidity"`
	Volume struct {
		H24 float64 `json:"h24"`
	} `json:"volume"`
	PriceUsd string `json:"priceUsd"`
}

// CoinMarketCap models (NEW)
type CoinMarketCapResponse struct {
	Data []CoinMarketCapCoin `json:"data"`
}

type CoinMarketCapCoin struct {
	ID                int      `json:"id"`
	Symbol            string   `json:"symbol"`
	Name              string   `json:"name"`
	Slug              string   `json:"slug"`
	CMCRank           int      `json:"cmc_rank"`
	Tags              []string `json:"tags"`
	CirculatingSupply float64  `json:"circulating_supply"`
	TotalSupply       float64  `json:"total_supply"`
	MaxSupply         float64  `json:"max_supply"`
	Quote             struct {
		USD struct {
			Price                 float64   `json:"price"`
			Volume24h             float64   `json:"volume_24h"`
			VolumeChange24h       float64   `json:"volume_change_24h"`
			PercentChange1h       float64   `json:"percent_change_1h"`
			PercentChange24h      float64   `json:"percent_change_24h"`
			PercentChange7d       float64   `json:"percent_change_7d"`
			PercentChange30d      float64   `json:"percent_change_30d"` // Added
			PercentChange90d      float64   `json:"percent_change_90d"` // Added
			MarketCap             float64   `json:"market_cap"`
			MarketCapDominance    float64   `json:"market_cap_dominance"`
			FullyDilutedMarketCap float64   `json:"fully_diluted_market_cap"`
			Sparkline             []float64 `json:"sparkline"` // NEW: CMC Sparkline data
		} `json:"USD"`
	} `json:"quote"`
}

// Messari models (NEW)
type MessariResponse struct {
	Data []MessariAsset `json:"data"`
}

type MessariAsset struct {
	Symbol     string `json:"symbol"`
	Name       string `json:"name"`
	MarketData struct {
		PriceUSD  float64 `json:"price_usd"`
		Volume24h float64 `json:"volume_last_24_hours"`
		MarketCap float64 `json:"market_cap_usd"`
	} `json:"market_data"`
}

// FilterParams represents query parameters for filtering tokens
type FilterParams struct {
	MinMcap  float64
	MaxMcap  float64
	Category string
	Limit    int
}

// MarketStats represents global market statistics
type MarketStats struct {
	TotalMarketCap      float64            `json:"total_market_cap"`
	TotalVolume         float64            `json:"total_volume"`
	MarketCapPercentage map[string]float64 `json:"market_cap_percentage"`
	MarketCapChange24h  float64            `json:"market_cap_change_24h"`
	BTCDominance        float64            `json:"btc_dominance"`
	ETHDominance        float64            `json:"eth_dominance"`
}

// CoinGeckoGlobalResponse represents the global data from CoinGecko
type CoinGeckoGlobalResponse struct {
	Data struct {
		TotalMarketCap         map[string]float64 `json:"total_market_cap"`
		TotalVolume            map[string]float64 `json:"total_volume"`
		MarketCapPercentage    map[string]float64 `json:"market_cap_percentage"`
		MarketCapChangePercent float64            `json:"market_cap_change_percentage_24h_usd"`
	} `json:"data"`
}
