<script setup>
import { ref, computed, nextTick, watch, onMounted, onUnmounted } from 'vue'

const props = defineProps({
  messages: Array,      // Received from App.vue (Shared State)
  isConnected: Boolean  // Received from App.vue
})

const emit = defineEmits(['send-message'])

// DATA
const inputContent = ref('')
const activeTab = ref('ALL')      // What messages we SEE
const transmitTarget = ref('GLOBAL') // What channel we SEND to
const isCollapsed = ref(false)
const logContainer = ref(null)
const showChannelMenu = ref(false) // Dropdown state

// TOOLTIP STATE
const tooltipData = ref(null) 
const tooltipStyle = ref({ top: '0px', left: '0px' })

// --- COMPUTED ---
const filteredMessages = computed(() => {
  if (activeTab.value === 'ALL') return props.messages
  return props.messages.filter(m => m.type === activeTab.value)
})

// Auto-scroll logic
watch(() => props.messages.length, () => {
    if (!isCollapsed.value) {
        nextTick(() => {
            if (logContainer.value) {
                logContainer.value.scrollTop = logContainer.value.scrollHeight
            }
        })
    }
})

// --- ACTIONS ---

// Tab switching now ONLY changes the view, not the target
const setTab = (tab) => {
  activeTab.value = tab
}

const toggleCollapse = () => {
  isCollapsed.value = !isCollapsed.value
}

// Channel Selector Logic
const toggleChannelMenu = () => {
    showChannelMenu.value = !showChannelMenu.value
}

const setTransmitChannel = (channel) => {
    transmitTarget.value = channel
    showChannelMenu.value = false
}

// Tooltip Logic
const showTooltip = (msg, event) => {
  if (msg.type === 'LOG' || !msg.sender) return
  const rect = event.target.getBoundingClientRect()
  tooltipData.value = {
    name: msg.sender,
    ship: msg.sender_ship || 'Unknown Signal',
    system: msg.sender_system || 'Unknown Sector'
  }
  tooltipStyle.value = {
    top: `${rect.top - 80}px`, 
    left: `${rect.left}px`
  }
  setTimeout(() => { tooltipData.value = null }, 3000)
}

const closeUI = () => {
  tooltipData.value = null
  showChannelMenu.value = false
}

const sendMessage = () => {
  if (!inputContent.value.trim()) return
  
  // Delegate sending to Parent (App.vue)
  emit('send-message', { 
      type: transmitTarget.value, 
      content: inputContent.value 
  })
  inputContent.value = ''
}

// Global click listener to close tooltips/menus
onMounted(() => { 
  document.addEventListener('click', closeUI)
})

onUnmounted(() => { 
  document.removeEventListener('click', closeUI)
})
</script>

