<template>
  <div class="glass-card p-5 rounded-2xl bg-[#13151a] relative group hover:-translate-y-1 transition-transform duration-300">
    <!-- Icon Background -->
    <div class="absolute top-4 right-4 text-2xl opacity-20 group-hover:opacity-40 transition-opacity grayscale group-hover:grayscale-0">
      {{ icon }}
    </div>

    <div class="flex flex-col h-full justify-between relative z-10">
      <div class="text-xs font-bold text-gray-400 uppercase tracking-widest mb-2">{{ label }}</div>
      
      <div class="flex items-end gap-2">
        <div class="text-xl font-bold text-white tracking-tight">{{ value }}</div>
        
        <div
          v-if="change !== undefined && change !== null && !isNaN(parseFloat(change))"
          class="text-xs font-bold px-1.5 py-0.5 rounded flex items-center gap-0.5"
          :class="isPositive ? 'text-green-400 bg-green-900/30' : 'text-red-400 bg-red-900/30'"
        >
          {{ isPositive ? '↑' : '↓' }} {{ Math.abs(parseFloat(change)).toFixed(1) }}%
        </div>
      </div>

      <!-- Optional Sub-label or mini-progress -->
      <div v-if="subLabel" class="mt-2 text-[10px] text-gray-500 font-medium uppercase tracking-wider flex justify-between items-center">
        <span>{{ subLabel }}</span>
      </div>

      <div v-if="progress !== undefined" class="mini-progress-bg">
        <div 
          class="mini-progress-fill" 
          :class="progressClass"
          :style="{ width: `${Math.min(100, Math.max(0, progress))}%` }"
        ></div>
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
    default: 'cyan'
  }
});

const isPositive = computed(() => {
  if (!props.change) return true;
  return parseFloat(props.change) >= 0;
});

const progressClass = computed(() => {
  if (props.progressColor === 'cyan') return 'bg-cyan-500';
  if (props.progressColor === 'green') return 'bg-emerald-500';
  if (props.progressColor === 'red') return 'bg-rose-500';
  return 'bg-blue-500';
});
</script>
