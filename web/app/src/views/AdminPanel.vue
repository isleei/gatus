<template>
  <div class="container mx-auto px-4 py-6 max-w-7xl space-y-4">
    <div>
      <h1 class="text-3xl font-bold tracking-tight">{{ t('adminV2.title') }}</h1>
      <p class="text-muted-foreground mt-1">{{ t('adminV2.subtitle') }}</p>
    </div>

    <div class="p-3 rounded-lg border bg-card text-sm text-muted-foreground">
      <span class="font-medium text-foreground">{{ t('admin.overlayPath') }}</span>
      <span class="ml-2">{{ overlayPath || t('common.noData') }}</span>
    </div>

    <div class="grid grid-cols-2 lg:grid-cols-4 gap-3">
      <Card>
        <CardContent class="p-4">
          <p class="text-xs text-muted-foreground">{{ t('adminV2.kpiTotal') }}</p>
          <p class="text-2xl font-semibold">{{ monitorsKpi.total }}</p>
        </CardContent>
      </Card>
      <Card>
        <CardContent class="p-4">
          <p class="text-xs text-muted-foreground">{{ t('adminV2.kpiUnhealthy') }}</p>
          <p class="text-2xl font-semibold text-red-600">{{ monitorsKpi.unhealthy }}</p>
        </CardContent>
      </Card>
      <Card>
        <CardContent class="p-4">
          <p class="text-xs text-muted-foreground">{{ t('adminV2.kpiDisabled') }}</p>
          <p class="text-2xl font-semibold">{{ monitorsKpi.disabled }}</p>
        </CardContent>
      </Card>
      <Card>
        <CardContent class="p-4">
          <p class="text-xs text-muted-foreground">{{ t('adminV2.kpiUnknown') }}</p>
          <p class="text-2xl font-semibold">{{ monitorsKpi.unknown }}</p>
        </CardContent>
      </Card>
    </div>

    <div class="flex flex-wrap gap-2">
      <Button :variant="activeTab === 'monitors' ? 'default' : 'outline'" size="sm" @click="activeTab = 'monitors'">{{ t('adminV2.tabMonitors') }}</Button>
      <Button :variant="activeTab === 'audit' ? 'default' : 'outline'" size="sm" @click="activeTab = 'audit'">{{ t('adminV2.tabAudit') }}</Button>
      <Button :variant="activeTab === 'advanced' ? 'default' : 'outline'" size="sm" @click="activeTab = 'advanced'">{{ t('adminV2.tabAdvanced') }}</Button>
    </div>

    <template v-if="activeTab === 'monitors'">
      <Card>
        <CardContent class="p-4 space-y-3">
          <div class="grid gap-3 md:grid-cols-2 xl:grid-cols-4">
            <Input v-model="filters.q" :placeholder="t('adminV2.searchPlaceholder')" @keyup.enter="refreshMonitors" />
            <Select v-model="filters.entityType" :options="entityTypeOptions" @update:model-value="onFilterChange" />
            <Select v-model="filters.status" :options="statusFilterOptions" @update:model-value="onFilterChange" />
            <Select v-model="filters.enabled" :options="enabledFilterOptions" @update:model-value="onFilterChange" />
            <Select v-model="filters.group" :options="groupFilterOptions" @update:model-value="onFilterChange" />
            <Select v-model="filters.sortBy" :options="sortByOptions" @update:model-value="onFilterChange" />
            <Select v-model="filters.sortDir" :options="sortDirOptions" @update:model-value="onFilterChange" />
            <Select v-model="filters.pageSize" :options="pageSizeOptions" @update:model-value="onFilterChange" />
          </div>
          <div class="flex flex-wrap gap-2">
            <Button size="sm" variant="outline" @click="refreshMonitors" :disabled="loadingMonitors">{{ t('common.refresh') }}</Button>
            <Button size="sm" variant="secondary" @click="openCreate('endpoint')">{{ t('adminV2.newEndpoint') }}</Button>
            <Button size="sm" variant="secondary" @click="openCreate('suite')">{{ t('adminV2.newSuite') }}</Button>
            <Button size="sm" variant="secondary" @click="openCreate('external')">{{ t('adminV2.newExternal') }}</Button>
            <Button size="sm" variant="outline" @click="openImportExport">{{ t('adminV2.importExport') }}</Button>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardContent class="p-0">
          <div v-if="loadingMonitors" class="p-6 text-sm text-muted-foreground">{{ t('common.loading') }}</div>
          <div v-else>
            <div class="hidden lg:block overflow-x-auto">
              <table class="w-full text-sm">
                <thead class="bg-muted/50">
                  <tr>
                    <th class="p-3 text-left w-10">
                      <input type="checkbox" :checked="allPageSelected" @change="toggleAllPageSelection($event.target.checked)" />
                    </th>
                    <th class="p-3 text-left">{{ t('common.status') }}</th>
                    <th class="p-3 text-left">{{ t('common.name') }}</th>
                    <th class="p-3 text-left">{{ t('common.group') }}</th>
                    <th class="p-3 text-left">{{ t('common.type') }}</th>
                    <th class="p-3 text-left">{{ t('common.url') }}</th>
                    <th class="p-3 text-left">{{ t('common.interval') }}</th>
                    <th class="p-3 text-left">{{ t('adminV2.lastCheck') }}</th>
                    <th class="p-3 text-left">{{ t('common.duration') }}</th>
                    <th class="p-3 text-left">{{ t('adminV2.notifications') }}</th>
                    <th class="p-3 text-left">{{ t('adminV2.actions') }}</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="monitor in monitors" :key="monitor.entityType + '-' + monitor.key" class="border-t">
                    <td class="p-3 align-top">
                      <input type="checkbox" :checked="selectedKeySet.has(monitor.key)" @change="toggleSelection(monitor.key, $event.target.checked)" />
                    </td>
                    <td class="p-3 align-top">
                      <span class="inline-flex items-center rounded-full px-2 py-0.5 text-xs" :class="statusClass(monitor.status)">
                        {{ monitor.status }}
                      </span>
                    </td>
                    <td class="p-3 align-top font-medium">{{ monitor.name }}</td>
                    <td class="p-3 align-top">{{ monitor.group || '-' }}</td>
                    <td class="p-3 align-top">{{ monitor.type }}</td>
                    <td class="p-3 align-top break-all">{{ monitor.steps > 0 ? `${monitor.steps} steps` : monitor.url || '-' }}</td>
                    <td class="p-3 align-top">{{ monitor.interval || '-' }}</td>
                    <td class="p-3 align-top">{{ formatDateTime(monitor.lastCheck) }}</td>
                    <td class="p-3 align-top">{{ monitor.duration || '-' }}</td>
                    <td class="p-3 align-top">{{ (monitor.notificationTypes || []).join(', ') || '-' }}</td>
                    <td class="p-3 align-top">
                      <div class="flex flex-wrap gap-1">
                        <Button size="sm" variant="outline" @click="openEdit(monitor)">{{ t('adminV2.edit') }}</Button>
                        <Button size="sm" variant="outline" @click="openCopy(monitor)">{{ t('adminV2.copy') }}</Button>
                        <Button size="sm" variant="outline" @click="toggleEnabled(monitor)">{{ monitor.enabled ? t('adminV2.disable') : t('adminV2.enable') }}</Button>
                        <Button size="sm" variant="destructive" @click="deleteMonitor(monitor)">{{ t('adminV2.delete') }}</Button>
                      </div>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>

            <div class="lg:hidden space-y-3 p-3">
              <div v-for="monitor in monitors" :key="monitor.entityType + '-' + monitor.key" class="rounded-lg border p-3 space-y-2">
                <div class="flex items-center justify-between gap-2">
                  <div class="flex items-center gap-2">
                    <input type="checkbox" :checked="selectedKeySet.has(monitor.key)" @change="toggleSelection(monitor.key, $event.target.checked)" />
                    <p class="font-medium">{{ monitor.name }}</p>
                  </div>
                  <span class="inline-flex items-center rounded-full px-2 py-0.5 text-xs" :class="statusClass(monitor.status)">
                    {{ monitor.status }}
                  </span>
                </div>
                <p class="text-xs text-muted-foreground">{{ monitor.group || '-' }} · {{ monitor.type }}</p>
                <p class="text-xs break-all">{{ monitor.steps > 0 ? `${monitor.steps} steps` : monitor.url || '-' }}</p>
                <div class="flex flex-wrap gap-1">
                  <Button size="sm" variant="outline" @click="openEdit(monitor)">{{ t('adminV2.edit') }}</Button>
                  <Button size="sm" variant="outline" @click="openCopy(monitor)">{{ t('adminV2.copy') }}</Button>
                  <Button size="sm" variant="outline" @click="toggleEnabled(monitor)">{{ monitor.enabled ? t('adminV2.disable') : t('adminV2.enable') }}</Button>
                  <Button size="sm" variant="destructive" @click="deleteMonitor(monitor)">{{ t('adminV2.delete') }}</Button>
                </div>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      <Card v-if="selectedKeys.length > 0">
        <CardContent class="p-4 space-y-3">
          <p class="text-sm text-muted-foreground">{{ t('adminV2.selectedCount', { count: selectedKeys.length }) }}</p>
          <div class="grid gap-3 md:grid-cols-3">
            <Select v-model="batch.action" :options="batchActionOptions" />
            <Input v-if="batch.action === 'set-group'" v-model="batch.group" :placeholder="t('adminV2.groupPlaceholder')" />
            <Input v-if="batch.action === 'set-interval'" v-model="batch.interval" :placeholder="t('adminV2.intervalPlaceholder')" />
            <Input v-if="batch.action === 'set-alert-types'" v-model="batch.alertTypes" :placeholder="t('adminV2.alertTypesPlaceholder')" />
          </div>
          <div class="flex flex-wrap gap-2">
            <Button size="sm" @click="executeBatch(false)" :disabled="batchLoading">{{ t('adminV2.applyBatch') }}</Button>
            <Button size="sm" variant="outline" @click="executeBatch(true)" :disabled="batchLoading">{{ t('adminV2.dryRunBatch') }}</Button>
            <Button size="sm" variant="outline" @click="clearSelection">{{ t('adminV2.clearSelection') }}</Button>
          </div>
          <p v-if="batchMessage" class="text-sm text-muted-foreground">{{ batchMessage }}</p>
        </CardContent>
      </Card>

      <div class="flex items-center justify-between text-sm">
        <Button size="sm" variant="outline" :disabled="filters.page <= 1" @click="filters.page--; refreshMonitors()">{{ t('pagination.previous') }}</Button>
        <span>{{ t('pagination.pageOf', { current: filters.page, total: monitorTotalPages }) }}</span>
        <Button size="sm" variant="outline" :disabled="filters.page >= monitorTotalPages" @click="filters.page++; refreshMonitors()">{{ t('pagination.next') }}</Button>
      </div>
    </template>

    <template v-if="activeTab === 'audit'">
      <Card>
        <CardContent class="p-4 space-y-3">
          <div class="grid gap-3 md:grid-cols-2 xl:grid-cols-5">
            <Input v-model="auditFilters.actor" :placeholder="t('adminV2.actor')" @keyup.enter="refreshAuditLogs" />
            <Input v-model="auditFilters.action" :placeholder="t('adminV2.action')" @keyup.enter="refreshAuditLogs" />
            <Select v-model="auditFilters.result" :options="auditResultOptions" />
            <Input v-model="auditFilters.q" :placeholder="t('adminV2.searchPlaceholder')" @keyup.enter="refreshAuditLogs" />
            <Button variant="outline" @click="refreshAuditLogs" :disabled="loadingAudit">{{ t('common.refresh') }}</Button>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardContent class="p-0 overflow-x-auto">
          <table class="w-full text-sm">
            <thead class="bg-muted/50">
              <tr>
                <th class="p-3 text-left">{{ t('common.timestamp') }}</th>
                <th class="p-3 text-left">{{ t('adminV2.actor') }}</th>
                <th class="p-3 text-left">{{ t('adminV2.action') }}</th>
                <th class="p-3 text-left">{{ t('common.type') }}</th>
                <th class="p-3 text-left">{{ t('common.status') }}</th>
                <th class="p-3 text-left">{{ t('common.errors') }}</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="item in auditLogs" :key="item.id" class="border-t">
                <td class="p-3">{{ formatDateTime(item.timestamp) }}</td>
                <td class="p-3">{{ item.actor }}</td>
                <td class="p-3">{{ item.action }} {{ item.entityType }} {{ item.entityKey || '' }}</td>
                <td class="p-3">{{ item.entityType }}</td>
                <td class="p-3">{{ item.result }}</td>
                <td class="p-3 text-xs text-muted-foreground break-all">{{ item.error || '-' }}</td>
              </tr>
            </tbody>
          </table>
        </CardContent>
      </Card>

      <div class="flex items-center justify-between text-sm">
        <Button size="sm" variant="outline" :disabled="auditFilters.page <= 1" @click="auditFilters.page--; refreshAuditLogs()">{{ t('pagination.previous') }}</Button>
        <span>{{ t('pagination.pageOf', { current: auditFilters.page, total: auditTotalPages }) }}</span>
        <Button size="sm" variant="outline" :disabled="auditFilters.page >= auditTotalPages" @click="auditFilters.page++; refreshAuditLogs()">{{ t('pagination.next') }}</Button>
      </div>
    </template>

    <template v-if="activeTab === 'advanced'">
      <Card>
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
    </template>

    <div v-if="drawerOpen" class="fixed inset-0 z-50 flex justify-end bg-black/40" @click.self="closeDrawer">
      <div class="h-full w-full md:w-[720px] bg-background border-l shadow-xl p-4 overflow-y-auto space-y-4">
        <div class="flex items-center justify-between gap-2">
          <h3 class="text-lg font-semibold">{{ drawerTitle }}</h3>
          <Button size="sm" variant="outline" @click="closeDrawer">{{ t('adminV2.close') }}</Button>
        </div>

        <div class="grid gap-3 md:grid-cols-2">
          <div>
            <label class="text-sm font-medium">{{ t('common.type') }}</label>
            <Select v-model="drawer.entityType" :options="entityTypeEditableOptions" :disabled="drawerMode === 'edit'" />
          </div>
          <div>
            <label class="text-sm font-medium">{{ t('adminV2.mode') }}</label>
            <Input :model-value="drawerMode" disabled />
          </div>
        </div>

        <template v-if="drawer.entityType === 'endpoint'">
          <div class="grid gap-3 md:grid-cols-2">
            <div>
              <label class="text-sm font-medium">{{ t('common.name') }}</label>
              <Input v-model="endpointForm.name" />
            </div>
            <div>
              <label class="text-sm font-medium">{{ t('common.group') }}</label>
              <Input v-model="endpointForm.group" />
            </div>
            <div class="md:col-span-2">
              <label class="text-sm font-medium">{{ t('common.url') }}</label>
              <Input v-model="endpointForm.url" />
            </div>
            <div>
              <label class="text-sm font-medium">{{ t('common.method') }}</label>
              <Select v-model="endpointForm.method" :options="methodOptions" />
            </div>
            <div>
              <label class="text-sm font-medium">{{ t('common.interval') }}</label>
              <Input v-model="endpointForm.interval" />
            </div>
          </div>
          <div class="grid gap-3 md:grid-cols-2">
            <div class="space-y-2 rounded-md border p-3">
              <p class="text-sm font-medium">{{ t('adminV2.probeSettings') }}</p>
              <div>
                <label class="text-xs text-muted-foreground">{{ t('adminV2.clientTimeout') }}</label>
                <Input v-model="endpointForm.clientTimeout" />
              </div>
              <label class="inline-flex items-center gap-2 text-sm">
                <input v-model="endpointForm.clientInsecure" type="checkbox" class="h-4 w-4 rounded border-input" />
                {{ t('adminV2.clientInsecure') }}
              </label>
              <label class="inline-flex items-center gap-2 text-sm">
                <input v-model="endpointForm.clientIgnoreRedirect" type="checkbox" class="h-4 w-4 rounded border-input" />
                {{ t('adminV2.clientIgnoreRedirect') }}
              </label>
            </div>
            <div class="space-y-2 rounded-md border p-3">
              <p class="text-sm font-medium">{{ t('adminV2.certificateSettings') }}</p>
              <label class="inline-flex items-center gap-2 text-sm">
                <input v-model="endpointForm.certificateEnabled" type="checkbox" class="h-4 w-4 rounded border-input" />
                {{ t('adminV2.certificateEnabled') }}
              </label>
              <div>
                <label class="text-xs text-muted-foreground">{{ t('adminV2.certificateThreshold') }}</label>
                <Input v-model="endpointForm.certificateExpirationThreshold" :disabled="!endpointForm.certificateEnabled" />
              </div>
              <p class="text-xs text-muted-foreground">{{ t('adminV2.certificateExpirationHint') }}</p>
            </div>
          </div>
          <div class="space-y-2 rounded-md border p-3">
            <p class="text-sm font-medium">{{ t('adminV2.tamperSettings') }}</p>
            <div class="grid gap-3 md:grid-cols-4">
              <label class="inline-flex items-center gap-2 text-sm md:col-span-4">
                <input v-model="endpointForm.tamperEnabled" type="checkbox" class="h-4 w-4 rounded border-input" />
                {{ t('adminV2.tamperEnabled') }}
              </label>
              <div>
                <label class="text-xs text-muted-foreground">{{ t('adminV2.tamperBaselineSamples') }}</label>
                <Input v-model="endpointForm.tamperBaselineSamples" type="number" min="1" :disabled="!endpointForm.tamperEnabled" />
              </div>
              <div>
                <label class="text-xs text-muted-foreground">{{ t('adminV2.tamperDriftThreshold') }}</label>
                <Input v-model="endpointForm.tamperDriftThresholdPercent" type="number" min="1" :disabled="!endpointForm.tamperEnabled" />
              </div>
              <div>
                <label class="text-xs text-muted-foreground">{{ t('adminV2.tamperConsecutiveBreaches') }}</label>
                <Input v-model="endpointForm.tamperConsecutiveBreaches" type="number" min="1" :disabled="!endpointForm.tamperEnabled" />
              </div>
              <div class="md:col-span-2">
                <label class="text-xs text-muted-foreground">{{ t('adminV2.tamperRequiredSubstrings') }}</label>
                <textarea
                  v-model="endpointForm.tamperRequiredSubstringsText"
                  class="mt-1 w-full min-h-[90px] rounded-md border bg-background p-2 text-sm font-mono"
                  spellcheck="false"
                  :disabled="!endpointForm.tamperEnabled"
                />
              </div>
              <div class="md:col-span-2">
                <label class="text-xs text-muted-foreground">{{ t('adminV2.tamperForbiddenSubstrings') }}</label>
                <textarea
                  v-model="endpointForm.tamperForbiddenSubstringsText"
                  class="mt-1 w-full min-h-[90px] rounded-md border bg-background p-2 text-sm font-mono"
                  spellcheck="false"
                  :disabled="!endpointForm.tamperEnabled"
                />
              </div>
            </div>
          </div>
          <div>
            <label class="text-sm font-medium">{{ t('admin.conditionsOnePerLine') }}</label>
            <textarea v-model="endpointForm.conditionsText" class="w-full min-h-[140px] rounded-md border bg-background p-2 text-sm font-mono" spellcheck="false" />
          </div>
          <div>
            <label class="text-sm font-medium">{{ t('admin.headersKeyValue') }}</label>
            <textarea v-model="endpointForm.headersText" class="w-full min-h-[120px] rounded-md border bg-background p-2 text-sm font-mono" spellcheck="false" />
          </div>
          <div>
            <label class="text-sm font-medium">{{ t('common.body') }}</label>
            <textarea v-model="endpointForm.body" class="w-full min-h-[100px] rounded-md border bg-background p-2 text-sm font-mono" spellcheck="false" />
          </div>
          <div class="grid gap-3 md:grid-cols-2">
            <div>
              <label class="text-sm font-medium">{{ t('adminV2.alertsJson') }}</label>
              <textarea
                v-model="endpointForm.alertsJson"
                class="w-full min-h-[120px] rounded-md border bg-background p-2 text-sm font-mono"
                :placeholder="t('adminV2.alertsJsonPlaceholder')"
                spellcheck="false"
              />
            </div>
            <div class="space-y-2 mt-1">
              <label class="inline-flex items-center gap-2 text-sm">
                <input v-model="endpointForm.resolveSuccessfulConditions" type="checkbox" class="h-4 w-4 rounded border-input" />
                {{ t('adminV2.resolveSuccessfulConditions') }}
              </label>
              <label class="inline-flex items-center gap-2 text-sm">
                <input v-model="endpointForm.dontResolveFailedConditions" type="checkbox" class="h-4 w-4 rounded border-input" />
                {{ t('adminV2.dontResolveFailedConditions') }}
              </label>
              <label class="inline-flex items-center gap-2 text-sm">
                <input v-model="endpointForm.enabled" type="checkbox" class="h-4 w-4 rounded border-input" />
                {{ t('admin.enabled') }}
              </label>
            </div>
          </div>
        </template>

        <template v-if="drawer.entityType === 'suite'">
          <div class="grid gap-3 md:grid-cols-2">
            <div>
              <label class="text-sm font-medium">{{ t('common.name') }}</label>
              <Input v-model="suiteForm.name" />
            </div>
            <div>
              <label class="text-sm font-medium">{{ t('common.group') }}</label>
              <Input v-model="suiteForm.group" />
            </div>
            <div>
              <label class="text-sm font-medium">{{ t('common.interval') }}</label>
              <Input v-model="suiteForm.interval" />
            </div>
            <div>
              <label class="text-sm font-medium">{{ t('adminV2.timeout') }}</label>
              <Input v-model="suiteForm.timeout" />
            </div>
          </div>
          <div>
            <label class="text-sm font-medium">{{ t('adminV2.contextJson') }}</label>
            <textarea v-model="suiteForm.contextJson" class="w-full min-h-[100px] rounded-md border bg-background p-2 text-sm font-mono" spellcheck="false" />
          </div>
          <div>
            <label class="text-sm font-medium">{{ t('adminV2.suiteEndpointsJson') }}</label>
            <textarea v-model="suiteForm.endpointsJson" class="w-full min-h-[220px] rounded-md border bg-background p-2 text-sm font-mono" spellcheck="false" />
          </div>
          <label class="inline-flex items-center gap-2 text-sm">
            <input v-model="suiteForm.enabled" type="checkbox" class="h-4 w-4 rounded border-input" />
            {{ t('admin.enabled') }}
          </label>
        </template>

        <template v-if="drawer.entityType === 'external'">
          <div class="grid gap-3 md:grid-cols-2">
            <div>
              <label class="text-sm font-medium">{{ t('common.name') }}</label>
              <Input v-model="externalForm.name" />
            </div>
            <div>
              <label class="text-sm font-medium">{{ t('common.group') }}</label>
              <Input v-model="externalForm.group" />
            </div>
            <div class="md:col-span-2">
              <label class="text-sm font-medium">{{ t('adminV2.token') }}</label>
              <Input v-model="externalForm.token" />
            </div>
            <div>
              <label class="text-sm font-medium">{{ t('adminV2.heartbeatInterval') }}</label>
              <Input v-model="externalForm.heartbeatInterval" />
            </div>
            <div>
              <label class="text-sm font-medium">{{ t('adminV2.alertsJson') }}</label>
              <textarea
                v-model="externalForm.alertsJson"
                class="w-full min-h-[120px] rounded-md border bg-background p-2 text-sm font-mono"
                :placeholder="t('adminV2.alertsJsonPlaceholder')"
                spellcheck="false"
              />
            </div>
          </div>
          <label class="inline-flex items-center gap-2 text-sm">
            <input v-model="externalForm.enabled" type="checkbox" class="h-4 w-4 rounded border-input" />
            {{ t('admin.enabled') }}
          </label>
        </template>

        <div class="flex flex-wrap gap-2">
          <Button @click="saveDrawer" :disabled="drawerSaving">{{ t('adminV2.save') }}</Button>
          <Button variant="outline" @click="closeDrawer" :disabled="drawerSaving">{{ t('adminV2.cancel') }}</Button>
          <Button v-if="drawerMode === 'edit'" variant="destructive" @click="deleteDrawer" :disabled="drawerSaving">{{ t('adminV2.delete') }}</Button>
        </div>
        <p v-if="drawerMessage" class="text-sm text-muted-foreground">{{ drawerMessage }}</p>
      </div>
    </div>

    <div v-if="importExportOpen" class="fixed inset-0 z-50 flex items-center justify-center bg-black/40" @click.self="importExportOpen = false">
      <div class="w-full max-w-5xl bg-background border rounded-lg shadow-xl p-4 max-h-[90vh] overflow-y-auto space-y-4">
        <div class="flex items-center justify-between gap-2">
          <h3 class="text-lg font-semibold">{{ t('adminV2.importExport') }}</h3>
          <Button size="sm" variant="outline" @click="importExportOpen = false">{{ t('adminV2.close') }}</Button>
        </div>

        <div class="grid gap-3 md:grid-cols-4">
          <Select v-model="importExport.entityType" :options="entityTypeOptions" />
          <Select v-model="importExport.mode" :options="importModeOptions" />
          <Button variant="outline" @click="loadExportData" :disabled="importExportLoading">{{ t('adminV2.exportData') }}</Button>
          <Button @click="runImport(true)" :disabled="importExportLoading">{{ t('adminV2.dryRunImport') }}</Button>
        </div>
        <div class="flex gap-2">
          <Button @click="runImport(false)" :disabled="importExportLoading">{{ t('adminV2.applyImport') }}</Button>
          <Button variant="outline" @click="copyExportToImport">{{ t('adminV2.copyExportToImport') }}</Button>
        </div>

        <div class="grid gap-3 lg:grid-cols-2">
          <div>
            <label class="text-sm font-medium">{{ t('adminV2.exportJson') }}</label>
            <textarea v-model="importExport.exportJson" class="w-full min-h-[320px] rounded-md border bg-background p-2 text-xs font-mono" spellcheck="false" />
          </div>
          <div>
            <label class="text-sm font-medium">{{ t('adminV2.importJson') }}</label>
            <textarea v-model="importExport.importJson" class="w-full min-h-[320px] rounded-md border bg-background p-2 text-xs font-mono" spellcheck="false" />
          </div>
        </div>
        <p v-if="importExportMessage" class="text-sm text-muted-foreground">{{ importExportMessage }}</p>
        <pre v-if="importPreview" class="text-xs bg-muted/40 p-3 rounded-md overflow-x-auto">{{ importPreview }}</pre>
      </div>
    </div>
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

