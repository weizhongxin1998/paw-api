# Phase 2: HTTP Request Debugging Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Enable full HTTP request debugging workflow: collection tree management, request editing (method/URL/params/headers/body/auth), sending requests, and viewing responses.

**Architecture:** Go HTTP client wraps `net/http` with timeout and redirect control. Frontend uses reactive Pinia stores to hold the current request state. RequestEditor emits to a composable that calls the Wails-bound handler.

**Tech Stack:** Go 1.23 `net/http`, Vue 3 `<script setup>`, Naive UI (`NTree`, `NDataTable`, `NDynamicInput`, `NInput`, `NButton`, `NTabs`), Pinia

---

## File Structure

```
paw-api/
├── pkg/
│   └── httpclient/
│       └── client.go          # Rewrite: real HTTP client with Do()
├── services/
│   └── request_service.go     # Modify: add Send()
├── handlers/
│   └── request_handler.go     # Modify: add SendRequest()
│
├── frontend/src/
│   ├── components/
│   │   ├── AppSidebar.vue     # Rewrite: collection tree with NTree + CRUD
│   │   ├── RequestEditor.vue  # Rewrite: full request editing
│   │   ├── ResponseViewer.vue # Rewrite: response display
│   │   └── KeyValueEditor.vue # New: reusable key-value pair editor
│   ├── stores/
│   │   ├── request.ts         # Modify: full request state
│   │   └── project.ts         # Modify: collection CRUD actions
│   └── composables/
│       └── useRequest.ts      # Modify: send logic + state mgmt
```

---

### Task 1: Go — HTTP client and send endpoint

**Files:**
- Rewrite: `pkg/httpclient/client.go`
- Modify: `services/request_service.go` (add `Send`)
- Modify: `handlers/request_handler.go` (add `SendRequest`)

- [ ] **Step 1: Rewrite `pkg/httpclient/client.go`**

```go
package httpclient

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	hc *http.Client
}

type Request struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    string
}

type Response struct {
	Status     int                 `json:"status"`
	StatusText string              `json:"status_text"`
	Headers    map[string][]string `json:"headers"`
	Body       string              `json:"body"`
	DurationMs int64               `json:"duration_ms"`
}

func NewClient() *Client {
	return &Client{
		hc: &http.Client{
			Timeout:   30 * time.Second,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if len(via) >= 10 {
					return http.ErrUseLastResponse
				}
				return nil
			},
		},
	}
}

func (c *Client) Do(req *Request) (*Response, error) {
	var bodyReader io.Reader
	if req.Body != "" {
		bodyReader = bytes.NewBufferString(req.Body)
	}

	httpReq, err := http.NewRequest(req.Method, req.URL, bodyReader)
	if err != nil {
		return nil, err
	}

	for k, v := range req.Headers {
		httpReq.Header.Set(k, v)
	}

	start := time.Now()
	httpResp, err := c.hc.Do(httpReq)
	duration := time.Since(start).Milliseconds()
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}

	statusText := strings.TrimPrefix(httpResp.Status, httpResp.Status+http.StatusText(httpResp.StatusCode))

	return &Response{
		Status:     httpResp.StatusCode,
		StatusText: statusText,
		Headers:    httpResp.Header,
		Body:       string(respBody),
		DurationMs: duration,
	}, nil
}
```

- [ ] **Step 2: Add `Send` to `services/request_service.go`**

Import `paw-api/pkg/httpclient` and add:

```go
func (s *RequestService) Send(method, url string, headers map[string]string, body string) (*httpclient.Response, error) {
	client := httpclient.NewClient()
	return client.Do(&httpclient.Request{
		Method:  method,
		URL:     url,
		Headers: headers,
		Body:    body,
	})
}
```

- [ ] **Step 3: Add `SendRequest` to `handlers/request_handler.go`**

```go
type SendRequestInput struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    string
}

func (h *RequestHandler) SendRequest(input SendRequestInput) (*httpclient.Response, error) {
	return h.service.Send(input.Method, input.URL, input.Headers, input.Body)
}
```

Also update the import:

```go
import (
	"paw-api/models"
	"paw-api/pkg/httpclient"
	"paw-api/services"
)
```

- [ ] **Step 4: Verify Go compiles**

Run: `cd D:\javap\paw-api && go build ./...`
Expected: No errors

---

### Task 2: Frontend — KeyValueEditor component

