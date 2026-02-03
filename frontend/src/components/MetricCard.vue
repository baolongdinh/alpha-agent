<template>
  <div class="glass-card p-6 rounded-[24px] bg-gradient-to-br from-[#16181d] to-[#0a0b0f] border border-white/5 relative group hover:border-primary/30 transition-all duration-500 overflow-hidden shadow-lg hover:shadow-primary/5">
    <!-- Animated Glow -->
    <div class="absolute -top-12 -right-12 w-24 h-24 bg-primary/10 blur-3xl rounded-full group-hover:bg-primary/20 transition-all duration-700"></div>

    <!-- Icon Background -->
    <div class="absolute top-6 right-6 text-3xl opacity-40 group-hover:opacity-80 transition-all duration-500 transform group-hover:scale-110 group-hover:rotate-6 grayscale-[50%] group-hover:grayscale-0">
      {{ icon }}
    </div>

    <div class="flex flex-col h-full justify-between relative z-10">
      <div>
        <div class="text-[10px] font-black text-gray-500 uppercase tracking-[0.3em] mb-4 flex items-center gap-2">
          <span class="w-1 h-3 rounded-full bg-primary/40"></span>
          {{ label }}
        </div>
        
        <div class="flex items-baseline gap-3">
          <div class="text-3xl font-black text-white tracking-tighter tabular-nums drop-shadow-sm">{{ value }}</div>
          
          <div
            v-if="change !== undefined && change !== null && !isNaN(parseFloat(change))"
            class="text-[10px] font-black px-2 py-1 rounded-xl flex items-center gap-1 shadow-sm border"
            :class="isPositive ? 'text-green-400 bg-green-500/10 border-green-500/20' : 'text-rose-400 bg-rose-500/10 border-rose-500/20'"
          >
            <span class="text-sm leading-none">{{ isPositive ? '↑' : '↓' }}</span> 
            <span>{{ Math.abs(parseFloat(change)).toFixed(1) }}%</span>
          </div>
        </div>
      </div>

      <div class="mt-6 space-y-4">
        <!-- Optional Sub-label or mini-progress -->
        <div v-if="subLabel" class="text-[10px] text-gray-400 font-bold uppercase tracking-[0.2em] flex justify-between items-center opacity-70 group-hover:opacity-100 transition-opacity">
          <span>{{ subLabel }}</span>
        </div>

        <div v-if="progress !== undefined" class="relative">
          <div class="h-1.5 w-full bg-white/5 rounded-full overflow-hidden border border-white/5 p-[1px]">
            <div 
              class="h-full rounded-full transition-all duration-1000 ease-out shadow-[0_0_8px_rgba(var(--tw-shadow-color),0.5)]" 
              :class="progressClass"
              :style="{ width: `${Math.min(100, Math.max(0, progress))}%` }"
            ></div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from "vue";

const props = defineProps({
  label: String,
  value: [String, Number],
  change: [Number, String],
  icon: String, // Emoji or SVG string
  subLabel: String,
  progress: Number, // 0 to 100
  progressColor: {
    type: String,
    default: 'primary'
  }
});

const isPositive = computed(() => {
  if (props.change === undefined || props.change === null) return true;
  return parseFloat(props.change) >= 0;
});

const progressClass = computed(() => {
  if (props.progressColor === 'cyan') return 'bg-cyan-500 shadow-cyan-500/50';
  if (props.progressColor === 'green' || props.progressColor === 'emerald') return 'bg-emerald-500 shadow-emerald-500/50';
  if (props.progressColor === 'red' || props.progressColor === 'rose') return 'bg-rose-500 shadow-rose-500/50';
  if (props.progressColor === 'amber') return 'bg-amber-500 shadow-amber-500/50';
  return 'bg-primary shadow-primary/50';
});
</script>

<style scoped>
.glass-card {
  backdrop-filter: blur(10px);
}
</style>
