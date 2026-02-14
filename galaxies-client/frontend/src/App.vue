<script setup lang="ts">
import { reactive, onMounted, ref } from 'vue'
import { 
  GetShipState, GetAvailableContracts, Travel, AcceptJob, 
  DropJob, GetPlanets, Refuel, GetModules, BuyModule 
} from '../wailsjs/go/main/App'

import StarMap from './components/StarMap.vue'
import ShipStatus from './components/ShipStatus.vue'
import OperationsPanel from './components/OperationsPanel.vue'
import ChatLog from './components/ChatLog.vue'

// --- STATE ---
const state = reactive({
  ship: {} as any,
  jobs: [] as any[],
  planets: [] as any[],
  modules: [] as any[], 
  chatMessages: [] as any[],
  loading: false
})

const chatCollapsed = ref(false)

// --- WEBSOCKET & LOGIC (Kept same as before) ---
let socket: WebSocket | null = null
function connectWS() {
    socket = new WebSocket("ws://localhost:8081/ws")
    socket.onmessage = (event) => {
        try {
            const msg = JSON.parse(event.data);
            if (msg.type === "chat_global") pushMessage(msg);
            if (msg.type === "market_pulse") handleMarketPulse(msg);
        } catch (e) { console.error(e); }
    }
    socket.onclose = () => setTimeout(connectWS, 5000)
}

function pushMessage(msg: any) {
    const msgWithTime = { ...msg, timestamp: new Date().toLocaleTimeString([], {hour: '2-digit', minute:'2-digit'}) }
    state.chatMessages = [...state.chatMessages, msgWithTime];
    if (state.chatMessages.length > 50) state.chatMessages.shift();
}

function handleMarketPulse(msg: any) {
    const updatedKeys = msg.updated_planets || [];
    if (updatedKeys.includes(state.ship.location_key)) {
        GetAvailableContracts().then(jobs => state.jobs = jobs);
    }
    const names = updatedKeys.map((key: string) => {
        const p = state.planets.find(x => x.key === key);
        return p ? p.name : key;
    });
    if (names.length > 0) {
        pushMessage({ type: "system_alert", sender: "NET_UPLINK", payload: `MARKET REFRESH: ${names.join(', ')}` });
    }
}

function handleSendMessage(text: string) {
    if (!socket || socket.readyState !== WebSocket.OPEN) return
    const msg = { type: "chat_global", sender: state.ship.name || 'PILOT', payload: text }
    socket.send(JSON.stringify(msg))
}

async function refreshAll() {
  state.loading = true
  try {
    state.ship = await GetShipState()
    state.jobs = await GetAvailableContracts()
    if (state.planets.length === 0) state.planets = await GetPlanets() || []
    state.modules = await GetModules() || []
  } catch (e) { console.error(e) } 
  finally { state.loading = false }
}

const actions = {
  travel: async (dest: string) => { state.ship = await Travel(dest); await refreshAll() },
  accept: async (id: string) => { await AcceptJob(id); await refreshAll() },
  drop: async (id: string) => { await DropJob(id); await refreshAll() },
  refuel: async () => { await Refuel(); await refreshAll() },
  buyModule: async (key: string) => { await BuyModule(key); await refreshAll() }
}

onMounted(() => { refreshAll(); connectWS() })
</script>

<template>
  <div class="terminal-grid">
    <div class="scanline"></div>

    <div class="col-left">
      
      <div class="section-status">
        <ShipStatus 
          :ship="state.ship" 
          :planets="state.planets"
          :loading="state.loading"
          @refuel="actions.refuel" 
        />
      </div>

      <div class="section-ops">
        <OperationsPanel 
          :ship="state.ship"
          :jobs="state.jobs"
          :planets="state.planets"
          :modules="state.modules"
          @accept="actions.accept"
          @drop="actions.drop"
          @buy="actions.buyModule"
        />
      </div>

      <div class="section-comms" :class="{ 'collapsed': chatCollapsed }">
        <div class="panel-header" @click="chatCollapsed = !chatCollapsed">
          <span>:: COMMS_UV_LINK ::</span>
          <span class="toggle-icon">{{ chatCollapsed ? '▲' : '▼' }}</span>
        </div>
        <div class="comms-content" v-show="!chatCollapsed">
          <ChatLog 
            :messages="state.chatMessages" 
            :shipName="state.ship.name" 
            @sendMessage="handleSendMessage"
          />
        </div>
      </div>

    </div>

    <div class="col-right">
      <StarMap 
        :universe="state.planets"
        :currentLocation="state.ship.location_key"
        :ship="state.ship"
        @travel="actions.travel"
      />
    </div>

  </div>
</template>

<style>
/* FIXED LAYOUT - Prevents "Planets Disappearing" */
.terminal-grid {
  display: grid;
  grid-template-columns: 450px 1fr; /* Fixed left, Flex right */
  height: 100vh;
  width: 100vw;
  background: #050505;
  color: #00ff41;
  overflow: hidden; /* Critical */
}

.col-left {
  display: flex;
  flex-direction: column;
  border-right: 2px solid #004400;
  background: rgba(0,20,0,0.3);
  min-width: 0; /* Critical for Flexbox nesting */
  height: 100%;
  overflow: hidden;
}

.col-right {
  position: relative;
  background: #000;
  min-width: 0;
  height: 100%;
  overflow: hidden;
}

/* SECTIONS */
.section-status {
  flex-shrink: 0;
  border-bottom: 1px solid #004400;
}

.section-ops {
  flex: 1; /* Takes all remaining height */
  overflow-y: auto; /* SCROLLBAR FIX */
  border-bottom: 1px solid #004400;
  min-height: 0; /* Firefox fix */
}

.section-comms {
  height: 200px; /* Default Height */
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  transition: height 0.3s ease;
  border-top: 1px solid #004400;
}

.section-comms.collapsed {
  height: 30px; /* Header only */
}

.panel-header {
  background: #002200;
  color: #008f11;
  padding: 4px 10px;
  font-size: 0.75rem;
  letter-spacing: 2px;
  cursor: pointer;
  display: flex;
  justify-content: space-between;
}
.panel-header:hover { background: #003300; color: #fff; }

.comms-content {
  flex: 1;
  overflow: hidden;
  background: rgba(0,0,0,0.8);
}

.scanline {
  position: fixed; top: 0; left: 0; width: 100%; height: 100%;
  background: linear-gradient(to bottom, rgba(255,255,255,0), rgba(255,255,255,0) 50%, rgba(0,0,0,0.1) 50%, rgba(0,0,0,0.1));
  background-size: 100% 4px; pointer-events: none; z-index: 9999; opacity: 0.3;
}

/* SCROLLBARS */
::-webkit-scrollbar { width: 5px; }
::-webkit-scrollbar-track { background: #001100; }
::-webkit-scrollbar-thumb { background: #004400; }
::-webkit-scrollbar-thumb:hover { background: #008f11; }
</style>