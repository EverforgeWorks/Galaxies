<script setup>
import { ref, onMounted, computed } from 'vue'
import TerminalLayout from './components/TerminalLayout.vue'
import LoginView from './components/LoginView.vue'
import DashboardView from './components/DashboardView.vue'

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

// -- ACTIONS --
const logout = () => {
  localStorage.removeItem('token')
  token.value = ''
  player.value = null
  if (socket) {
    socket.close()
    socket = null
  }
  window.location.href = '/'
}

const addLog = (msg) => logs.value.unshift(`${new Date().toLocaleTimeString()}: ${msg}`)

const connectWS = () => {
  if (!token.value || socket) return
  
  status.value = 'CONNECTING...'
  const protocol = window.location.protocol === 'https:' ? 'wss' : 'ws'
  // Vite proxy targets port 8080, but in prod/docker we hit nginx at port 80
  const host = window.location.host 
  const wsUrl = `${protocol}://${host}/ws?token=${token.value}`
  
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
  
  socket.onclose = (e) => {
    status.value = 'DISCONNECTED'
    socket = null
    // If we get a 401 (Close Code 1006 usually in browser), it means bad token
    addLog('Connection Lost. Token may be expired.')
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
    
    <div v-else-if="currentView === 'LOADING'" class="flex flex-col items-center mt-20 space-y-4">
      <div class="animate-pulse text-terminal-green text-xl">
        > ESTABLISHING NEURAL LINK...
      </div>
      <button @click="logout" class="text-xs text-red-500 hover:text-red-400 underline cursor-pointer">
        [ ABORT / RESET UPLINK ]
      </button>
    </div>
    
    <DashboardView 
      v-else 
      :player="player" 
      :star="currentStar" 
      :logs="logs" 
    />

  </TerminalLayout>
</template>
