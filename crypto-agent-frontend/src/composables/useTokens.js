import { ref, computed } from 'vue'
import { api } from '../services/api'

// Shared state (singleton pattern)
const tokens = ref([])
const loading = ref(false)
const error = ref(null)
const lastFetch = ref(null)

// Filter state
const filters = ref({
  search: '',
  minMcap: null,
  maxMcap: null,
  minScore: null,
  maxScore: null,
  minPrice: null,
  maxPrice: null,
  minChange: null,
  maxChange: null,
  category: '',
  limit: 5000,
})

const marketStats = ref(null)

export function useTokens() {
  /**
   * Fetch tokens from API 
   */
  const fetchTokens = async () => {
    loading.value = true
    error.value = null

    try {
      console.log('🔍 Fetching base market data...')
      
      const [tokensResp, statsResp] = await Promise.all([
        api.fetchMarketData(), // Fetch full set, let FE handle the filtering for speed
        api.fetchMarketStats()
      ])

      tokens.value = tokensResp.data || []
      marketStats.value = statsResp
      lastFetch.value = new Date()
      console.log(`✅ Loaded ${tokens.value.length} tokens into memory`)

    } catch (err) {
      error.value = err.message || 'Failed to fetch tokens'
      console.error('❌ Error fetching tokens:', error.value)
      tokens.value = []
    } finally {
      loading.value = false
    }
  }

  /**
   * Fetch detailed token data
   */
  const fetchTokenDetail = async (id) => {
    try {
       const response = await api.fetchTokenDetail(id)
       return response
    } catch (e) {
       console.error("Error fetching detail:", e)
       return null
    }
  }

  /**
   * Update filters (Frontend Reactive)
   */
  const updateFilters = (newFilters) => {
    filters.value = { ...filters.value, ...newFilters }
  }

  /**
   * Reset filters to default
   */
  const resetFilters = () => {
    filters.value = {
      search: '',
      minMcap: null,
      maxMcap: null,
      minScore: null,
      maxScore: null,
      minPrice: null,
      maxPrice: null,
      minChange: null,
      maxChange: null,
      category: '',
      limit: 5000,
    }
  }

  /**
   * Refresh data
   */
  const refresh = () => {
    fetchTokens()
  }

  // Reactive Filter logic
  const filteredTokens = computed(() => {
    let result = [...tokens.value]
    const f = filters.value

    // 1. Search (Symbol or Name)
    if (f.search && f.search.trim()) {
      const term = f.search.toLowerCase().trim()
      result = result.filter(t => 
        t.symbol.toLowerCase().includes(term) || 
        t.name.toLowerCase().includes(term)
      )
    }

    // 2. Category / Sector
    if (f.category && f.category !== 'All Sectors' && f.category !== '') {
      result = result.filter(t => t.category === f.category)
    }

    // 3. Market Cap Range
    if (f.minMcap != null && f.minMcap !== '') {
      result = result.filter(t => t.market_cap >= Number(f.minMcap))
    }
    if (f.maxMcap != null && f.maxMcap !== '') {
      result = result.filter(t => t.market_cap <= Number(f.maxMcap))
    }

    // 4. Score Range
    if (f.minScore != null && f.minScore !== '') {
      result = result.filter(t => t.trust_score >= Number(f.minScore))
    }
    if (f.maxScore != null && f.maxScore !== '') {
      result = result.filter(t => t.trust_score <= Number(f.maxScore))
    }

    // 5. Price Range
    if (f.minPrice != null && f.minPrice !== '') {
      result = result.filter(t => t.price >= Number(f.minPrice))
    }
    if (f.maxPrice != null && f.maxPrice !== '') {
      result = result.filter(t => t.price <= Number(f.maxPrice))
    }

    // 6. 24h Change Range
    if (f.minChange != null && f.minChange !== '') {
      result = result.filter(t => t.change_24h >= Number(f.minChange))
    }
    if (f.maxChange != null && f.maxChange !== '') {
      result = result.filter(t => t.change_24h <= Number(f.maxChange))
    }

    return result
  })

  // Computed properties for UI
  const hasTokens = computed(() => filteredTokens.value.length > 0)
  const isEmpty = computed(() => !loading.value && filteredTokens.value.length === 0)
  const topTokens = computed(() => filteredTokens.value.slice(0, 10))

  return {
    // State
    tokens: filteredTokens, // Expose filtered list as primary data source
    rawTokens: tokens,
    loading,
    error,
    lastFetch,
    filters,
    marketStats,
    
    // Computed
    hasTokens,
    isEmpty,
    topTokens,
    
    // Methods
    fetchTokens,
    fetchTokenDetail,
    updateFilters,
    resetFilters,
    refresh,
  }
}
