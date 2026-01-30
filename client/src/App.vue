<script setup>
/**
 * App.vue - Main Application Controller
 * * Handles global state, WebSocket connection, routing logic (Login -> Register -> Dashboard),
 * and high-level event coordination between the backend and UI components.
 */

import { ref, onMounted, computed } from 'vue'
import TerminalLayout from './components/TerminalLayout.vue'
import LoginView from './components/LoginView.vue'
import DashboardView from './components/DashboardView.vue'
import RegistrationView from './components/RegistrationView.vue'

// -- GLOBAL STATE --
const token = ref(localStorage.getItem('token') || '')
const status = ref('DISCONNECTED') // WebSocket status: CONNECTING, CONNECTED, DISCONNECTED
const logs = ref([])               // System log history
const player = ref(null)           // Current player entity
const currentStar = ref(null)      // Current star location entity
const isRegistering = ref(false)   // Flag to interrupt flow for new pilot registration
let socket = null                  // WebSocket instance

// -- VIEW ROUTING LOGIC --
// Determines which main component to display based on auth and data state.
const currentView = computed(() => {
  // 1. No token? User must log in via OAuth.
  if (!token.value) return 'LOGIN'
  
  // 2. Registration flag set? User needs to confirm callsign.
  if (isRegistering.value) return 'REGISTER'
  
  // 3. Token exists but no player data yet? Show loading/connecting state.
  if (token.value && !player.value) return 'LOADING' 
  
  // 4. All systems go.
  return 'DASHBOARD'
})

// -- ACTIONS --

/**
 * handleRegistration
 * Called when the user submits their callsign in RegistrationView.
 * Updates the local player state and clears the registration flag to reveal the dashboard.
 * @param {string} callsign - The chosen user name
 */
const handleRegistration = (callsign) => {
    // In a future version, this would send an 'UPDATE_PROFILE' message to the backend.
    // For now, we optimistically update the local state.
    if (player.value) {
        player.value.name = callsign
    }
    isRegistering.value = false
    addLog(`Identity Confirmed: Pilot ${callsign} active.`)
}

/**
 * logout
 * Clears authentication tokens, closes connections, and resets state.
 * Acts as an "Emergency Reset" if the client gets into a bad state.
 */
const logout = () => {
  localStorage.removeItem('token')
  token.value = ''
  player.value = null
  isRegistering.value = false
  
  if (socket) {
    socket.close()
    socket = null
  }
  
  // Hard reload to clear any memory artifacts
  window.location.href = '/'
}

/**
 * login
 * Redirects the user to the backend OAuth endpoint.
 * @param {string} provider - 'github' or 'discord'
 */
const login = (provider) => {
  window.location.href = `/auth/${provider}`
}

/**
 * addLog
 * Appends a message to the system log with a timestamp.
 */
const addLog = (msg) => {
    logs.value.unshift(`${new Date().toLocaleTimeString()}: ${msg}`)
}

/**
 * connectWS
 * Establishes the WebSocket connection using the JWT token.
 */
const connectWS = () => {
  if (!token.value || socket) return
  
  status.value = 'CONNECTING...'
  
  // Determine protocol (wss for https, ws for http)
  const protocol = window.location.protocol === 'https:' ? 'wss' : 'ws'
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
      console.error("Protocol Error: Failed to parse message", err)
    }
  }
  
  socket.onclose = (e) => {
    status.value = 'DISCONNECTED'
    socket = null
    // Code 1006 usually indicates abnormal closure (e.g., auth failure)
    if (e.code === 1006) {
        addLog('Connection Lost. Token may be expired.')
    } else {
        addLog('Connection terminated.')
    }
  }
}

/**
 * handleMessage
 * Routes incoming WebSocket messages to the appropriate state handlers.
 */
const handleMessage = (msg) => {
  switch (msg.type) {
    case 'PLAYER_UPDATE':
      // LOGIC: Check if this is a "new" user who needs to register.
      // We check if 'last_login' is very recent (less than 60 seconds ago).
      // Note: A more robust method would be a specific "is_new" flag from the backend.
      const lastLoginTime = new Date(msg.payload.last_login).getTime()
      const now = new Date().getTime()
      const isNewUser = (now - lastLoginTime) < 60000 // 60s threshold
      
      // If new and we haven't loaded player data yet, trigger registration flow
      if (isNewUser && !player.value) {
         isRegistering.value = true
      }
      
      player.value = msg.payload
      break
      
    case 'STAR_UPDATE':
      currentStar.value = msg.payload
      break
      
    case 'CHAT_MESSAGE':
      addLog(`[COMM] ${msg.payload.sender_name}: ${msg.payload.content}`)
      break
      
    default:
      console.warn("Unknown message type received:", msg.type)
  }
}

// -- LIFECYCLE --
onMounted(() => {
  // Check for JWT token in URL query params (returned from OAuth callback)
  const params = new URLSearchParams(window.location.search)
  const urlToken = params.get('token')
  
  if (urlToken) {
    token.value = urlToken
    localStorage.setItem('token', urlToken)
    // Remove token from URL for cleaner UI
    window.history.replaceState({}, document.title, "/")
  }

  // If we have a token (from URL or LocalStorage), connect immediately
  if (token.value) {
      connectWS()
  }
})
</script>

<template>
  <TerminalLayout :status="status">
    
    <LoginView 
      v-if="currentView === 'LOGIN'" 
      @login="login" 
    />
    
    <RegistrationView
      v-else-if="currentView === 'REGISTER'"
      :default-name="player?.name"
      @submit-name="handleRegistration"
    />
    
    <div v-else-if="currentView === 'LOADING'" class="flex flex-col items-center mt-20 space-y-4">
      <div class="animate-pulse text-terminal-green text-xl">
        > ESTABLISHING NEURAL LINK...
      </div>
      <button 
        @click="logout" 
        class="text-xs text-red-500 hover:text-red-400 underline cursor-pointer hover:bg-red-900/20 px-2 py-1 rounded"
      >
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
