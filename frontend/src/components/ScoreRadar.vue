<template>
  <div class="radar-container relative">
    <Radar :data="chartData" :options="chartOptions" />
  </div>
</template>

<script setup>
import { computed } from "vue";
import {
  Chart as ChartJS,
  RadialLinearScale,
  PointElement,
  LineElement,
  Filler,
  Tooltip,
  Legend,
} from "chart.js";
import { Radar } from "vue-chartjs";

ChartJS.register(
  RadialLinearScale,
  PointElement,
  LineElement,
  Filler,
  Tooltip,
  Legend,
);

const props = defineProps({
  scores: {
    type: Object,
    required: true,
    // Expected format: { liquidity: 80, volume: 70, trend: 60, tvl: 90, market: 75, risk: 85 }
  },
});

const chartData = computed(() => {
  return {
    labels: [
      "LIQUIDITY", 
      "VOLUME", 
      "TREND", 
      "SOCIAL", 
      "MARKET", 
      "STABILITY"
    ],
    datasets: [
      {
        label: "Score",
        backgroundColor: "rgba(59, 130, 246, 0.4)",
        borderColor: "#3b82f6",
        borderWidth: 3,
        pointBackgroundColor: "#3b82f6",
        pointBorderColor: "#fff",
        pointHoverBackgroundColor: "#fff",
        pointHoverBorderColor: "#3b82f6",
        pointRadius: 4,
        pointHoverRadius: 6,
        data: [
          props.scores.liquidity || 0,
          props.scores.volume || 0,
          props.scores.trend || 0,
          props.scores.social || 0,
          props.scores.market || 0,
          props.scores.stability || 0,
        ],
        fill: true,
      },
      // Subtle Glow Layer
      {
        label: "Glow",
        backgroundColor: "rgba(59, 130, 246, 0.1)",
        borderColor: "transparent",
        data: [
          (props.scores.liquidity || 0) + 5,
          (props.scores.volume || 0) + 5,
          (props.scores.trend || 0) + 5,
          (props.scores.social || 0) + 5,
          (props.scores.market || 0) + 5,
          (props.scores.stability || 0) + 5,
        ],
        fill: true,
        pointRadius: 0,
      }
    ],
  };
});

const chartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  layout: {
    padding: 30
  },
  scales: {
    r: {
      angleLines: {
        color: "rgba(255, 255, 255, 0.15)",
        lineWidth: 1
      },
      grid: {
        color: "rgba(255, 255, 255, 0.1)",
        circular: true
      },
      pointLabels: {
        color: "#64748b",
        font: {
          size: 11,
          weight: 'bold',
          family: "'Inter', sans-serif"
        },
        padding: 25
      },
      ticks: {
        display: false,
        stepSize: 20
      },
      min: 0,
      max: 100,
      beginAtZero: true
    },
  },
  plugins: {
    legend: {
      display: false,
    },
    tooltip: {
      backgroundColor: '#1a1c24',
      titleColor: '#fff',
      bodyColor: '#94a3b8',
      borderColor: '#3b82f6',
      borderWidth: 1,
      padding: 10,
      displayColors: false
    }
  },
};
</script>

<style scoped>
.radar-container {
  height: 100%;
  width: 100%;
}
</style>
