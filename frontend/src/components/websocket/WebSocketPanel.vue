<template>
  <div class="ws-panel">
    <div class="ws-toolbar">
      <div class="ws-url-row">
        <n-input
          v-model:value="wsUrl"
          placeholder="ws://localhost:8080/ws"
          :disabled="connected"
          size="small"
          class="ws-url-input"
        />
        <n-button
          :type="connected ? 'error' : 'primary'"
          size="small"
          @click="toggleConnection"
          :loading="connecting"
        >
          {{ connected ? 'Disconnect' : 'Connect' }}
        </n-button>
        <span class="ws-status" :class="{ connected: connected }">
          <span class="ws-dot"></span>
          {{ connected ? 'Connected' : 'Disconnected' }}
        </span>
      </div>
    </div>

    <div class="ws-messages" ref="messagesRef">
      <div v-if="messages.length === 0" class="ws-empty">
        <n-empty description="No messages yet" size="small" />
      </div>
      <div
        v-for="(msg, idx) in messages"
        :key="idx"
        class="ws-message"
        :class="msg.type"
      >
        <div class="ws-msg-time">{{ msg.time }}</div>
        <div class="ws-msg-content">{{ msg.content }}</div>
      </div>
    </div>

    <div class="ws-input-row">
      <n-input
        v-model:value="inputMessage"
        placeholder="Type a message (Ctrl+Enter to send)"
        size="small"
        :disabled="!connected"
        @keydown="onKeydown"
      />
      <n-button
        size="small"
        type="primary"
        :disabled="!connected || !inputMessage.trim()"
        @click="sendMessage"
      >
        Send
      </n-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, nextTick, onUnmounted, watch } from 'vue'
import { NInput, NButton, NEmpty } from 'naive-ui'
import { EventsOn, EventsOff } from '../../../wailsjs/runtime/runtime'
import {
  WSConnect,
  WSSend,
  WSDisconnect,
} from '../../../wailsjs/go/main/App'

interface WSMessage {
  type: 'sent' | 'received' | 'system'
  content: string
  time: string
}

const wsUrl = ref('ws://localhost:8080/ws')
const connected = ref(false)
const connecting = ref(false)
const inputMessage = ref('')
const messages = ref<WSMessage[]>([])
const messagesRef = ref<HTMLElement | null>(null)

let unsubMessage: (() => void) | null = null
let unsubError: (() => void) | null = null
let unsubClosed: (() => void) | null = null
let unsubConnected: (() => void) | null = null

function nowTime() {
  const d = new Date()
  return d.toLocaleTimeString()
}

function addSystem(msg: string) {
  messages.value.push({ type: 'system', content: msg, time: nowTime() })
  scrollBottom()
}

function scrollBottom() {
  nextTick(() => {
    if (messagesRef.value) {
      messagesRef.value.scrollTop = messagesRef.value.scrollHeight
    }
  })
}

async function toggleConnection() {
  if (connected.value) {
    try {
      await WSDisconnect(wsUrl.value)
    } catch (_) {
      // ignore
    }
    connected.value = false
    addSystem('Disconnected')
   } else {
    connecting.value = true
    try {
      await WSConnect(wsUrl.value, '{}')
      connected.value = true
      addSystem('Connected to ' + wsUrl.value)
      registerListeners()
    } catch (err: any) {
      addSystem('Connection failed: ' + (err?.message || err))
    } finally {
      connecting.value = false
    }
  }
}

function registerListeners() {
  unsubMessage = EventsOn('ws:message', (url: string, data: string) => {
    if (url === wsUrl.value) {
      messages.value.push({ type: 'received', content: data, time: nowTime() })
      scrollBottom()
    }
  })
  unsubError = EventsOn('ws:error', (url: string, err: string) => {
    if (url === wsUrl.value) {
      messages.value.push({
        type: 'system',
        content: 'Error: ' + err,
        time: nowTime(),
      })
      scrollBottom()
    }
  })
  unsubClosed = EventsOn('ws:closed', (url: string) => {
    if (url === wsUrl.value) {
      connected.value = false
      messages.value.push({
        type: 'system',
        content: 'Connection closed',
        time: nowTime(),
      })
      scrollBottom()
      cleanupListeners()
    }
  })
  unsubConnected = EventsOn('ws:connected', (url: string) => {
    if (url === wsUrl.value) {
      connected.value = true
    }
  })
}

function cleanupListeners() {
  EventsOff('ws:message', 'ws:error', 'ws:closed', 'ws:connected')
  unsubMessage = null
  unsubError = null
  unsubClosed = null
  unsubConnected = null
}

async function sendMessage() {
  if (!inputMessage.value.trim() || !connected.value) return

  const msg = inputMessage.value
  inputMessage.value = ''

  try {
    await WSSend(wsUrl.value, msg)
    messages.value.push({ type: 'sent', content: msg, time: nowTime() })
    scrollBottom()
  } catch (err: any) {
    addSystem('Send failed: ' + (err?.message || err))
  }
}

function onKeydown(e: KeyboardEvent) {
  if (e.ctrlKey && e.key === 'Enter') {
    e.preventDefault()
    sendMessage()
  }
}

onUnmounted(() => {
  cleanupListeners()
  if (connected.value) {
    WSDisconnect(wsUrl.value).catch(() => {})
  }
})
</script>

<style scoped>
.ws-panel {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-height: 0;
}

.ws-toolbar {
  padding: 8px 10px;
  border-bottom: 1px solid #e8e8e8;
  flex-shrink: 0;
}

.ws-url-row {
  display: flex;
  gap: 6px;
  align-items: center;
}

.ws-url-input {
  flex: 1;
}

.ws-status {
  font-size: 12px;
  color: #999;
  display: flex;
  align-items: center;
  gap: 4px;
  white-space: nowrap;
}

.ws-status.connected {
  color: #18a058;
}

.ws-dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: #ddd;
}

.ws-status.connected .ws-dot {
  background: #18a058;
}

.ws-messages {
  flex: 1;
  overflow-y: auto;
  padding: 10px;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.ws-empty {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
}

.ws-message {
  max-width: 80%;
  padding: 6px 10px;
  border-radius: 8px;
  font-size: 13px;
  line-height: 1.4;
  word-break: break-word;
}

.ws-message.sent {
  align-self: flex-end;
  background: #d4edda;
  color: #155724;
}

.ws-message.received {
  align-self: flex-start;
  background: #e9e9e9;
  color: #333;
}

.ws-message.system {
  align-self: center;
  background: transparent;
  color: #999;
  font-size: 11px;
  padding: 2px 6px;
}

.ws-msg-time {
  font-size: 10px;
  opacity: 0.6;
  margin-bottom: 2px;
}

.ws-input-row {
  padding: 8px 10px;
  border-top: 1px solid #e8e8e8;
  display: flex;
  gap: 6px;
  flex-shrink: 0;
}
</style>
