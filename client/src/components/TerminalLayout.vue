<script setup>
defineProps({
  status: String
})
</script>

<template>
  <div class="min-h-screen bg-black text-terminal-green font-mono p-4 selection:bg-terminal-green selection:text-black">
    <div class="pointer-events-none fixed inset-0 bg-[linear-gradient(rgba(18,16,16,0)_50%,rgba(0,0,0,0.25)_50%),linear-gradient(90deg,rgba(255,0,0,0.06),rgba(0,255,0,0.02),rgba(0,0,255,0.06))] z-50 bg-[length:100%_2px,3px_100%]"></div>
    
    <header class="flex justify-between items-center border-b border-terminal-green pb-2 mb-6 uppercase">
      <h1 class="text-2xl font-bold tracking-widest text-shadow-glow">GALAXIES: BURN RATE</h1>
      <div class="text-xs">
        NET_STATUS: <span :class="status === 'CONNECTED' ? 'text-terminal-green' : 'text-red-500'">{{ status }}</span>
      </div>
    </header>

    <main class="relative z-10">
      <slot></slot>
    </main>
  </div>
</template>

<style scoped>
.text-shadow-glow {
  text-shadow: 0 0 5px rgba(51, 255, 0, 0.5);
}
</style>
