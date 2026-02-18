<template>
  <div class="container mx-auto px-4 py-8 max-w-7xl">
    <div class="mb-6">
      <h1 class="text-3xl font-bold tracking-tight">Monitor Management</h1>
      <p class="text-muted-foreground mt-2">
        Endpoints and notification channels are now UI-driven. Suites and external endpoints can still be managed via the JSON overlay editor below.
      </p>
    </div>

    <div class="p-4 rounded-lg border bg-card mb-6">
      <p class="text-sm">
        <span class="font-semibold">Overlay Path:</span> {{ overlayPath || 'N/A' }}
      </p>
      <p class="text-muted-foreground text-sm mt-1">
        Saved changes are auto-applied by Gatus within a few seconds.
      </p>
    </div>

    <div class="grid gap-6 lg:grid-cols-[22rem,1fr]">
      <Card>
        <CardHeader>
          <div class="flex items-center justify-between gap-2">
            <CardTitle>Endpoints</CardTitle>
            <div class="flex items-center gap-2">
              <Button variant="outline" size="sm" @click="startCreate" :disabled="savingEndpoint">New</Button>
              <Button variant="outline" size="sm" @click="loadEndpoints" :disabled="loadingEndpoints || savingEndpoint">Refresh</Button>
            </div>
          </div>
        </CardHeader>
        <CardContent>
          <p v-if="loadingEndpoints" class="text-sm text-muted-foreground">Loading endpoints...</p>
          <p v-else-if="endpoints.length === 0" class="text-sm text-muted-foreground">No endpoint found in the managed configuration.</p>
          <div v-else class="space-y-2">
            <button
              v-for="endpoint in endpoints"
              :key="endpoint.key"
              class="w-full text-left rounded-md border p-3 transition-colors hover:bg-accent"
              :class="selectedKey === endpoint.key ? 'bg-accent border-primary' : ''"
              @click="selectEndpoint(endpoint)"
            >
              <div class="flex items-center justify-between gap-2">
                <p class="font-medium truncate">{{ endpoint.group ? `${endpoint.group} / ${endpoint.name}` : endpoint.name }}</p>
                <span class="text-xs text-muted-foreground uppercase">{{ endpoint.type }}</span>
              </div>
              <p class="text-xs text-muted-foreground truncate mt-1">{{ endpoint.url }}</p>
              <p class="text-xs text-muted-foreground mt-1">Key: {{ endpoint.key }}</p>
            </button>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>{{ isEditing ? `Edit Endpoint: ${selectedKey}` : 'Create Endpoint' }}</CardTitle>
        </CardHeader>
        <CardContent class="space-y-4">
          <div class="grid gap-4 md:grid-cols-2">
            <div class="space-y-1">
              <label class="text-sm font-medium">Name</label>
              <Input v-model="form.name" placeholder="frontend" />
            </div>
            <div class="space-y-1">
              <label class="text-sm font-medium">Group</label>
              <Input v-model="form.group" placeholder="core" />
            </div>
            <div class="space-y-1 md:col-span-2">
              <label class="text-sm font-medium">URL</label>
              <Input v-model="form.url" placeholder="https://example.org/health" />
            </div>
            <div class="space-y-1">
              <label class="text-sm font-medium">Method</label>
              <Select v-model="form.method" :options="methodOptions" />
            </div>
            <div class="space-y-1">
              <label class="text-sm font-medium">Interval</label>
              <Input v-model="form.interval" placeholder="30s / 1m / 5m" />
            </div>
          </div>

          <div class="grid gap-4 md:grid-cols-2">
            <div class="space-y-1">
              <label class="text-sm font-medium">Conditions (one per line)</label>
              <textarea
                v-model="form.conditionsText"
                class="w-full min-h-[180px] rounded-md border bg-background p-3 text-sm font-mono"
                spellcheck="false"
              />
            </div>
            <div class="space-y-1">
              <label class="text-sm font-medium">Headers (Key: Value)</label>
              <textarea
                v-model="form.headersText"
                class="w-full min-h-[180px] rounded-md border bg-background p-3 text-sm font-mono"
                spellcheck="false"
                placeholder="Authorization: Bearer token"
              />
            </div>
          </div>

          <div class="space-y-1">
            <label class="text-sm font-medium">Body</label>
            <textarea
              v-model="form.body"
              class="w-full min-h-[120px] rounded-md border bg-background p-3 text-sm font-mono"
              spellcheck="false"
            />
          </div>

          <div class="flex flex-wrap items-center gap-4">
            <label class="inline-flex items-center gap-2 text-sm">
              <input v-model="form.enabled" type="checkbox" class="h-4 w-4 rounded border-input" />
              Enabled
            </label>
            <label class="inline-flex items-center gap-2 text-sm">
              <input v-model="form.graphql" type="checkbox" class="h-4 w-4 rounded border-input" />
              GraphQL Body
            </label>
          </div>

          <div class="flex flex-wrap items-center gap-2">
            <Button @click="saveEndpoint" :disabled="savingEndpoint || loadingEndpoints">
              {{ isEditing ? 'Save Endpoint' : 'Create Endpoint' }}
            </Button>
            <Button variant="outline" @click="startCreate" :disabled="savingEndpoint">Clear</Button>
            <Button
              v-if="isEditing"
              variant="destructive"
              @click="deleteEndpoint"
              :disabled="savingEndpoint"
            >
              Delete Endpoint
            </Button>
            <span v-if="endpointMessage" class="text-sm text-muted-foreground">{{ endpointMessage }}</span>
          </div>
        </CardContent>
      </Card>
    </div>

    <Card class="mt-6">
      <CardHeader>
        <div class="flex items-center justify-between gap-2">
          <CardTitle>Notification Channels</CardTitle>
          <Button variant="outline" size="sm" @click="loadNotifications" :disabled="loadingNotifications || savingNotification">Refresh</Button>
        </div>
      </CardHeader>
      <CardContent>
        <div class="grid gap-6 lg:grid-cols-[22rem,1fr]">
          <div class="space-y-2">
            <p v-if="loadingNotifications" class="text-sm text-muted-foreground">Loading notifications...</p>
            <p v-else-if="notifications.length === 0" class="text-sm text-muted-foreground">No notification provider type available.</p>
            <button
              v-for="notification in notifications"
              :key="notification.type"
              class="w-full text-left rounded-md border p-3 transition-colors hover:bg-accent"
              :class="selectedNotificationType === notification.type ? 'bg-accent border-primary' : ''"
              @click="selectNotification(notification)"
            >
              <div class="flex items-center justify-between gap-2">
                <p class="font-medium truncate">{{ notification.type }}</p>
                <span class="text-xs" :class="notification.configured ? 'text-green-600' : 'text-muted-foreground'">
                  {{ notification.configured ? 'Configured' : 'Not Configured' }}
                </span>
              </div>
              <p class="text-xs text-muted-foreground mt-1">
                Used by {{ notification.usedByEndpoints }} endpoint(s), {{ notification.usedByExternalEndpoints }} external endpoint(s)
              </p>
            </button>
          </div>

          <div class="space-y-3">
            <p v-if="!selectedNotification" class="text-sm text-muted-foreground">Select a notification type from the left panel.</p>
            <template v-else>
              <div class="space-y-1">
                <label class="text-sm font-medium">Type</label>
                <Input :model-value="selectedNotification.type" disabled />
              </div>
              <div class="space-y-1">
                <label class="text-sm font-medium">Provider Config (JSON)</label>
                <textarea
                  v-model="notificationConfigText"
                  class="w-full min-h-[240px] rounded-md border bg-background p-3 text-sm font-mono"
                  spellcheck="false"
                />
              </div>
              <div class="flex flex-wrap items-center gap-2">
                <Button @click="saveNotification" :disabled="savingNotification || loadingNotifications">Save Notification</Button>
                <Button
                  variant="destructive"
                  @click="deleteNotification"
                  :disabled="savingNotification || !selectedNotification.configured"
                >
                  Delete Notification
                </Button>
                <span v-if="notificationMessage" class="text-sm text-muted-foreground">{{ notificationMessage }}</span>
              </div>
            </template>
          </div>
        </div>
      </CardContent>
    </Card>

    <Card class="mt-6">
      <CardHeader>
        <div class="flex items-center justify-between gap-2">
          <CardTitle>Advanced JSON Overlay</CardTitle>
          <div class="flex items-center gap-2">
            <Button variant="outline" size="sm" @click="loadManagedConfig" :disabled="loadingManaged || savingManaged">Reload</Button>
            <Button size="sm" @click="saveManagedConfig" :disabled="loadingManaged || savingManaged">Save Overlay</Button>
            <Button variant="destructive" size="sm" @click="resetOverlay" :disabled="loadingManaged || savingManaged">Reset Overlay</Button>
          </div>
        </div>
      </CardHeader>
      <CardContent class="space-y-3">
        <textarea
          v-model="jsonText"
          class="w-full min-h-[420px] rounded-md border bg-background p-3 font-mono text-xs"
          spellcheck="false"
        />
        <span v-if="managedMessage" class="text-sm text-muted-foreground">{{ managedMessage }}</span>
      </CardContent>
    </Card>
  </div>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Select } from '@/components/ui/select'

