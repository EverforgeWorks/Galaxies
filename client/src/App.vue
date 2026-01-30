<script setup>
import { ref, onMounted, computed } from 'vue'
import TerminalLayout from './components/TerminalLayout.vue'
import LoginView from './views/LoginView.vue'
import RegistrationView from './views/RegistrationView.vue'
import DashboardView from './views/DashboardView.vue'

// -- STATE --
const token = ref(localStorage.getItem('token') || '')
const status = ref('DISCONNECTED')
const logs = ref([])
const player = ref(null)
const currentStar = ref(null)
let socket = null

// -- COMPUTED --
// Logic: If no token -> Login. If token but no player data yet -> Loading (or Reg). If player data -> Dashboard.
const currentView = computed(() => {
  if (!token.value) return 'LOGIN'
  if (token.value && !player.value) return 'LOADING' // Or Registration if we add that flag
  return 'DASHBOARD'
})

// -- LOGIC --
const addLog = (msg) => logs.value.unshift(`${new Date().toLocaleTimeString()}: ${msg}`)

const connectWS = () => {
  if (!token.value || socket) return
  
  status.value = 'CONNECTING...'
  const protocol = window.location.protocol === 'https:' ? 'wss' : 'ws'
  const wsUrl = `${protocol}://${window.location.host}/ws?token=${token.value}`
  
  socket = new WebSocket(wsUrl)
  
  socket.onopen = () => {
    status.value = 'CONNECTED'
    addLog('Uplink established.')
  }
  
  socket.onmessage = (e) => {
    try {
      const msg = JSON.parse(e.data)
      handleMessage(msg)
    } catch (err) {
      console.error(err)
    }
  }
  
  socket.onclose = () => {
    status.value = 'DISCONNECTED'
    socket = null
  }
}

const handleMessage = (msg) => {
  switch (msg.type) {
    case 'PLAYER_UPDATE':
      player.value = msg.payload
      break
    case 'STAR_UPDATE':
      currentStar.value = msg.payload
      break
    case 'CHAT_MESSAGE':
      addLog(`[COMM] ${msg.payload.sender_name}: ${msg.payload.content}`)
      break
  }
}

const login = (provider) => {
  window.location.href = `/auth/${provider}`
}

// TODO: Wire this up to a backend message type "UPDATE_NAME" later
const handleRegistration = (name) => {
  addLog(`Requesting callsign: ${name}...`)
  // socket.send(JSON.stringify({ type: 'UPDATE_NAME', payload: { name } }))
}

// -- LIFECYCLE --
onMounted(() => {
  // Check for OAuth callback token
  const params = new URLSearchParams(window.location.search)
  const urlToken = params.get('token')
  
  if (urlToken) {
    token.value = urlToken
    localStorage.setItem('token', urlToken)
    window.history.replaceState({}, document.title, "/")
  }

  if (token.value) connectWS()
})
</script>

<template>
  <TerminalLayout :status="status">
    <LoginView v-if="currentView === 'LOGIN'" @login="login" />
    
    <div v-else-if="currentView === 'LOADING'" class="text-center mt-20 animate-pulse">
      ESTABLISHING NEURAL LINK...
    </div>
    
    <DashboardView 
      v-else 
      :player="player" 
      :star="currentStar" 
      :logs="logs" 
    />
  </TerminalLayout>
</template>
