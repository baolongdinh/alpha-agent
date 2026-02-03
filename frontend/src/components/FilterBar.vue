<template>
  <div class="mb-8 space-y-4">
    <!-- Main Bar -->
    <div class="flex flex-col md:flex-row gap-4 items-center justify-between">
      <div class="flex items-center gap-3 w-full md:w-auto">
        <!-- Search -->
        <div class="relative w-full md:w-96 group">
          <input
            type="text"
            v-model="localFilters.search"
            @input="onSearchInput"
            placeholder="Search symbol or name..."
            class="w-full bg-white/5 border border-white/10 rounded-2xl py-3 pl-12 pr-10 text-white text-sm placeholder-gray-500 focus:outline-none focus:border-primary/50 focus:bg-white/10 transition-all shadow-2xl group-hover:border-white/20"
          />
          <svg class="w-5 h-5 absolute left-4 top-1/2 -translate-y-1/2 text-gray-500 group-hover:text-primary transition-colors pointer-events-none" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
        </div>

        <!-- Advanced Toggle -->
        <button
          @click="showAdvanced = !showAdvanced"
          class="flex items-center gap-2 px-5 py-3 rounded-2xl transition-all border shrink-0"
          :class="showAdvanced ? 'bg-primary text-black border-primary font-black' : 'bg-white/5 text-gray-400 border-white/10 hover:bg-white/10 font-bold'"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6V4m0 2a2 2 0 100 4m0-4a2 2 0 110 4m-6 8a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4m6 6v10m6-2a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4" />
          </svg>
          <span class="text-xs uppercase tracking-widest">{{ showAdvanced ? 'Hide Filters' : 'Advanced Filters' }}</span>
        </button>
      </div>

      <div class="flex items-center gap-3">
        <!-- Refresh -->
        <button
          @click="refresh"
          :disabled="loading"
          class="group flex items-center gap-2 px-5 py-3 rounded-2xl bg-white/5 text-gray-400 border border-white/10 hover:bg-white/10 hover:text-white transition-all disabled:opacity-50"
        >
          <svg class="w-4 h-4 group-hover:rotate-180 transition-transform duration-500" :class="{ 'animate-spin': loading }" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
          </svg>
          <span class="text-xs font-bold uppercase tracking-wider">Refresh</span>
        </button>

        <!-- Reset -->
        <button
          @click="resetFilters"
          class="px-5 py-3 rounded-2xl bg-rose-500/10 text-rose-500 border border-rose-500/20 hover:bg-rose-500/20 transition-all font-bold text-xs uppercase tracking-wider"
        >
          Reset
        </button>
      </div>
    </div>

    <!-- Advanced Panel -->
    <div 
      v-if="showAdvanced" 
      class="glass-panel p-8 rounded-[32px] bg-black/40 backdrop-blur-3xl border border-white/10 shadow-3xl animate-slide-in relative overflow-hidden"
    >
      <div class="absolute top-0 right-0 w-64 h-64 bg-primary/5 blur-[100px] rounded-full"></div>
      
      <div class="grid grid-cols-1 md:grid-cols-4 gap-8 relative z-10">
        <!-- Market Cap Filter -->
        <div class="space-y-3">
          <label class="text-[10px] font-black text-gray-500 uppercase tracking-[0.2em] flex items-center gap-2">
            <span class="w-1.5 h-1.5 bg-primary rounded-full"></span> Market Cap (USD)
          </label>
          <div class="grid grid-cols-2 gap-2">
            <input 
              type="number" 
              v-model="localFilters.minMcap" 
              @input="onFilterUpdate"
              placeholder="Min" 
              class="bg-white/5 border border-white/5 rounded-xl px-4 py-2 text-sm text-white focus:outline-none focus:border-primary/40 focus:bg-white/10 transition-all font-mono"
            />
            <input 
              type="number" 
              v-model="localFilters.maxMcap" 
              @input="onFilterUpdate"
              placeholder="Max" 
              class="bg-white/5 border border-white/5 rounded-xl px-4 py-2 text-sm text-white focus:outline-none focus:border-primary/40 focus:bg-white/10 transition-all font-mono"
            />
          </div>
          <div class="text-[9px] text-gray-600 font-bold italic">e.g. 1000000000 for 1B</div>
        </div>

        <!-- Trust Score Filter -->
        <div class="space-y-3">
          <label class="text-[10px] font-black text-gray-500 uppercase tracking-[0.2em] flex items-center gap-2">
            <span class="w-1.5 h-1.5 bg-green-500 rounded-full"></span> Alpha Trust Score
          </label>
          <div class="grid grid-cols-2 gap-2">
            <input 
              type="number" 
              v-model="localFilters.minScore" 
              @input="onFilterUpdate"
              placeholder="0" 
              max="100"
              class="bg-white/5 border border-white/5 rounded-xl px-4 py-2 text-sm text-white focus:outline-none focus:border-primary/40 focus:bg-white/10 transition-all font-mono"
            />
            <input 
              type="number" 
              v-model="localFilters.maxScore" 
              @input="onFilterUpdate"
              placeholder="100" 
              max="100"
              class="bg-white/5 border border-white/5 rounded-xl px-4 py-2 text-sm text-white focus:outline-none focus:border-primary/40 focus:bg-white/10 transition-all font-mono"
            />
          </div>
          <div class="flex justify-between text-[9px] text-gray-600 font-bold">
            <span>Minimum</span>
            <span>Maximum</span>
          </div>
        </div>

        <!-- Price Rank -->
        <div class="space-y-3">
          <label class="text-[10px] font-black text-gray-500 uppercase tracking-[0.2em] flex items-center gap-2">
            <span class="w-1.5 h-1.5 bg-purple-500 rounded-full"></span> Unit Price (USD)
          </label>
          <div class="grid grid-cols-2 gap-2">
            <input 
              type="number" 
              v-model="localFilters.minPrice" 
              @input="onFilterUpdate"
              placeholder="Min" 
              class="bg-white/5 border border-white/5 rounded-xl px-4 py-2 text-sm text-white focus:outline-none focus:border-primary/40 focus:bg-white/10 transition-all font-mono"
            />
            <input 
              type="number" 
              v-model="localFilters.maxPrice" 
              @input="onFilterUpdate"
              placeholder="Max" 
              class="bg-white/5 border border-white/5 rounded-xl px-4 py-2 text-sm text-white focus:outline-none focus:border-primary/40 focus:bg-white/10 transition-all font-mono"
            />
          </div>
        </div>

        <!-- Sector / Category -->
        <div class="space-y-3">
          <label class="text-[10px] font-black text-gray-500 uppercase tracking-[0.2em] flex items-center gap-2">
            <span class="w-1.5 h-1.5 bg-yellow-500 rounded-full"></span> Sector Classification
          </label>
          <select 
            v-model="localFilters.category" 
            @change="onFilterUpdate"
            class="w-full bg-white/5 border border-white/5 rounded-xl px-4 py-2 text-sm text-white focus:outline-none focus:border-primary/40 focus:bg-white/10 transition-all appearance-none cursor-pointer"
          >
            <option value="">All Sectors</option>
            <option value="L1/L2">L1/L2 Protocols</option>
            <option value="DeFi">DeFi Engines</option>
            <option value="AI">AI & Agents</option>
            <option value="Memecoins">High Volatility / Memes</option>
            <option value="Infrastructure">Web3 Infra</option>
          </select>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'

