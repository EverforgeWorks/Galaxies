<script setup>
import { ref, onMounted, computed, watch } from 'vue'

const props = defineProps({
  token: String,
  player: Object // <--- SOURCE OF TRUTH (Received from App.vue via WS)
})

const emit = defineEmits(['refresh'])

// --- STATE ---
const marketItems = ref([])
const loading = ref(true)
const transactionInProgress = ref(false)
const errorMsg = ref('')

// Filters
const searchQuery = ref('')
const selectedCategory = ref('ALL')
const expandedId = ref(null) 

// Transaction Input State: { "ItemName": Amount }
const tradeAmounts = ref({})

// --- AUTH HELPER ---
const getAuthHeaders = () => ({
  'Content-Type': 'application/json',
  'Authorization': `Bearer ${props.token}`
})

// --- COMPUTED HELPERS ---

const cargoUsed = computed(() => {
  if (!props.player?.ship?.cargo) return 0
  return props.player.ship.cargo.reduce((sum, item) => sum + item.quantity, 0)
})

const cargoMax = computed(() => props.player?.ship?.stats?.cargo_volume || 0)
const cargoSpace = computed(() => cargoMax.value - cargoUsed.value)

// System Modifiers
const sysStats = computed(() => props.player?.current_system?.stats || {})
const buyMult = computed(() => sysStats.value.market_buy_mult || 1.0)
const sellMult = computed(() => sysStats.value.market_sell_mult || 1.0)
const taxRate = computed(() => sysStats.value.tax_rate || 0.0)

// Categories (Union of Market and Cargo categories)
const categories = computed(() => {
  const cats = new Set()
  marketItems.value.forEach(i => cats.add(i.category))
  
  if (props.player?.ship?.cargo) {
    props.player.ship.cargo.forEach(i => cats.add(i.category))
  }
  return ['ALL', ...Array.from(cats).sort()]
})

// --- PRICING ---
const getBuyPrice = (base) => Math.ceil(base * buyMult.value * (1 + taxRate.value))
const getSellPrice = (base) => Math.floor(base * sellMult.value)

// --- FILTERING ---
const filterList = (list) => {
  if (!list) return []
  return list.filter(item => {
    const matchCat = selectedCategory.value === 'ALL' || item.category === selectedCategory.value
    const matchSearch = item.name.toLowerCase().includes(searchQuery.value.toLowerCase())
    return matchCat && matchSearch
  })
}

const filteredCargo = computed(() => filterList(props.player?.ship?.cargo))
const filteredMarket = computed(() => filterList(marketItems.value))

// --- ACTIONS ---

const fetchMarket = async () => {
  loading.value = true
  try {
    const res = await fetch('/api/market', {
      headers: getAuthHeaders()
    })
    const data = await res.json()
    marketItems.value = data.market || []
  } catch (e) { console.error(e) }
  loading.value = false
}

const toggleExpand = (id) => {
  if (expandedId.value === id) {
    expandedId.value = null
  } else {
    expandedId.value = id
    tradeAmounts.value = {} 
  }
}

// --- TRADE LOGIC ---

const adjustAmount = (name, delta, limit) => {
  const current = tradeAmounts.value[name] || 0
  const newVal = Math.max(0, Math.min(limit, current + delta))
  tradeAmounts.value[name] = newVal
}

const setMax = (name, limit) => {
  tradeAmounts.value[name] = limit
}

const getMaxBuy = (item) => {
  const price = getBuyPrice(item.base_value)
  if (price === 0) return 0
  const afford = Math.floor((props.player?.credits || 0) / price)
  const space = cargoSpace.value
  return Math.min(afford, space, item.quantity)
}

const executeTrade = async (endpoint, item, qty) => {
  if (qty <= 0) return
  transactionInProgress.value = true
  try {
    const res = await fetch(endpoint, {
      method: 'POST',
      headers: getAuthHeaders(),
      body: JSON.stringify({
        item_name: item.name,
        quantity: qty
      })
    })
    
    if (!res.ok) {
      const err = await res.json()
      errorMsg.value = err.error
      setTimeout(() => errorMsg.value = '', 3000)
    } else {
      tradeAmounts.value[item.name] = 0
      await fetchMarket() 
      emit('refresh') 
    }
  } catch(e) { console.error(e) }
  transactionInProgress.value = false
}

onMounted(fetchMarket)
</script>