const endpoints = ref([])
const selectedKey = ref('')
const overlayPath = ref('')

const loadingEndpoints = ref(false)
const savingEndpoint = ref(false)
const endpointMessage = ref('')

const notifications = ref([])
const selectedNotificationType = ref('')
const notificationConfigText = ref('{}')
const loadingNotifications = ref(false)
const savingNotification = ref(false)
const notificationMessage = ref('')

const loadingManaged = ref(false)
const savingManaged = ref(false)
const managedMessage = ref('')
const jsonText = ref('')

const methodOptions = [
  { label: 'GET', value: 'GET' },
  { label: 'POST', value: 'POST' },
  { label: 'PUT', value: 'PUT' },
  { label: 'PATCH', value: 'PATCH' },
  { label: 'DELETE', value: 'DELETE' },
  { label: 'HEAD', value: 'HEAD' },
  { label: 'OPTIONS', value: 'OPTIONS' },
]

const createDefaultForm = () => ({
  enabled: true,
  name: '',
  group: '',
  url: '',
  method: 'GET',
  interval: '1m',
  conditionsText: '[STATUS] == 200',
  headersText: '',
  body: '',
  graphql: false,
})

const form = ref(createDefaultForm())

const isEditing = computed(() => selectedKey.value.length > 0)
const selectedNotification = computed(() => notifications.value.find((notification) => notification.type === selectedNotificationType.value) || null)