**Files:**
- Create: `frontend/src/components/KeyValueEditor.vue`

- [ ] **Step 1: Create `KeyValueEditor.vue`**

```vue
<script lang="ts" setup>
import { NButton, NInput, NIcon } from 'naive-ui'
import { Add, Trash } from '@vicons/ionicons5'

export interface KeyValuePair {
  key: string
  value: string
  enabled: boolean
}

const props = withDefaults(defineProps<{
  modelValue: KeyValuePair[]
  keyPlaceholder?: string
  valuePlaceholder?: string
}>(), {
  keyPlaceholder: 'Key',
  valuePlaceholder: 'Value',
})

const emit = defineEmits<{
  (e: 'update:modelValue', val: KeyValuePair[]): void
}>()

function addRow() {
  emit('update:modelValue', [...props.modelValue, { key: '', value: '', enabled: true }])
}

function removeRow(index: number) {
  const next = props.modelValue.filter((_, i) => i !== index)
  emit('update:modelValue', next)
}

function updateRow(index: number, field: 'key' | 'value' | 'enabled', val: string | boolean) {
  const next = props.modelValue.map((pair, i) => {
    if (i !== index) return pair
    return { ...pair, [field]: val }
  })
  emit('update:modelValue', next)
}
</script>

<template>
  <div class="kv-editor">
    <div class="kv-rows">
      <div v-for="(pair, index) in modelValue" :key="index" class="kv-row">
        <input
          type="checkbox"
          :checked="pair.enabled"
          class="kv-checkbox"
          @change="updateRow(index, 'enabled', ($event.target as HTMLInputElement).checked)"
        />
        <NInput
          :value="pair.key"
          size="tiny"
          :placeholder="keyPlaceholder"
          class="kv-input"
          @update:value="updateRow(index, 'key', $event)"
        />
        <NInput
          :value="pair.value"
          size="tiny"
          :placeholder="valuePlaceholder"
          class="kv-input"
          @update:value="updateRow(index, 'value', $event)"
        />
        <NButton quaternary circle size="tiny" @click="removeRow(index)">
          <template #icon><NIcon><Trash /></NIcon></template>
        </NButton>
      </div>
    </div>
    <NButton size="tiny" quaternary @click="addRow" class="kv-add-btn">
      <template #icon><NIcon><Add /></NIcon></template>
      Add
    </NButton>
  </div>
</template>

<style scoped>
.kv-rows {
  display: flex;
  flex-direction: column;
  gap: 4px;
}
.kv-row {
  display: flex;
  align-items: center;
  gap: 4px;
}
.kv-checkbox {
  width: 14px;
  height: 14px;
  cursor: pointer;
  flex-shrink: 0;
}
.kv-input {
  flex: 1;
}
.kv-add-btn {
  margin-top: 4px;
}
</style>
```

---

### Task 3: Frontend — Collection tree sidebar

**Files:**
- Rewrite: `frontend/src/components/AppSidebar.vue`
- Modify: `frontend/src/stores/project.ts` (add collection CRUD actions)

- [ ] **Step 1: Update `frontend/src/stores/project.ts`**

```typescript
import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Project, Collection } from '../types/project'

export const useProjectStore = defineStore('project', () => {
  const projects = ref<Project[]>([])
  const currentProject = ref<Project | null>(null)
  const collections = ref<Collection[]>([])
  const selectedCollectionId = ref<string | null>(null)

  function setCurrentProject(p: Project) {
    currentProject.value = p
  }

  function setCollections(list: Collection[]) {
    collections.value = list
  }

  function addCollection(c: Collection) {
    collections.value.push(c)
  }

  function removeCollection(id: string) {
    collections.value = collections.value.filter(c => c.id !== id)
    if (selectedCollectionId.value === id) {
      selectedCollectionId.value = null
    }
  }

  function updateCollection(c: Collection) {
    const idx = collections.value.findIndex(x => x.id === c.id)
    if (idx !== -1) {
      collections.value[idx] = c
    }
  }

  function selectCollection(id: string | null) {
    selectedCollectionId.value = id
  }

  return {
    projects, currentProject, collections, selectedCollectionId,
    setCurrentProject, setCollections, addCollection, removeCollection,
    updateCollection, selectCollection,
  }
})
```

- [ ] **Step 2: Rewrite `AppSidebar.vue`**