const activeTab = ref('monitors')
const overlayPath = ref('')

const filters = ref({
  entityType: 'all',
  q: '',
  group: 'all',
  enabled: 'all',
  status: 'all',
  sortBy: 'updatedAt',
  sortDir: 'desc',
  page: 1,
  pageSize: '50',
})

const monitors = ref([])
const monitorTotal = ref(0)
const monitorsKpi = ref({ total: 0, unhealthy: 0, disabled: 0, unknown: 0 })
const monitorGroups = ref([])
const loadingMonitors = ref(false)

const selectedKeys = ref([])
const selectedKeySet = computed(() => new Set(selectedKeys.value))
const allPageSelected = computed(() => monitors.value.length > 0 && monitors.value.every((item) => selectedKeySet.value.has(item.key)))

const drawerOpen = ref(false)
const drawerMode = ref('create')
const drawerSaving = ref(false)
const drawerMessage = ref('')
const drawer = ref({
  entityType: 'endpoint',
  key: '',
})

const endpointForm = ref({
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
  alertsJson: '[]',
  resolveSuccessfulConditions: false,
  dontResolveFailedConditions: false,
  clientTimeout: '10s',
  clientInsecure: false,
  clientIgnoreRedirect: false,
  certificateEnabled: false,
  certificateExpirationThreshold: '72h',
  tamperEnabled: false,
  tamperBaselineSamples: 20,
  tamperDriftThresholdPercent: 20,
  tamperConsecutiveBreaches: 3,
  tamperRequiredSubstringsText: '',
  tamperForbiddenSubstringsText: '',
})

