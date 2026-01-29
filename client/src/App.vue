<script setup>
import { ref, onMounted, computed, onUnmounted } from 'vue'
import StarMap from './components/StarMap.vue'
import CommsPanel from './components/CommsPanel.vue'
import MarketPanel from './components/MarketPanel.vue'
import AdminConsole from './components/AdminConsole.vue'
import AdminPanel from './components/AdminPanel.vue'
import Cockpit from './components/Cockpit.vue'

// --- STATE ---
const token = ref(localStorage.getItem('token') || '')
const player = ref(null)
const nearbySystems = ref([])
const message = ref('')
const loading = ref(false)
const showAdminConsole = ref(false)

// Chat & Connection State
const isConnected = ref(false)
const chatMessages = ref([])
let socket = null

// New Pilot State
const isNewPilot = ref(false)
const starterShips = ref([])
const selectedShipIndex = ref(0)
const pilotName = ref('')

// View State
const viewMode = ref('COCKPIT') 
const universeSystems = ref([])

// --- AUTH HELPER ---
const getAuthHeaders = () => {
    return {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token.value}`
    }
}

// --- WEBSOCKET MANAGER ---
const connectWS = () => {
    if (!token.value) return

    const protocol = window.location.protocol === 'https:' ? 'wss' : 'ws'
    const wsUrl = `${protocol}://${window.location.host}/api/ws?token=${token.value}`
    
    socket = new WebSocket(wsUrl)

    socket.onopen = () => {
        console.log("UPLINK ESTABLISHED")
        isConnected.value = true
        message.value = "SYSTEM ONLINE"
    }

    socket.onclose = () => {
        console.log("UPLINK LOST")
        isConnected.value = false
        setTimeout(connectWS, 3000)
    }

    socket.onmessage = (event) => {
        try {
            const msg = JSON.parse(event.data)
            
            // 1. AUTOMATIC STATE SYNC
            if (msg.type === 'PLAYER_UPDATE') {
                player.value = msg.data
                if (viewMode.value === 'COCKPIT' && player.value.current_system) {
                    scanSystems()
                }
            }
            // 2. CHAT & SYSTEM MESSAGES
            else if (['GLOBAL', 'SYSTEM', 'LOG'].includes(msg.type)) {
                 // Determine if this is a "Server Alert" or "Player Chat"
                 const isServerAlert = (!msg.sender || msg.sender === 'SYSTEM');

                 // If it's a Server Alert, show it in the HUD area
                 if (isServerAlert) {
                     message.value = msg.text || msg.content;
                 }

                 // ALWAYS push to chat log, preserving original fields (ship, system, sender)
                 // This fixes the "Unknown Ship" tooltip issue
                 chatMessages.value.push({
                    ...msg,
                    timestamp: msg.timestamp || (Date.now() / 1000)
                 })
            }
        } catch (e) { console.error("Packet Corrupt", e) }
    }
}

const sendChatMessage = (payload) => {
    if (socket && socket.readyState === 1) {
        socket.send(JSON.stringify(payload))
    }
}

// --- INITIALIZATION ---
onMounted(async () => {
  const urlParams = new URLSearchParams(window.location.search)
  const urlToken = urlParams.get('token')
  const urlNew = urlParams.get('new')

  if (urlToken) {
    token.value = urlToken
    localStorage.setItem('token', urlToken)
    window.history.replaceState({}, document.title, "/")
    if (urlNew === 'true') {
      isNewPilot.value = true
      pilotName.value = "Unknown Pilot"
      await fetchStarters()
    } else {
      await fetchPlayerData()
      connectWS() 
    }
  } else if (token.value) {
    await fetchPlayerData()
    connectWS() 
  }

  window.addEventListener('keydown', (e) => {
    if (e.key === '`' || e.key === '~') {
      if (player.value?.is_admin) {
        showAdminConsole.value = !showAdminConsole.value
      }
    }
  })
})

onUnmounted(() => {
    if (socket) socket.close()
})

// --- NAVIGATION ---
const switchView = async (mode) => {
  viewMode.value = mode
  if (mode === 'MAP') await fetchUniverse()
}

