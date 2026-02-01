<template>
  <Modal @close="$emit('close')">
    <div class="p-0 overflow-hidden bg-[#0f1115] max-w-6xl w-full mx-auto rounded-3xl shadow-[0_0_50px_rgba(0,0,0,0.5)] border border-white/10 flex flex-col max-h-[90vh] relative">
      <!-- Header Section -->
      <div class="p-8 pb-12 bg-gradient-to-br from-[#1a1c24] to-[#0a0b10] relative">
        <button @click="$emit('close')" class="absolute top-6 right-6 p-2 rounded-full bg-black/50 hover:bg-black/80 text-white transition-colors z-20">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/></svg>
        </button>

        <div v-if="displayToken && displayToken.symbol">
          <div class="flex items-start justify-between">
            <div class="flex items-center gap-6">
              <div class="w-20 h-20 rounded-2xl bg-gradient-to-br from-gray-700 to-gray-800 flex items-center justify-center text-3xl font-bold border border-white/10 shadow-2xl overflow-hidden">
                <img v-if="displayToken.image" :src="displayToken.image" :alt="displayToken.name" class="w-full h-full object-cover" />
                <span v-else-if="displayToken.symbol">{{ displayToken.symbol[0] }}</span>
                <span v-else>?</span>
              </div>
              <div>
                <div class="flex items-center gap-3 mb-2">
                  <h2 class="text-4xl font-bold text-white tracking-tight">{{ displayToken.name }}</h2>
                  <span class="px-2 py-1 bg-white/10 rounded-md text-gray-400 text-xs font-mono uppercase tracking-widest">{{ displayToken.symbol }}</span>
                </div>
                <div class="flex items-center gap-2">
                  <div class="px-3 py-1 bg-blue-600/10 text-blue-400 text-xs font-bold rounded-full border border-blue-500/20">Rank #{{ displayToken.rank }}</div>
                  <div v-if="displayToken.category" class="px-3 py-1 bg-white/5 text-gray-400 text-xs font-bold rounded-full border border-white/10">{{ displayToken.category }}</div>
                </div>
              </div>
            </div>
            <div class="text-right hidden sm:block">
              <div class="text-5xl font-bold text-white tracking-tighter mb-2">{{ formatCurrency(displayToken.price) }}</div>
              <div 
                class="flex flex-col items-end gap-1 font-mono"
              >
                <!-- 7D Mini Chart (CMC Sourced) -->
                <div v-if="displayToken.sparkline && displayToken.sparkline.length" class="mb-2 h-10 w-28 opacity-80 group-hover:opacity-100 transition-opacity flex items-center justify-center">
                  <SparklineChart 
                    :data="displayToken.sparkline" 
                    width="110" 
                    height="34" 
                    :is-positive="displayToken.change_7d >= 0"
                  />
                </div>

                <div v-if="displayToken.change_1h !== undefined" :class="displayToken.change_1h >= 0 ? 'text-green-400' : 'text-red-400'">
                  <span class="text-[10px] bg-black/30 px-2 py-0.5 rounded-full border border-white/5 opacity-70">
                    {{ displayToken.change_1h >= 0 ? '↑' : '↓' }} 1h {{ Math.abs(displayToken.change_1h).toFixed(2) }}%
                  </span>
                </div>
                <div :class="displayToken.change_24h >= 0 ? 'text-green-400' : 'text-red-400'">
                  <span class="text-sm bg-black/30 px-2 py-0.5 rounded-full border" :class="displayToken.change_24h >= 0 ? 'border-green-500/30' : 'border-red-500/30'">
                    {{ displayToken.change_24h >= 0 ? '↑' : '↓' }} 24H {{ Math.abs(displayToken.change_24h).toFixed(2) }}%
                  </span>
                </div>
                <div v-if="displayToken.change_7d" :class="displayToken.change_7d >= 0 ? 'text-green-400' : 'text-red-400'">
                  <span class="text-[10px] bg-black/30 px-2 py-0.5 rounded-full border opacity-80" :class="displayToken.change_7d >= 0 ? 'border-green-500/20' : 'border-red-500/20'">
                    {{ displayToken.change_7d >= 0 ? '↑' : '↓' }} 7D {{ Math.abs(displayToken.change_7d).toFixed(2) }}%
                  </span>
                </div>
              </div>
            </div>
          </div>

          <!-- Supply Highlights (NEW) -->
          <div class="grid grid-cols-2 sm:grid-cols-4 gap-4 mt-8 pb-8 border-b border-white/5">
            <div class="space-y-1">
              <div class="text-[10px] text-gray-500 uppercase tracking-wider">Circulating Supply</div>
              <div class="text-sm font-mono text-white">{{ formatNumber(displayToken.circulating_supply) }} {{ displayToken.symbol }}</div>
            </div>
            <div class="space-y-1">
              <div class="text-[10px] text-gray-500 uppercase tracking-wider">Total Supply</div>
              <div class="text-sm font-mono text-white">{{ formatNumber(displayToken.total_supply) }} {{ displayToken.symbol }}</div>
            </div>
            <div class="space-y-1">
              <div class="text-[10px] text-gray-500 uppercase tracking-wider">Max Supply</div>
              <div class="text-sm font-mono text-white">{{ displayToken.max_supply > 0 ? formatNumber(displayToken.max_supply) : '∞' }}</div>
            </div>
            <div class="space-y-1">
              <div class="text-[10px] text-gray-500 uppercase tracking-wider">Market Dom</div>
              <div class="text-sm font-mono text-amber-400">{{ displayToken.market_cap_dominance?.toFixed(2) || '0.00' }}%</div>
            </div>
          </div>
        </div>
      </div>

      <!-- Main Content Area -->
      <div class="flex-1 overflow-y-auto custom-scrollbar bg-[#0d0e12] p-8 rounded-t-[40px] border-t border-white/5 relative z-10 -mt-6">
        <div class="grid lg:grid-cols-12 gap-8">
          
          <!-- LEFT: Scoring & Sentiment (4 cols) -->
          <div class="lg:col-span-4 space-y-6">
            <!-- Trust Score Card -->
            <div class="glass-bento p-6">
              <div class="text-[10px] font-bold text-gray-500 uppercase tracking-[0.2em] mb-4">Alpha Intelligence</div>
              <div class="flex items-center gap-6">
                <RatingBadge 
                  :grade="displayToken.score_breakdown?.grade || 'D'" 
                  class="text-4xl px-5 py-3 rounded-2xl shadow-xl"
                />
                <div>
                  <div class="text-4xl font-black text-white tracking-tighter">{{ displayToken.trust_score?.toFixed(0) || '—' }}</div>
                  <div class="text-[10px] text-gray-500 font-bold uppercase tracking-widest mt-1">Trust Score</div>
                </div>
              </div>
            </div>

            <!-- Sentiment Gauge Card -->
            <div class="glass-bento p-6 flex flex-col items-center justify-center">
              <div class="text-[10px] font-bold text-gray-500 uppercase tracking-[0.2em] mb-6 w-full text-left">Market Sentiment</div>
              <SentimentGauge :score="displayToken.sentiment || 0" />
              <div v-if="displayToken.market_cap" class="mt-4 text-[10px] text-gray-500 font-medium">
                Based on ${{ formatNumber(displayToken.market_cap) }} cap
              </div>
            </div>

            <!-- Radar Chart Container -->
            <div class="glass-bento p-6 flex flex-col items-center justify-center min-h-[350px]">
              <div class="w-full h-[300px] relative">
                <ScoreRadar v-if="displayToken.score_breakdown" :scores="{
                  liquidity: displayToken.score_breakdown.liquidity_score || 0,
                  volume: displayToken.score_breakdown.volume_score || 0,
                  trend: displayToken.score_breakdown.trend_score || 0,
                  social: Math.min((displayToken.market_cap_dominance || 0) * 10, 100),
                  market: displayToken.score_breakdown.market_health_score || 0,
                  stability: displayToken.score_breakdown.risk_score || 0
                }" />
                <div v-else class="text-gray-500 text-sm animate-pulse">Computing Matrix...</div>
              </div>
            </div>
          </div>

          <!-- RIGHT: High Density Bento Grid (8 cols) -->
          <div class="lg:col-span-8 space-y-8">
            
            <!-- Section: Performance Benchmarks -->
            <div>
              <div class="text-[10px] font-bold text-cyan-500/80 uppercase tracking-[0.2em] mb-4 px-2">Performance Benchmarks</div>
              <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
                <MetricCard 
                  label="Market Cap" 
                  :value="formatCurrencyCompact(displayToken.market_cap)" 
                  icon="💰"
                  :subLabel="`Rank #${displayToken.rank}`"
                />
                <MetricCard 
                   label="24h Volume" 
                   :value="formatCurrencyCompact(displayToken.volume_24h)" 
                   icon="📊"
                   :subLabel="displayToken.volume_change_24h ? (displayToken.volume_change_24h >= 0 ? '↑ ' : '↓ ') + displayToken.volume_change_24h.toFixed(2) + '% Momentum' : 'Vol / Cap Ratio'"
                   :progress="displayToken.market_cap > 0 ? (displayToken.volume_24h / displayToken.market_cap) * 100 : 0"
                />
                <MetricCard 
                  label="ATH Distance" 
                  :value="displayToken.ath_change ? displayToken.ath_change.toFixed(1) + '%' : '—'" 
                  :change="displayToken.ath_change"
                  icon="🏔️"
                  :progress="Math.abs(displayToken.ath_change || 0)"
                  progressColor="red"
                />
              </div>
            </div>

            <!-- Section: Advanced Analytics -->
            <div>
              <div class="text-[10px] font-bold text-amber-500/80 uppercase tracking-[0.2em] mb-4 px-2">Advanced Analytics</div>
              <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
                <MetricCard 
                   label="Market Strength" 
                   :value="displayToken.market_cap_dominance ? displayToken.market_cap_dominance.toFixed(2) + '%' : '—'" 
                   icon="⭐"
                   subLabel="Total Market Share"
                   :progress="displayToken.market_cap_dominance * 5"
                />
                <MetricCard 
                   label="7D Volatility" 
                   :value="displayToken.volatility_7d ? displayToken.volatility_7d.toFixed(2) + '%' : '—'" 
                   icon="📊"
                   subLabel="Risk Coefficient"
                   :progress="displayToken.volatility_7d"
                />
                <MetricCard 
                   v-if="displayToken.market_cap > 0 && displayToken.fdv > 0"
                   label="FDV / Mcap" 
                   :value="(displayToken.fdv / displayToken.market_cap).toFixed(2)" 
                   icon="⚖️"
                   subLabel="Circulation Ratio"
                />
                 <MetricCard 
                   v-if="displayToken.holder_count > 0"
                   label="Holders" 
                   :value="displayToken.holder_count.toLocaleString()" 
                   icon="👥"
                   subLabel="Verified Addresses"
                />
                <MetricCard 
                   v-if="displayToken.market_cap > 0 && (displayToken.liquidity > 0 || displayToken.tvl > 0)"
                   label="Liquidity Index" 
                   :value="displayToken.liquidity > 0 ? ((displayToken.liquidity / displayToken.market_cap) * 100).toFixed(2) + '%' : 'N/A'" 
                   icon="📈"
                   subLabel="Liq / MCap Index"
                   :progress="displayToken.liquidity > 0 ? Math.min((displayToken.liquidity / displayToken.market_cap) * 1000, 100) : 0"
                />
              </div>
            </div>

             <!-- AI Analysis Bento -->
             <div class="glass-bento p-8 relative overflow-hidden group">
               <!-- Decorative Gradient -->
               <div class="absolute -top-24 -right-24 w-64 h-64 bg-blue-600/5 blur-[80px] rounded-full group-hover:bg-blue-600/10 transition-colors"></div>
               
               <div class="flex justify-between items-center mb-6 relative z-10">
                  <div class="flex items-center gap-3">
                    <div class="w-8 h-8 rounded-lg bg-blue-600/20 flex items-center justify-center text-blue-400">✨</div>
                    <h3 class="text-sm font-bold text-white uppercase tracking-widest">Alpha Intelligence Report</h3>
                  </div>
                  <div v-if="isAnalyzing" class="flex items-center gap-2">
                    <span class="w-2 h-2 bg-blue-500 rounded-full animate-pulse"></span>
                    <span class="text-[10px] font-bold text-blue-400 uppercase tracking-tighter">Analyzing Real-time Data...</span>
                  </div>
               </div>
               
               <div v-if="analysis" class="prose prose-invert prose-sm max-w-none text-gray-300 relative z-10 leading-relaxed">
                  <div v-html="renderMarkdown(analysis)"></div>
               </div>
               <div v-else-if="aiError" class="text-rose-400 text-xs font-medium bg-rose-500/10 p-4 rounded-xl border border-rose-500/20 mb-3 relative z-10">
                  ⚠️ {{ aiError }}
               </div>
               
               <!-- Manual Trigger Button -->
               <div v-if="!analysis && !isAnalyzing" class="flex flex-col items-center justify-center py-10 text-center relative z-10">
                  <p class="text-gray-500 text-xs mb-8 max-w-xs uppercase tracking-widest font-medium">
                    Synthesizing multi-source signals...
                  </p>
                  <button 
                    @click="analyzeToken(displayToken)"
                    class="px-10 py-4 rounded-2xl bg-gradient-to-r from-blue-600 to-indigo-600 hover:from-blue-500 hover:to-indigo-500 text-white font-black text-sm uppercase tracking-widest transition-all shadow-[0_10px_30px_rgba(37,99,235,0.3)] hover:shadow-[0_15px_40px_rgba(37,99,235,0.5)] active:scale-95 flex items-center gap-4"
                  >
                    Generate Agent Report
                  </button>
               </div>
             </div>
          </div>
        </div>
      </div>
    </div>
  </Modal>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import Modal from './Modal.vue'