<template>
  <div 
    class="flex flex-col relative bg-[#050505] border border-[#333] transition-all duration-300 h-[300px]"
    :class="{ 'h-[40px] overflow-hidden border-b border-[#333]': isCollapsed }"
    @click.stop
  >
    
    <div 
      v-if="tooltipData" 
      class="fixed bg-[#0d1117] border border-[#00ff41] p-2.5 z-[9999] shadow-[0_0_10px_rgba(0,255,65,0.2)] min-w-[150px] font-mono text-sm"
      :style="tooltipStyle"
    >
        <div class="border-b border-[#333] mb-[5px] pb-[2px] text-white font-bold">{{ tooltipData.name }}</div>
        <div class="text-[#888] mb-[2px] text-xs">SHIP: <span class="text-[#e5c07b] float-right">{{ tooltipData.ship }}</span></div>
        <div class="text-[#888] mb-[2px] text-xs">LOC: <span class="text-[#e5c07b] float-right">{{ tooltipData.system }}</span></div>
    </div>

    <div 
      class="flex justify-between items-center px-2.5 py-[5px] bg-[#111] border-b border-[#333] cursor-pointer h-[40px] hover:bg-[#1a1a1a]"
      @click="toggleCollapse"
    >
      <div class="flex items-center gap-2.5">
        <div 
          class="w-2 h-2 rounded-full bg-[#330000]"
          :class="{ 'bg-[#00ff41] shadow-[0_0_5px_#00ff41]': isConnected }"
        ></div>
        <span class="font-bold text-[#888] font-mono">COMMS_NET</span>
      </div>
      
      <div class="flex gap-[15px]" @click.stop>
        <button 
          v-for="tab in ['ALL', 'GLOBAL', 'SYSTEM']" 
          :key="tab"
          class="bg-transparent border-none text-[#555] cursor-pointer font-mono font-bold pb-[2px] hover:text-[#888]"
          :class="{ 'text-white border-b-2 border-[#00ff41]': activeTab === tab }"
          @click="setTab(tab)"
        >{{ tab }}</button>
      </div>
      
      <div class="text-[#555] text-xs">{{ isCollapsed ? '▲' : '▼' }}</div>
    </div>

    <div v-show="!isCollapsed" class="flex-1 flex flex-col overflow-hidden">
      <div class="flex-1 overflow-y-auto p-2.5 font-mono text-[0.85rem]" ref="logContainer">
        <div 
          v-for="(msg, idx) in filteredMessages" :key="idx" 
          class="mb-1 leading-snug break-words"
          :class="{ 
            'text-[#61afef]': msg.type === 'GLOBAL', 
            'text-[#e5c07b]': msg.type === 'SYSTEM',
            'text-[#555] italic': msg.type === 'LOG'
          }"
        >
          <span class="text-[#444] mr-2 text-xs">[{{ new Date(msg.timestamp * 1000).toLocaleTimeString([], {hour: '2-digit', minute:'2-digit'}) }}]</span>
          
          <span v-if="msg.type === 'LOG'" class="text-[#e06c75] mr-[5px] font-bold">[LOG]</span>
          <span v-else-if="msg.type === 'GLOBAL'" class="text-[#61afef] mr-[5px] font-bold">[G]</span>
          <span v-else-if="msg.type === 'SYSTEM'" class="text-[#e5c07b] mr-[5px] font-bold">[S]</span>

          <span 
            v-if="msg.type !== 'LOG'" 
            class="cursor-pointer underline decoration-dotted hover:text-white hover:decoration-solid" 
            :class="msg.type === 'GLOBAL' ? 'text-[#61afef]' : 'text-[#e5c07b]'"
            @click.stop="showTooltip(msg, $event)"
          >{{ msg.sender }}:</span>
          
          <span v-else class="text-[#555]">{{ msg.sender }}:</span> 
          <span class="text-[#ccc]">{{ msg.content || msg.text }}</span>
        </div>
      </div>

      <div class="flex items-center px-2.5 py-[5px] border-t border-[#333] bg-[#0a0a0a] relative">
        <div class="relative">
            <span 
              class="text-xs font-bold mr-2.5 min-w-[60px] font-mono cursor-pointer hover:text-white select-none"
              :class="transmitTarget === 'GLOBAL' ? 'text-[#61afef]' : 'text-[#e5c07b]'"
              @click.stop="toggleChannelMenu"
            >{{ transmitTarget }} &gt;</span>
            
            <div v-if="showChannelMenu" class="absolute bottom-full left-0 mb-1 bg-[#111] border border-[#333] z-50 min-w-[100px] shadow-lg">
                <div 
                    class="px-3 py-1 cursor-pointer hover:bg-[#222] text-[#61afef] text-xs font-bold"
                    @click="setTransmitChannel('GLOBAL')"
                >GLOBAL</div>
                <div 
                    class="px-3 py-1 cursor-pointer hover:bg-[#222] text-[#e5c07b] text-xs font-bold"
                    @click="setTransmitChannel('SYSTEM')"
                >SYSTEM</div>
            </div>
        </div>

        <input 
          v-model="inputContent" 
          @keyup.enter="sendMessage"
          placeholder="Transmit..."
          class="flex-1 p-[5px] text-[0.9rem] border-none bg-transparent text-white font-mono focus:outline-none"
        />
      </div>
    </div>
  </div>
</template>
