<script setup>
import { ref, computed,  onMounted, watch } from 'vue'

const props = defineProps({
  token: String
})

const activeTable = ref('players')
const tableData = ref([])
const loading = ref(false)
const error = ref('')

const tables = [
  { id: 'players', label: 'PLAYERS' },
  { id: 'ships', label: 'SHIPS' },
  { id: 'systems', label: 'SYSTEMS' },
  { id: 'market_listings', label: 'MARKET_STOCK' },
  { id: 'ship_inventory', label: 'PLAYER_CARGO' }
]

const fetchTable = async (tableName) => {
  loading.value = true
  error.value = ""
  try {
    const res = await fetch('/api/admin/exec', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ 
        token: props.token, 
        command: `SELECT * FROM ${tableName} LIMIT 100;` 
      })
    })
    const data = await res.json()
    if (data.error) {
      error.value = data.error
    } else {
      tableData.value = data.result || []
    }
  } catch (e) {
    error.value = "Failed to reach sector authority."
  }
  loading.value = false
}

// Dynamically get headers from the first object in the array
const headers = computed(() => {
  if (tableData.value.length === 0) return []
  return Object.keys(tableData.value[0])
})

const formatVal = (val) => {
    if (val === null) return 'NULL'
    if (typeof val === 'object') return JSON.stringify(val)
    return val
}

onMounted(() => fetchTable(activeTable.value))
watch(activeTable, (newVal) => fetchTable(newVal))
</script>

<template>
  <div class="admin-panel">
    <div class="admin-sidebar">
      <button 
        v-for="t in tables" :key="t.id"
        @click="activeTable = t.id"
        :class="{ active: activeTable === t.id }"
      >
        {{ t.label }}
      </button>
    </div>

    <div class="admin-content">
      <div v-if="loading" class="status">QUERYING DATABASE...</div>
      <div v-else-if="error" class="status error">{{ error }}</div>
      
      <div v-else class="table-wrapper">
        <table v-if="tableData.length">
          <thead>
            <tr>
              <th v-for="h in headers" :key="h">{{ h.toUpperCase() }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(row, idx) in tableData" :key="idx">
              <td v-for="h in headers" :key="h">
                <div class="cell-content">{{ formatVal(row[h]) }}</div>
              </td>
            </tr>
          </tbody>
        </table>
        <div v-else class="status">TABLE IS EMPTY</div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.admin-panel { display: flex; height: 100%; background: #050505; border: 1px solid #ff0055; overflow: hidden; }

.admin-sidebar { width: 150px; background: #111; border-right: 1px solid #333; display: flex; flex-direction: column; }
.admin-sidebar button {
  background: transparent; border: none; border-bottom: 1px solid #222;
  color: #666; padding: 12px; text-align: left; cursor: pointer;
  font-family: monospace; font-size: 0.75rem;
}
.admin-sidebar button:hover { background: #1a1a1a; color: #fff; }
.admin-sidebar button.active { color: #ff0055; border-left: 3px solid #ff0055; background: #1a050a; }

.admin-content { flex: 1; display: flex; flex-direction: column; overflow: hidden; position: relative; }
.status { padding: 20px; text-align: center; font-family: monospace; color: #ff0055; }
.error { color: #fff; background: #500; }

.table-wrapper { flex: 1; overflow: auto; }
table { width: 100%; border-collapse: collapse; font-family: monospace; font-size: 0.7rem; }
th { 
  position: sticky; top: 0; background: #222; color: #aaa; 
  padding: 8px; text-align: left; border-bottom: 2px solid #ff0055;
  white-space: nowrap;
}
td { border: 1px solid #222; padding: 4px 8px; color: #ccc; max-width: 300px; }
tr:hover { background: #111; }

.cell-content { 
    max-height: 50px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; 
}
</style>