import RatingBadge from './RatingBadge.vue'
import ScoreRadar from './ScoreRadar.vue'
import MetricCard from './MetricCard.vue'
import SentimentGauge from './SentimentGauge.vue'
import SparklineChart from './SparklineChart.vue'
import { formatCurrency, formatPrice, formatNumber } from '../utils/formatters'
import { formatCurrencyCompact } from '../services/api'
import { useTokens } from '../composables/useTokens'
import { useAIAnalysis } from '../composables/useAIAnalysis'
import { marked } from 'marked'
import DOMPurify from 'dompurify'

const props = defineProps({
  token: {
    type: Object,
    required: true
  }
})

defineEmits(['close'])

const { fetchTokenDetail } = useTokens()
const { analysis, isAnalyzing, error: aiError, analyzeToken, clearAnalysis } = useAIAnalysis()

const detailData = ref(null)
const isLoading = ref(false)

const renderMarkdown = (text) => {
  if (!text) return ''
  const html = marked.parse(text)
  return DOMPurify.sanitize(html)
}

// Merge list data with detailed data
const displayToken = computed(() => {
  return { ...props.token, ...detailData.value }
})

onMounted(async () => {
  clearAnalysis() // Clear previous analysis
  if (props.token.id) {
    isLoading.value = true
    try {
      const data = await fetchTokenDetail(props.token.id)
      if (data) {
        detailData.value = data
      }
    } finally {
      isLoading.value = false
    }
  }
})
</script>
