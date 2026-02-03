# 🦅 AlphaAgent - Intelligent Crypto Market Tracker

AlphaAgent is an intelligent cryptocurrency market tracking and analysis platform, designed to provide deep insights (Alpha) through multi-source data aggregation and advanced AI analysis.

## ✨ Key Features

- **🔍 Multi-Source Aggregation**: Automatically aggregates data from CoinMarketCap, DexScreener, and DeFiLlama to provide the most comprehensive view of any token.
- **🤖 AI Market Analyst**: Integrates Google Gemini AI for in-depth token analysis, providing:
  - **Strategic Execution Plan**: Specific Buy/Sell plans (Entry, Take Profit, Stop Loss).
  - **Risk Assessment**: Multi-dimensional risk evaluation.
  - **Sentiment Shift**: Detection of community sentiment trends.
- **⚡ Real-time Metrics**: Real-time updates for price, volume, and market volatility.
- **📊 Advanced Visualization**: Professional Dashboard (Dark Mode) featuring Radar charts, Gauges, and Sparklines.
- **🛡️ Trust Score Engine**: A proprietary scoring system evaluating token reliability based on liquidity, holder distribution, and volatility.

## 🛠️ Tech Stack

### Backend (Golang)

- **Framework**: Gin Gonic
- **AI Integration**: Google Generative AI Cloud (Gemini)
- **Data Handling**: Concurrent fetching & merging (Goroutines)
- **Caching**: In-memory caching for optimized API performance.

### Frontend (Vue 3)

- **Framework**: Vue 3 (Composition API) + Vite
- **Styling**: TailwindCSS + Custom Glassmorphism UI
- **State Management**: Vue Reactivity
- **Charting**: Custom SVG Components (Radar, Gauge)

## 🚀 Installation & Setup

### Prerequisites

- Go 1.21+
- Node.js 18+
- API Keys: CoinMarketCap, Google Gemini (configure in `.env`)

### 1. Start Backend

```bash
cd backend

# Install dependencies
go mod download

# Configure Environment Variables
# Create a .env file and add:
# COINMARKETCAP_API_KEY=your_key
# GEMINI_API_KEY=your_key

# Run the server
go run .
```

### 2. Start Frontend

```bash
cd frontend

# Install dependencies
npm install

# Run dev server
npm run dev
```

[DEMO]
<img width="1911" height="910" alt="image" src="https://github.com/user-attachments/assets/a024a2e8-50a2-483e-903e-70a791d22f92" />

<img width="1920" height="2804" alt="image" src="https://github.com/user-attachments/assets/5abdaff2-a186-4a92-b13e-7b742b636ffa" />
