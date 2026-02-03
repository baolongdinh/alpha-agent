<template>
  <div
    class="rating-badge inline-flex items-center justify-center font-bold rounded px-2 py-0.5 text-xs tracking-wider"
    :style="{ backgroundColor: bgColor, color: '#fff' }"
  >
    {{ grade }}
  </div>
</template>

<script setup>
import { computed } from "vue";

const props = defineProps({
  grade: {
    type: String,
    required: true,
    default: "D",
  },
});

const bgColor = computed(() => {
  const g = props.grade.toUpperCase();
  if (g.startsWith("A")) return "var(--rating-aaa)";
  if (g.startsWith("B")) return "var(--rating-b)";
  if (g === "C") return "var(--rating-c)";
  return "var(--rating-d)";
});
</script>

<style scoped>
.rating-badge {
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.2);
  min-width: 40px;
}
</style>