// --- API FETCHERS ---
const fetchPlayerData = async () => {
  try {
    const res = await fetch(`/api/me`, { headers: getAuthHeaders() })
    if (res.ok) {
      player.value = await res.json()
      if (player.value.current_system) {
        await scanSystems()
      } else {
        if (!player.value.ship) {
             isNewPilot.value = true
             await fetchStarters()
        }
      }
    } else {
      logout()
    }
  } catch (e) { console.error(e) }
}

const scanSystems = async () => {
  try {
    const res = await fetch(`/api/scan`, { headers: getAuthHeaders() })
    const data = await res.json()
    nearbySystems.value = data.systems || []
  } catch (e) { console.error(e) }
}

const fetchUniverse = async () => {
  try {
    const res = await fetch(`/api/map`, { headers: getAuthHeaders() })
    const data = await res.json()
    universeSystems.value = data.systems || []
  } catch (e) { console.error(e) }
}

// --- ACTIONS ---
const warp = async (targetId) => {
  if (!targetId) return
  message.value = "Calculations underway..."
  try {
    const res = await fetch('/api/warp', {
      method: 'POST',
      headers: getAuthHeaders(),
      body: JSON.stringify({ target_system_id: targetId })
    })
    if (res.ok) {
      message.value = "Jump Complete."
      // Server will push update via WS
    } else {
      const err = await res.json()
      message.value = "Jump Failed: " + err.error
    }
  } catch (e) { message.value = "Comms Down." }
}

const handleWarp = async (targetId) => {
  await warp(targetId)
  if (!message.value.includes('Failed')) {
      switchView('COCKPIT')
  }
}

const refuel = async () => {
  try {
    const res = await fetch('/api/refuel', {
      method: 'POST',
      headers: getAuthHeaders(),
      body: JSON.stringify({}) 
    })
    if (res.ok) {
      message.value = "Refueling Complete."
    } else {
      const err = await res.json()
      message.value = "Refuel Error: " + err.error
    }
  } catch (e) { message.value = "Connection Error." }
}

// --- DEBUG TOOLS ---
const debugRefuel = async () => {
  await fetch('/api/debug/refuel', { method: 'POST', headers: getAuthHeaders(), body: JSON.stringify({}) })
}

const debugCredits = async () => {
  await fetch('/api/debug/credits', { method: 'POST', headers: getAuthHeaders(), body: JSON.stringify({ amount: 10000 }) })
}

// --- ONBOARDING ---
const fetchStarters = async () => {
  const res = await fetch('/api/starters')
  const data = await res.json()
  starterShips.value = data.ships
  selectedShipIndex.value = 0
}

const commissionShip = async () => {
  loading.value = true
  const res = await fetch('/api/onboard', {
    method: 'POST',
    headers: getAuthHeaders(),
    body: JSON.stringify({ name: pilotName.value, ship_index: selectedShipIndex.value })
  })
  if (res.ok) {
    isNewPilot.value = false
    await fetchPlayerData()
    connectWS()
  } else {
    message.value = "Commission Failed."
  }
  loading.value = false
}

// --- AUTH & LOGOUT ---
const logout = async () => {
  if (token.value) {
    try {
        await fetch('/api/save', { method: 'POST', headers: getAuthHeaders(), body: JSON.stringify({}) })
    } catch(e) { console.log("Save failed on logout", e) }
  }
  
  if (socket) socket.close()
  token.value = ''
  player.value = null
  localStorage.removeItem('token')
  window.location.href = '/'
}

// --- COMPUTED ---
const selectedShip = computed(() => starterShips.value[selectedShipIndex.value] || {})

const sortedSystems = computed(() => {
  if (!player.value || !player.value.current_system) return []
  const cx = player.value.current_system.x
  const cy = player.value.current_system.y
  
  return [...nearbySystems.value].sort((a, b) => {
    const da = Math.sqrt(Math.pow(a.x - cx, 2) + Math.pow(a.y - cy, 2))
    const db = Math.sqrt(Math.pow(b.x - cx, 2) + Math.pow(b.y - cy, 2))
    return da - db
  })
})
</script>

