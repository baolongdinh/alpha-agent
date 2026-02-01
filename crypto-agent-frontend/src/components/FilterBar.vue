<template>
  <div class="flex flex-col md:flex-row gap-4 items-center justify-between mb-6">
    <!-- Search Bar -->
    <div class="relative w-full md:w-96 group">
      <input
        type="text"
        v-model="localFilters.search"
        @input="onSearchInput"
        placeholder="Search tokens..."
        class="w-full bg-white/5 border border-white/10 rounded-full py-2.5 pl-10 pr-10 text-white text-sm placeholder-gray-500 focus:outline-none focus:border-primary/50 focus:bg-white/10 transition-all shadow-sm group-hover:border-white/20"
      />
      <!-- Search Icon -->
      <svg class="w-4 h-4 absolute left-3.5 top-1/2 -translate-y-1/2 text-gray-500 group-hover:text-gray-400 transition-colors pointer-events-none" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
      </svg>
    </div>

    <!-- Refresh Button -->
    <button
      @click="refresh"
      :disabled="loading"
      class="flex items-center gap-2 px-4 py-2 rounded-full bg-primary/10 text-primary border border-primary/20 hover:bg-primary/20 transition-all disabled:opacity-50"
    >
      <svg class="w-4 h-4" :class="{ 'animate-spin': loading }" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
      </svg>
      <span class="text-xs font-bold uppercase tracking-wider">Refresh</span>
    </button>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'

const props = defineProps({
  filters: {
    type: Object,
    required: true
  },
  showFavorites: {
    type: Boolean,
    default: false
  },
  loading: {
    type: Boolean,
    default: false
  },
  lastFetch: {
    type: Date,
    default: null
  }
})

const emit = defineEmits(['update:filters', 'refresh', 'reset'])

const localFilters = ref({ ...props.filters })
let debounceTimeout = null

// Debounce function
const debounce = (fn, delay) => {
  return (...args) => {
    if (debounceTimeout) clearTimeout(debounceTimeout)
    debounceTimeout = setTimeout(() => {
      fn(...args)
    }, delay)
  }
}

// Watchers
watch(() => props.filters, (newFilters) => {
  // Only update if not currently typing deeply (simple sync)
  if (newFilters.search !== localFilters.value.search) {
     localFilters.value = { ...newFilters }
  } else {
     // merge other fields
     localFilters.value = { ...localFilters.value, ...newFilters }
  }
}, { deep: true })

const applyFilters = () => {
  emit('update:filters', { ...localFilters.value })
}

const selectCategory = (cat) => {
  localFilters.value.category = cat
  applyFilters()
}

// Debounced Search
const debouncedSearch = debounce(() => {
  applyFilters()
}, 300)

const onSearchInput = () => {
  debouncedSearch()
}

const clearSearch = () => {
  localFilters.value.search = ''
  applyFilters()
}

const resetFilters = () => {
  emit('reset')
}

const refresh = () => {
  emit('refresh')
}
</script>
