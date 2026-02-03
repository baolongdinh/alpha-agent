<template>
  <div class="space-y-6">
    <!-- Filter Bar -->
    <FilterBar
      :filters="filters"
      :loading="loading"
      :last-fetch="lastFetch"
      :show-favorites="showFavorites"
      @update:filters="updateFilters"
      @refresh="refresh"
      @reset="resetFilters"
      @toggle:favorites="toggleFavoritesFilter"
    />

    <!-- Error State -->
    <div v-if="error" class="glass rounded-xl p-6 border-2 border-red-500/30">
      <div class="flex items-center gap-3 text-red-400">
        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <span class="font-medium">{{ error }}</span>
      </div>
    </div>

    <!-- Table Container -->
    <div class="glass-panel rounded-xl overflow-hidden shadow-2xl">
      <!-- Scrollable Area containing both Header and Rows -->
      <div 
        ref="scrollContainer"
        class="custom-scrollbar max-h-[calc(100vh-350px)] overflow-y-auto relative"
      >
        <!-- Sticky Table Header -->
        <div class="sticky top-0 z-30 bg-[#0f1115] border-b border-white/10 px-6 py-4">
          <div class="grid grid-cols-12 gap-4 text-xs font-bold text-gray-400 uppercase tracking-wider">
            <div 
              class="col-span-1 cursor-pointer hover:text-white flex items-center gap-1 transition-colors select-none whitespace-nowrap"
              @click="toggleSort('rank')"
            >
              Rank
              <span v-if="sortField === 'rank'" class="text-primary">{{ sortDesc ? '↓' : '↑' }}</span>
            </div>
            <div class="col-span-3 whitespace-nowrap">Token Project</div>
            <div 
              class="col-span-1 text-right cursor-pointer hover:text-white flex items-center justify-end gap-1 transition-colors select-none whitespace-nowrap"
              @click="toggleSort('price')"
            >
              Price
              <span v-if="sortField === 'price'" class="text-primary">{{ sortDesc ? '↓' : '↑' }}</span>
            </div>
            
            <div 
              class="col-span-1 text-right cursor-pointer hover:text-white flex items-center justify-end gap-1 transition-colors select-none whitespace-nowrap"
              @click="toggleSort('change_24h')"
            >
              24h %
              <span v-if="sortField === 'change_24h'" class="text-primary">{{ sortDesc ? '↓' : '↑' }}</span>
            </div>
            
            <div class="col-span-1 text-right whitespace-nowrap">7d %</div>

            <div 
              class="col-span-2 text-right cursor-pointer hover:text-white flex items-center justify-end gap-1 transition-colors select-none whitespace-nowrap"
              @click="toggleSort('market_cap')"
            >
              Market Cap
              <span v-if="sortField === 'market_cap'" class="text-primary">{{ sortDesc ? '↓' : '↑' }}</span>
            </div>
            
            <div 
              class="col-span-1 text-right cursor-pointer hover:text-white flex items-center justify-end gap-1 transition-colors select-none whitespace-nowrap"
              @click="toggleSort('volume_24h')"
            >
              Volume
              <span v-if="sortField === 'volume_24h'" class="text-primary">{{ sortDesc ? '↓' : '↑' }}</span>
            </div>

            <div 
              class="col-span-2 text-right cursor-pointer hover:text-white flex items-center justify-end gap-1 transition-colors select-none whitespace-nowrap"
              @click="toggleSort('trust_score')"
            >
              Score
              <span v-if="sortField === 'trust_score'" class="text-primary">{{ sortDesc ? '↓' : '↑' }}</span>
            </div>
          </div>
        </div>

        <!-- Table Body -->
        <!-- Loading State -->
        <LoadingSkeleton v-if="loading" :rows="10" />

        <!-- Empty State -->
        <div v-else-if="isEmpty" class="py-24 text-center">
          <svg class="w-16 h-16 mx-auto text-gray-700 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
          <p class="text-gray-400 text-lg font-light">No tokens found matching your criteria</p>
          <button 
            @click="$emit('reset')"
            class="mt-4 text-primary hover:text-accent font-medium text-sm transition-colors"
          >
            Clear Filters
          </button>
        </div>

        <!-- Token Rows -->
        <div v-else class="divide-y divide-white/5">
          <TokenRow
            v-for="token in sortedTokens"
            :key="token.symbol"
            :token="token"
            class="hover:bg-white/5 transition-colors cursor-pointer"
          />
          
          <!-- Infinite Scroll Sentinel -->
          <div ref="sentinel" class="py-10 flex justify-center items-center">
            <div v-if="loading && tokens.length > 0" class="flex items-center gap-3 text-primary text-sm font-medium">
              <svg class="animate-spin h-5 w-5" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" fill="none"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              <span>Loading more tokens...</span>
            </div>
            <div v-else-if="!hasMore && tokens.length > 0" class="text-gray-500 text-xs uppercase tracking-widest font-bold opacity-50">
              End of market pulse
            </div>
          </div>
        </div>
      </div>
    </div>


  </div>
