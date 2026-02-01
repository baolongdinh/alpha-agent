<template>
  <div class="min-h-screen bg-[#050510] bg-[radial-gradient(ellipse_at_top,_var(--tw-gradient-stops))] from-blue-900/20 via-[#050510] to-[#050510]">
    <!-- Top Navigation -->
    <header class="sticky top-0 z-40 w-full glass-panel border-b border-white/5 backdrop-blur-xl">
      <div class="container mx-auto px-4 h-16 flex items-center justify-between">
        <!-- Logo -->
        <div class="flex items-center gap-2">
          <div class="w-8 h-8 rounded-lg bg-gradient-to-tr from-blue-500 to-cyan-500 flex items-center justify-center text-white font-black shadow-lg shadow-blue-500/20">
            AA
          </div>
          <span class="text-xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-white to-gray-400 tracking-tight">
            AlphaAgent
          </span>
        </div>

        <div class="flex items-center gap-4">
           <div class="flex items-center gap-2 px-3 py-1 rounded-full bg-emerald-500/10 border border-emerald-500/20 transition-all duration-300">
              <span class="w-1.5 h-1.5 rounded-full bg-emerald-500"></span>
              <span class="text-[10px] font-bold text-emerald-400 uppercase tracking-widest px-1">Market Real-time</span>
           </div>
        </div>
      </div>
    </header>

    <!-- Global Ticker (Marquee) -->
    <div v-if="hasTokens && marketStats" class="bg-black/60 border-b border-white/5 overflow-hidden h-10 flex items-center relative transition-opacity duration-500" :class="{ 'opacity-0': loading && !hasTokens }">
      <div class="flex animate-marquee hover:pause whitespace-nowrap gap-16 text-[10px] uppercase font-bold tracking-widest text-gray-500 px-8">
        <!-- Duplicate content for seamless loop -->
        <span v-for="i in 4" :key="i" class="flex gap-12 items-center">
          <!-- Individual Tokens -->
          <span v-for="token in topTokens.slice(0, 5)" :key="token.symbol" class="flex gap-2 items-center">
            <span class="text-white/40">{{ token.symbol }}</span>
            <span :class="token.change_24h >= 0 ? 'text-emerald-400' : 'text-rose-400'" class="font-mono">
              ${{ formatPrice(token.price) }}
            </span>
          </span>

          <!-- Global Benchmarks -->
          <span class="flex gap-4 items-center px-4 border-x border-white/10">
            <span class="flex gap-2">
              <span class="text-white/40">BTC DOM:</span>
              <span class="text-white font-mono">{{ marketStats.btc_dominance?.toFixed(1) }}%</span>
            </span>
            <span class="flex gap-2 text-cyan-500/80">
              <span class="text-white/40">ETH DOM:</span>
              <span class="text-cyan-400 font-mono">{{ marketStats.eth_dominance?.toFixed(1) }}%</span>
            </span>
            <span class="flex gap-2">
              <span class="text-white/40">GLOBAL CAP:</span>
              <span class="text-blue-400 font-mono">{{ formatCurrencyCompact(marketStats.total_market_cap) }}</span>
            </span>
            <span class="flex gap-2">
              <span class="text-white/40">24H VOL:</span>
              <span class="text-purple-400 font-mono">{{ formatCurrencyCompact(marketStats.total_volume) }}</span>
            </span>
          </span>
        </span>
      </div>
    </div>

    <!-- Main Layout -->
    <div class="container mx-auto px-4 py-8">
      <div class="grid grid-cols-12 gap-8">
        <!-- Main Content -->
        <main class="col-span-12 xl:col-span-10 xl:col-start-2">
          <!-- Page Header -->
          <div class="mb-8 flex justify-between items-end border-b border-white/5 pb-6">
            <div>
              <h1 class="text-4xl font-black text-white mb-2 tracking-tighter">Market Pulse</h1>
              <p class="text-gray-500 text-sm font-medium">Proprietary AI scoring for cross-chain liquidity & token health.</p>
            </div>
            <div class="text-right hidden sm:block">
               <div class="text-[10px] text-gray-500 font-bold uppercase tracking-widest mb-1">Last Data Sync</div>
               <div class="text-xs text-white font-mono">{{ lastFetch ? lastFetch.toLocaleTimeString() : 'connecting...' }}</div>
            </div>
          </div>

          <TokenTable />
        </main>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted, onUnmounted, ref } from 'vue'
import TokenTable from './components/TokenTable.vue'
import { useTokens } from './composables/useTokens'
import { formatPrice } from './utils/formatters'
import { formatCurrencyCompact } from './services/api'

const { refresh: refreshTokens, marketStats, topTokens, hasTokens, lastFetch, loading } = useTokens()

onMounted(() => {
  console.log('🚀 AlphaAgent initialized')
  refreshTokens() // Initial fetch only
})

onUnmounted(() => {
  if (timer) clearInterval(timer)
})
</script>

<style scoped>
/* Ensure the marquee animation exists in style.css or define here if needed, 
   but it's usually in global CSS */
</style>

