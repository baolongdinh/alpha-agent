<template>
  <div class="container mx-auto px-4 py-8">
    <div class="grid grid-cols-12 gap-8">
      <!-- Main Content -->
      <main class="col-span-12 xl:col-span-10 xl:col-start-2">
        <!-- Page Header -->
        <div class="mb-8 flex justify-between items-end border-b border-white/5 pb-6">
          <div>
            <h1 class="text-4xl font-black text-white mb-2 tracking-tighter">Market Pulse</h1>
            <p class="text-gray-500 text-sm font-medium">Proprietary AI scoring for cross-chain liquidity & token health.</p>
          </div>
          <div class="text-right hidden sm:block">
             <div class="text-[10px] text-gray-500 font-bold uppercase tracking-widest mb-1">Last Data Sync</div>
             <div class="text-xs text-white font-mono">{{ lastFetch ? lastFetch.toLocaleTimeString() : 'connecting...' }}</div>
          </div>
        </div>

        <TokenTable />
      </main>
    </div>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import TokenTable from '../components/TokenTable.vue'
import { useTokens } from '../composables/useTokens'

const { refresh: refreshTokens, lastFetch } = useTokens()

onMounted(() => {
  console.log('🏠 HomePage mounted')
  // Initial fetch handled by TokenTable component
})
</script>
