<script setup>
import { ref, computed, onMounted, watch, nextTick } from 'vue'

const props = defineProps({
  systems: Array,
  currentSystemId: String,
  shipRange: Number
})

const emit = defineEmits(['warp'])

const selectedSystem = ref(null)
const viewport = ref(null)

// --- MAP CONFIG & ZOOM ---
const GRID_TILES = 40
const BASE_CELL_SIZE = 60 
const ZOOM_LEVELS = [1.2, 0.8, 0.5, 0.3] 
const zoomIndex = ref(1) 

const currentScale = computed(() => ZOOM_LEVELS[zoomIndex.value])
const cellSize = computed(() => BASE_CELL_SIZE * currentScale.value)
const mapSize = computed(() => GRID_TILES * cellSize.value)
const centerOffset = computed(() => mapSize.value / 2)

// --- VISUAL GENERATION ---
const STAR_TYPES = [
  { color: '#9bb0ff', glow: 'rgba(155, 176, 255, 0.6)', size: 1.2 }, 
  { color: '#aabfff', glow: 'rgba(170, 191, 255, 0.6)', size: 1.0 }, 
  { color: '#cad7ff', glow: 'rgba(202, 215, 255, 0.6)', size: 1.0 }, 
  { color: '#f8f7ff', glow: 'rgba(248, 247, 255, 0.6)', size: 0.9 }, 
  { color: '#fff4ea', glow: 'rgba(255, 244, 234, 0.6)', size: 0.9 }, 
  { color: '#ffd2a1', glow: 'rgba(255, 210, 161, 0.6)', size: 0.8 }, 
  { color: '#ffcc6f', glow: 'rgba(255, 204, 111, 0.6)', size: 0.7 }  
]

const seededRandom = (seed) => {
  let hash = 0
  for (let i = 0; i < seed.length; i++) {
    hash = seed.charCodeAt(i) + ((hash << 5) - hash)
  }
  const x = Math.sin(hash) * 10000
  return x - Math.floor(x)
}

const visualSystems = computed(() => {
  return props.systems.map(sys => {
    const rngX = seededRandom(sys.id + 'x')
    const rngY = seededRandom(sys.id + 'y')
    const rngType = seededRandom(sys.id + 'type')
    
    // Deterministic offset within the grid cell
    const offsetX = 0.2 + (rngX * 0.6)
    const offsetY = 0.2 + (rngY * 0.6)
    
    const typeIndex = Math.floor(rngType * STAR_TYPES.length)
    const starVisuals = STAR_TYPES[typeIndex]

    return {
      ...sys,
      visual: {
        gridX: sys.x + offsetX, 
        gridY: sys.y + offsetY,
        color: starVisuals.color,
        glow: starVisuals.glow,
        baseSize: starVisuals.size
      }
    }
  })
})

const currentSystemObj = computed(() => 
  visualSystems.value.find(s => s.id === props.currentSystemId)
)

// Drag State
const isDragging = ref(false)
const startX = ref(0)
const startY = ref(0)
const scrollLeft = ref(0)
const scrollTop = ref(0)
const wasDragAction = ref(false)

// --- INITIALIZATION ---
const centerMap = () => {
  if (viewport.value && currentSystemObj.value) {
    const sys = currentSystemObj.value
    const pxX = (sys.visual.gridX * cellSize.value) + centerOffset.value
    const pxY = (sys.visual.gridY * cellSize.value) + centerOffset.value
    
    viewport.value.scrollLeft = pxX - (viewport.value.clientWidth / 2)
    viewport.value.scrollTop = pxY - (viewport.value.clientHeight / 2)
  }
}

onMounted(() => {
  nextTick(centerMap)
})

watch(zoomIndex, () => nextTick(centerMap))

// --- ZOOM HANDLERS ---
const zoomIn = () => {
  if (zoomIndex.value > 0) zoomIndex.value--
}

const zoomOut = () => {
  if (zoomIndex.value < ZOOM_LEVELS.length - 1) zoomIndex.value++
}

// --- DRAG HANDLERS ---
const startDrag = (e) => {
  isDragging.value = true
  wasDragAction.value = false
  startX.value = e.pageX - viewport.value.offsetLeft
  startY.value = e.pageY - viewport.value.offsetTop
  scrollLeft.value = viewport.value.scrollLeft
  scrollTop.value = viewport.value.scrollTop
}