<template>
  <div class="flex flex-col h-full bg-[#050505] text-[#ccc] font-mono relative">
    
    <div class="bg-[#111] border-b border-[#333] p-2.5">
      <div class="flex justify-between items-baseline mb-2.5">
        <h3 class="m-0 text-white font-bold">{{ player?.current_system?.name }} EXCHANGE</h3>
        <div class="font-bold">
          <span class="text-[#e5c07b] mr-4">{{ player?.credits?.toLocaleString() }} CR</span>
          <span :class="{ 'text-[#e06c75]': cargoUsed >= cargoMax }">CARGO: {{ cargoUsed }} / {{ cargoMax }}</span>
        </div>
      </div>
      
      <div class="flex gap-2.5 flex-wrap">
        <input 
          v-model="searchQuery" 
          placeholder="SEARCH COMMODITIES..." 
          class="bg-black border border-[#444] text-white py-1 px-2.5 font-mono flex-1 min-w-[150px] outline-none focus:border-[#00ff41]" 
        />
        <div class="flex gap-[5px] overflow-x-auto">
          <button 
            v-for="cat in categories" :key="cat"
            class="bg-[#222] border border-[#444] text-[#888] px-2 py-1 cursor-pointer whitespace-nowrap text-xs hover:text-white hover:border-white transition-colors"
            :class="{ 'border-[#e5c07b] text-[#e5c07b] bg-[#221a00]': selectedCategory === cat }"
            @click="selectedCategory = cat"
          >{{ cat }}</button>
        </div>
      </div>
    </div>

    <div class="flex flex-1 overflow-hidden">
      
      <div class="flex-1 flex flex-col border-r border-[#222]">
        <div class="bg-[#080808] text-center p-2 font-bold text-[#555] border-b border-[#222] text-xs">SHIP INVENTORY</div>
        <div class="flex-1 overflow-y-auto p-[5px]">
          <div v-if="filteredCargo.length === 0" class="p-5 text-center text-[#444] italic">NO ITEMS MATCH</div>
          
          <div 
            v-for="item in filteredCargo" :key="item.id" 
            class="border border-[#222] mb-1 bg-[#0a0a0a] transition-all duration-200 hover:border-[#444]"
            :class="{ 'border-[#666] bg-[#111]': expandedId === item.id }"
          >
            <div class="p-2 cursor-pointer flex justify-between items-center" @click="toggleExpand(item.id)">
              <div class="flex flex-col">
                <span class="font-bold text-[#eee]">{{ item.name }}</span>
                <span class="text-xs text-[#666]">x{{ item.quantity }}</span>
              </div>
              <div class="text-right">
                <span class="font-bold text-[#00ff41]">{{ getSellPrice(item.base_value) }} CR</span>
              </div>
            </div>

            <div v-if="expandedId === item.id" class="p-2.5 border-t border-dashed border-[#333] bg-[#0e0e0e] animate-slideDown">
              <div class="flex justify-between text-xs text-[#666] mb-2.5">
                <span>AVG COST: {{ item.avg_cost.toFixed(0) }}</span>
                <span>PROFIT: <span :class="getSellPrice(item.base_value) >= item.avg_cost ? 'text-[#00ff41]' : 'text-[#e06c75]'">
                  {{ (getSellPrice(item.base_value) - item.avg_cost).toFixed(0) }}
                </span></span>
              </div>
              
              <div class="flex items-center gap-[5px] mb-[5px]">
                <button class="bg-[#222] border border-[#444] text-white w-7 h-7 cursor-pointer flex items-center justify-center hover:bg-[#333]" @click="adjustAmount(item.name, -1, item.quantity)">-</button>
                <span class="flex-1 text-center bg-black border border-[#333] h-7 flex items-center justify-center font-bold">{{ tradeAmounts[item.name] || 0 }}</span>
                <button class="bg-[#222] border border-[#444] text-white w-7 h-7 cursor-pointer flex items-center justify-center hover:bg-[#333]" @click="adjustAmount(item.name, 1, item.quantity)">+</button>
                <button class="bg-[#222] border border-[#444] text-white h-7 px-2 cursor-pointer flex items-center justify-center hover:bg-[#333] text-xs w-auto" @click="setMax(item.name, item.quantity)">ALL</button>
              </div>

              <button 
                class="w-full p-2 font-bold font-mono border-none cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed bg-[#00ff41] text-black hover:bg-[#88ffaa]"
                :disabled="!tradeAmounts[item.name] || transactionInProgress"
                @click="executeTrade('/api/sell', item, tradeAmounts[item.name])"
              >
                SELL FOR {{ (tradeAmounts[item.name] || 0) * getSellPrice(item.base_value) }} CR
              </button>
            </div>
          </div>
        </div>
      </div>

      <div class="flex-1 flex flex-col">
        <div class="bg-[#080808] text-center p-2 font-bold text-[#555] border-b border-[#222] text-xs">LOCAL MARKET</div>
        <div class="flex-1 overflow-y-auto p-[5px]">
          <div v-if="filteredMarket.length === 0" class="p-5 text-center text-[#444] italic">NO COMMODITIES FOUND</div>

          <div 
            v-for="item in filteredMarket" :key="item.id" 
            class="border border-[#222] mb-1 bg-[#0a0a0a] transition-all duration-200 hover:border-[#444]"
            :class="{ 
              'border-[#666] bg-[#111]': expandedId === item.id, 
              'border-l-[3px] border-l-[#e06c75]': item.is_illegal 
            }"
          >
            <div class="p-2 cursor-pointer flex justify-between items-center" @click="toggleExpand(item.id)">
              <div class="flex flex-col">
                <span class="font-bold text-[#eee]">
                  {{ item.name }}
                  <span v-if="item.is_illegal">💀</span>
                </span>
                <span class="text-xs text-[#666]">STOCK: {{ item.quantity }}</span>
              </div>
              <div class="text-right">
                <span class="font-bold text-[#e06c75]">{{ getBuyPrice(item.base_value) }} CR</span>
              </div>
            </div>

            <div v-if="expandedId === item.id" class="p-2.5 border-t border-dashed border-[#333] bg-[#0e0e0e] animate-slideDown">
              <div class="flex justify-between text-xs text-[#666] mb-2.5">
                <span>BASE: {{ item.base_value }}</span>
                <span>TAX: {{ (taxRate * 100).toFixed(0) }}%</span>
                <span>MULT: x{{ buyMult.toFixed(2) }}</span>
              </div>

              <div class="flex items-center gap-[5px] mb-[5px]">
                <button class="bg-[#222] border border-[#444] text-white w-7 h-7 cursor-pointer flex items-center justify-center hover:bg-[#333]" @click="adjustAmount(item.name, -1, getMaxBuy(item))">-</button>
                <span class="flex-1 text-center bg-black border border-[#333] h-7 flex items-center justify-center font-bold">{{ tradeAmounts[item.name] || 0 }}</span>
                <button class="bg-[#222] border border-[#444] text-white w-7 h-7 cursor-pointer flex items-center justify-center hover:bg-[#333]" @click="adjustAmount(item.name, 1, getMaxBuy(item))">+</button>
                <button class="bg-[#222] border border-[#444] text-white h-7 px-2 cursor-pointer flex items-center justify-center hover:bg-[#333] text-xs w-auto" @click="setMax(item.name, getMaxBuy(item))">MAX</button>
              </div>

              <div class="text-[0.7rem] text-[#555] text-center mb-2.5">
                MAX BUY: {{ getMaxBuy(item) }} (Afford/Space/Stock)
              </div>

              <button 
                class="w-full p-2 font-bold font-mono border-none cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed bg-[#e06c75] text-black hover:bg-[#ff8892]"
                :disabled="!tradeAmounts[item.name] || transactionInProgress"
                @click="executeTrade('/api/buy', item, tradeAmounts[item.name])"
              >
                BUY FOR {{ (tradeAmounts[item.name] || 0) * getBuyPrice(item.base_value) }} CR
              </button>
            </div>
          </div>
        </div>
      </div>

    </div>

    <div v-if="errorMsg" class="absolute bottom-5 left-1/2 -translate-x-1/2 bg-[#500] text-white py-2.5 px-5 border border-[#f00] shadow-[0_0_10px_#000]">
      {{ errorMsg }}
    </div>
  </div>
</template>

<style scoped>
/* Keeping this animation here as it's cleaner than a tailwind config extension for a single component */
@keyframes slideDown { 
  from { opacity: 0; transform: translateY(-5px); } 
  to { opacity: 1; transform: translateY(0); } 
}
.animate-slideDown { 
  animation: slideDown 0.2s ease-out; 
}
</style>
