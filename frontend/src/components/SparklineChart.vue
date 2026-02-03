<template>
  <svg :width="width" :height="height" class="overflow-visible">
    <!-- Gradient Defs -->
    <defs>
      <linearGradient :id="'gradient-' + id" x1="0" x2="0" y1="0" y2="1">
        <stop offset="0%" :stop-color="color" stop-opacity="0.2" />
        <stop offset="100%" :stop-color="color" stop-opacity="0" />
      </linearGradient>
    </defs>
    
    <!-- Area Path -->
    <path
      :d="areaPath"
      :fill="'url(#gradient-' + id + ')'"
      stroke="none"
    />
    
    <!-- Line Path -->
    <path
      :d="linePath"
      :stroke="color"
      fill="none"
      stroke-width="1.5"
      stroke-linecap="round"
      stroke-linejoin="round"
    />
  </svg>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  data: {
    type: Array,
    required: true,
    default: () => []
  },
  width: {
    type: [Number, String],
    default: 100
  },
  height: {
    type: [Number, String],
    default: 40
  },
  isPositive: {
    type: Boolean,
    default: true
  }
})

const id = Math.random().toString(36).substr(2, 9)
const color = computed(() => props.isPositive ? '#4ade80' : '#f87171')

const points = computed(() => {
  if (!props.data.length) return []
  
  const min = Math.min(...props.data)
  const max = Math.max(...props.data)
  const range = max - min || 1
  
  const width = parseFloat(props.width)
  const height = parseFloat(props.height)
  const step = width / (props.data.length - 1)
  
  return props.data.map((val, i) => {
    const x = i * step
    const normalizedVal = (val - min) / range
    // Increase internal padding to 6px for better vertical balance
    const y = height - (normalizedVal * (height - 12)) - 6 
    return { x, y }
  })
})

const linePath = computed(() => {
  if (!points.value.length) return ''
  return `M ${points.value.map(p => `${p.x},${p.y}`).join(' L ')}`
})

const areaPath = computed(() => {
  if (!points.value.length) return ''
  const h = parseFloat(props.height)
  const first = points.value[0]
  const last = points.value[points.value.length - 1]
  return `M ${first.x},${h} L ${points.value.map(p => `${p.x},${p.y}`).join(' L ')} L ${last.x},${h} Z`
})
</script>
