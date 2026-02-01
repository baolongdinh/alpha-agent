import axios from 'axios'
// import mockData from './mockData.json' // Removed to prevent fake data fallback

// Configuration
const API_BASE_URL = 'http://localhost:8080/api'
const COINGECKO_API = 'https://api.coingecko.com/api/v3'
const DEXSCREENER_API = 'https://api.dexscreener.com/latest/dex/tokens'
const GOPLUS_API = 'https://api.goplusecurity.com/api/v1/token_security'

// Cache Store
const cache = {
  marketData: {
    data: null,
    timestamp: 0,
    duration: 10000 // 10s (backend has its own cache)
  },
  tokenDetails: new Map() // Key: id, Value: { data, timestamp }
}

export const api = {
  /**
   * Fetch Market Data (From Backend)
   */
  async fetchMarketData() {
    // 1. Check Cache
    const now = Date.now()
    if (cache.marketData.data && (now - cache.marketData.timestamp < cache.marketData.duration)) {
      console.log('⚡ Using cached market data')
      return { data: cache.marketData.data }
    }

    try {
      console.log('🌍 Fetching from Backend:', `${API_BASE_URL}/tokens`)
      const response = await axios.get(`${API_BASE_URL}/tokens`, {
         params: {
             limit: 1000 // Increased limit to fetch more data
         }
      })
      
      // Backend returns { status: "success", data: [...], ... }
      if (response.data && response.data.status === 'success') {
          const mappedData = response.data.data.map(transformBackendToken)
          
          // Update Cache
          cache.marketData = {
            data: mappedData,
            timestamp: now
          }
          return { data: mappedData }
      } else {
          throw new Error('Backend returned error status')
      }

    } catch (error) {
      console.error('❌ Backend fetch failed:', error)
      // Return empty or throw, do NOT use mock data
      throw error 
    }
  },

  /**
   * Fetch DexScreener Data
   */
  async fetchDexPairs(chainId, address) {
      if (!chainId || !address) return null
      // DexScreener uses different chain IDs e.g. 'ethereum', 'bsc', 'solana'
      // My getChainAndAddress helper returns '1', '56' etc.
      // Need mapping or just try. DexScreener endpoint: /latest/dex/tokens/{tokenAddresses}
      // It returns pairs across all chains for that address.
      try {
        const response = await axios.get(`${DEXSCREENER_API}/${address}`)
        return response.data.pairs || []
      } catch (e) {
        console.warn('DexScreener fetch failed', e)
        return null
      }
  },

  /**
   * Fetch Token Detail (CoinGecko + GoPlus + DexScreener)
   */
  /**
   * AI Analysis
   */
  async analyzeToken(payload) {
    try {
      console.log('🤖 Sending AI Analysis Request:', payload)
      // Call backend AI service
      const response = await axios.post(`${API_BASE_URL}/analyze`, payload)
      console.log('✅ AI Analysis Response:', response.data)
      return response.data
    } catch (e) {
      console.error('❌ AI Analysis failed:', e.response ? e.response.data : e.message)
      throw e
    }
  },

  /**
   * Fetch Token Detail (CoinGecko + GoPlus + DexScreener)
   */
  async fetchTokenDetail(id) {
    const now = Date.now()
    const cached = cache.tokenDetails.get(id)
    if (cached && (now - cached.timestamp < 30000)) {
        return cached.data
    }

    try {
        console.log(`🌍 Fetching Token Detail for ${id} from Backend...`)
        // The 'id' here is usually the symbol or name passed from transformBackendToken
        const response = await axios.get(`${API_BASE_URL}/tokens/${id}`)
        
        if (response.data) {
           const data = response.data
           // No need to transform again since backend returns the full Token model
           cache.tokenDetails.set(id, { data, timestamp: now })
           return data
        }
        return null
    } catch (e) {
        console.error('❌ Failed to fetch backend token detail', e)
        return null
    }
  },

  /**
   * Fetch Global Market Stats
   */
  async fetchMarketStats() {
    try {
      const response = await axios.get(`${API_BASE_URL}/market/stats`)
      return response.data
    } catch (e) {
      console.error('Failed to fetch market stats', e)
      return null
    }
  },

  // ...
}

function transformBackendToken(token) {
    // Map Backend model to Frontend UI expectations
    return {
        id: (token.id || token.name).toLowerCase().replace(/\s+/g, '-'),
        
        rank: token.rank,
        name: token.name,
        symbol: token.symbol,
        image: token.image,
        price: token.price,
        change_24h: token.change_24h,
        change_7d: token.change_7d,
        market_cap: token.market_cap,
        volume_24h: token.volume_24h,
        tvl: token.tvl,
        liquidity: token.liquidity,
        ath: token.ath,
        ath_change: token.ath_change,
        fdv: token.fdv,
        sparkline: token.sparkline, // Backend returns array
        trust_score: token.trust_score,
        score_breakdown: token.score_breakdown, // Pass through the breakdown!
        category: token.category
    }
}

/**
 * Compact Number Formatter (e.g., 1.25B, 450M)
 */
export const formatCompact = (val) => {
  if (!val) return '—'
  const formatter = Intl.NumberFormat('en-US', {
    notation: 'compact',
    maximumFractionDigits: 2
  })
  return formatter.format(val)
}

/**
 * Currency Formatter with Compact option
 */
export const formatCurrencyCompact = (val) => {
  if (!val && val !== 0) return '—'
  const formatter = Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: 'USD',
    notation: 'compact',
    maximumFractionDigits: 2
  })
  return formatter.format(val)
}

function getChainAndAddress(platforms) {
  if (!platforms) return { chainId: null, address: null }
  
  // Priority map: CoinGecko Platform -> GoPlus Chain ID
  // 1: Ethereum, 56: BSC, 137: Polygon, 42161: Arbitrum, 43114: Avalanche, 10: Optimism
  // Solana: GoPlus uses specialized API usually, but let's check basic EVM first or use Sol
  // For Solana (solana), GoPlus might have different endpoint or ID. 
  // Simplify to EVM for this demo phase + Solana if easy.
  
  if (platforms.ethereum) return { chainId: '1', address: platforms.ethereum }
  if (platforms['binance-smart-chain']) return { chainId: '56', address: platforms['binance-smart-chain'] }
  if (platforms['polygon-pos']) return { chainId: '137', address: platforms['polygon-pos'] }
  if (platforms.solana) return { chainId: 'solana', address: platforms.solana } // GoPlus might handle solana differently
  
  return { chainId: null, address: null }
}
