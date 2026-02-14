<script setup lang="ts">
import { ref, onUpdated, nextTick } from 'vue'

const props = defineProps({
  messages: Array as () => any[],
  shipName: String
})

const emit = defineEmits(['sendMessage'])
const newMessage = ref('')
const scrollContainer = ref<HTMLElement | null>(null)

// Auto-scroll to bottom when new messages arrive
onUpdated(() => {
    nextTick(() => {
        if (scrollContainer.value) {
            scrollContainer.value.scrollTop = scrollContainer.value.scrollHeight
        }
    })
})

const send = () => {
    if (!newMessage.value.trim()) return
    emit('sendMessage', newMessage.value)
    newMessage.value = ''
}
</script>

<template>
  <div class="chat-wrapper">
    <div class="message-list" ref="scrollContainer">
      <div v-for="(msg, i) in messages" :key="i" class="chat-row">
        
        <template v-if="msg.type === 'system_alert'">
            <div class="sys-msg">
                <span class="timestamp">[{{ msg.timestamp }}]</span>
                <span class="sys-sender">:: {{ msg.sender }} ::</span>
                <span class="sys-text">{{ msg.payload }}</span>
            </div>
        </template>

        <template v-else>
            <div>
                <span class="timestamp">[{{ msg.timestamp }}]</span>
                <span class="sender" :class="{ 'is-me': msg.sender === shipName }">{{ msg.sender }}:</span>
                <span class="text">{{ msg.payload }}</span>
            </div>
        </template>

      </div>
      <div v-if="!messages?.length" class="empty-chat">-- NO SIGNALS DETECTED --</div>
    </div>

    <div class="input-area">
      <input 
        v-model="newMessage" 
        @keyup.enter="send"
        placeholder="ENTER COMMS..."
        maxlength="140"
      />
      <button @click="send">SEND</button>
    </div>
  </div>
</template>

<style scoped>
.chat-wrapper { display: flex; flex-direction: column; height: 100%; background: rgba(0, 5, 0, 0.9); font-family: 'Courier New', monospace; }

.message-list {
  flex: 1;
  overflow-y: auto;
  padding: 10px;
  display: flex;
  flex-direction: column;
  gap: 6px;
  font-size: 0.8rem;
}

.chat-row { line-height: 1.2; word-break: break-all; }

/* STANDARD CHAT */
.timestamp { color: #004400; margin-right: 6px; font-size: 0.7rem; }
.sender { color: #008f11; font-weight: bold; margin-right: 6px; }
.sender.is-me { color: #fff; }
.text { color: #00ff41; }

/* SYSTEM MESSAGE STYLES */
.sys-msg { color: #00ffff; font-style: italic; border-top: 1px dashed #004400; padding-top: 4px; margin-top: 2px; }
.sys-sender { color: #008f11; font-weight: bold; margin-right: 6px; }
.sys-text { color: #aaffff; text-shadow: 0 0 2px #00ffff; }

.input-area {
  display: flex;
  padding: 5px;
  border-top: 1px solid #004400;
  background: #000;
}

input {
  flex: 1;
  background: transparent;
  border: none;
  color: #00ff41;
  font-family: 'Courier New', monospace;
  outline: none;
  padding: 5px;
  font-size: 0.8rem;
}

button {
  background: #004400;
  color: #00ff41;
  border: 1px solid #008f11;
  padding: 0 15px;
  cursor: pointer;
  font-family: 'Courier New';
  font-weight: bold;
  font-size: 0.8rem;
}
button:hover { background: #00ff41; color: #000; }

.empty-chat { text-align: center; color: #004400; margin-top: 20px; font-style: italic; }
</style>