const normalizePayload = (payload) => {
  const safe = payload && typeof payload === 'object' ? payload : {}
  const normalizedAlerting = safe.alerting && typeof safe.alerting === 'object' && !Array.isArray(safe.alerting) ? safe.alerting : null
  return {
    alerting: normalizedAlerting,
    endpoints: Array.isArray(safe.endpoints) ? safe.endpoints : [],
    externalEndpoints: Array.isArray(safe.externalEndpoints) ? safe.externalEndpoints : [],
    suites: Array.isArray(safe.suites) ? safe.suites : [],
  }
}

const endpointToForm = (endpoint) => {
  const headers = endpoint && endpoint.headers && typeof endpoint.headers === 'object' ? endpoint.headers : {}
  const headerLines = Object.entries(headers).map(([name, value]) => `${name}: ${value}`)
  const conditions = Array.isArray(endpoint.conditions) ? endpoint.conditions : []
  return {
    enabled: endpoint.enabled !== false,
    name: endpoint.name || '',
    group: endpoint.group || '',
    url: endpoint.url || '',
    method: endpoint.method || 'GET',
    interval: endpoint.interval || '1m',
    conditionsText: conditions.join('\n'),
    headersText: headerLines.join('\n'),
    body: endpoint.body || '',
    graphql: endpoint.graphql === true,
  }
}

const parseHeadersText = (headersText) => {
  const headers = {}
  for (const rawLine of headersText.split('\n')) {
    const line = rawLine.trim()
    if (!line) continue
    const separatorIndex = line.indexOf(':')
    if (separatorIndex < 1) {
      throw new Error(`Invalid header format: ${line}`)
    }
    const name = line.slice(0, separatorIndex).trim()
    const value = line.slice(separatorIndex + 1).trim()
    if (!name) {
      throw new Error(`Invalid header name in line: ${line}`)
    }
    headers[name] = value
  }
  return headers
}

const buildPayloadFromForm = () => {
  const name = form.value.name.trim()
  const url = form.value.url.trim()
  const conditions = form.value.conditionsText
    .split('\n')
    .map((condition) => condition.trim())
    .filter((condition) => condition.length > 0)
  if (!name || !url) {
    throw new Error('Name and URL are required.')
  }
  if (conditions.length === 0) {
    throw new Error('At least one condition is required.')
  }
  return {
    enabled: form.value.enabled,
    name,
    group: form.value.group.trim(),
    url,
    method: form.value.method,
    interval: form.value.interval.trim(),
    conditions,
    headers: parseHeadersText(form.value.headersText),
    body: form.value.body,
    graphql: form.value.graphql,
  }
}

