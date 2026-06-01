# Phase 3: Environment Variables Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add multi-environment management with `{{variable}}` interpolation support across URL, headers, and body.

**Architecture:** Environment CRUD backed by existing Go service/handler. Active environment selected via frontend dropdown. Variable interpolation performed client-side before sending requests, using a composable utility.

**Tech Stack:** Go 1.23, Vue 3 `<script setup>`, Naive UI (`NSelect`, `NModal`, `NDynamicInput`, `NButton`, `NTag`), Pinia

---

## File Structure

```
paw-api/
├── services/
│   └── environment_service.go     # Modify: add Update, Delete methods
├── handlers/
│   └── environment_handler.go     # Modify: add UpdateEnvironment, DeleteEnvironment
│
├── frontend/src/
│   ├── components/
│   │   ├── EnvSelector.vue        # New: dropdown to select active environment + manage button
│   │   └── EnvManager.vue         # New: modal dialog for environment CRUD + variable editing
│   │   └── AppSidebar.vue         # Modify: add EnvSelector to header
│   ├── stores/
│   │   └── environment.ts         # Rewrite: full CRUD actions + active env state
│   ├── composables/
│   │   └── useVariableResolver.ts # New: resolve {{var}} patterns using active env
│   └── components/
│       └── RequestEditor.vue      # Modify: resolve variables before sending
```

---

### Task 1: Go — Add Update/Delete to environment service and handler

**Files:**
- Modify: `services/environment_service.go`
- Modify: `handlers/environment_handler.go`

- [ ] **Step 1: Add Update and Delete to `services/environment_service.go`**

Add before the closing brace:

```go
func (s *EnvironmentService) Update(id, name, variables string) (*models.Environment, error) {
	if name == "" {
		return nil, errors.New("environment name is required")
	}
	e, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if e == nil {
		return nil, errors.New("environment not found")
	}
	e.Name = name
	e.Variables = variables
	e.UpdatedAt = time.Now()
	return e, s.repo.Update(e)
}

func (s *EnvironmentService) Delete(id string) error {
	return s.repo.Delete(id)
}
```

- [ ] **Step 2: Add handlers to `handlers/environment_handler.go`**

Add before the closing brace:

```go
func (h *EnvironmentHandler) UpdateEnvironment(id, name, variables string) (*models.Environment, error) {
	return h.service.Update(id, name, variables)
}

func (h *EnvironmentHandler) DeleteEnvironment(id string) error {
	return h.service.Delete(id)
}
```

- [ ] **Step 3: Verify Go compiles**

Run: `cd D:\javap\paw-api && go build ./...`
Expected: No errors

---

### Task 2: Frontend — Environment store

**Files:**
- Rewrite: `frontend/src/stores/environment.ts`

- [ ] **Step 1: Rewrite `frontend/src/stores/environment.ts`**

```typescript
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Environment } from '../types/environment'

export interface EnvVariable {
  key: string
  value: string
  enabled: boolean
}

export const useEnvironmentStore = defineStore('environment', () => {
  const environments = ref<Environment[]>([])
  const activeEnvironment = ref<Environment | null>(null)

  const activeVariables = computed<EnvVariable[]>(() => {
    if (!activeEnvironment.value) return []
    try {
      return JSON.parse(activeEnvironment.value.variables)
    } catch {
      return []
    }
  })

  function setEnvironments(list: Environment[]) {
    environments.value = list
  }

  function addEnvironment(env: Environment) {
    environments.value.push(env)
  }

  function removeEnvironment(id: string) {
    environments.value = environments.value.filter(e => e.id !== id)
    if (activeEnvironment.value?.id === id) {
      activeEnvironment.value = null
    }
  }

  function updateEnvironment(env: Environment) {
    const idx = environments.value.findIndex(e => e.id === env.id)
    if (idx !== -1) {
      environments.value[idx] = env
    }
    if (activeEnvironment.value?.id === env.id) {
      activeEnvironment.value = env
    }
  }

  function setActiveEnvironment(env: Environment | null) {
    activeEnvironment.value = env
  }

  return {
    environments, activeEnvironment, activeVariables,
    setEnvironments, addEnvironment, removeEnvironment,
    updateEnvironment, setActiveEnvironment,
  }
})
```

---

### Task 3: Frontend — Variable resolver composable

**Files:**
- Create: `frontend/src/composables/useVariableResolver.ts`

- [ ] **Step 1: Create `useVariableResolver.ts`**

