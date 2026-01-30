<script setup>
const emit = defineEmits(['login'])
</script>

<template>
  <div class="flex flex-col items-center justify-center h-[60vh] space-y-8 animate-fade-in">
    <div class="text-center space-y-2">
      <p class="text-lg animate-pulse">>>> UPLINK REQUIRED</p>
      <p class="text-xs opacity-70">SECURE CONNECTION PROTOCOL INITIATED</p>
    </div>

    <div class="flex flex-col space-y-4 w-64">
      <button @click="$emit('login', 'github')" class="btn btn-outline border-terminal-green text-terminal-green hover:bg-terminal-green hover:text-black font-mono uppercase rounded-none">
        [ ACCESS VIA GITHUB ]
      </button>
      <button @click="$emit('login', 'discord')" class="btn btn-outline border-indigo-500 text-indigo-500 hover:bg-indigo-500 hover:text-white font-mono uppercase rounded-none">
        [ ACCESS VIA DISCORD ]
      </button>
    </div>
  </div>
</template>
