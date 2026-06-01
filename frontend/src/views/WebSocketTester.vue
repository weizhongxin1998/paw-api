<script lang="ts" setup>
import { computed, onMounted, onUnmounted, nextTick, ref } from 'vue'
import { NInput, NButton, NIcon, NTag } from 'naive-ui'
import { Send } from '@vicons/ionicons5'
import { Connect, Send as WsSend, Disconnect } from '../../wailsjs/go/handlers/WebSocketHandler'
import { EventsOn, EventsOff } from '../../wailsjs/runtime/runtime'
import { useTabsStore, type WsMessage } from '../stores/tabs'

const tabsStore = useTabsStore()
const tab = computed(() => tabsStore.activeTab)
const wsData = computed(() => tab.value?.wsData)
const inputMsg = ref('')
const messagesEnd = ref<HTMLDivElement | null>(null)
const connecting = ref(false)

function onWSMessage(raw: string) {
  try {
    const msg: WsMessage = JSON.parse(raw)
    if (!tab.value || !wsData.value) return
    tabsStore.updateWsData({ messages: [...wsData.value.messages, msg] })
    nextTick(() => messagesEnd.value?.scrollIntoView({ behavior: 'smooth' }))
  } catch {}
}

async function doConnect() {
  const data = wsData.value
  if (!data || !data.url.trim()) return
  connecting.value = true
  try {
    await Connect(data.url.trim())
    tabsStore.updateWsData({ connected: true, messages: [...data.messages, { type: 'system', content: 'Connected', time: Date.now() }] })
  } catch (e: any) {
    tabsStore.updateWsData({ messages: [...data.messages, { type: 'system', content: `Error: ${e.message || e}`, time: Date.now() }] })
  } finally { connecting.value = false }
}

async function doDisconnect() {
  const data = wsData.value
  if (!data) return
  await Disconnect()
  tabsStore.updateWsData({ connected: false, messages: [...data.messages, { type: 'system', content: 'Disconnected', time: Date.now() }] })
}

async function sendMessage() {
  const data = wsData.value
  if (!data || !inputMsg.value.trim()) return
  try {
    await WsSend(inputMsg.value.trim())
    inputMsg.value = ''
  } catch (e: any) {
    tabsStore.updateWsData({ messages: [...data.messages, { type: 'system', content: `Send error: ${e.message || e}`, time: Date.now() }] })
  }
}

function formatTime(ts: number): string {
  return new Date(ts).toLocaleTimeString()
}

onMounted(() => { EventsOn('ws-message', onWSMessage) })
onUnmounted(() => { EventsOff('ws-message') })
</script>

<template>
  <div v-if="wsData" class="ws-view">
    <div class="ws-toolbar">
      <NInput v-model:value="wsData.url" placeholder="ws://localhost:8080/ws" size="small" class="url-input" />
      <NButton v-if="!wsData.connected" size="small" type="primary" :loading="connecting" @click="doConnect">Connect</NButton>
      <NButton v-else size="small" type="error" @click="doDisconnect">Disconnect</NButton>
      <NTag v-if="wsData.connected" type="success" size="small">Connected</NTag>
      <NTag v-else size="small">Disconnected</NTag>
    </div>
    <div class="ws-messages">
      <div v-for="(msg, i) in wsData.messages" :key="i" class="ws-msg" :class="msg.type">
        <span class="msg-time">{{ formatTime(msg.time) }}</span>
        <span v-if="msg.type === 'system'" class="msg-system">{{ msg.content }}</span>
        <span v-else class="msg-content">{{ msg.content }}</span>
      </div>
      <div ref="messagesEnd" />
    </div>
    <div class="ws-input-row">
      <NInput v-model:value="inputMsg" placeholder="Type a message..." size="small" class="msg-input" :disabled="!wsData.connected" @keydown.enter.prevent="sendMessage" />
      <NButton size="small" type="primary" :disabled="!wsData.connected" @click="sendMessage">
        <template #icon><NIcon><Send /></NIcon></template>
        Send
      </NButton>
    </div>
  </div>
</template>

<style scoped>
.ws-view { display: flex; flex-direction: column; padding: 12px 16px; border-bottom: 1px solid var(--border-color); }
.ws-toolbar { display: flex; align-items: center; gap: 8px; margin-bottom: 8px; }
.url-input { flex: 1; }
.ws-messages { height: 200px; overflow-y: auto; background: var(--tab-bar-bg); border: 1px solid var(--border-color); border-radius: 4px; padding: 8px; margin-bottom: 8px; font-family: monospace; font-size: 13px; }
.ws-msg { padding: 2px 0; display: flex; gap: 8px; }
.msg-time { color: #999; font-size: 11px; white-space: nowrap; min-width: 70px; }
.msg-content { color: #333; word-break: break-all; }
.ws-msg.sent .msg-content { color: #2080f0; }
.ws-msg.received .msg-content { color: #18a058; }
.msg-system { color: #999; font-style: italic; }
.ws-input-row { display: flex; gap: 8px; }
.msg-input { flex: 1; }
</style>