const startCreate = () => {
  selectedKey.value = ''
  form.value = createDefaultForm()
  endpointMessage.value = ''
}

const selectEndpoint = (endpoint) => {
  selectedKey.value = endpoint.key
  form.value = endpointToForm(endpoint)
  endpointMessage.value = ''
}

const selectNotification = (notification) => {
  selectedNotificationType.value = notification.type
  const configObject = notification.config && typeof notification.config === 'object' ? notification.config : {}
  notificationConfigText.value = JSON.stringify(configObject, null, 2)
  notificationMessage.value = ''
}

const loadEndpoints = async () => {
  loadingEndpoints.value = true
  endpointMessage.value = ''
  try {
    const response = await fetch('/api/v1/admin/endpoints', {
      credentials: 'include',
    })
    const data = await response.json()
    if (!response.ok) {
      throw new Error(data.error || 'Failed to load endpoints')
    }
    endpoints.value = Array.isArray(data.endpoints) ? data.endpoints : []
    overlayPath.value = data.overlayPath || overlayPath.value
    if (selectedKey.value) {
      const selectedEndpoint = endpoints.value.find((endpoint) => endpoint.key === selectedKey.value)
      if (selectedEndpoint) {
        form.value = endpointToForm(selectedEndpoint)
      } else {
        startCreate()
      }
    }
  } catch (error) {
    endpointMessage.value = error.message
  } finally {
    loadingEndpoints.value = false
  }
}

const saveEndpoint = async () => {
  savingEndpoint.value = true
  endpointMessage.value = ''
  try {
    const payload = buildPayloadFromForm()
    const targetURL = isEditing.value ? `/api/v1/admin/endpoints/${selectedKey.value}` : '/api/v1/admin/endpoints'
    const method = isEditing.value ? 'PUT' : 'POST'
    const response = await fetch(targetURL, {
      method,
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(payload),
    })
    const data = await response.json()
    if (!response.ok) {
      throw new Error(data.error || 'Failed to save endpoint')
    }
    selectedKey.value = data.key || selectedKey.value
    await loadEndpoints()
    endpointMessage.value = isEditing.value ? 'Endpoint saved.' : 'Endpoint created.'
  } catch (error) {
    endpointMessage.value = error.message
  } finally {
    savingEndpoint.value = false
  }
}

const deleteEndpoint = async () => {
  if (!isEditing.value) return
  const confirmed = window.confirm(`Delete endpoint ${selectedKey.value}?`)
  if (!confirmed) return
  savingEndpoint.value = true
  endpointMessage.value = ''
  try {
    const response = await fetch(`/api/v1/admin/endpoints/${selectedKey.value}`, {
      method: 'DELETE',
      credentials: 'include',
    })
    if (!response.ok && response.status !== 204) {
      let message = 'Failed to delete endpoint'
      try {
        const data = await response.json()
        message = data.error || message
      } catch {
        // noop
      }
      throw new Error(message)
    }
    startCreate()
    await loadEndpoints()
    endpointMessage.value = 'Endpoint deleted.'
  } catch (error) {
    endpointMessage.value = error.message
  } finally {
    savingEndpoint.value = false
  }
}

const loadNotifications = async () => {
  loadingNotifications.value = true
  notificationMessage.value = ''
  try {
    const response = await fetch('/api/v1/admin/notifications', {
      credentials: 'include',
    })
    const data = await response.json()
    if (!response.ok) {
      throw new Error(data.error || 'Failed to load notifications')
    }
    notifications.value = Array.isArray(data.notifications) ? data.notifications : []
    overlayPath.value = data.overlayPath || overlayPath.value
    if (selectedNotificationType.value) {
      const selected = notifications.value.find((notification) => notification.type === selectedNotificationType.value)
      if (selected) {
        selectNotification(selected)
      } else if (notifications.value.length > 0) {
        selectNotification(notifications.value[0])
      } else {
        selectedNotificationType.value = ''
        notificationConfigText.value = '{}'
      }
    } else if (notifications.value.length > 0) {
      selectNotification(notifications.value[0])
    }
  } catch (error) {
    notificationMessage.value = error.message
  } finally {
    loadingNotifications.value = false
  }
}

