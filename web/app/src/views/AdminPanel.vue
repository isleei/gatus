<template>
  <div class="container mx-auto px-4 py-8 max-w-7xl">
    <div class="mb-6">
      <h1 class="text-3xl font-bold tracking-tight">{{ t('admin.title') }}</h1>
      <p class="text-muted-foreground mt-2">
        {{ t('admin.subtitle') }}
      </p>
    </div>

    <div class="p-4 rounded-lg border bg-card mb-6">
      <p class="text-sm">
        <span class="font-semibold">{{ t('admin.overlayPath') }}</span> {{ overlayPath || t('common.noData') }}
      </p>
      <p class="text-muted-foreground text-sm mt-1">
        {{ t('admin.autoAppliedHint') }}
      </p>
    </div>

    <div class="grid gap-6 lg:grid-cols-[22rem,1fr]">
      <Card>
        <CardHeader>
          <div class="flex items-center justify-between gap-2">
            <CardTitle>{{ t('common.endpoints') }}</CardTitle>
            <div class="flex items-center gap-2">
              <Button variant="outline" size="sm" @click="startCreate" :disabled="savingEndpoint">{{ t('admin.new') }}</Button>
              <Button variant="outline" size="sm" @click="loadEndpoints" :disabled="loadingEndpoints || savingEndpoint">{{ t('common.refresh') }}</Button>
            </div>
          </div>
        </CardHeader>
        <CardContent>
          <p v-if="loadingEndpoints" class="text-sm text-muted-foreground">{{ t('admin.loadingEndpoints') }}</p>
          <p v-else-if="endpoints.length === 0" class="text-sm text-muted-foreground">{{ t('admin.noEndpoints') }}</p>
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
              <p class="text-xs text-muted-foreground mt-1">{{ t('admin.endpointKey', { key: endpoint.key }) }}</p>
            </button>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>{{ isEditing ? t('admin.editEndpoint', { key: selectedKey }) : t('admin.createEndpoint') }}</CardTitle>
        </CardHeader>
        <CardContent class="space-y-4">
          <div class="grid gap-4 md:grid-cols-2">
            <div class="space-y-1">
              <label class="text-sm font-medium">{{ t('common.name') }}</label>
              <Input v-model="form.name" placeholder="frontend" />
            </div>
            <div class="space-y-1">
              <label class="text-sm font-medium">{{ t('common.group') }}</label>
              <Input v-model="form.group" placeholder="core" />
            </div>
            <div class="space-y-1 md:col-span-2">
              <label class="text-sm font-medium">{{ t('common.url') }}</label>
              <Input v-model="form.url" placeholder="https://example.org/health" />
            </div>
            <div class="space-y-1">
              <label class="text-sm font-medium">{{ t('common.method') }}</label>
              <Select v-model="form.method" :options="methodOptions" />
            </div>
            <div class="space-y-1">
              <label class="text-sm font-medium">{{ t('common.interval') }}</label>
              <Input v-model="form.interval" placeholder="30s / 1m / 5m" />
            </div>
          </div>

          <div class="grid gap-4 md:grid-cols-2">
            <div class="space-y-1">
              <label class="text-sm font-medium">{{ t('admin.conditionsOnePerLine') }}</label>
              <textarea
                v-model="form.conditionsText"
                class="w-full min-h-[180px] rounded-md border bg-background p-3 text-sm font-mono"
                spellcheck="false"
              />
            </div>
            <div class="space-y-1">
              <label class="text-sm font-medium">{{ t('admin.headersKeyValue') }}</label>
              <textarea
                v-model="form.headersText"
                class="w-full min-h-[180px] rounded-md border bg-background p-3 text-sm font-mono"
                spellcheck="false"
                placeholder="Authorization: Bearer token"
              />
            </div>
          </div>

          <div class="space-y-1">
            <label class="text-sm font-medium">{{ t('common.body') }}</label>
            <textarea
              v-model="form.body"
              class="w-full min-h-[120px] rounded-md border bg-background p-3 text-sm font-mono"
              spellcheck="false"
            />
          </div>

          <div class="flex flex-wrap items-center gap-4">
            <label class="inline-flex items-center gap-2 text-sm">
              <input v-model="form.enabled" type="checkbox" class="h-4 w-4 rounded border-input" />
              {{ t('admin.enabled') }}
            </label>
            <label class="inline-flex items-center gap-2 text-sm">
              <input v-model="form.graphql" type="checkbox" class="h-4 w-4 rounded border-input" />
              {{ t('admin.graphqlBody') }}
            </label>
          </div>

          <div class="flex flex-wrap items-center gap-2">
            <Button @click="saveEndpoint" :disabled="savingEndpoint || loadingEndpoints">
              {{ isEditing ? t('admin.saveEndpoint') : t('admin.createEndpointAction') }}
            </Button>
            <Button variant="outline" @click="startCreate" :disabled="savingEndpoint">{{ t('admin.clear') }}</Button>
            <Button
              v-if="isEditing"
              variant="destructive"
              @click="deleteEndpoint"
              :disabled="savingEndpoint"
            >
              {{ t('admin.deleteEndpoint') }}
            </Button>
            <span v-if="endpointMessage" class="text-sm text-muted-foreground">{{ endpointMessage }}</span>
          </div>
        </CardContent>
      </Card>
    </div>

    <Card class="mt-6">
      <CardHeader>
        <div class="flex items-center justify-between gap-2">
          <CardTitle>{{ t('admin.notificationChannels') }}</CardTitle>
          <Button variant="outline" size="sm" @click="loadNotifications" :disabled="loadingNotifications || savingNotification">{{ t('common.refresh') }}</Button>
        </div>
      </CardHeader>
      <CardContent>
        <div class="grid gap-6 lg:grid-cols-[22rem,1fr]">
          <div class="space-y-2">
            <p v-if="loadingNotifications" class="text-sm text-muted-foreground">{{ t('admin.loadingNotifications') }}</p>
            <p v-else-if="notifications.length === 0" class="text-sm text-muted-foreground">{{ t('admin.noNotificationProvider') }}</p>
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
                  {{ notification.configured ? t('admin.configured') : t('admin.notConfigured') }}
                </span>
              </div>
              <p class="text-xs text-muted-foreground mt-1">
                {{ t('admin.usedBy', { endpoints: notification.usedByEndpoints, externalEndpoints: notification.usedByExternalEndpoints }) }}
              </p>
            </button>
          </div>

          <div class="space-y-3">
            <p v-if="!selectedNotification" class="text-sm text-muted-foreground">{{ t('admin.selectNotificationType') }}</p>
            <template v-else>
              <div class="space-y-1">
                <label class="text-sm font-medium">{{ t('common.type') }}</label>
                <Input :model-value="selectedNotification.type" disabled />
              </div>
              <div class="space-y-1">
                <label class="text-sm font-medium">{{ t('admin.providerConfigJson') }}</label>
                <textarea
                  v-model="notificationConfigText"
                  class="w-full min-h-[240px] rounded-md border bg-background p-3 text-sm font-mono"
                  spellcheck="false"
                />
              </div>
              <div class="flex flex-wrap items-center gap-2">
                <Button @click="saveNotification" :disabled="savingNotification || loadingNotifications">{{ t('admin.saveNotification') }}</Button>
                <Button
                  variant="destructive"
                  @click="deleteNotification"
                  :disabled="savingNotification || !selectedNotification.configured"
                >
                  {{ t('admin.deleteNotification') }}
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
          <CardTitle>{{ t('admin.advancedJsonOverlay') }}</CardTitle>
          <div class="flex items-center gap-2">
            <Button variant="outline" size="sm" @click="loadManagedConfig" :disabled="loadingManaged || savingManaged || applyingManaged">{{ t('admin.reload') }}</Button>
            <Button size="sm" @click="saveManagedConfig" :disabled="loadingManaged || savingManaged || applyingManaged">{{ t('admin.saveOverlay') }}</Button>
            <Button variant="secondary" size="sm" @click="applyManagedConfig" :disabled="loadingManaged || savingManaged || applyingManaged">{{ t('admin.applyNow') }}</Button>
            <Button variant="destructive" size="sm" @click="resetOverlay" :disabled="loadingManaged || savingManaged || applyingManaged">{{ t('admin.resetOverlay') }}</Button>
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
import { useI18n } from '@/i18n'

