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
  category: '',
  limit: 1000,
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
      console.log('🔍 Fetching market data...')
      
      const [tokensResp, statsResp] = await Promise.all([
        api.fetchMarketData(),
        api.fetchMarketStats()
      ])

      let fetchedTokens = tokensResp.data || []
      marketStats.value = statsResp
      
      // 1. Filter by Category
      if (filters.value.category && filters.value.category !== 'All Sectors') {
         fetchedTokens = fetchedTokens.filter(t => t.category === filters.value.category)
      }

      // 2. Filter by Search
      if (filters.value.search && filters.value.search.trim()) {
        const searchTerm = filters.value.search.toLowerCase().trim()
        fetchedTokens = fetchedTokens.filter(token => 
          token.symbol.toLowerCase().includes(searchTerm) ||
          token.name.toLowerCase().includes(searchTerm)
        )
      }
      
      // 3. Limit - Removed frontend slice to respect backend limit fully
      // if (filters.value.limit) {
      //    fetchedTokens = fetchedTokens.slice(0, filters.value.limit)
      // }
      
      tokens.value = fetchedTokens
      lastFetch.value = new Date()
      console.log(`✅ Fetched ${tokens.value.length} tokens`)

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
       // Response is already the data object from api.js
       return response
    } catch (e) {
       console.error("Error fetching detail:", e)
       return null
    }
  }

  /**
   * Update filters and refetch
   */
  const updateFilters = (newFilters) => {
    filters.value = { ...filters.value, ...newFilters }
    fetchTokens()
  }

  /**
   * Reset filters to default
   */
  const resetFilters = () => {
    filters.value = {
      search: '',
      minMcap: null,
      maxMcap: null,
      category: '',
      limit: 50,
    }
    fetchTokens()
  }

  /**
   * Refresh data (force refetch)
   */
  const refresh = () => {
    fetchTokens()
  }

  // Computed properties
  const hasTokens = computed(() => tokens.value.length > 0)
  const isEmpty = computed(() => !loading.value && tokens.value.length === 0)
  const topTokens = computed(() => tokens.value.slice(0, 10))

  return {
    // State
    tokens,
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
