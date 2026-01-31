<script setup>
import { ref } from 'vue'

const emit = defineEmits(['commission'])

const shipName = ref('')
const error = ref('')

function submitName() {
  const name = shipName.value.trim()
  if (name.length < 3) {
    error.value = "ERROR: MINIMUM LENGTH 3 CHARS"
    return
  }
  if (name.length > 25) {
    error.value = "ERROR: MAXIMUM LENGTH 25 CHARS"
    return
  }
  
  emit('commission', name)
}
</script>

<template>
  <div class="fixed inset-0 bg-black/90 z-50 flex items-center justify-center font-mono backdrop-blur-sm">
    <div class="border-2 border-green-500 bg-black p-8 max-w-md w-full shadow-[0_0_50px_rgba(34,197,94,0.2)]">
      
      <h2 class="text-2xl text-green-500 mb-2 uppercase tracking-widest text-center border-b-2 border-green-900 pb-4">
        Vessel Commission
      </h2>
      
      <p class="text-green-400 mb-6 text-sm leading-relaxed mt-4">
        <span class="text-green-700 block mb-2">> PENDING AUTHORIZATION...</span>
        Welcome, Commander. GalCom has issued you a standard <span class="text-white font-bold">Scout Class</span> vessel. 
        Registration protocols require a unique designation before launch clearance is granted.
      </p>

      <div class="mb-6">
        <label class="block text-green-700 text-xs uppercase mb-1">Enter Vessel Designation</label>
        <input 
          v-model="shipName"
          @keyup.enter="submitName"
          type="text" 
          placeholder="NAME YOUR SHIP..."
          class="w-full bg-black border border-green-800 text-green-400 p-3 focus:outline-none focus:border-green-400 focus:shadow-[0_0_10px_rgba(34,197,94,0.2)] placeholder-green-900 transition-all text-center uppercase tracking-wider"
          autofocus
        />
        <p v-if="error" class="text-red-500 text-xs mt-2 border-l-2 border-red-500 pl-2">{{ error }}</p>
      </div>

      <button 
        @click="submitName"
        class="w-full bg-green-900/20 hover:bg-green-500/20 text-green-400 border border-green-600 py-3 uppercase tracking-wider transition-all hover:shadow-[0_0_15px_rgba(34,197,94,0.4)] group"
      >
        <span class="group-hover:animate-pulse">Initialize Systems</span>
      </button>
    </div>
  </div>
</template>