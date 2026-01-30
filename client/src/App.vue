<script setup>
import { ref, onMounted, computed } from 'vue'

// Import the components you created
import TerminalLayout from './components/TerminalLayout.vue'
import LoginView from './components/LoginView.vue'
import DashboardView from './components/DashboardView.vue'
// import RegistrationView from './components/RegistrationView.vue' // Enable when ready

// -- STATE --
const token = ref(localStorage.getItem('token') || '')
const status = ref('DISCONNECTED')
const logs = ref([])
const player = ref(null)
const currentStar = ref(null)
let socket = null

// -- VIEW LOGIC --
const currentView = computed(() => {
  if (!token.value) return 'LOGIN'
  if (token.value && !player.value) return 'LOADING' 
  return 'DASHBOARD'
})

// -- NETWORK LOGIC --
const addLog = (msg) => logs.value.unshift(`${new Date().toLocaleTimeString()}: ${msg}`)

const connectWS = () => {
  if (!token.value || socket) return
  
  status.value = 'CONNECTING...'
  const protocol = window.location.protocol === 'https:' ? 'wss' : 'ws'
  const wsUrl = `${protocol}://${window.location.host}/ws?token=${token.value}`
  
  socket = new WebSocket(wsUrl)
  
  socket.onopen = () => {
    status.value = 'CONNECTED'
    addLog('System Link Established.')
  }
  
  socket.onmessage = (e) => {
    try {
      const msg = JSON.parse(e.data)
      handleMessage(msg)
    } catch (err) {
      console.error("Protocol Error:", err)
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

// -- LIFECYCLE --
onMounted(() => {
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
    
    <LoginView 
      v-if="currentView === 'LOGIN'" 
      @login="login" 
    />
    
    <div v-else-if="currentView === 'LOADING'" class="text-center mt-20 animate-pulse text-terminal-green">
      > ESTABLISHING NEURAL LINK...
    </div>
    
    <DashboardView 
      v-else 
      :player="player" 
      :star="currentStar" 
      :logs="logs" 
    />

  </TerminalLayout>
</template>