const { t } = useI18n()

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
const applyingManaged = ref(false)
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
      throw new Error(t('admin.invalidHeader', { line }))
    }
    const name = line.slice(0, separatorIndex).trim()
    const value = line.slice(separatorIndex + 1).trim()
    if (!name) {
      throw new Error(t('admin.invalidHeaderName', { line }))
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
    throw new Error(t('admin.nameUrlRequired'))
  }
  if (conditions.length === 0) {
    throw new Error(t('admin.atLeastOneCondition'))
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
      throw new Error(data.error || t('admin.failedLoadEndpoints'))
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
      throw new Error(data.error || t('admin.failedSaveEndpoint'))
    }
    selectedKey.value = data.key || selectedKey.value
    await loadEndpoints()
    endpointMessage.value = isEditing.value ? t('admin.endpointSaved') : t('admin.endpointCreated')
  } catch (error) {
    endpointMessage.value = error.message
  } finally {
    savingEndpoint.value = false
  }
}

const deleteEndpoint = async () => {
  if (!isEditing.value) return
  const confirmed = window.confirm(t('admin.confirmDeleteEndpoint', { key: selectedKey.value }))
  if (!confirmed) return
  savingEndpoint.value = true
  endpointMessage.value = ''
  try {
    const response = await fetch(`/api/v1/admin/endpoints/${selectedKey.value}`, {
      method: 'DELETE',
      credentials: 'include',
    })
    if (!response.ok && response.status !== 204) {
      let message = t('admin.failedDeleteEndpoint')
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
    endpointMessage.value = t('admin.endpointDeleted')
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
      throw new Error(data.error || t('admin.failedLoadNotifications'))
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
    notificationMessage.value = t('admin.selectNotificationFirst')
    return
  }
  savingNotification.value = true
  notificationMessage.value = ''
  try {
    const parsed = JSON.parse(notificationConfigText.value || '{}')
    if (!parsed || typeof parsed !== 'object' || Array.isArray(parsed)) {
      throw new Error(t('admin.notificationConfigJsonObject'))
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
      throw new Error(data.error || t('admin.failedSaveNotification'))
    }
    const currentType = selectedNotification.value.type
    await loadNotifications()
    const selected = notifications.value.find((notification) => notification.type === currentType)
    if (selected) {
      selectNotification(selected)
    }
    notificationMessage.value = t('admin.notificationSaved')
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
  const confirmed = window.confirm(t('admin.confirmDeleteNotification', { type: selectedNotification.value.type }))
  if (!confirmed) return
  savingNotification.value = true
  notificationMessage.value = ''
  try {
    const response = await fetch(`/api/v1/admin/notifications/${selectedNotification.value.type}`, {
      method: 'DELETE',
      credentials: 'include',
    })
    if (!response.ok && response.status !== 204) {
      let message = t('admin.failedDeleteNotification')
      try {
        const data = await response.json()
        message = data.error || message
      } catch {
        // noop
      }
      throw new Error(message)
    }
    await loadNotifications()
    notificationMessage.value = t('admin.notificationDeleted')
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
      throw new Error(data.error || t('admin.failedLoadManagedConfig'))
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
      throw new Error(data.error || t('admin.failedSaveManagedConfig'))
    }
    overlayPath.value = data.overlayPath || overlayPath.value
    managedMessage.value = data.message || t('admin.overlaySaved')
    await Promise.all([loadEndpoints(), loadNotifications()])
  } catch (error) {
    managedMessage.value = error.message
  } finally {
    savingManaged.value = false
  }
}

const applyManagedConfig = async () => {
  applyingManaged.value = true
  managedMessage.value = ''
  try {
    const response = await fetch('/api/v1/admin/reload', {
      method: 'POST',
      credentials: 'include',
    })
    const data = await response.json()
    if (!response.ok) {
      throw new Error(data.error || t('admin.failedApplyImmediately'))
    }
    managedMessage.value = data.message || t('admin.immediateReloadRequested')
    await Promise.all([loadEndpoints(), loadNotifications()])
  } catch (error) {
    managedMessage.value = error.message
  } finally {
    applyingManaged.value = false
  }
}

const resetOverlay = async () => {
  const confirmed = window.confirm(t('admin.confirmResetOverlay'))
  if (!confirmed) return
  savingManaged.value = true
  managedMessage.value = ''
  try {
    const response = await fetch('/api/v1/admin/managed-config', {
      method: 'DELETE',
      credentials: 'include',
    })
    if (!response.ok && response.status !== 204) {
      let message = t('admin.failedResetOverlay')
      try {
        const data = await response.json()
        message = data.error || message
      } catch {
        // noop
      }
      throw new Error(message)
    }
    managedMessage.value = t('admin.overlayDeletedHint')
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
