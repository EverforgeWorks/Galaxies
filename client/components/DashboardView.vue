<script setup>
defineProps({
  player: Object,
  star: Object,
  logs: Array
})
</script>

<template>
  <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
    <div class="border border-terminal-green p-4 bg-black/50 min-h-[200px]">
      <h3 class="border-b border-terminal-green pb-1 mb-3 opacity-50">PILOT_DATA</h3>
      <div v-if="player">
        <div class="flex justify-between my-1"><span>ID:</span> <span>{{ player.name }}</span></div>
        <div class="flex justify-between my-1"><span>CR:</span> <span class="text-yellow-400">{{ player.credits }}</span></div>
        <div class="flex justify-between my-1"><span>RNK:</span> <span>{{ player.is_admin ? 'ADM' : 'PLT' }}</span></div>
      </div>
    </div>

    <div class="border border-terminal-green p-4 bg-black/50 min-h-[200px]">
      <h3 class="border-b border-terminal-green pb-1 mb-3 opacity-50">NAV_COMPUTER</h3>
      <div v-if="star" class="text-center">
        <div class="text-4xl font-bold my-4">{{ star.name }}</div>
        <div class="text-xs tracking-widest">Global Coords: [{{ star.x }}, {{ star.y }}]</div>
      </div>
      <div v-else class="text-center animate-pulse mt-8">CALCULATING JUMP...</div>
    </div>

    <div class="border border-terminal-green p-4 bg-black/50 h-[400px] flex flex-col">
      <h3 class="border-b border-terminal-green pb-1 mb-3 opacity-50">SYS_LOGS</h3>
      <div class="flex-1 overflow-y-auto font-light text-sm space-y-1">
        <div v-for="(log, i) in logs" :key="i" class="opacity-80">> {{ log }}</div>
      </div>
    </div>
  </div>
</template>
