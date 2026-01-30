<script setup>
import { ref, onMounted } from 'vue'

const token = ref(localStorage.getItem('token') || '')
const status = ref('DISCONNECTED')
const logs = ref([])
let socket = null

const addLog = (msg) => {
  logs.value.push(`${new Date().toLocaleTimeString()}: ${msg}`)
}

const connectWS = () => {
  if (!token.value) return
  
  status.value = 'CONNECTING...'
  const protocol = window.location.protocol === 'https:' ? 'wss' : 'ws'
  // Pointing to your new Go backend route
  const wsUrl = `${protocol}://${window.location.host}/ws?token=${token.value}`
  
  socket = new WebSocket(wsUrl)
  
  socket.onopen = () => {
    status.value = 'CONNECTED'
    addLog('Uplink established.')
  }
  
  socket.onmessage = (e) => addLog(`Received: ${e.data}`)
  socket.onclose = () => {
    status.value = 'DISCONNECTED'
    addLog('Uplink lost.')
  }
}

onMounted(() => {
  // Capture token from OAuth redirect: http://localhost?token=XYZ
  const params = new URLSearchParams(window.location.search)
  const urlToken = params.get('token')
  
  if (urlToken) {
    token.value = urlToken
    localStorage.setItem('token', urlToken)
    window.history.replaceState({}, document.title, "/") // Clean URL
  }

  if (token.value) connectWS()
})

const login = (provider) => {
  window.location.href = `/auth/${provider}`
}
</script>

<template>
  <div class="p-8 max-w-2xl mx-auto font-mono">
    <h1 class="text-3xl font-bold mb-4 text-primary">GALAXIES_REBUILD_V1</h1>
    
    <div v-if="!token" class="space-x-4">
      <button @click="login('github')" class="btn btn-outline">Login with GitHub</button>
      <button @click="login('discord')" class="btn btn-outline btn-secondary">Login with Discord</button>
    </div>

    <div v-else class="space-y-4">
      <div class="stats shadow border border-base-300 w-full">
        <div class="stat">
          <div class="stat-title">Uplink Status</div>
          <div class="stat-value text-sm" :class="status === 'CONNECTED' ? 'text-success' : 'text-error'">
            {{ status }}
          </div>
        </div>
      </div>

      <div class="bg-black p-4 rounded border border-neutral h-64 overflow-y-auto">
        <div v-for="log in logs" :key="log" class="text-xs text-success mb-1">
          > {{ log }}
        </div>
      </div>
    </div>
  </div>
</template>