const props = defineProps({
  filters: {
    type: Object,
    required: true
  },
  loading: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:filters', 'refresh', 'reset'])

const localFilters = ref({ ...props.filters })
const showAdvanced = ref(false)

let debounceTimeout = null
const debounce = (fn, delay) => {
  return (...args) => {
    if (debounceTimeout) clearTimeout(debounceTimeout)
    debounceTimeout = setTimeout(() => fn(...args), delay)
  }
}

watch(() => props.filters, (newVal) => {
  localFilters.value = { ...newVal }
}, { deep: true })

const onFilterUpdate = () => {
  emit('update:filters', { ...localFilters.value })
}

const onSearchInput = debounce(() => {
  onFilterUpdate()
}, 300)

const resetFilters = () => {
  emit('reset')
}

const refresh = () => {
  emit('refresh')
}
</script>

<style scoped>
.animate-slide-in {
  animation: slideIn 0.4s cubic-bezier(0.4, 0, 0.2, 1);
}

@keyframes slideIn {
  from { opacity: 0; transform: translateY(-10px); }
  to { opacity: 1; transform: translateY(0); }
}

select {
  background-image: url("data:image/svg+xml,%3csvg xmlns='http://www.w3.org/2000/svg' fill='none' viewBox='0 0 20 20'%3e%3cpath stroke='%236b7280' stroke-linecap='round' stroke-linejoin='round' stroke-width='1.5' d='M6 8l4 4 4-4'/%3e%3c/svg%3e");
  background-position: right 0.5rem center;
  background-repeat: no-repeat;
  background-size: 1.5em 1.5em;
}

input::-webkit-outer-spin-button,
input::-webkit-inner-spin-button {
  -webkit-appearance: none;
  margin: 0;
}
</style>