```vue
<script lang="ts" setup>
import { ref, computed } from 'vue'
import { NLayoutSider, NTree, NButton, NIcon, NInput, NModal, NForm, NFormItem, NSpace, useMessage } from 'naive-ui'
import { Add } from '@vicons/ionicons5'
import { useRouter, useRoute } from 'vue-router'
import { useProjectStore } from '../stores/project'
import { CreateCollection } from '../../wailsjs/go/handlers/CollectionHandler'

const router = useRouter()
const route = useRoute()
const projectStore = useProjectStore()
const message = useMessage()

const showAddModal = ref(false)
const newCollectionName = ref('')
const addingParentId = ref<string | null>(null)
const editingId = ref<string | null>(null)
const editingName = ref('')

const menuItems = [
  { label: 'Workspace', key: '/workspace' },
  { label: 'History', key: '/history' },
  { label: 'Docs', key: '/docs' },
  { label: 'Test Runner', key: '/test-runner' },
]

function navigateTo(path: string) {
  router.push(path)
}

const treeData = computed(() => {
  function buildTree(parentId: string | null): any[] {
    return projectStore.collections
      .filter(c => c.parent_id === parentId)
      .sort((a, b) => a.sort_order - b.sort_order)
      .map(c => ({
        label: c.name,
        key: c.id,
        isLeaf: false,
        children: buildTree(c.id),
      }))
  }
  return [
    ...menuItems.map(m => ({
      label: m.label,
      key: m.key,
      isLeaf: true,
    })),
    {
      label: 'Collections',
      key: 'collections-header',
      isLeaf: true,
      disabled: true,
    },
    ...buildTree(null),
  ]
})

function handleNodeSelect(keys: string[]) {
  if (keys.length === 0) return
  const key = keys[0]
  const menuItem = menuItems.find(m => m.key === key)
  if (menuItem) {
    router.push(key)
    return
  }
  projectStore.selectCollection(key)
}

import { CreateCollection, ListCollections, DeleteCollection } from '../../wailsjs/go/handlers/CollectionHandler'

function startAdd(parentId: string | null) {
  addingParentId.value = parentId
  newCollectionName.value = ''
  showAddModal.value = true
}

async function confirmAdd() {
  if (!newCollectionName.value.trim()) return
  try {
    const col = await CreateCollection(
      projectStore.currentProject?.id ?? '',
      addingParentId.value ?? '',
      newCollectionName.value.trim(),
      0,
    )
    projectStore.addCollection(col)
    showAddModal.value = false
    message.success('Collection created')
  } catch (e: any) {
    message.error(e.message || 'Failed to create collection')
  }
}

function handleContextMenu(e: MouseEvent, key: string) {
  // For now, simple double-click rename
  if (key.startsWith('/')) return
  const col = projectStore.collections.find(c => c.id === key)
  if (col) {
    editingId.value = key
    editingName.value = col.name
  }
}

function confirmRename() {
  if (!editingId.value || !editingName.value.trim()) {
    editingId.value = null
    return
  }
  const col = projectStore.collections.find(c => c.id === editingId.value)
  if (col) {
    col.name = editingName.value.trim()
    projectStore.updateCollection(col)
  }
  editingId.value = null
}
</script>

<template>
  <NLayoutSider bordered width="220" :native-scrollbar="false" class="app-sidebar">
    <div class="sidebar-header">
      <span class="sidebar-title">Paw API</span>
    </div>
    <NTree
      :data="treeData"
      :default-expand-all="true"
      block-line
      selectable
      @update:selected-keys="handleNodeSelect"
    />
    <div class="sidebar-actions">
      <NButton size="tiny" quaternary @click="startAdd(null)">
        <template #icon><NIcon><Add /></NIcon></template>
        New Collection
      </NButton>
    </div>

    <NModal v-model:show="showAddModal" title="New Collection" preset="card" style="width: 360px">
      <NForm>
        <NFormItem label="Name">
          <NInput v-model:value="newCollectionName" placeholder="Collection name" />
        </NFormItem>
      </NForm>
      <template #footer>
        <NSpace justify="end">
          <NButton @click="showAddModal = false">Cancel</NButton>
          <NButton type="primary" @click="confirmAdd">Create</NButton>
        </NSpace>
      </template>
    </NModal>
  </NLayoutSider>
</template>

<style scoped>
.app-sidebar {
  height: 100%;
  display: flex;
  flex-direction: column;
}
.sidebar-header {
  padding: 16px 16px 12px;
  border-bottom: 1px solid var(--border-color);
}
.sidebar-title {
  font-size: 16px;
  font-weight: 700;
}
.sidebar-actions {
  padding: 8px;
  border-top: 1px solid var(--border-color);
  margin-top: auto;
}
</style>
```