const suiteForm = ref({
  enabled: true,
  name: '',
  group: '',
  interval: '10m',
  timeout: '5m',
  contextJson: '{}',
  endpointsJson: '[\n  {\n    "name": "step-1",\n    "url": "https://example.org/health",\n    "method": "GET",\n    "conditions": ["[STATUS] == 200"]\n  }\n]',
})

const externalForm = ref({
  enabled: true,
  name: '',
  group: '',
  token: '',
  heartbeatInterval: '1m',
  alertsJson: '[]',
})

const batch = ref({
  action: 'enable',
  group: '',
  interval: '1m',
  alertTypes: '',
})
const batchLoading = ref(false)
const batchMessage = ref('')

const auditFilters = ref({
  page: 1,
  pageSize: 20,
  actor: '',
  action: '',
  result: 'all',
  q: '',
})
const auditLogs = ref([])
const auditTotal = ref(0)
const loadingAudit = ref(false)

const loadingManaged = ref(false)
const savingManaged = ref(false)
const applyingManaged = ref(false)
const managedMessage = ref('')
const jsonText = ref('')

const importExportOpen = ref(false)
const importExportLoading = ref(false)
const importExport = ref({
  entityType: 'all',
  mode: 'merge',
  exportJson: '',
  importJson: '',
})
const importPreview = ref('')
const importExportMessage = ref('')