const saveNotification = async () => {
  if (!selectedNotification.value) {
    notificationMessage.value = 'Please select a notification type first.'
    return
  }
  savingNotification.value = true
  notificationMessage.value = ''
  try {
    const parsed = JSON.parse(notificationConfigText.value || '{}')
    if (!parsed || typeof parsed !== 'object' || Array.isArray(parsed)) {
      throw new Error('Notification config must be a JSON object.')
    }
    const response = await fetch(`/api/v1/admin/notifications/${selectedNotification.value.type}`, {
      method: 'PUT',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(parsed),
    })
    const data = await response.json()
    if (!response.ok) {
      throw new Error(data.error || 'Failed to save notification')
    }
    const currentType = selectedNotification.value.type
    await loadNotifications()
    const selected = notifications.value.find((notification) => notification.type === currentType)
    if (selected) {
      selectNotification(selected)
    }
    notificationMessage.value = 'Notification saved.'
  } catch (error) {
    notificationMessage.value = error.message
  } finally {
    savingNotification.value = false
  }
}

const deleteNotification = async () => {
  if (!selectedNotification.value || !selectedNotification.value.configured) {
    return
  }
  const confirmed = window.confirm(`Delete notification type ${selectedNotification.value.type}?`)
  if (!confirmed) return
  savingNotification.value = true
  notificationMessage.value = ''
  try {
    const response = await fetch(`/api/v1/admin/notifications/${selectedNotification.value.type}`, {
      method: 'DELETE',
      credentials: 'include',
    })
    if (!response.ok && response.status !== 204) {
      let message = 'Failed to delete notification'
      try {
        const data = await response.json()
        message = data.error || message
      } catch {
        // noop
      }
      throw new Error(message)
    }
    await loadNotifications()
    notificationMessage.value = 'Notification deleted.'
  } catch (error) {
    notificationMessage.value = error.message
  } finally {
    savingNotification.value = false
  }
}

const loadManagedConfig = async () => {
  loadingManaged.value = true
  managedMessage.value = ''
  try {
    const response = await fetch('/api/v1/admin/managed-config', {
      credentials: 'include',
    })
    const data = await response.json()
    if (!response.ok) {
      throw new Error(data.error || 'Failed to load managed configuration')
    }
    overlayPath.value = data.overlayPath || overlayPath.value
    jsonText.value = JSON.stringify(normalizePayload(data), null, 2)
  } catch (error) {
    managedMessage.value = error.message
  } finally {
    loadingManaged.value = false
  }
}

const saveManagedConfig = async () => {
  savingManaged.value = true
  managedMessage.value = ''
  try {
    const parsed = JSON.parse(jsonText.value || '{}')
    const payload = normalizePayload(parsed)
    const response = await fetch('/api/v1/admin/managed-config', {
      method: 'PUT',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(payload),
    })
    const data = await response.json()
    if (!response.ok) {
      throw new Error(data.error || 'Failed to save managed configuration')
    }
    overlayPath.value = data.overlayPath || overlayPath.value
    managedMessage.value = data.message || 'Overlay saved.'
    await Promise.all([loadEndpoints(), loadNotifications()])
  } catch (error) {
    managedMessage.value = error.message
  } finally {
    savingManaged.value = false
  }
}

const resetOverlay = async () => {
  const confirmed = window.confirm('Delete managed overlay and revert to base YAML configuration?')
  if (!confirmed) return
  savingManaged.value = true
  managedMessage.value = ''
  try {
    const response = await fetch('/api/v1/admin/managed-config', {
      method: 'DELETE',
      credentials: 'include',
    })
    if (!response.ok && response.status !== 204) {
      let message = 'Failed to reset overlay'
      try {
        const data = await response.json()
        message = data.error || message
      } catch {
        // noop
      }
      throw new Error(message)
    }
    managedMessage.value = 'Overlay deleted. Base YAML configuration will be active after reload.'
    startCreate()
    await Promise.all([loadEndpoints(), loadNotifications(), loadManagedConfig()])
  } catch (error) {
    managedMessage.value = error.message
  } finally {
    savingManaged.value = false
  }
}

onMounted(async () => {
  await Promise.all([loadEndpoints(), loadNotifications(), loadManagedConfig()])
})
</script>