const onDrag = (e) => {
  if (!isDragging.value) return
  e.preventDefault()
  
  const x = e.pageX - viewport.value.offsetLeft
  const y = e.pageY - viewport.value.offsetTop
  const walkX = (x - startX.value) * 1.5 
  const walkY = (y - startY.value) * 1.5

  viewport.value.scrollLeft = scrollLeft.value - walkX
  viewport.value.scrollTop = scrollTop.value - walkY
  
  if (Math.abs(walkX) > 5 || Math.abs(walkY) > 5) {
    wasDragAction.value = true
  }
}

const stopDrag = () => {
  isDragging.value = false
}

// --- LOGIC ---
const selectSystem = (sys) => {
  if (!wasDragAction.value) {
    selectedSystem.value = sys
  }
}

const getDistance = (target) => {
  const current = currentSystemObj.value
  if (!current || !target) return 0
  const dx = current.x - target.x 
  const dy = current.y - target.y
  return Math.sqrt(dx*dx + dy*dy).toFixed(1)
}

// Check if a system is within warp range (used for visual dimming)
const isInRange = (target) => {
  if (!currentSystemObj.value) return false
  if (target.id === props.currentSystemId) return true
  const dx = currentSystemObj.value.x - target.x
  const dy = currentSystemObj.value.y - target.y
  return (dx*dx + dy*dy) <= (props.shipRange * props.shipRange)
}

const canWarp = computed(() => {
  if (!selectedSystem.value) return false
  const dist = getDistance(selectedSystem.value)
  return dist <= props.shipRange
})
</script>

