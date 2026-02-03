<template>
  <div class="flex items-center justify-center">
    <div class="relative" :style="{ width: sizeValue, height: sizeValue }">
      <!-- Background circle -->
      <svg class="transform -rotate-90" :width="sizeValue" :height="sizeValue">
        <circle
          :cx="radius + strokeWidth"
          :cy="radius + strokeWidth"
          :r="radius"
          :stroke-width="strokeWidth"
          fill="none"
          class="stroke-white/10"
        />
        <!-- Progress circle -->
        <circle
          :cx="radius + strokeWidth"
          :cy="radius + strokeWidth"
          :r="radius"
          :stroke-width="strokeWidth"
          fill="none"
          :class="scoreColor.text"
          :stroke-dasharray="circumference"
          :stroke-dashoffset="dashOffset"
          class="transition-all duration-1000 ease-out"
          stroke-linecap="round"
        />
      </svg>
      
      <!-- Score text -->
      <div class="absolute inset-0 flex flex-col items-center justify-center">
        <span :class="[textSizeClass, 'font-bold', scoreColor.text]">
          {{ Math.round(score) }}
        </span>
        <span class="text-xs text-gray-400" v-if="showLabel">Score</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { getScoreColor } from '../utils/colors'

const props = defineProps({
  score: {
    type: Number,
    required: true,
    default: 0
  },
  size: {
    type: String,
    default: 'md', // sm, md, lg
    validator: (value) => ['sm', 'md', 'lg'].includes(value)
  },
  showLabel: {
    type: Boolean,
    default: false
  }
})

// Size configurations
const sizeConfig = {
  sm: { size: 60, stroke: 4, text: 'text-lg' },
  md: { size: 80, stroke: 6, text: 'text-2xl' },
  lg: { size: 120, stroke: 8, text: 'text-4xl' }
}

const config = computed(() => sizeConfig[props.size])
const sizeValue = computed(() => `${config.value.size}px`)
const strokeWidth = computed(() => config.value.stroke)
const textSizeClass = computed(() => config.value.text)
const radius = computed(() => (config.value.size - config.value.stroke * 2) / 2)
const circumference = computed(() => 2 * Math.PI * radius.value)
const dashOffset = computed(() => {
  const progress = Math.min(Math.max(props.score, 0), 100) / 100
  return circumference.value * (1 - progress)
})

const scoreColor = computed(() => getScoreColor(props.score))
</script>
