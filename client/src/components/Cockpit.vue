<script setup>
import { ref } from 'vue'

const props = defineProps({
  player: Object,
  systems: Array, // Expects sorted list of systems
  message: String,
  loading: Boolean
})

const emit = defineEmits(['warp', 'refuel', 'debug-refuel', 'debug-credits'])

const selectedNavTarget = ref('')

const getDist = (sys) => {
  if (!props.player) return 0
  const cx = props.player.current_system.x
  const cy = props.player.current_system.y
  return Math.sqrt(Math.pow(sys.x - cx, 2) + Math.pow(sys.y - cy, 2)).toFixed(1)
}

const handleWarp = () => {
    emit('warp', selectedNavTarget.value)
    // Optional: clear selection after jump?
    // selectedNavTarget.value = '' 
}
</script>

<template>
  <div class="h-full font-mono text-[#00ff41]">
    <div class="mb-4">
        <div>
            <span class="text-[#888] font-bold">PILOT:</span> {{ player.name }} <span class="mx-2 text-[#333]">|</span>
            <span class="text-[#888] font-bold">CREDITS:</span> {{ player.credits }}
        </div>
    </div>
    
    <div class="grid grid-cols-1 md:grid-cols-2 gap-8 mt-8">
        <div class="border border-[#333] p-6 bg-[#050505]">
            <h2 class="text-xl border-b-2 border-[#30363d] pb-2 mb-4 tracking-widest text-[#00ff41]">Ship Status</h2>
            <p class="text-2xl font-bold text-[#00ff41]">{{ player.ship.name }}</p>
            <p class="text-sm text-[#888] mb-4">{{ player.ship.model_name }}</p>
            
            <div class="mb-4">
                <label class="text-[#00ff41]">FUEL ({{ Math.floor(player.ship.current_fuel) }}/{{ player.ship.stats.max_fuel }})</label>
                <div class="bg-[#111] h-2.5 w-full border border-[#333] mt-1">
                    <div class="h-full bg-[#e5c07b]" :style="{ width: (player.ship.current_fuel / player.ship.stats.max_fuel * 100) + '%' }"></div>
                </div>
            </div>

            <div class="mb-4">
                <label class="text-[#00ff41]">HULL ({{ Math.floor(player.ship.current_hull) }}/{{ player.ship.stats.max_hull }})</label>
                <div class="bg-[#111] h-2.5 w-full border border-[#333] mt-1">
                    <div class="h-full bg-[#e06c75]" :style="{ width: (player.ship.current_hull / player.ship.stats.max_hull * 100) + '%' }"></div>
                </div>
            </div>

            <div class="mb-4">
                <label class="text-[#00ff41]">SHIELD ({{ Math.floor(player.ship.current_shield) }}/{{ player.ship.stats.max_shield }})</label>
                <div class="bg-[#111] h-2.5 w-full border border-[#333] mt-1">
                    <div class="h-full bg-[#61afef]" :style="{ width: (player.ship.current_shield / player.ship.stats.max_shield * 100) + '%' }"></div>
                </div>
            </div>

            <div class="grid grid-cols-3 gap-2 mt-4 pt-4 border-t border-dotted border-[#333]">
                <div class="text-xs text-[#888]">
                    <strong class="text-white block mb-1">FITTING</strong>
                    <p>High Slots: {{ player.ship.stats.high_slots }}</p>
                    <p>Mid Slots: {{ player.ship.stats.mid_slots }}</p>
                    <p>Low Slots: {{ player.ship.stats.low_slots }}</p>
                    <p>Power: {{ player.ship.stats.max_power_grid }} MW</p>
                </div>
                <div class="text-xs text-[#888]">
                    <strong class="text-white block mb-1">CAPACITY</strong>
                    <p>Cargo: {{ player.ship.stats.cargo_volume }} units</p>
                    <p>Cabins: {{ player.ship.stats.passenger_cabins }}</p>
                    <p>Bunks: {{ player.ship.stats.crew_bunks }}</p>
                </div>
                <div class="text-xs text-[#888]">
                    <strong class="text-white block mb-1">COMBAT</strong>
                    <p>Regen: {{ player.ship.stats.shield_regen }}/s</p>
                    <p>Stealth: {{ player.ship.stats.stealth_rating }}</p>
                    <p>Dmg Bonus: +{{ (player.ship.stats.damage_bonus * 100 - 100).toFixed(0) }}%</p>
                </div>
            </div>

            <div class="mt-4">
                <button 
                    v-if="player.current_system.stats.has_refueling"
                    class="w-full bg-[#e5c07b] text-black font-bold py-2 px-4 cursor-pointer disabled:bg-[#333] disabled:text-[#666] disabled:cursor-not-allowed hover:bg-white transition-colors"
                    @click="emit('refuel')"
                    :disabled="player.ship.current_fuel >= player.ship.stats.max_fuel || loading"
                >
                    {{ player.ship.current_fuel >= player.ship.stats.max_fuel ? 'TANK FULL' : 'PURCHASE FUEL' }}
                </button>
                <div v-else class="text-[#333] italic text-center text-sm mt-2">NO REFUELING AVAILABLE</div>
            </div>
        </div>

        <div class="border border-[#333] p-6 bg-[#050505]">
            <h2 class="text-xl border-b-2 border-[#30363d] pb-2 mb-4 tracking-widest text-[#00ff41]">Navigation</h2>
            <p class="mb-4 text-[#00ff41]">SYSTEM: <strong>{{ player.current_system.name }}</strong></p>
            
            <div class="mb-8">
                <div class="flex justify-between border-b border-[#333] py-1"><span class="text-[#888]">GOV:</span> {{ player.current_system.political }}</div>
                <div class="flex justify-between border-b border-[#333] py-1"><span class="text-[#888]">ECO:</span> {{ player.current_system.economic }}</div>
                <div class="flex justify-between border-b border-[#333] py-1"><span class="text-[#888]">SOC:</span> {{ player.current_system.social }}</div>
            </div>

            <div class="mt-4">
                <div v-if="message" class="text-[#e5c07b] border border-[#e5c07b] p-2 mb-4">{{ message }}</div>
                
                <label class="block mb-2 font-bold text-[#00ff41]">JUMP TARGET SELECTION:</label>
                <select v-model="selectedNavTarget" class="w-full bg-black border border-[#00ff41] text-[#00ff41] p-2 text-lg font-mono focus:outline-none cursor-pointer mb-4">
                    <option value="" disabled>-- SELECT DESTINATION --</option>
                    <option v-for="sys in systems" :key="sys.id" :value="sys.id">
                        {{ sys.name }} ({{ getDist(sys) }} LY)
                    </option>
                </select>
                
                <button 
                    class="w-full p-4 bg-[#00ff41] text-black border-none text-xl font-bold font-mono cursor-pointer hover:bg-white disabled:opacity-50 disabled:cursor-not-allowed disabled:hover:bg-[#00ff41] transition-colors" 
                    @click="handleWarp"
                    :disabled="!selectedNavTarget || loading"
                >
                INITIATE JUMP
                </button>
            </div>
        </div>
    </div>

    <div class="mt-12 pt-4 border-t border-dashed border-[#333] text-right">
        <span class="text-[#888] text-sm mr-2">DEBUG TOOLS:</span>
        <button @click="emit('debug-refuel')" class="bg-[#220000] text-[#ff5555] border border-[#550000] font-mono cursor-pointer ml-2 px-2 py-0.5 hover:bg-[#550000] hover:text-white transition-colors">[CHEAT_FUEL]</button>
        <button @click="emit('debug-credits')" class="bg-[#220000] text-[#ff5555] border border-[#550000] font-mono cursor-pointer ml-2 px-2 py-0.5 hover:bg-[#550000] hover:text-white transition-colors">[CHEAT_CREDITS]</button>
    </div>
  </div>
</template>
