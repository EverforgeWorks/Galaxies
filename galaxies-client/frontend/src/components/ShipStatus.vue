<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps({ 
  ship: Object as () => any,
  planets: Array as () => any[],
  loading: Boolean
})

const emit = defineEmits(['refuel'])

const planetName = computed(() => {
    if (!props.planets || !props.ship) return "UNKNOWN"
    const p = props.planets.find(x => x.key === props.ship.location_key)
    return p ? p.name : props.ship.location_key
})

const fuelPct = computed(() => props.ship ? (props.ship.fuel / props.ship.max_fuel) * 100 : 0)
const fuelCost = computed(() => props.ship ? Math.floor((props.ship.max_fuel - props.ship.fuel) / 100 * 4) : 0)
</script>

<template>
  <div class="status-panel">
    <div class="row">
        <span class="ship-name">{{ ship?.name || 'NO_SIGNAL' }}</span>
        <span class="loc">@{{ planetName }}</span>
    </div>

    <div class="stats-grid">
        <div class="stat-cell">
            <span class="label">CREDITS</span>
            <span class="value val-cr">{{ ship?.credits?.toLocaleString() }}</span>
        </div>
        <div class="stat-cell">
            <span class="label">MASS</span>
            <span class="value">{{ ship?.base_mass }}kg</span>
        </div>
        <div class="stat-cell">
             <span class="label">FUEL ({{ (ship?.fuel/100).toFixed(0) }}%)</span>
             <button class="btn-refuel" :disabled="fuelPct > 99" @click="emit('refuel')">
                {{ fuelPct > 99 ? 'FULL' : `FILL -${fuelCost}` }}
             </button>
        </div>
    </div>
    
    <div class="fuel-track">
        <div class="fuel-bar" :style="{ width: fuelPct + '%' }" :class="{ 'crit': fuelPct < 20 }"></div>
    </div>
  </div>
</template>

<style scoped>
.status-panel { font-family: 'Courier New', monospace; font-size: 0.8rem; background: #001100; padding: 10px; }

.row { display: flex; justify-content: space-between; align-items: baseline; margin-bottom: 6px; }
.ship-name { font-weight: bold; color: #fff; font-size: 0.9rem; }
.loc { color: #00ff41; }

.stats-grid { display: grid; grid-template-columns: 1fr 1fr 1.2fr; gap: 10px; margin-bottom: 5px; }
.stat-cell { display: flex; flex-direction: column; }
.label { font-size: 0.65rem; color: #006600; margin-bottom: 2px; }
.value { font-size: 0.85rem; color: #00aa00; }
.val-cr { color: #fff; font-weight: bold; }

.btn-refuel {
    background: #004400; border: none; color: #00ff41; 
    font-size: 0.7rem; cursor: pointer; padding: 1px 4px;
}
.btn-refuel:hover:not(:disabled) { background: #00ff41; color: #000; }
.btn-refuel:disabled { opacity: 0.3; cursor: default; }

.fuel-track { height: 4px; background: #002200; margin-top: 4px; }
.fuel-bar { height: 100%; background: #00ff41; transition: width 0.5s; }
.fuel-bar.crit { background: #ff0000; }
</style>