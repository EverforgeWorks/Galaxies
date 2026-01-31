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
import CommissionView from './components/CommissionView.vue' // NEW IMPORT

// -- GLOBAL STATE --
const token = ref(localStorage.getItem('token') || '')
const status = ref('DISCONNECTED') // WebSocket status
const logs = ref([])               // System log history
const player = ref(null)           // Current player entity
const currentStar = ref(null)      // Current star location entity
const isRegistering = ref(false)   // Flag to interrupt flow for new pilot registration
let socket = null                  // WebSocket instance

// -- VIEW ROUTING LOGIC --
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

// -- NEW: Check if ship needs naming --
// This runs only when we are in the 'DASHBOARD' view
const needsCommission = computed(() => {
  return player.value?.ship?.name === "Uncommissioned Scout"
})

// -- ACTIONS --

/**
 * handleRegistration (Pilot Name)
 * Updates the local player state and clears the registration flag.
 */
const handleRegistration = (callsign) => {
    if (player.value) {
        player.value.name = callsign
    }
    isRegistering.value = false
    addLog(`Identity Confirmed: Pilot ${callsign} active.`)
}

/**
 * handleCommission (Ship Name)
 * NEW: Sends the rename command to the server.
 */
const handleCommission = (name) => {
    if (!socket) return
    socket.send(JSON.stringify({
        type: "COMMISSION_SHIP",
        payload: name
    }))
    addLog(`Transmitting vessel registration: ${name}...`)
}

const logout = () => {
  localStorage.removeItem('token')
  token.value = ''
  player.value = null
  isRegistering.value = false
  
  if (socket) {
    socket.close()
    socket = null
  }
  window.location.href = '/'
}

const login = (provider) => {
  window.location.href = `/auth/${provider}`
}

const addLog = (msg) => {
    logs.value.unshift(`${new Date().toLocaleTimeString()}: ${msg}`)
}

const connectWS = () => {
  if (!token.value || socket) return
  
  status.value = 'CONNECTING...'
  
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
      // PRESERVED LOGIC: 60s threshold for new users
      const lastLoginTime = new Date(msg.payload.last_login).getTime()
      const now = new Date().getTime()
      const isNewUser = (now - lastLoginTime) < 60000 // 60s threshold
      
      if (isNewUser && !player.value) {
         isRegistering.value = true
      }
      
      player.value = msg.payload
      break


      // NEW: Handle the universe dump
    case 'UNIVERSE_STATE':
      universe.value = msg.payload
      // Optional: Log it (can be verbose)
      // addLog(`Nav System: ${msg.payload.length} sectors mapped.`)
      break
    
    // -- NEW HANDLERS --
    case 'SHIP_UPDATED':
      if (player.value) {
          player.value.ship = msg.payload
          addLog(`Vessel Commissioned: ${msg.payload.name} ready for launch.`)
      }
      break

    case 'ERROR':
        addLog(`[ERROR] ${msg.payload}`)
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
  const params = new URLSearchParams(window.location.search)
  const urlToken = params.get('token')
  
  if (urlToken) {
    token.value = urlToken
    localStorage.setItem('token', urlToken)
    window.history.replaceState({}, document.title, "/")
  }

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
    
    <div v-else class="h-full relative">
        
        <CommissionView 
            v-if="needsCommission" 
            @commission="handleCommission" 
        />
        
        <DashboardView 
            :player="player" 
            :star="currentStar" 
            :logs="logs" 
            :universe="universe"
        />
    </div>

  </TerminalLayout>
</template>