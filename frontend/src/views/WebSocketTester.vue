<script lang="ts" setup>
import { ref, onMounted, onUnmounted, nextTick } from 'vue'
import { useI18n } from 'vue-i18n'
import { NInput, NButton, NIcon, NTag } from 'naive-ui'
import { Send } from '@vicons/ionicons5'
import { Connect, Send as WsSend, Disconnect, IsConnected } from '../../wailsjs/go/handlers/WebSocketHandler'
import { EventsOn, EventsOff } from '../../wailsjs/runtime/runtime'

const { t } = useI18n()
const url = ref('ws://localhost:8080')
const message = ref('')
const messages = ref<Array<{ type: string; content: string; time: number }>>([])
const connected = ref(false)
const connecting = ref(false)
const messagesEnd = ref<HTMLDivElement | null>(null)

async function doConnect() {
  if (!url.value.trim()) return
  connecting.value = true
  try {
    await Connect(url.value.trim())
    connected.value = true
    messages.value.push({ type: 'system', content: t('ws.systemConnected'), time: Date.now() })
  } catch (e: any) {
    messages.value.push({ type: 'system', content: `${t('ws.error')}: ${e.message || e}`, time: Date.now() })
  } finally { connecting.value = false }
}

async function doDisconnect() {
  await Disconnect()
  connected.value = false
  messages.value.push({ type: 'system', content: t('ws.systemDisconnected'), time: Date.now() })
}

async function sendMessage() {
  if (!message.value.trim()) return
  try {
    await WsSend(message.value.trim())
    message.value = ''
  } catch (e: any) {
    messages.value.push({ type: 'system', content: `${t('ws.sendError')}: ${e.message || e}`, time: Date.now() })
  }
}

function onWSMessage(raw: string) {
  try { const msg = JSON.parse(raw); messages.value.push(msg); nextTick(() => messagesEnd.value?.scrollIntoView({ behavior: 'smooth' })) } catch {}
}

function formatTime(ts: number): string { return new Date(ts).toLocaleTimeString() }

onMounted(async () => { connected.value = await IsConnected(); EventsOn('ws-message', onWSMessage) })
onUnmounted(() => { EventsOff('ws-message') })
</script>

<template>
  <div class="ws-view">
    <div class="ws-toolbar">
      <NInput v-model:value="url" :placeholder="t('ws.placeholder')" size="small" class="url-input" />
      <NButton v-if="!connected" size="small" type="primary" :loading="connecting" @click="doConnect">{{ t('ws.connect') }}</NButton>
      <NButton v-else size="small" type="error" @click="doDisconnect">{{ t('ws.disconnect') }}</NButton>
      <NTag v-if="connected" type="success" size="small">{{ t('ws.connected') }}</NTag>
      <NTag v-else size="small">{{ t('ws.disconnected') }}</NTag>
    </div>
    <div class="ws-messages">
      <div v-for="(msg, i) in messages" :key="i" class="ws-msg" :class="msg.type">
        <span class="msg-time">{{ formatTime(msg.time) }}</span>
        <span v-if="msg.type === 'system'" class="msg-system">{{ msg.content }}</span>
        <span v-else class="msg-content">{{ msg.content }}</span>
      </div>
      <div ref="messagesEnd" />
    </div>
    <div class="ws-input-row">
      <NInput v-model:value="message" :placeholder="t('ws.messagePlaceholder')" size="small" class="msg-input" :disabled="!connected" @keydown.enter.prevent="sendMessage" />
      <NButton size="small" type="primary" :disabled="!connected" @click="sendMessage">
        <template #icon><NIcon><Send /></NIcon></template>
        {{ t('ws.send') }}
      </NButton>
    </div>
  </div>
</template>

<style scoped>
.ws-view { display: flex; flex-direction: column; height: 100%; padding: 12px 16px; }
.ws-toolbar { display: flex; align-items: center; gap: 8px; margin-bottom: 8px; }
.url-input { flex: 1; }
.ws-messages { flex: 1; overflow-y: auto; background: var(--tab-bar-bg); border: 1px solid var(--border-color); border-radius: 4px; padding: 8px; margin-bottom: 8px; font-family: monospace; font-size: 13px; }
.ws-msg { padding: 2px 0; display: flex; gap: 8px; }
.msg-time { color: #999; font-size: 11px; white-space: nowrap; min-width: 70px; }
.msg-content { color: #333; word-break: break-all; }
.ws-msg.sent .msg-content { color: #2080f0; }
.ws-msg.received .msg-content { color: #18a058; }
.msg-system { color: #999; font-style: italic; }
.ws-input-row { display: flex; gap: 8px; }
.msg-input { flex: 1; }
</style>
