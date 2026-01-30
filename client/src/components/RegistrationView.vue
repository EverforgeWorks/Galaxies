<script setup>
import { ref } from 'vue'

const props = defineProps({
  defaultName: String
})

const emit = defineEmits(['submit-name'])
const callsign = ref(props.defaultName || '')

const submit = () => {
  if (callsign.value.length > 2) {
    emit('submit-name', callsign.value)
  }
}
</script>

<template>
  <div class="max-w-md mx-auto border border-terminal-green p-6 bg-black/80 mt-10">
    <h2 class="text-xl mb-4 border-b border-gray-800 pb-2">NEW PILOT REGISTRATION</h2>
    <p class="text-sm mb-6 text-gray-400">
      Identity verification incomplete. Please assign a permanent callsign for the GALAXIES register.
    </p>

    <div class="form-control w-full">
      <label class="label">
        <span class="label-text text-terminal-green">DESIRED CALLSIGN</span>
      </label>
      <input 
        v-model="callsign" 
        type="text" 
        placeholder="ENTER NAME..." 
        class="input input-bordered border-terminal-green bg-black text-terminal-green w-full rounded-none focus:outline-none focus:ring-1 focus:ring-terminal-green"
        @keyup.enter="submit"
      />
    </div>

    <button 
      @click="submit" 
      class="btn btn-block mt-6 bg-terminal-green text-black hover:bg-white hover:text-black rounded-none border-none"
    >
      CONFIRM IDENTITY
    </button>
  </div>
</template>
