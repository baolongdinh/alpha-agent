<template>
  <div class="sentiment-gauge flex flex-col items-center">
    <div class="relative w-24 h-12 overflow-hidden">
      <!-- Gauge Background (Semi-circle) -->
      <svg class="w-24 h-24 transform -rotate-90" viewBox="00 100 100">
        <circle
          cx="50"
          cy="50"
          r="45"
          fill="none"
          stroke="rgba(255,255,255,0.1)"
          stroke-width="8"
          stroke-dasharray="141.37"
          stroke-dashoffset="141.37"
          stroke-linecap="round"
        />
        <!-- Gauge Fill -->
        <circle
          cx="50"
          cy="50"
          r="45"
          fill="none"
          :stroke="gaugeColor"
          stroke-width="8"
          :stroke-dasharray="141.37"
          :stroke-dashoffset="dashOffset"
          stroke-linecap="round"
          class="transition-all duration-1000 ease-out"
        />
      </svg>
      
      <!-- Center Needle or Indicator -->
      <div 
        class="absolute bottom-0 left-1/2 w-1 h-8 bg-white origin-bottom transition-transform duration-1000 ease-out"
        :style="{ transform: `translateX(-50%) rotate(${needleRotation}deg)` }"
      ></div>
    </div>
    <div class="mt-2 text-sm font-medium" :class="sentimentClass">
      {{ sentimentLabel }} ({{ (score * 100).toFixed(0) }}%)
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue';

const props = defineProps({
  score: {
    type: Number,
    required: true, // -1 to 1
    default: 0
  }
});

// Normalize score from [-1, 1] to [0, 1] for the gauge
const normalizedScore = computed(() => (props.score + 1) / 2);

// Semi-circle circumference is PI * r = 3.1415 * 45 = 141.37
const dashOffset = computed(() => 141.37 * (1 - normalizedScore.value));

const needleRotation = computed(() => (normalizedScore.value * 180) - 90);

const gaugeColor = computed(() => {
  if (props.score > 0.2) return '#10b981'; // Positive
  if (props.score < -0.2) return '#ef4444'; // Negative
  return '#94a3b8'; // Neutral
});

const sentimentLabel = computed(() => {
  if (props.score > 0.5) return 'Very Bullish';
  if (props.score > 0.1) return 'Bullish';
  if (props.score < -0.5) return 'Very Bearish';
  if (props.score < -0.1) return 'Bearish';
  return 'Neutral';
});

const sentimentClass = computed(() => {
  if (props.score > 0.1) return 'sentiment-positive';
  if (props.score < -0.1) return 'sentiment-negative';
  return 'sentiment-neutral';
});
</script>

<style scoped>
.sentiment-gauge {
  user-select: none;
}
</style>
