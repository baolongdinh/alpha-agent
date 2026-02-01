<template>
  <div class="fixed inset-0 z-50 flex items-center justify-center p-4">
    <!-- Backdrop -->
    <div 
      class="absolute inset-0 bg-black/60 backdrop-blur-sm"
      @click="$emit('close')"
    ></div>
    
    <!-- Modal Content -->
    <div 
      class="relative bg-[#181a20] border border-[rgba(255,255,255,0.1)] rounded-xl w-full max-w-4xl max-h-[90vh] overflow-y-auto shadow-2xl animate-fade-in"
    >
      <button 
        class="absolute top-4 right-4 text-gray-400 hover:text-white"
        @click="$emit('close')"
      >
        ✕
      </button>
      
      <slot />
    </div>
  </div>
</template>

<script setup>
defineEmits(['close'])
</script>