```typescript
import { useEnvironmentStore, type EnvVariable } from '../stores/environment'

export function useVariableResolver() {
  function resolve(text: string, variables: EnvVariable[]): string {
    let result = text
    for (const v of variables) {
      if (!v.enabled) continue
      const pattern = new RegExp(`\\{\\{\\s*${escapeRegex(v.key)}\\s*\\}\\}`, 'g')
      result = result.replace(pattern, v.value)
    }
    return result
  }

  function escapeRegex(str: string): string {
    return str.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
  }

  return { resolve }
}
```

---

### Task 4: Frontend — Environment selector component

**Files:**
- Create: `frontend/src/components/EnvSelector.vue`

- [ ] **Step 1: Create `EnvSelector.vue`**

```vue
<script lang="ts" setup>
import { computed } from 'vue'
import { NSelect, NButton, NIcon, NSpace } from 'naive-ui'
import { Settings } from '@vicons/ionicons5'
import { useEnvironmentStore } from '../stores/environment'
import { useProjectStore } from '../stores/project'

const envStore = useEnvironmentStore()
const projectStore = useProjectStore()

const emit = defineEmits<{
  (e: 'manage'): void
}>()

const envOptions = computed(() => envStore.environments.map(e => ({
  label: e.name,
  value: e.id,
})))

async function handleChange(envId: string | null) {
  if (!envId || !projectStore.currentProject) return
  try {
    const env = await SetActiveEnvironment(envId, projectStore.currentProject.id)
    envStore.setActiveEnvironment(env)
  } catch (e: any) {
    console.error('Failed to set active environment', e)
  }
}
</script>

<template>
  <NSpace align="center" size="small">
    <NSelect
      :options="envOptions"
      :value="envStore.activeEnvironment?.id ?? null"
      placeholder="No environment"
      size="tiny"
      style="width: 140px"
      clearable
      @update:value="handleChange"
    />
    <NButton quaternary size="tiny" @click="emit('manage')">
      <template #icon><NIcon><Settings /></NIcon></template>
    </NButton>
  </NSpace>
</template>
```

---

### Task 5: Frontend — Environment manager modal

**Files:**
- Create: `frontend/src/components/EnvManager.vue`

- [ ] **Step 1: Create `EnvManager.vue`**

