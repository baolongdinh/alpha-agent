<template>
  <div class="space-y-2">
    <div v-for="i in rows" :key="i" class="glass rounded-lg p-4">
      <div class="grid grid-cols-12 gap-4 items-center">
        <!-- Rank -->
        <div class="col-span-1">
          <div class="shimmer h-6 w-12 rounded"></div>
        </div>

        <!-- Token Info -->
        <div class="col-span-3">
          <div class="shimmer h-6 w-24 rounded mb-2"></div>
          <div class="shimmer h-4 w-32 rounded"></div>
        </div>

        <!-- Score -->
        <div class="col-span-1 flex justify-center">
          <div class="shimmer h-16 w-16 rounded-full"></div>
        </div>

        <!-- Price -->
        <div class="col-span-2">
          <div class="shimmer h-5 w-24 rounded"></div>
        </div>

        <!-- Market Cap -->
        <div class="col-span-2">
          <div class="shimmer h-5 w-28 rounded"></div>
        </div>

        <!-- Volume -->
        <div class="col-span-2">
          <div class="shimmer h-5 w-24 rounded"></div>
        </div>

        <!-- Change -->
        <div class="col-span-1">
          <div class="shimmer h-5 w-16 rounded"></div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  rows: {
    type: Number,
    default: 8
  }
})
</script>