---

### Task 4: Frontend — Request editor

**Files:**
- Rewrite: `frontend/src/components/RequestEditor.vue`

- [ ] **Step 1: Rewrite `RequestEditor.vue`**

```vue
<script lang="ts" setup>
import { ref, reactive } from 'vue'
import { NInput, NSelect, NButton, NTabs, NTabPane, NSpace, NIcon, useMessage } from 'naive-ui'
import { Send } from '@vicons/ionicons5'
import KeyValueEditor, { type KeyValuePair } from './KeyValueEditor.vue'
import { useRequestStore } from '../stores/request'

const httpMethods = ['GET', 'POST', 'PUT', 'DELETE', 'PATCH', 'HEAD', 'OPTIONS']
const requestStore = useRequestStore()
const message = useMessage()

const method = ref('GET')
const url = ref('')
const params = ref<KeyValuePair[]>([])
const headers = ref<KeyValuePair[]>([{ key: 'Content-Type', value: 'application/json', enabled: true }])
const body = ref('')
const bodyType = ref('none')
const auth = ref<KeyValuePair[]>([])

const sending = ref(false)

function buildHeadersMap(): Record<string, string> {
  const map: Record<string, string> = {}
  for (const h of headers.value) {
    if (h.enabled && h.key) {
      map[h.key] = h.value
    }
  }
  return map
}

import { SendRequest } from '../../wailsjs/go/handlers/RequestHandler'

async function handleSend() {
  if (!url.value.trim()) {
    message.warning('Please enter a URL')
    return
  }
  sending.value = true
  try {
    const resp = await SendRequest({
      Method: method.value,
      URL: url.value.trim(),
      Headers: buildHeadersMap(),
      Body: bodyType.value === 'none' ? '' : body.value,
    })
    requestStore.setLastResponse(resp)
    message.success(`${resp.status} ${resp.status_text} — ${resp.duration_ms}ms`)
  } catch (e: any) {
    message.error(e.message || 'Request failed')
  } finally {
    sending.value = false
  }
}
</script>

<template>
  <div class="request-editor">
    <div class="url-row">
      <NSelect
        :options="httpMethods.map(m => ({ label: m, value: m }))"
        v-model:value="method"
        style="width: 110px"
        size="small"
      />
      <NInput
        v-model:value="url"
        placeholder="https://api.example.com/endpoint"
        size="small"
        class="url-input"
      />
      <NButton type="primary" size="small" :loading="sending" @click="handleSend">
        <template #icon><NIcon><Send /></NIcon></template>
        Send
      </NButton>
    </div>
    <NTabs type="line" size="small" class="editor-tabs">
      <NTabPane name="params" tab="Params">
        <KeyValueEditor
          v-model="params"
          key-placeholder="Parameter name"
          value-placeholder="Value"
        />
      </NTabPane>
      <NTabPane name="headers" tab="Headers">
        <KeyValueEditor v-model="headers" />
      </NTabPane>
      <NTabPane name="body" tab="Body">
        <div class="body-controls">
          <NSelect
            :options="[
              { label: 'None', value: 'none' },
              { label: 'JSON', value: 'json' },
              { label: 'Text', value: 'text' },
              { label: 'Form-encoded', value: 'form' },
            ]"
            v-model:value="bodyType"
            size="tiny"
            style="width: 120px; margin-bottom: 8px;"
          />
          <NInput
            v-if="bodyType !== 'none'"
            v-model:value="body"
            type="textarea"
            :rows="6"
            placeholder="Request body"
            class="body-input"
          />
        </div>
      </NTabPane>
      <NTabPane name="auth" tab="Auth">
        <KeyValueEditor
          v-model="auth"
          key-placeholder="Auth key"
          value-placeholder="Value"
        />
      </NTabPane>
    </NTabs>
  </div>
</template>

<style scoped>
.request-editor {
  padding: 12px 16px;
  border-bottom: 1px solid var(--border-color);
}
.url-row {
  display: flex;
  gap: 8px;
  margin-bottom: 8px;
}
.url-input {
  flex: 1;
}
.editor-tabs {
  margin-top: 4px;
}
.body-controls {
  padding: 8px 0;
}
.body-input {
  font-family: 'Cascadia Code', 'Fira Code', 'Consolas', monospace;
  font-size: 13px;
}
</style>
```