```vue
<script lang="ts" setup>
import { ref, onMounted } from 'vue'
import { NModal, NButton, NInput, NIcon, NForm, NFormItem, NSpace, NTabPane, NTabs, useMessage } from 'naive-ui'
import { Add, Trash } from '@vicons/ionicons5'
import { useEnvironmentStore, type EnvVariable } from '../stores/environment'
import {
  ListEnvironments,
  CreateEnvironment,
  UpdateEnvironment,
  DeleteEnvironment,
  SetActiveEnvironment,
} from '../../wailsjs/go/handlers/EnvironmentHandler'
import { useProjectStore } from '../stores/project'

const props = defineProps<{
  show: boolean
}>()

const emit = defineEmits<{
  (e: 'update:show', val: boolean): void
}>()

const envStore = useEnvironmentStore()
const projectStore = useProjectStore()
const message = useMessage()

const editingEnvId = ref<string | null>(null)
const envName = ref('')
const variables = ref<EnvVariable[]>([])

function resetForm() {
  editingEnvId.value = null
  envName.value = ''
  variables.value = []
}

function editEnv(env: typeof envStore.environments[0]) {
  editingEnvId.value = env.id
  envName.value = env.name
  try {
    variables.value = JSON.parse(env.variables)
  } catch {
    variables.value = []
  }
}

function addVariable() {
  variables.value.push({ key: '', value: '', enabled: true })
}

function removeVariable(index: number) {
  variables.value.splice(index, 1)
}

async function save() {
  if (!envName.value.trim()) {
    message.warning('Environment name is required')
    return
  }
  if (!projectStore.currentProject) {
    message.error('No project selected')
    return
  }
  const varsJSON = JSON.stringify(variables.value)
  try {
    if (editingEnvId.value) {
      const updated = await UpdateEnvironment(editingEnvId.value, envName.value.trim(), varsJSON)
      envStore.updateEnvironment(updated)
      message.success('Environment updated')
    } else {
      const created = await CreateEnvironment(
        projectStore.currentProject.id,
        envName.value.trim(),
        varsJSON,
        envStore.environments.length === 0,
      )
      envStore.addEnvironment(created)
      if (created.is_active) {
        envStore.setActiveEnvironment(created)
      }
      message.success('Environment created')
    }
    resetForm()
  } catch (e: any) {
    message.error(e.message || 'Failed to save environment')
  }
}

async function removeEnv(id: string) {
  try {
    await DeleteEnvironment(id)
    envStore.removeEnvironment(id)
    message.success('Environment deleted')
  } catch (e: any) {
    message.error(e.message || 'Failed to delete environment')
  }
}

async function setActive(id: string) {
  if (!projectStore.currentProject) return
  try {
    const env = await SetActiveEnvironment(id, projectStore.currentProject.id)
    envStore.setActiveEnvironment(env)
    message.success(`${env.name} is now active`)
  } catch (e: any) {
    message.error(e.message || 'Failed to set active environment')
  }
}

onMounted(async () => {
  if (projectStore.currentProject) {
    try {
      const envs = await ListEnvironments(projectStore.currentProject.id)
      envStore.setEnvironments(envs)
      const active = envs.find(e => e.is_active)
      if (active) envStore.setActiveEnvironment(active)
    } catch (e: any) {
      console.error('Failed to load environments', e)
    }
  }
})
</script>

<template>
  <NModal
    :show="props.show"
    @update:show="emit('update:show', $event)"
    title="Environment Manager"
    preset="card"
    style="width: 520px"
  >
    <NTabs type="line">
      <NTabPane name="list" tab="Environments">
        <div v-if="envStore.environments.length === 0" class="empty-hint">
          No environments yet. Add one to start using variables.
        </div>
        <div v-for="env in envStore.environments" :key="env.id" class="env-row">
          <div class="env-info" @click="editEnv(env)">
            <span class="env-name">{{ env.name }}</span>
            <span v-if="envStore.activeEnvironment?.id === env.id" class="active-badge">ACTIVE</span>
          </div>
          <NSpace>
            <NButton size="tiny" quaternary @click.stop="setActive(env.id)">Activate</NButton>
            <NButton size="tiny" quaternary @click.stop="editEnv(env)">Edit</NButton>
            <NButton size="tiny" quaternary @click.stop="removeEnv(env.id)">
              <template #icon><NIcon><Trash /></NIcon></template>
            </NButton>
          </NSpace>
        </div>
      </NTabPane>
      <NTabPane name="edit" tab="Edit">
        <NForm>
          <NFormItem label="Name">
            <NInput v-model:value="envName" placeholder="e.g. Development" />
          </NFormItem>
          <NFormItem label="Variables">
            <div class="var-list">
              <div v-for="(v, i) in variables" :key="i" class="var-row">
                <input type="checkbox" v-model="v.enabled" class="var-checkbox" />
                <NInput v-model:value="v.key" size="tiny" placeholder="Key" class="var-input" />
                <NInput v-model:value="v.value" size="tiny" placeholder="Value" class="var-input" />
                <NButton quaternary circle size="tiny" @click="removeVariable(i)">
                  <template #icon><NIcon><Trash /></NIcon></template>
                </NButton>
              </div>
              <NButton size="tiny" quaternary @click="addVariable">
                <template #icon><NIcon><Add /></NIcon></template>
                Add Variable
              </NButton>
            </div>
          </NFormItem>
        </NForm>
        <NSpace justify="end" style="margin-top: 12px;">
          <NButton @click="resetForm">Reset</NButton>
          <NButton type="primary" @click="save">Save</NButton>
        </NSpace>
      </NTabPane>
    </NTabs>
  </NModal>
</template>

<style scoped>
.env-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 4px;
  border-bottom: 1px solid var(--border-color);
}
.env-info {
  cursor: pointer;
  flex: 1;
}
.env-name {
  font-weight: 500;
}
.active-badge {
  margin-left: 8px;
  font-size: 10px;
  background: #18a058;
  color: white;
  padding: 1px 6px;
  border-radius: 3px;
}
.empty-hint {
  color: #999;
  font-size: 13px;
  padding: 24px 0;
  text-align: center;
}
.var-list {
  width: 100%;
}
.var-row {
  display: flex;
  align-items: center;
  gap: 4px;
  margin-bottom: 4px;
}
.var-checkbox {
  width: 14px;
  height: 14px;
  flex-shrink: 0;
}
.var-input {
  flex: 1;
}
</style>
```

---

### Task 6: Frontend — Integrate into sidebar and RequestEditor

**Files:**
- Modify: `frontend/src/components/AppSidebar.vue`
- Modify: `frontend/src/components/RequestEditor.vue`

- [ ] **Step 1: Add EnvSelector to `AppSidebar.vue` header**

Replace the sidebar-header section:

