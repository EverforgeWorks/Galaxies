<script setup>
import { computed } from 'vue'

const props = defineProps({
  player: Object,
  star: Object,
  logs: Array
})

// Helper to format credits like "1,000 CR"
const fmtMoney = (amount) => {
  return new Intl.NumberFormat('en-US').format(amount) + ' CR'
}
</script>

<template>
  <div class="h-full grid grid-cols-12 grid-rows-6 gap-4 p-4">
    
    <div class="col-span-12 md:col-span-4 row-span-3 border-2 border-green-500/50 p-4 relative bg-black/50">
      <div class="absolute -top-3 left-4 bg-black px-2 text-green-400 font-bold uppercase tracking-widest">
        Pilot_Data
      </div>
      <div class="flex flex-col space-y-4 mt-2 text-xl">
        <div class="flex justify-between border-b border-green-900 pb-1">
          <span class="text-green-700">ID:</span>
          <span>{{ player?.name || 'UNKNOWN' }}</span>
        </div>
        <div class="flex justify-between border-b border-green-900 pb-1">
          <span class="text-green-700">RANK:</span>
          <span>{{ player?.is_admin ? 'ADMIN' : 'PLT-1' }}</span>
        </div>
        <div class="flex justify-between border-b border-green-900 pb-1">
          <span class="text-green-700">CREDITS:</span>
          <span class="text-yellow-400">{{ fmtMoney(player?.credits || 0) }}</span>
        </div>
        <div class="flex justify-between">
          <span class="text-green-700">LOC:</span>
          <span>{{ star?.name || 'DEEP SPACE' }} [{{ star?.x }}, {{ star?.y }}]</span>
        </div>
      </div>
    </div>

    <div class="col-span-12 md:col-span-8 row-span-3 border-2 border-green-500/50 p-4 relative bg-black/50 flex flex-col items-center justify-center">
      <div class="absolute -top-3 left-4 bg-black px-2 text-green-400 font-bold uppercase tracking-widest">
        Nav_Computer
      </div>
      
      <div class="text-center space-y-2 opacity-80">
        <div class="text-4xl animate-pulse">Scanning Sector...</div>
        <div class="text-sm text-green-700">NO ADJACENT SYSTEMS DETECTED</div>
      </div>
    </div>

    <div class="col-span-12 row-span-3 border-2 border-green-500/50 p-4 relative bg-black/50 flex flex-col">
      <div class="absolute -top-3 left-4 bg-black px-2 text-green-400 font-bold uppercase tracking-widest">
        Sys_Logs
      </div>
      <div class="flex-1 overflow-y-auto font-mono text-sm space-y-1 mt-2">
        <div v-for="(log, i) in logs" :key="i" class="border-l-2 border-green-800 pl-2">
          <span class="opacity-50 mr-2">></span> {{ log }}
        </div>
        <div v-if="logs.length === 0" class="opacity-30 italic">No activity recorded.</div>
      </div>
    </div>

  </div>
</template>