const monitorTotalPages = computed(() => Math.max(1, Math.ceil(monitorTotal.value / Number(filters.value.pageSize || 50))))
const auditTotalPages = computed(() => Math.max(1, Math.ceil(auditTotal.value / auditFilters.value.pageSize)))

const methodOptions = [
  { label: 'GET', value: 'GET' },
  { label: 'POST', value: 'POST' },
  { label: 'PUT', value: 'PUT' },
  { label: 'PATCH', value: 'PATCH' },
  { label: 'DELETE', value: 'DELETE' },
  { label: 'HEAD', value: 'HEAD' },
  { label: 'OPTIONS', value: 'OPTIONS' },
]

const entityTypeOptions = [
  { label: t('adminV2.entityAll'), value: 'all' },
  { label: t('adminV2.entityEndpoint'), value: 'endpoint' },
  { label: t('adminV2.entitySuite'), value: 'suite' },
  { label: t('adminV2.entityExternal'), value: 'external' },
]

const entityTypeEditableOptions = [
  { label: t('adminV2.entityEndpoint'), value: 'endpoint' },
  { label: t('adminV2.entitySuite'), value: 'suite' },
  { label: t('adminV2.entityExternal'), value: 'external' },
]

const statusFilterOptions = [
  { label: t('search.none'), value: 'all' },
  { label: t('common.healthy'), value: 'healthy' },
  { label: t('common.unhealthy'), value: 'unhealthy' },
  { label: t('common.unknown'), value: 'unknown' },
  { label: t('adminV2.statusDisabled'), value: 'disabled' },
]

