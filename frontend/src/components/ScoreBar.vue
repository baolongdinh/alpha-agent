<template>
  <div class="mb-2">
    <div class="flex items-center justify-between text-sm mb-1">
      <span class="text-gray-400">{{ label }}</span>
      <span class="text-gray-300 font-medium">{{ score.toFixed(1) }} / {{ max }}</span>
    </div>
    <div class="h-2 bg-white/10 rounded-full overflow-hidden">
      <div 
        class="h-full transition-all duration-1000 ease-out"
        :class="barColor"
        :style="{ width: `${percentage}%` }"
      ></div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { getScoreGradient } from '../utils/colors'

const props = defineProps({
  label: {
    type: String,
    required: true
  },
  score: {
    type: Number,
    required: true
  },
  max: {
    type: Number,
    required: true
  }
})

const percentage = computed(() => {
  return Math.min((props.score / props.max) * 100, 100)
})

const normalizedScore = computed(() => {
  // Convert to 0-100 scale for color
  return (props.score / props.max) * 100
})

const barColor = computed(() => {
  return getScoreGradient(normalizedScore.value)
})
</script>