<template>
  <div class="flex h-[600px] border border-[#30363d] bg-black relative overflow-hidden select-none">
    
    <div class="absolute top-2.5 left-2.5 z-[300] flex flex-col gap-1">
      <button class="bg-[#111] border border-[#333] text-white w-8 h-8 cursor-pointer font-bold hover:bg-[#333] disabled:opacity-50 disabled:cursor-default" @click="zoomIn" :disabled="zoomIndex === 0">+</button>
      <button class="bg-[#111] border border-[#333] text-white w-8 h-8 cursor-pointer font-bold hover:bg-[#333] disabled:opacity-50 disabled:cursor-default" @click="zoomOut" :disabled="zoomIndex === ZOOM_LEVELS.length - 1">-</button>
    </div>

    <div 
      class="flex-1 overflow-auto cursor-grab active:cursor-grabbing bg-[#050505]" 
      ref="viewport"
      @mousedown="startDrag"
      @mouseleave="stopDrag"
      @mouseup="stopDrag"
      @mousemove="onDrag"
    >
      <div 
        class="relative transition-all duration-300 ease-out" 
        :style="{ 
          width: mapSize + 'px', 
          height: mapSize + 'px',
          backgroundImage: `linear-gradient(#1a1a1a 1px, transparent 1px), linear-gradient(90deg, #1a1a1a 1px, transparent 1px)`,
          backgroundSize: `${cellSize}px ${cellSize}px`
        }"
      >
        <div class="absolute w-5 h-5 border-l border-t border-[#333] -translate-x-px -translate-y-px pointer-events-none" :style="{ left: centerOffset + 'px', top: centerOffset + 'px' }"></div>
        
        <div 
          v-for="sys in visualSystems" 
          :key="sys.id"
          class="absolute -translate-x-1/2 -translate-y-1/2 cursor-pointer flex flex-col items-center justify-center transition-all duration-300 group min-w-[50px] min-h-[50px]" 
          :class="{ 'z-[100]': selectedSystem?.id === sys.id }"
          :style="{ 
            left: (sys.visual.gridX * cellSize) + centerOffset + 'px', 
            top: (sys.visual.gridY * cellSize) + centerOffset + 'px' 
          }"
          @click="selectSystem(sys)"
        >
          <div 
            class="rounded-full transition-all duration-300 shadow-star relative"
            :class="{
              'ring-1 ring-white': selectedSystem?.id === sys.id && sys.id !== currentSystemId,
              'opacity-40 grayscale': !isInRange(sys) && sys.id !== currentSystemId
            }"
            :style="{
              width: (12 * sys.visual.baseSize * currentScale) + 'px',
              height: (12 * sys.visual.baseSize * currentScale) + 'px',
              backgroundColor: sys.visual.color,
              boxShadow: isInRange(sys) || sys.id === currentSystemId ? `0 0 ${10 * currentScale}px ${sys.visual.glow}, 0 0 ${20 * currentScale}px ${sys.visual.glow}` : 'none'
            }"
          >
            <div 
                v-if="sys.id === currentSystemId" 
                class="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 rounded-full border border-[#00ff41] shadow-[0_0_8px_#00ff41]"
                style="width: 220%; height: 220%;"
            ></div>
          </div>

          <div 
            class="absolute top-full mt-1 text-[0.6rem] whitespace-nowrap bg-black/80 px-1 py-0.5 border border-[#333] text-[#aaa] pointer-events-none"
            :class="{ 
                'block text-[#00ff41] border-[#00ff41]': sys.id === currentSystemId,
                'block text-white': selectedSystem?.id === sys.id,
                'hidden group-hover:block': selectedSystem?.id !== sys.id && sys.id !== currentSystemId,
                'text-[#555]': !isInRange(sys)
            }"
            :style="{ fontSize: (0.7 * currentScale) + 'rem' }"
          >
            {{ sys.name }}
          </div>
        </div>
      </div>
    </div>

    <div class="w-[300px] bg-[#0d1117]/95 border-l border-[#30363d] p-6 flex flex-col z-[200] shadow-[-5px_0_15px_rgba(0,0,0,0.5)]">
      <h3 class="text-xl mb-4 text-[#00ff41] font-bold border-b border-[#30363d] pb-2">NAV_CHART</h3>
      
      <div v-if="selectedSystem" class="flex-1">
        <h4 class="text-lg text-white font-bold border-b border-[#333] pb-2 mb-4">{{ selectedSystem.name }}</h4>
        <div class="text-sm text-[#aaa] leading-relaxed">
          <p>COORDS: [{{ selectedSystem.x }}, {{ selectedSystem.y }}]</p>
          <p>DIST: <span :class="{ 'text-[#e06c75] font-bold': !canWarp, 'text-[#00ff41]': canWarp }">{{ getDistance(selectedSystem) }} LY</span></p>
          <hr class="border-t border-[#333] my-3">
          <p>GOV: <span class="text-[#e5c07b]">{{ selectedSystem.political }}</span></p>
          <p>ECO: <span class="text-[#e5c07b]">{{ selectedSystem.economic }}</span></p>
          <p>SOC: <span class="text-[#e5c07b]">{{ selectedSystem.social }}</span></p>
          <hr class="border-t border-[#333] my-3">
          <p class="mb-2">FACILITIES:</p>
          <ul class="space-y-1">
             <li v-if="selectedSystem.stats.has_shipyard"><span class="text-[#238636] mr-2">•</span>SHIPYARD</li>
             <li v-if="selectedSystem.stats.has_refueling"><span class="text-[#238636] mr-2">•</span>REFUELING</li>
             <li v-if="selectedSystem.stats.has_outfitter"><span class="text-[#238636] mr-2">•</span>OUTFITTER</li>
             <li v-if="selectedSystem.stats.has_black_market" class="text-[#e06c75] font-bold"><span class="text-[#e06c75] mr-2">•</span>BLACK MARKET</li>
             <li v-if="selectedSystem.stats.has_mission_board"><span class="text-[#238636] mr-2">•</span>MISSIONS</li>
             <li v-if="selectedSystem.stats.has_cantina"><span class="text-[#238636] mr-2">•</span>CANTINA</li>
             <li v-if="!selectedSystem.stats.has_refueling && !selectedSystem.stats.has_shipyard" class="text-[#333] text-xs italic">NONE DETECTED</li>
          </ul>
        </div>
        
        <div class="mt-8">
            <button 
            v-if="selectedSystem.id !== currentSystemId"
            class="w-full py-4 bg-[#238636] border border-[#2ea043] text-white font-bold font-mono tracking-widest cursor-pointer hover:bg-[#2ea043] disabled:bg-[#222] disabled:border-[#444] disabled:text-[#666] disabled:cursor-not-allowed" 
            :disabled="!canWarp"
            @click="$emit('warp', selectedSystem.id)"
            >
            {{ canWarp ? 'INITIATE JUMP' : 'OUT OF RANGE' }}
            </button>
            <div v-else class="text-center text-[#00ff41] border border-[#00ff41] bg-[#00ff41]/10 p-3 font-bold">
                CURRENT LOCATION
            </div>
        </div>
      </div>
      
      <div v-else class="text-center text-[#444] mt-8">
        <p class="mb-4">SELECT A STAR SYSTEM</p>
        <p class="text-xs text-[#333]">GRID: {{ cellSize.toFixed(0) }}px = 1 LY</p>
        <p class="text-xs text-[#333]">DRAG TO PAN</p>
      </div>
    </div>

  </div>
</template>

<style scoped>
/* Force layer promotion for performance on dragging */
.universe-grid {
    will-change: transform; 
}
</style>