<template>
  <div class="max-w-[1000px] mx-auto my-0 md:my-8 p-4 md:p-8 border-0 md:border md:border-[#30363d] bg-[#050505] relative overflow-hidden min-h-screen md:min-h-[80vh] flex flex-col font-mono text-[#00ff41]">
    <div class="absolute top-0 left-0 w-full h-[5px] bg-[#00ff41]/10 opacity-40 pointer-events-none z-[100] animate-scan"></div>
    
    <div v-if="!isConnected && player && !isNewPilot" class="absolute top-0 left-0 w-full h-full bg-black/85 z-[9999] flex flex-col justify-center items-center text-[#ff0055] font-bold tracking-widest">
        <div class="animate-blink mb-2">SIGNAL LOST</div>
        <div class="text-xs text-[#666]">ATTEMPTING RECONNECT...</div>
    </div>

    <AdminConsole v-if="showAdminConsole" :token="token" @close="showAdminConsole = false" />

    <div class="flex flex-col md:flex-row justify-between items-start md:items-center border-b-2 border-[#30363d] mb-8 pb-2">
        <h1 class="text-xl md:text-2xl m-0 p-0 tracking-widest border-b-0">Galaxies: Burn Rate</h1>
        
        <div v-if="player && !isNewPilot" class="mt-2 md:mt-0 flex flex-wrap gap-2 w-full md:w-auto justify-between">
            <button 
                @click="switchView('COCKPIT')" 
                class="flex-1 md:flex-none px-4 py-1 bg-transparent border border-[#30363d] text-[#666] font-bold hover:text-white hover:border-white transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                :class="{ 'bg-[#238636] text-white border-[#238636]': viewMode === 'COCKPIT' }"
                :disabled="loading"
            >Cockpit</button>
            
            <button 
                @click="switchView('MAP')" 
                class="flex-1 md:flex-none px-4 py-1 bg-transparent border border-[#30363d] text-[#666] font-bold hover:text-white hover:border-white transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                :class="{ 'bg-[#238636] text-white border-[#238636]': viewMode === 'MAP' }"
                :disabled="loading"
            >Star Map</button>
            
            <button 
                @click="switchView('MARKET')" 
                class="flex-1 md:flex-none px-4 py-1 bg-transparent border border-[#30363d] text-[#666] font-bold hover:text-white hover:border-white transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                :class="{ 'bg-[#238636] text-white border-[#238636]': viewMode === 'MARKET' }"
                :disabled="loading"
            >Market</button>
            
            <button v-if="player.is_admin" 
                    @click="switchView('ADMIN')" 
                    class="flex-1 md:flex-none px-4 py-1 bg-transparent border border-[#ff0055] text-[#ff0055] font-bold hover:bg-[#ff0055] hover:text-black transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                    :class="{ 'bg-[#ff0055] text-black': viewMode === 'ADMIN' }"
                    :disabled="loading">
                ADMIN
            </button>
            
            <button @click="logout" class="flex-1 md:flex-none px-4 py-1 bg-transparent border border-[#550000] text-[#cc6666] font-bold hover:bg-[#550000] hover:text-white transition-colors">Logout</button>
        </div>
    </div>
    
    <div v-if="!player && !isNewPilot" class="text-center py-16">
      <p class="animate-blink mb-4">IDENTITY_REQUIRED</p>
      
      <div class="flex justify-center gap-4 mt-4">
        <a href="/auth/github" class="flex items-center justify-center gap-2 px-6 py-3 bg-[#24292e] text-white font-bold border border-white/20 hover:bg-[#2f363d] hover:border-white transition-all">GITHUB</a>
        <a href="/auth/discord" class="flex items-center justify-center gap-2 px-6 py-3 bg-[#5865F2] text-white font-bold border border-white/20 hover:bg-[#7983f5] hover:border-white transition-all">DISCORD</a>
      </div>
    </div>

    <div v-else-if="isNewPilot" class="max-w-[600px] mx-auto w-full">
      <h2 class="text-xl border-b-2 border-[#30363d] pb-2 mb-4 tracking-widest">PILOT COMMISSIONING PROTOCOL</h2>
      <div class="mb-8">
        <label class="block mb-2">PILOT_CALLSIGN:</label>
        <input v-model="pilotName" class="w-full bg-black border border-[#00ff41] text-[#00ff41] p-2 text-lg font-mono focus:outline-none" />
      </div>
      
      <div class="mb-4">
        <h3 class="border-b-2 border-[#30363d] pb-2 mb-4 tracking-widest">SELECT VESSEL CLASS:</h3>
        <div class="flex gap-2 mb-4">
          <button 
            v-for="(ship, idx) in starterShips" 
            :key="idx" 
            @click="selectedShipIndex = idx" 
            class="flex-1 p-2 bg-[#111] border border-[#30363d] text-[#666] font-mono cursor-pointer hover:text-white hover:border-white"
            :class="{ 'border-[#00ff41] text-[#00ff41] bg-[#002200]': selectedShipIndex === idx }"
          >{{ ship.model_name }}</button>
        </div>
        
        <div v-if="selectedShip" class="bg-[#0a0a0a] p-4 border border-dashed border-[#333] leading-relaxed">
          <p><strong class="text-white">CLASS:</strong> {{ selectedShip.model_name }}</p>
          <p><strong class="text-white">HULL:</strong> {{ selectedShip.stats?.max_hull }}</p>
          <p><strong class="text-white">SHIELD:</strong> {{ selectedShip.stats?.max_shield }}</p>
          <p><strong class="text-white">RANGE:</strong> {{ selectedShip.stats?.jump_range }} LY</p>
        </div>
      </div>
      
      <button 
        @click="commissionShip" 
        class="w-full p-4 mt-4 bg-[#00ff41] text-black border-none text-xl font-bold font-mono cursor-pointer hover:bg-white disabled:opacity-50 disabled:cursor-not-allowed" 
        :disabled="loading"
      >
        {{ loading ? 'COMMISSIONING...' : 'ACCEPT COMMAND' }}
      </button>
      <p class="text-[#e5c07b] border border-[#e5c07b] p-2 mt-4" v-if="message">{{ message }}</p>
    </div>

    <div v-else class="flex flex-col flex-1 gap-8 relative">
        
        <div class="flex-1 min-h-[400px]">
            <div v-if="viewMode === 'COCKPIT'" class="h-full">
                <Cockpit 
                    :player="player"
                    :systems="nearbySystems"
                    :message="message"
                    :loading="loading"
                    @warp="handleWarp"
                    @refuel="refuel"
                    @debug-refuel="debugRefuel"
                    @debug-credits="debugCredits"
                />
            </div>

            <div v-else-if="viewMode === 'MARKET'" class="h-full border border-[#333]">
                <MarketPanel 
                    :token="token" 
                    :player="player" 
                    @refresh="fetchPlayerData" 
                />
            </div>

            <div v-else-if="viewMode === 'ADMIN'" class="flex-1 border border-[#ff0055] overflow-hidden flex flex-col bg-black h-full">
                <AdminPanel :token="token" />
            </div>

            <div v-else class="border border-[#333] p-2 bg-black">
                <StarMap 
                    :systems="universeSystems"
                    :current-system-id="player.current_system.id"
                    :ship-range="player.ship.stats.jump_range"
                    @warp="handleWarp"
                />
            </div>
        </div>

        <div class="mt-8 border-t border-[#333]">
             <CommsPanel 
                :messages="chatMessages"
                :is-connected="isConnected"
                @send-message="sendChatMessage" 
             />
        </div>

    </div>

    <div class="fixed bottom-1 left-2 text-[0.7rem] text-[#333] pointer-events-none z-[999]">Galaxies: Burn Rate version 0.0.1-alpha</div>
  </div>
</template>

<style>
@keyframes scan { 0% { top: 0%; } 100% { top: 100%; } }
.animate-scan { animation: scan 6s linear infinite; }

@keyframes blink { 50% { opacity: 0; } }
.animate-blink { animation: blink 1s step-end infinite; }
</style>