const enabledFilterOptions = [
  { label: t('search.none'), value: 'all' },
  { label: t('common.yes'), value: 'true' },
  { label: t('common.no'), value: 'false' },
]

const sortByOptions = [
  { label: t('adminV2.sortUpdatedAt'), value: 'updatedAt' },
  { label: t('search.name'), value: 'name' },
  { label: t('search.group'), value: 'group' },
  { label: t('search.health'), value: 'status' },
  { label: t('common.interval'), value: 'interval' },
  { label: t('common.duration'), value: 'duration' },
]

const sortDirOptions = [
  { label: t('adminV2.sortDesc'), value: 'desc' },
  { label: t('adminV2.sortAsc'), value: 'asc' },
]

const pageSizeOptions = [
  { label: '20', value: '20' },
  { label: '50', value: '50' },
  { label: '100', value: '100' },
  { label: '200', value: '200' },
]

const batchActionOptions = [
  { label: t('adminV2.batchEnable'), value: 'enable' },
  { label: t('adminV2.batchDisable'), value: 'disable' },
  { label: t('adminV2.batchSetGroup'), value: 'set-group' },
  { label: t('adminV2.batchSetInterval'), value: 'set-interval' },
  { label: t('adminV2.batchSetAlertTypes'), value: 'set-alert-types' },
  { label: t('adminV2.batchDelete'), value: 'delete' },
]

const auditResultOptions = [
  { label: t('search.none'), value: 'all' },
  { label: t('adminV2.resultSuccess'), value: 'success' },
  { label: t('adminV2.resultFailure'), value: 'failure' },
]

const importModeOptions = [
  { label: t('adminV2.importModeMerge'), value: 'merge' },
  { label: t('adminV2.importModeReplace'), value: 'replace' },
]

const groupFilterOptions = computed(() => {
  const options = [{ label: t('search.none'), value: 'all' }]
  monitorGroups.value.forEach((group) => {
    options.push({ label: group, value: group })
  })
  return options
})

const drawerTitle = computed(() => {
  const mode = drawerMode.value === 'edit' ? t('adminV2.edit') : drawerMode.value === 'copy' ? t('adminV2.copy') : t('adminV2.create')
  return `${mode} ${drawer.value.entityType}`
})

const onFilterChange = () => {
  filters.value.page = 1
  refreshMonitors()
}

const refreshMonitors = async () => {
  loadingMonitors.value = true
  try {
    const search = new URLSearchParams({
      entityType: filters.value.entityType,
      q: filters.value.q,
      group: filters.value.group,
      enabled: filters.value.enabled,
      status: filters.value.status,
      sortBy: filters.value.sortBy,
      sortDir: filters.value.sortDir,
      page: String(filters.value.page),
      pageSize: String(filters.value.pageSize),
    })
    const response = await fetch(`/api/v1/admin/monitors?${search.toString()}`, { credentials: 'include' })
    const data = await response.json()
    if (!response.ok) {
      throw new Error(data.error || t('admin.failedLoadEndpoints'))
    }
    monitors.value = Array.isArray(data.items) ? data.items : []
    monitorTotal.value = Number(data.total || 0)
    monitorsKpi.value = data.kpi || { total: 0, unhealthy: 0, disabled: 0, unknown: 0 }
    monitorGroups.value = Array.isArray(data.groups) ? data.groups : []
    if (filters.value.page > monitorTotalPages.value) {
      filters.value.page = monitorTotalPages.value
    }
  } catch (error) {
    batchMessage.value = error.message
  } finally {
    loadingMonitors.value = false
  }
}