</template>

<script setup>
import { ref, onMounted, computed, watch } from 'vue'
import FilterBar from './FilterBar.vue'
import LoadingSkeleton from './LoadingSkeleton.vue'
import TokenRow from './TokenRow.vue'
import { useTokens } from '../composables/useTokens'
import { useWatchlist } from '../composables/useWatchlist'
import { timeAgo } from '../utils/formatters'

const {
  tokens,
  loading,
  error,
  lastFetch,
  filters,
  hasTokens,
  isEmpty,
  hasMore,
  fetchTokens,
  loadMore,
  updateFilters,
  resetFilters,
  refresh
} = useTokens()

// Infinite Scroll Observer
const sentinel = ref(null)
const scrollContainer = ref(null)
let observer = null

const setupObserver = () => {
  if (observer) observer.disconnect()
  
  // Use a slight delay to ensure DOM is settled
  setTimeout(() => {
    if (!sentinel.value || !scrollContainer.value) {
      console.warn('⚠️ Sentinel or Container not found for observer')
      return
    }

    observer = new IntersectionObserver((entries) => {
      const entry = entries[0]
      if (entry.isIntersecting && hasMore.value && !loading.value) {
        console.log('📜 Load more triggered by scroll')
        loadMore()
      }
    }, {
      root: scrollContainer.value,
      rootMargin: '150px', 
      threshold: 0.01
    })

    observer.observe(sentinel.value)
    console.log('👀 Observer attached to sentinel')
  }, 100)
}

// Watch for data/loading changes to re-bind observer if needed
watch([tokens, loading], () => {
  if (!loading.value && tokens.value.length > 0) {
    setupObserver()
  }
}, { deep: true })

// Sorting State
const sortField = ref('rank') // Default sort by rank
const sortDesc = ref(false)   // Default asc (1, 2, 3...)
const showFavorites = ref(false) // Watchlist toggle

const { isFavorite } = useWatchlist()

const toggleSort = (field) => {
  if (sortField.value === field) {
    sortDesc.value = !sortDesc.value
  } else {
    sortField.value = field
    sortDesc.value = true // Default to descending for new metrics (Price, Score)
    if (field === 'rank') sortDesc.value = false // Rank is usually Ascending
  }
}

const toggleFavoritesFilter = () => {
  showFavorites.value = !showFavorites.value
}

const sortedTokens = computed(() => {
  if (!tokens.value) return []
  
  let result = [...tokens.value]

  // Filter by Favorites
  if (showFavorites.value) {
    result = result.filter(t => isFavorite(t.symbol))
  }

  return result.sort((a, b) => {
    let aVal = a[sortField.value]
    let bVal = b[sortField.value]

    // Handle nested score properties if needed, or mapped fields
    if (sortField.value === 'price') {
      aVal = a.price
      bVal = b.price
    } else if (sortField.value === 'change_24h') {
      aVal = a.change_24h
      bVal = b.change_24h
    } else if (sortField.value === 'trust_score') {
      aVal = a.trust_score || 0
      bVal = b.trust_score || 0
    } else if (sortField.value === 'rank') {
      aVal = a.rank
      bVal = b.rank
    }

    if (aVal === bVal) return 0
    const comparison = aVal > bVal ? 1 : -1
    return sortDesc.value ? -comparison : comparison
  })
})

// Fetch tokens on mount
onMounted(() => {
  fetchTokens()
  setupObserver()
})
</script>