- [ ] **Step 2: Update `frontend/src/stores/request.ts`** to add `lastResponse`

```typescript
import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Request } from '../types/request'

export interface ResponseData {
  status: number
  status_text: string
  headers: Record<string, string[]>
  body: string
  duration_ms: number
}

export const useRequestStore = defineStore('request', () => {
  const requests = ref<Request[]>([])
  const currentRequest = ref<Request | null>(null)
  const lastResponse = ref<ResponseData | null>(null)

  function setCurrentRequest(r: Request) {
    currentRequest.value = r
  }

  function setLastResponse(resp: ResponseData) {
    lastResponse.value = resp
  }

  return { requests, currentRequest, lastResponse, setCurrentRequest, setLastResponse }
})
```

---

### Task 5: Frontend — Response viewer

**Files:**
- Rewrite: `frontend/src/components/ResponseViewer.vue`

- [ ] **Step 1: Rewrite `ResponseViewer.vue`**

```vue
<script lang="ts" setup>
import { computed, ref } from 'vue'
import { NTag, NInput, NTabs, NTabPane, NEmpty } from 'naive-ui'
import { useRequestStore } from '../stores/request'

const requestStore = useRequestStore()
const response = computed(() => requestStore.lastResponse)
const statusType = computed(() => {
  if (!response.value) return 'default'
  if (response.value.status >= 200 && response.value.status < 300) return 'success'
  if (response.value.status >= 300 && response.value.status < 400) return 'warning'
  if (response.value.status >= 400) return 'error'
  return 'default'
})

const headersText = computed(() => {
  if (!response.value) return ''
  return Object.entries(response.value.headers)
    .map(([k, v]) => `${k}: ${v.join(', ')}`)
    .join('\n')
})

function formatJSON(text: string): string {
  try {
    return JSON.stringify(JSON.parse(text), null, 2)
  } catch {
    return text
  }
}
</script>

<template>
  <div class="response-viewer">
    <div v-if="response" class="response-header">
      <NTag :type="statusType" size="small" class="status-tag">
        {{ response.status }} {{ response.status_text }}
      </NTag>
      <span class="duration">{{ response.duration_ms }}ms</span>
    </div>
    <div v-if="!response" class="response-empty">
      <NEmpty description="Send a request to see the response" />
    </div>
    <div v-else class="response-body">
      <NTabs type="line" size="small">
        <NTabPane name="body" tab="Body">
          <NInput
            :value="formatJSON(response.body)"
            type="textarea"
            :rows="12"
            readonly
            class="resp-body-input"
          />
        </NTabPane>
        <NTabPane name="headers" tab="Headers">
          <NInput
            :value="headersText"
            type="textarea"
            :rows="8"
            readonly
            class="resp-headers-input"
          />
        </NTabPane>
      </NTabs>
    </div>
  </div>
</template>

<style scoped>
.response-viewer {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}
.response-header {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  border-bottom: 1px solid var(--border-color);
  font-size: 13px;
}
.duration {
  color: #999;
  font-size: 12px;
}
.response-empty {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
}
.response-body {
  flex: 1;
  padding: 8px 16px;
  overflow: auto;
}
.resp-body-input, .resp-headers-input {
  font-family: 'Cascadia Code', 'Fira Code', 'Consolas', monospace;
  font-size: 13px;
}
</style>
```

---

### Task 6: Build verification

**Files:**
- No file changes

- [ ] **Step 1: Generate Wails TypeScript bindings**

Run: `cd D:\javap\paw-api && wails generate module`
Expected: Generates `frontend/wailsjs/go/handlers/*.js` and `.d.ts` files

- [ ] **Step 2: Frontend type-check and build**

Run: `cd D:\javap\paw-api\frontend && npm run build`
Expected: `vue-tsc --noEmit` passes, `vite build` succeeds

- [ ] **Step 3: Go build**

Run: `cd D:\javap\paw-api && go build ./...`
Expected: No errors

- [ ] **Step 4: Wails build**

Run: `cd D:\javap\paw-api && wails build`
Expected: Builds successfully, outputs to `build/bin/`