const refreshAuditLogs = async () => {
  loadingAudit.value = true
  try {
    const search = new URLSearchParams({
      page: String(auditFilters.value.page),
      pageSize: String(auditFilters.value.pageSize),
      actor: auditFilters.value.actor,
      action: auditFilters.value.action,
      result: auditFilters.value.result,
      q: auditFilters.value.q,
    })
    const response = await fetch(`/api/v1/admin/audit-logs?${search.toString()}`, { credentials: 'include' })
    const data = await response.json()
    if (!response.ok) {
      throw new Error(data.error || 'Failed to load audit logs')
    }
    auditLogs.value = Array.isArray(data.items) ? data.items : []
    auditTotal.value = Number(data.total || 0)
  } catch (error) {
    managedMessage.value = error.message
  } finally {
    loadingAudit.value = false
  }
}

const toggleSelection = (key, checked) => {
  const set = new Set(selectedKeys.value)
  if (checked) {
    set.add(key)
  } else {
    set.delete(key)
  }
  selectedKeys.value = [...set]
}

const toggleAllPageSelection = (checked) => {
  if (!checked) {
    const set = new Set(selectedKeys.value)
    monitors.value.forEach((monitor) => set.delete(monitor.key))
    selectedKeys.value = [...set]
    return
  }
  const set = new Set(selectedKeys.value)
  monitors.value.forEach((monitor) => set.add(monitor.key))
  selectedKeys.value = [...set]
}

const clearSelection = () => {
  selectedKeys.value = []
}

const openCreate = (entityType) => {
  drawer.value = { entityType, key: '' }
  drawerMode.value = 'create'
  resetForms()
  drawerMessage.value = ''
  drawerOpen.value = true
}

const openCopy = async (monitor) => {
  await openEdit(monitor, true)
}

const openEdit = async (monitor, asCopy = false) => {
  drawer.value = { entityType: monitor.entityType, key: monitor.key }
  drawerMode.value = asCopy ? 'copy' : 'edit'
  drawerMessage.value = ''
  drawerOpen.value = true
  await loadFormFromServer(monitor.entityType, monitor.key)
  if (asCopy) {
    if (monitor.entityType === 'endpoint') endpointForm.value.name = `${endpointForm.value.name}-copy`
    if (monitor.entityType === 'suite') suiteForm.value.name = `${suiteForm.value.name}-copy`
    if (monitor.entityType === 'external') externalForm.value.name = `${externalForm.value.name}-copy`
  }
}

const closeDrawer = () => {
  drawerOpen.value = false
  drawerSaving.value = false
  drawerMessage.value = ''
}

const resetForms = () => {
  endpointForm.value = {
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
    alertsJson: '[]',
    resolveSuccessfulConditions: false,
    dontResolveFailedConditions: false,
    clientTimeout: '10s',
    clientInsecure: false,
    clientIgnoreRedirect: false,
    certificateEnabled: false,
    certificateExpirationThreshold: '72h',
    tamperEnabled: false,
    tamperBaselineSamples: 20,
    tamperDriftThresholdPercent: 20,
    tamperConsecutiveBreaches: 3,
    tamperRequiredSubstringsText: '',
    tamperForbiddenSubstringsText: '',
  }
  suiteForm.value = {
    enabled: true,
    name: '',
    group: '',
    interval: '10m',
    timeout: '5m',
    contextJson: '{}',
    endpointsJson: '[\n  {\n    "name": "step-1",\n    "url": "https://example.org/health",\n    "method": "GET",\n    "conditions": ["[STATUS] == 200"]\n  }\n]',
  }
  externalForm.value = {
    enabled: true,
    name: '',
    group: '',
    token: '',
    heartbeatInterval: '1m',
    alertsJson: '[]',
  }
}

const loadFormFromServer = async (entityType, key) => {
  try {
    if (entityType === 'endpoint') {
      const response = await fetch('/api/v1/admin/endpoints', { credentials: 'include' })
      const data = await response.json()
      const endpoint = (data.endpoints || []).find((item) => item.key === key)
      if (!endpoint) return
      endpointForm.value = {
        enabled: endpoint.enabled !== false,
        name: endpoint.name || '',
        group: endpoint.group || '',
        url: endpoint.url || '',
        method: endpoint.method || 'GET',
        interval: endpoint.interval || '1m',
        conditionsText: (endpoint.conditions || []).join('\n'),
        headersText: Object.entries(endpoint.headers || {}).map(([headerName, headerValue]) => `${headerName}: ${headerValue}`).join('\n'),
        body: endpoint.body || '',
        graphql: endpoint.graphql === true,
        alertsJson: JSON.stringify(endpoint.alerts || [], null, 2),
        resolveSuccessfulConditions: endpoint.ui?.resolveSuccessfulConditions === true,
        dontResolveFailedConditions: endpoint.ui?.dontResolveFailedConditions === true,
        clientTimeout: endpoint.client?.timeout || '10s',
        clientInsecure: endpoint.client?.insecure === true,
        clientIgnoreRedirect: endpoint.client?.ignoreRedirect === true,
        certificateEnabled: endpoint.certificate?.enabled === true,
        certificateExpirationThreshold: endpoint.certificate?.expirationThreshold || '72h',
        tamperEnabled: endpoint.tamper?.enabled === true,
        tamperBaselineSamples: endpoint.tamper?.baselineSamples || 20,
        tamperDriftThresholdPercent: endpoint.tamper?.driftThresholdPercent || 20,
        tamperConsecutiveBreaches: endpoint.tamper?.consecutiveBreaches || 3,
        tamperRequiredSubstringsText: (endpoint.tamper?.requiredSubstrings || []).join('\n'),
        tamperForbiddenSubstringsText: (endpoint.tamper?.forbiddenSubstrings || []).join('\n'),
      }
      overlayPath.value = data.overlayPath || overlayPath.value
      return
    }
    if (entityType === 'suite') {
      const response = await fetch('/api/v1/admin/suites', { credentials: 'include' })
      const data = await response.json()
      const suite = (data.suites || []).find((item) => item.key === key)
      if (!suite) return
      suiteForm.value = {
        enabled: suite.enabled !== false,
        name: suite.name || '',
        group: suite.group || '',
        interval: suite.interval || '10m',
        timeout: suite.timeout || '5m',
        contextJson: JSON.stringify(suite.context || {}, null, 2),
        endpointsJson: JSON.stringify(suite.endpoints || [], null, 2),
      }
      overlayPath.value = data.overlayPath || overlayPath.value
      return
    }
    if (entityType === 'external') {
      const response = await fetch('/api/v1/admin/external-endpoints', { credentials: 'include' })
      const data = await response.json()
      const externalEndpoint = (data.externalEndpoints || []).find((item) => item.key === key)
      if (!externalEndpoint) return
      externalForm.value = {
        enabled: externalEndpoint.enabled !== false,
        name: externalEndpoint.name || '',
        group: externalEndpoint.group || '',
        token: externalEndpoint.token || '',
        heartbeatInterval: externalEndpoint.heartbeatInterval || '1m',
        alertsJson: JSON.stringify(externalEndpoint.alerts || [], null, 2),
      }
      overlayPath.value = data.overlayPath || overlayPath.value
    }
  } catch (error) {
    drawerMessage.value = error.message
  }
}

