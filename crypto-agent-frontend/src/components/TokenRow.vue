<template>
  <div 
    class="block animate-fade-in group active:scale-[0.99] transition-transform duration-100 ease-out cursor-pointer"
    @click="router.push(`/token/${token.id}`)"
  >
    <div class="grid grid-cols-12 gap-4 items-center py-4 px-6 relative z-10">
      <!-- Rank & Star -->
      <div class="col-span-1 flex items-center gap-2">
        <button 
          @click.stop="toggleFavorite(token.symbol)"
          class="text-gray-600 hover:text-yellow-400 transition-colors focus:outline-none"
        >
          <svg v-if="isFavorite(token.symbol)" class="w-4 h-4 text-yellow-400 fill-current" viewBox="0 0 20 20"><path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z"/></svg>
          <svg v-else class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.197-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z" /></svg>
        </button>
        <span class="text-gray-500 font-mono text-sm group-hover:text-white transition-colors">
          {{ token.rank || '#' }}
        </span>
      </div>

      <!-- Token Info -->
      <div class="col-span-3">
        <div class="flex items-center gap-3">
          <div class="w-8 h-8 rounded-full bg-gradient-to-br from-gray-700 to-gray-800 flex items-center justify-center text-xs font-bold ring-1 ring-white/10 overflow-hidden">
            <img v-if="token.image" :src="token.image" :alt="token.symbol" class="w-full h-full object-cover" />
            <span v-else>{{ token.symbol[0] }}</span>
          </div>
          <div>
            <div class="font-bold text-white leading-none">{{ token.symbol }}</div>
            <div class="text-xs text-gray-500 mt-1 truncate max-w-[120px]">{{ token.name }}</div>
          </div>
          <div v-if="token.category" class="hidden sm:inline-block px-2 py-0.5 bg-white/5 text-gray-400 text-[10px] uppercase tracking-wide rounded border border-white/5">
            {{ token.category }}
          </div>
        </div>
      </div>

      <!-- Price -->
      <div class="col-span-2 text-right">
        <div class="text-white font-mono font-medium tracking-tight">{{ formatCurrency(token.price) }}</div>
      </div>

      <!-- 24h Change -->
      <div class="col-span-1 text-right">
        <div 
          class="font-mono text-sm"
          :class="token.change_24h >= 0 ? 'text-green-400' : 'text-red-400'"
        >
          {{ Math.abs(token.change_24h).toFixed(1) }}%
        </div>
      </div>

      <!-- 7d Trend (Sparkline) -->
      <div class="col-span-3 flex justify-center items-center h-10">
        <SparklineChart 
          v-if="token.sparkline && token.sparkline.length"
          :data="token.sparkline" 
          width="120" 
          height="40"
          :is-positive="token.change_7d >= 0"
        />
        <div v-else class="text-xs text-gray-600 self-center">No Data</div>
      </div>

      <!-- Trust Score -->
      <div class="col-span-2 text-right flex items-center justify-end gap-3">
        <div class="text-right">
          <div class="font-bold text-white text-sm">{{ token.trust_score?.toFixed(0) || '-' }}</div>
          <div class="text-[10px] text-gray-500 uppercase">Score</div>
        </div>
        <RatingBadge :grade="token.score_breakdown?.grade || 'D'" />
      </div>
    </div>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
import RatingBadge from './RatingBadge.vue'
import SparklineChart from './SparklineChart.vue'
import { formatNumber, formatCurrency } from '../utils/formatters'
import { useWatchlist } from '../composables/useWatchlist'

const router = useRouter()
const { isFavorite, toggleFavorite } = useWatchlist()

defineProps({
  token: {
    type: Object,
    required: true
  }
})

const getScoreColor = (score) => {
  if (score >= 70) return 'var(--rating-aaa)' // Green
  if (score >= 50) return 'var(--rating-b)'   // Yellow
  if (score >= 35) return 'var(--rating-c)'   // Orange
  return 'var(--rating-d)'                    // Red
}
</script>
