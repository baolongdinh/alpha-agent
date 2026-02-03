import { ref, computed } from 'vue'
import { api } from '../services/api'

// Shared state (singleton pattern)
const tokens = ref([])
const loading = ref(false)
const error = ref(null)
const lastFetch = ref(null)
const hasMore = ref(true)
const offset = ref(0)
const PAGE_SIZE = 50

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
  limit: PAGE_SIZE, // Changed from 200 to PAGE_SIZE
})

const marketStats = ref(null)

export function useTokens() {
  /**
   * Fetch tokens from API 
   */
  const fetchTokens = async (isLoadMore = false) => {
    if (!isLoadMore) {
      loading.value = true
      offset.value = 0
      hasMore.value = true
    }

    error.value = null

    try {
      console.log(`🔍 Fetching market data (isLoadMore: ${isLoadMore}, offset: ${offset.value})...`)

      const [tokensResp, statsResp] = await Promise.all([
        api.fetchMarketData(offset.value, filters.value),
        !isLoadMore ? api.fetchMarketStats() : Promise.resolve(marketStats.value)
      ])

      const newData = tokensResp.data || []
      if (isLoadMore) {
        // Append data
        const existingIds = new Set(tokens.value.map(t => t.symbol))
        const filteredNewData = newData.filter(t => !existingIds.has(t.symbol))
        tokens.value = [...tokens.value, ...filteredNewData]
      } else {
        tokens.value = newData
      }

      hasMore.value = tokensResp.hasMore
      marketStats.value = statsResp
      lastFetch.value = new Date()
      console.log(`✅ Total tokens in memory: ${tokens.value.length}, HasMore: ${hasMore.value}`)

    } catch (err) {
      error.value = err.message || 'Failed to fetch tokens'
      console.error('❌ Error fetching tokens:', error.value)
      if (!isLoadMore) tokens.value = []
    } finally {
      loading.value = false
    }
  }

  /**
   * Load More (Incremental)
   */
  const loadMore = async () => {
    if (loading.value || !hasMore.value) return

    offset.value += PAGE_SIZE
    await fetchTokens(true)
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
    fetchTokens() // Refetch when filters change
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
      limit: PAGE_SIZE,
    }
    offset.value = 0
    hasMore.value = true
  }

  /**
   * Refresh data
   */
  const refresh = () => {
    fetchTokens()
  }

  // Reactive Filter logic is now mostly handled by Backend
  const filteredTokens = computed(() => {
    // Sorting is still FE reactive for instant feel, but data subset is BE filtered
    let result = [...tokens.value]
    return result
  })

  // Computed properties for UI
  const hasTokens = computed(() => filteredTokens.value.length > 0)
  const isEmpty = computed(() => !loading.value && filteredTokens.value.length === 0)
  const topTokens = computed(() => filteredTokens.value.slice(0, 10))

  return {
    // State
    tokens: filteredTokens, // Expose filtered list as primary data source
    rawTokens: tokens, // Expose the raw data for internal use if needed
    loading,
    error,
    lastFetch,
    filters,
    marketStats,
    hasMore,

    // Computed
    hasTokens,
    isEmpty,
    topTokens,

    // Methods
    fetchTokens,
    loadMore,
    fetchTokenDetail,
    updateFilters,
    resetFilters,
    refresh,
  }
}