const parseAlertsJSON = (text) => {
  const trimmed = (text || '').trim()
  if (!trimmed) return []
  let parsed
  try {
    parsed = JSON.parse(trimmed)
  } catch {
    throw new Error(t('adminV2.alertsJsonInvalid'))
  }
  if (!Array.isArray(parsed)) {
    throw new Error(t('adminV2.alertsJsonInvalid'))
  }
  return parsed
}

const parseHeaders = (text) => {
  const headers = {}
  text.split('\n').forEach((rawLine) => {
    const line = rawLine.trim()
    if (!line) return
    const index = line.indexOf(':')
    if (index <= 0) throw new Error(t('admin.invalidHeader', { line }))
    const headerName = line.slice(0, index).trim()
    const headerValue = line.slice(index + 1).trim()
    headers[headerName] = headerValue
  })
  return headers
}

const buildDrawerPayload = () => {
  if (drawer.value.entityType === 'endpoint') {
    const conditions = endpointForm.value.conditionsText.split('\n').map((item) => item.trim()).filter((item) => item)
    if (!endpointForm.value.name.trim() || !endpointForm.value.url.trim() || conditions.length === 0) {
      throw new Error(t('admin.nameUrlRequired'))
    }
    const tamperBaselineSamples = Math.max(1, Number.parseInt(endpointForm.value.tamperBaselineSamples, 10) || 20)
    const tamperDriftThresholdPercent = Math.max(1, Number.parseInt(endpointForm.value.tamperDriftThresholdPercent, 10) || 20)
    const tamperConsecutiveBreaches = Math.max(1, Number.parseInt(endpointForm.value.tamperConsecutiveBreaches, 10) || 3)
    const tamperRequiredSubstrings = endpointForm.value.tamperRequiredSubstringsText
      .split('\n')
      .map((item) => item.trim())
      .filter((item) => item)
    const tamperForbiddenSubstrings = endpointForm.value.tamperForbiddenSubstringsText
      .split('\n')
      .map((item) => item.trim())
      .filter((item) => item)
    return {
      enabled: endpointForm.value.enabled,
      name: endpointForm.value.name.trim(),
      group: endpointForm.value.group.trim(),
      url: endpointForm.value.url.trim(),
      method: endpointForm.value.method,
      interval: endpointForm.value.interval.trim(),
      conditions,
      headers: parseHeaders(endpointForm.value.headersText),
      body: endpointForm.value.body,
      graphql: endpointForm.value.graphql,
      alerts: parseAlertsJSON(endpointForm.value.alertsJson),
      ui: {
        resolveSuccessfulConditions: endpointForm.value.resolveSuccessfulConditions,
        dontResolveFailedConditions: endpointForm.value.dontResolveFailedConditions,
      },
      client: {
        timeout: endpointForm.value.clientTimeout.trim(),
        insecure: endpointForm.value.clientInsecure,
        ignoreRedirect: endpointForm.value.clientIgnoreRedirect,
      },
      certificate: {
        enabled: endpointForm.value.certificateEnabled,
        expirationThreshold: endpointForm.value.certificateExpirationThreshold.trim(),
      },
      tamper: {
        enabled: endpointForm.value.tamperEnabled,
        baselineSamples: tamperBaselineSamples,
        driftThresholdPercent: tamperDriftThresholdPercent,
        consecutiveBreaches: tamperConsecutiveBreaches,
        requiredSubstrings: tamperRequiredSubstrings,
        forbiddenSubstrings: tamperForbiddenSubstrings,
      },
    }
  }

  if (drawer.value.entityType === 'suite') {
    const context = JSON.parse(suiteForm.value.contextJson || '{}')
    const endpoints = JSON.parse(suiteForm.value.endpointsJson || '[]')
    if (!Array.isArray(endpoints) || endpoints.length === 0) {
      throw new Error(t('adminV2.suiteEndpointsRequired'))
    }
    return {
      enabled: suiteForm.value.enabled,
      name: suiteForm.value.name.trim(),
      group: suiteForm.value.group.trim(),
      interval: suiteForm.value.interval.trim(),
      timeout: suiteForm.value.timeout.trim(),
      context,
      endpoints,
    }
  }

  return {
    enabled: externalForm.value.enabled,
    name: externalForm.value.name.trim(),
    group: externalForm.value.group.trim(),
    token: externalForm.value.token.trim(),
    heartbeatInterval: externalForm.value.heartbeatInterval.trim(),
    alerts: parseAlertsJSON(externalForm.value.alertsJson),
  }
}

const saveDrawer = async () => {
  drawerSaving.value = true
  drawerMessage.value = ''
  try {
    const payload = buildDrawerPayload()
    let basePath = '/api/v1/admin/endpoints'
    if (drawer.value.entityType === 'suite') basePath = '/api/v1/admin/suites'
    if (drawer.value.entityType === 'external') basePath = '/api/v1/admin/external-endpoints'

    const isEdit = drawerMode.value === 'edit'
    const method = isEdit ? 'PUT' : 'POST'
    const targetURL = isEdit ? `${basePath}/${drawer.value.key}` : basePath

    const response = await fetch(targetURL, {
      method,
      credentials: 'include',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload),
    })
    const data = await response.json().catch(() => ({}))
    if (!response.ok) {
      throw new Error(data.error || t('admin.failedSaveEndpoint'))
    }
    drawerMessage.value = t('admin.endpointSaved')
    await refreshMonitors()
    await refreshAuditLogs()
    if (drawerMode.value !== 'edit') {
      closeDrawer()
    }
  } catch (error) {
    drawerMessage.value = error.message
  } finally {
    drawerSaving.value = false
  }
}

const deleteDrawer = async () => {
  if (drawerMode.value !== 'edit') return
  if (!window.confirm(t('adminV2.confirmDelete'))) return
  drawerSaving.value = true
  drawerMessage.value = ''
  try {
    let basePath = '/api/v1/admin/endpoints'
    if (drawer.value.entityType === 'suite') basePath = '/api/v1/admin/suites'
    if (drawer.value.entityType === 'external') basePath = '/api/v1/admin/external-endpoints'
    const response = await fetch(`${basePath}/${drawer.value.key}`, {
      method: 'DELETE',
      credentials: 'include',
    })
    if (!response.ok && response.status !== 204) {
      const data = await response.json().catch(() => ({}))
      throw new Error(data.error || t('admin.failedDeleteEndpoint'))
    }
    closeDrawer()
    await refreshMonitors()
    await refreshAuditLogs()
  } catch (error) {
    drawerMessage.value = error.message
  } finally {
    drawerSaving.value = false
  }
}

