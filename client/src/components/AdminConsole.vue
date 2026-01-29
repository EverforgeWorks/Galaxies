<script setup>
import { ref, onMounted } from 'vue'

const props = defineProps({
  token: String
})

const output = ref(['> ADMIN ACCESS GRANTED', '> SYSTEM READY'])
const command = ref('')
const activeTab = ref('SQL') // SQL, PLAYERS, SYSTEMS

const execute = async () => {
  if (!command.value.trim()) return
  
  const cmd = command.value
  output.value.push(`> ${cmd}`)
  command.value = ''

  try {
    const res = await fetch('/api/admin/exec', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ token: props.token, command: cmd })
    })
    const data = await res.json()
    
    if (data.error) {
      output.value.push(`ERROR: ${data.error}`)
    } else {
      output.value.push(JSON.stringify(data.result, null, 2))
    }
  } catch (e) {
    output.value.push(`COMM ERROR: ${e.message}`)
  }
}
</script>

<template>
  <div class="admin-overlay">
    <div class="terminal-window">
      <div class="term-header">
        <span>DEV_CONSOLE // {{ activeTab }}</span>
      </div>
      <div class="term-output">
        <div v-for="(line, i) in output" :key="i" class="line">{{ line }}</div>
      </div>
      <div class="term-input">
        <span class="prompt">$</span>
        <input 
          v-model="command" 
          @keyup.enter="execute" 
          placeholder="ENTER SQL OR COMMAND..." 
          autofocus 
        />
      </div>
    </div>
  </div>
</template>

<style scoped>
.admin-overlay {
  position: fixed; top: 0; left: 0; width: 100vw; height: 100vh;
  background: rgba(0, 0, 0, 0.8); z-index: 9999;
  display: flex; justify-content: center; align-items: center;
}
.terminal-window {
  width: 80%; height: 80%; background: #111; border: 2px solid #ff0055;
  display: flex; flex-direction: column; font-family: monospace;
  box-shadow: 0 0 20px rgba(255, 0, 85, 0.2);
}
.term-header {
  background: #ff0055; color: #000; padding: 5px 10px; font-weight: bold;
}
.term-output {
  flex: 1; padding: 10px; overflow-y: auto; color: #ddd; white-space: pre-wrap;
}
.term-input {
  display: flex; border-top: 1px solid #333; padding: 10px;
}
.prompt { color: #ff0055; margin-right: 10px; }
input {
  flex: 1; background: transparent; border: none; color: #fff; 
  font-family: monospace; outline: none;
}
</style>