```vue
<script lang="ts" setup>
import { ref } from 'vue'
// ... existing imports ...
import EnvSelector from './EnvSelector.vue'
import EnvManager from './EnvManager.vue'

const showEnvManager = ref(false)
</script>

<template>
  <NLayoutSider bordered width="220" :native-scrollbar="false" class="app-sidebar">
    <div class="sidebar-header">
      <div class="sidebar-title-row">
        <span class="sidebar-title">Paw API</span>
      </div>
      <EnvSelector @manage="showEnvManager = true" />
    </div>
    <!-- ... rest of sidebar ... -->
    <EnvManager v-model:show="showEnvManager" />
  </NLayoutSider>
</template>

<style scoped>
.sidebar-header {
  padding: 12px 12px 8px;
  border-bottom: 1px solid var(--border-color);
  display: flex;
  flex-direction: column;
  gap: 6px;
}
.sidebar-title-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.sidebar-title {
  font-size: 16px;
  font-weight: 700;
}
</style>
```

Full modified `AppSidebar.vue`:

```vue
<script lang="ts" setup>
import { ref, computed } from 'vue'
import { NLayoutSider, NTree, NButton, NIcon, NInput, NModal, NForm, NFormItem, NSpace, useMessage } from 'naive-ui'
import { Add } from '@vicons/ionicons5'
import { useRouter, useRoute } from 'vue-router'
import { useProjectStore } from '../stores/project'
import { CreateCollection } from '../../wailsjs/go/handlers/CollectionHandler'
import EnvSelector from './EnvSelector.vue'
import EnvManager from './EnvManager.vue'

const router = useRouter()
const route = useRoute()
const projectStore = useProjectStore()
const message = useMessage()

const showAddModal = ref(false)
const newCollectionName = ref('')
const addingParentId = ref<string | null>(null)
const showEnvManager = ref(false)

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

function startAdd(parentId: string | null) {
  addingParentId.value = parentId
  newCollectionName.value = ''
  showAddModal.value = true
}

async function confirmAdd() {
  if (!newCollectionName.value.trim()) return
  if (!projectStore.currentProject) {
    message.error('No project selected')
    return
  }
  try {
    const col = await CreateCollection(
      projectStore.currentProject.id,
      addingParentId.value ?? '',
      newCollectionName.value.trim(),
      projectStore.collections.length,
    )
    projectStore.addCollection(col)
    showAddModal.value = false
    message.success('Collection created')
  } catch (e: any) {
    message.error(e.message || 'Failed to create collection')
  }
}
</script>

<template>
  <NLayoutSider bordered width="220" :native-scrollbar="false" class="app-sidebar">
    <div class="sidebar-header">
      <div class="sidebar-title-row">
        <span class="sidebar-title">Paw API</span>
      </div>
      <EnvSelector @manage="showEnvManager = true" />
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

    <EnvManager v-model:show="showEnvManager" />
  </NLayoutSider>
</template>

<style scoped>
.app-sidebar {
  height: 100%;
  display: flex;
  flex-direction: column;
}
.sidebar-header {
  padding: 12px 12px 8px;
  border-bottom: 1px solid var(--border-color);
  display: flex;
  flex-direction: column;
  gap: 6px;
}
.sidebar-title-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
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

- [ ] **Step 2: Add variable resolution to `RequestEditor.vue`**

Add import and integrate into `handleSend`:

```typescript
import { useVariableResolver } from '../composables/useVariableResolver'
import { useEnvironmentStore } from '../stores/environment'

// ... inside setup ...
const envStore = useEnvironmentStore()
const { resolve } = useVariableResolver()

async function handleSend() {
  // ... existing validation ...
  sending.value = true
  try {
    const resolvedURL = resolve(appendParamsToURL(url.value.trim()), envStore.activeVariables)
    const resolvedBody = bodyType.value === 'none' ? '' : resolve(body.value, envStore.activeVariables)
    const resolvedHeaders = buildHeadersMap()
    for (const k of Object.keys(resolvedHeaders)) {
      resolvedHeaders[k] = resolve(resolvedHeaders[k], envStore.activeVariables)
    }

    const resp = await SendRequest({
      Method: method.value,
      URL: resolvedURL,
      Headers: resolvedHeaders,
      Body: resolvedBody,
    })
    // ... rest ...
  }
}
```

---

### Task 7: Build verification

- [ ] **Step 1: Generate Wails bindings**

Run: `cd D:\javap\paw-api && wails generate module`
Expected: Generates updated binding files

- [ ] **Step 2: Frontend build**

Run: `cd D:\javap\paw-api\frontend && npm run build`
Expected: `vue-tsc --noEmit` passes, `vite build` succeeds

- [ ] **Step 3: Go build**

Run: `cd D:\javap\paw-api && go build ./...`
Expected: No errors

- [ ] **Step 4: Wails build**

Run: `cd D:\javap\paw-api && wails build`
Expected: Builds successfully, outputs to `build/bin/`