const deleteMonitor = async (monitor) => {
  if (!window.confirm(t('adminV2.confirmDelete'))) return
  try {
    let basePath = '/api/v1/admin/endpoints'
    if (monitor.entityType === 'suite') basePath = '/api/v1/admin/suites'
    if (monitor.entityType === 'external') basePath = '/api/v1/admin/external-endpoints'
    const response = await fetch(`${basePath}/${monitor.key}`, {
      method: 'DELETE',
      credentials: 'include',
    })
    if (!response.ok && response.status !== 204) {
      const data = await response.json().catch(() => ({}))
      throw new Error(data.error || t('admin.failedDeleteEndpoint'))
    }
    await refreshMonitors()
    await refreshAuditLogs()
  } catch (error) {
    batchMessage.value = error.message
  }
}

const toggleEnabled = async (monitor) => {
  try {
    await openEdit(monitor)
    if (drawer.value.entityType === 'endpoint') {
      endpointForm.value.enabled = !endpointForm.value.enabled
    } else if (drawer.value.entityType === 'suite') {
      suiteForm.value.enabled = !suiteForm.value.enabled
    } else {
      externalForm.value.enabled = !externalForm.value.enabled
    }
    await saveDrawer()
    closeDrawer()
  } catch (error) {
    batchMessage.value = error.message
  }
}

const executeBatch = async (dryRun) => {
  batchLoading.value = true
  batchMessage.value = ''
  try {
    const payload = {}
    if (batch.value.action === 'set-group') payload.group = batch.value.group
    if (batch.value.action === 'set-interval') payload.interval = batch.value.interval
    if (batch.value.action === 'set-alert-types') {
      payload.alertTypes = batch.value.alertTypes.split(',').map((item) => item.trim()).filter((item) => item)
    }

    const response = await fetch('/api/v1/admin/monitors/batch', {
      method: 'POST',
      credentials: 'include',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        entityType: 'all',
        keys: selectedKeys.value,
        action: batch.value.action,
        payload,
        dryRun,
      }),
    })
    const data = await response.json()
    if (!response.ok) {
      throw new Error(data.error || t('admin.failedSaveEndpoint'))
    }
    batchMessage.value = `${dryRun ? t('adminV2.dryRunDone') : t('adminV2.batchDone')} (${data.success}/${data.total})`
    await refreshMonitors()
    await refreshAuditLogs()
    if (!dryRun) clearSelection()
  } catch (error) {
    batchMessage.value = error.message
  } finally {
    batchLoading.value = false
  }
}

const openImportExport = () => {
  importExportOpen.value = true
  importExportMessage.value = ''
  importPreview.value = ''
}

const loadExportData = async () => {
  importExportLoading.value = true
  importExportMessage.value = ''
  try {
    const response = await fetch(`/api/v1/admin/export?entityType=${importExport.value.entityType}`, {
      credentials: 'include',
    })
    const data = await response.json()
    if (!response.ok) {
      throw new Error(data.error || 'Export failed')
    }
    importExport.value.exportJson = JSON.stringify(data, null, 2)
    overlayPath.value = data.overlayPath || overlayPath.value
  } catch (error) {
    importExportMessage.value = error.message
  } finally {
    importExportLoading.value = false
  }
}

const copyExportToImport = () => {
  importExport.value.importJson = importExport.value.exportJson
}

const runImport = async (dryRun) => {
  importExportLoading.value = true
  importExportMessage.value = ''
  importPreview.value = ''
  try {
    const parsed = JSON.parse(importExport.value.importJson || '{}')
    const requestBody = {
      entityType: importExport.value.entityType,
      mode: importExport.value.mode,
      dryRun,
      data: {
        alerting: parsed.alerting || null,
        endpoints: Array.isArray(parsed.endpoints) ? parsed.endpoints : [],
        externalEndpoints: Array.isArray(parsed.externalEndpoints) ? parsed.externalEndpoints : [],
        suites: Array.isArray(parsed.suites) ? parsed.suites : [],
      },
    }
    const response = await fetch('/api/v1/admin/import', {
      method: 'POST',
      credentials: 'include',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(requestBody),
    })
    const data = await response.json()
    if (!response.ok) {
      throw new Error(data.error || 'Import failed')
    }
    importPreview.value = JSON.stringify(data, null, 2)
    importExportMessage.value = dryRun ? t('adminV2.dryRunDone') : t('adminV2.importDone')
    if (!dryRun) {
      await refreshMonitors()
      await refreshAuditLogs()
      await loadManagedConfig()
    }
  } catch (error) {
    importExportMessage.value = error.message
  } finally {
    importExportLoading.value = false
  }
}

const loadManagedConfig = async () => {
  loadingManaged.value = true
  managedMessage.value = ''
  try {
    const response = await fetch('/api/v1/admin/managed-config', { credentials: 'include' })
    const data = await response.json()
    if (!response.ok) {
      throw new Error(data.error || t('admin.failedLoadManagedConfig'))
    }
    overlayPath.value = data.overlayPath || overlayPath.value
    jsonText.value = JSON.stringify({
      alerting: data.alerting || null,
      endpoints: data.endpoints || [],
      externalEndpoints: data.externalEndpoints || [],
      suites: data.suites || [],
    }, null, 2)
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
    const payload = JSON.parse(jsonText.value || '{}')
    const response = await fetch('/api/v1/admin/managed-config', {
      method: 'PUT',
      credentials: 'include',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload),
    })
    const data = await response.json()
    if (!response.ok) {
      throw new Error(data.error || t('admin.failedSaveManagedConfig'))
    }
    overlayPath.value = data.overlayPath || overlayPath.value
    managedMessage.value = data.message || t('admin.overlaySaved')
    await refreshMonitors()
    await refreshAuditLogs()
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
    await refreshAuditLogs()
  } catch (error) {
    managedMessage.value = error.message
  } finally {
    applyingManaged.value = false
  }
}

const resetOverlay = async () => {
  if (!window.confirm(t('admin.confirmResetOverlay'))) return
  savingManaged.value = true
  managedMessage.value = ''
  try {
    const response = await fetch('/api/v1/admin/managed-config', {
      method: 'DELETE',
      credentials: 'include',
    })
    if (!response.ok && response.status !== 204) {
      const data = await response.json().catch(() => ({}))
      throw new Error(data.error || t('admin.failedResetOverlay'))
    }
    managedMessage.value = t('admin.overlayDeletedHint')
    await loadManagedConfig()
    await refreshMonitors()
    await refreshAuditLogs()
  } catch (error) {
    managedMessage.value = error.message
  } finally {
    savingManaged.value = false
  }
}

const formatDateTime = (value) => {
  if (!value) return '-'
  const parsed = new Date(value)
  if (Number.isNaN(parsed.getTime())) return '-'
  return parsed.toLocaleString()
}

const statusClass = (status) => {
  if (status === 'healthy') return 'bg-green-100 text-green-800'
  if (status === 'unhealthy') return 'bg-red-100 text-red-800'
  if (status === 'disabled') return 'bg-slate-100 text-slate-700'
  return 'bg-amber-100 text-amber-800'
}

onMounted(async () => {
  await Promise.all([refreshMonitors(), refreshAuditLogs(), loadManagedConfig()])
})
</script>
