<template>
  <div class="dashboard-container bg-background">
    <div class="container mx-auto px-4 py-8 max-w-7xl space-y-4">
      <div class="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
        <div>
          <h1 class="text-3xl font-bold tracking-tight">{{ t('certificates.title') }}</h1>
          <p class="text-muted-foreground mt-1">{{ t('certificates.subtitle') }}</p>
        </div>
        <Button variant="ghost" size="icon" @click="fetchData" :title="t('common.refreshData')" :disabled="loading">
          <RefreshCw :class="['h-5 w-5', loading && 'animate-spin']" />
        </Button>
      </div>

      <div class="grid grid-cols-2 lg:grid-cols-4 gap-3">
        <Card>
          <CardContent class="p-4">
            <p class="text-xs text-muted-foreground">{{ t('certificates.total') }}</p>
            <p class="text-2xl font-semibold">{{ certificateStats.total }}</p>
          </CardContent>
        </Card>
        <Card>
          <CardContent class="p-4">
            <p class="text-xs text-muted-foreground">{{ t('certificates.expired') }}</p>
            <p class="text-2xl font-semibold text-red-600">{{ certificateStats.expired }}</p>
          </CardContent>
        </Card>
        <Card>
          <CardContent class="p-4">
            <p class="text-xs text-muted-foreground">{{ t('certificates.expiringSoon') }}</p>
            <p class="text-2xl font-semibold text-amber-600">{{ certificateStats.expiringSoon }}</p>
          </CardContent>
        </Card>
        <Card>
          <CardContent class="p-4">
            <p class="text-xs text-muted-foreground">{{ t('certificates.healthy') }}</p>
            <p class="text-2xl font-semibold text-green-600">{{ certificateStats.healthy }}</p>
          </CardContent>
        </Card>
      </div>

      <Card>
        <CardContent class="p-4">
          <Input
            v-model="searchQuery"
            type="text"
            :placeholder="t('certificates.searchPlaceholder')"
          />
        </CardContent>
      </Card>

      <Card>
        <CardContent class="p-0">
          <div v-if="loading" class="p-8 flex items-center justify-center">
            <Loading size="lg" />
          </div>
          <div v-else-if="sortedCertificateEndpoints.length === 0" class="p-8 text-center text-muted-foreground">
            {{ t('certificates.noCertificates') }}
          </div>
          <div v-else-if="filteredCertificateEndpoints.length === 0" class="p-8 text-center text-muted-foreground">
            {{ t('certificates.noSearchResult') }}
          </div>
          <div v-else class="overflow-x-auto">
            <table class="w-full text-sm">
              <thead class="bg-muted/50">
                <tr>
                  <th class="p-3 text-left">{{ t('common.name') }}</th>
                  <th class="p-3 text-left">{{ t('common.group') }}</th>
                  <th class="p-3 text-left">{{ t('certificates.hostname') }}</th>
                  <th class="p-3 text-left">{{ t('certificates.expiresAt') }}</th>
                  <th class="p-3 text-left">{{ t('certificates.remaining') }}</th>
                  <th class="p-3 text-left">{{ t('common.status') }}</th>
                </tr>
              </thead>
              <tbody>
                <tr
                  v-for="item in filteredCertificateEndpoints"
                  :key="item.key"
                  class="border-t hover:bg-accent/30 cursor-pointer"
                  @click="goToEndpoint(item.key)"
                >
                  <td class="p-3">
                    <p class="font-medium text-foreground">{{ item.name }}</p>
                    <p class="text-xs text-muted-foreground">{{ item.key }}</p>
                  </td>
                  <td class="p-3 text-muted-foreground">{{ item.group }}</td>
                  <td class="p-3 text-muted-foreground">{{ item.hostname }}</td>
                  <td class="p-3 font-mono">{{ formatExpirationAt(item.expirationTimestamp) }}</td>
                  <td class="p-3" :class="remainingClass(item.remainingSeconds)">{{ formatRemaining(item.remainingSeconds) }}</td>
                  <td class="p-3">
                    <StatusBadge :status="item.monitorStatus" />
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </CardContent>
      </Card>
    </div>

    <Settings @refreshData="fetchData" />
  </div>
</template>

<script setup>
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { RefreshCw } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Card, CardContent } from '@/components/ui/card'
import Loading from '@/components/Loading.vue'
import Settings from '@/components/Settings.vue'
import StatusBadge from '@/components/StatusBadge.vue'
import { useI18n, getCurrentLocale } from '@/i18n'
import { formatDurationFromSeconds } from '@/utils/time'

const CERT_EXPIRING_SOON_SECONDS = 7 * 24 * 3600

const router = useRouter()
const { t } = useI18n()

const loading = ref(false)
const endpointStatuses = ref([])
const searchQuery = ref('')
const nowTimestamp = ref(Date.now())

let nowTicker = null

const sortedCertificateEndpoints = computed(() => {
  const now = nowTimestamp.value

  return endpointStatuses.value
    .map((endpoint) => {
      if (!endpoint?.results || endpoint.results.length === 0) {
        return null
      }

      const latestResult = endpoint.results[endpoint.results.length - 1]
      const rawCertificateExpirationSeconds = latestResult?.certificateExpirationSeconds
      const certificateExpirationSeconds = Number(rawCertificateExpirationSeconds)
      if (!Number.isFinite(certificateExpirationSeconds)) {
        return null
      }

      const parsedTimestamp = latestResult?.timestamp ? new Date(latestResult.timestamp).getTime() : NaN
      const resultTimestamp = Number.isFinite(parsedTimestamp) ? parsedTimestamp : now
      const expirationTimestamp = resultTimestamp + Math.trunc(certificateExpirationSeconds) * 1000
      const remainingSeconds = Math.trunc((expirationTimestamp - now) / 1000)

      return {
        key: endpoint.key,
        name: endpoint.name,
        group: endpoint.group || t('home.noGroup'),
        hostname: latestResult?.hostname || t('common.noData'),
        monitorStatus: latestResult?.success === true ? 'healthy' : latestResult?.success === false ? 'unhealthy' : 'unknown',
        expirationTimestamp,
        remainingSeconds,
      }
    })
    .filter((item) => item !== null)
    .sort((a, b) => {
      if (a.expirationTimestamp !== b.expirationTimestamp) {
        return a.expirationTimestamp - b.expirationTimestamp
      }
      return a.name.localeCompare(b.name)
    })
})

const filteredCertificateEndpoints = computed(() => {
  const query = searchQuery.value.trim().toLowerCase()
  if (!query) {
    return sortedCertificateEndpoints.value
  }
  return sortedCertificateEndpoints.value.filter((item) => {
    return (
      item.name.toLowerCase().includes(query) ||
      item.group.toLowerCase().includes(query) ||
      item.hostname.toLowerCase().includes(query) ||
      item.key.toLowerCase().includes(query)
    )
  })
})

const certificateStats = computed(() => {
  const total = sortedCertificateEndpoints.value.length
  const expired = sortedCertificateEndpoints.value.filter((item) => item.remainingSeconds < 0).length
  const expiringSoon = sortedCertificateEndpoints.value.filter((item) => item.remainingSeconds >= 0 && item.remainingSeconds <= CERT_EXPIRING_SOON_SECONDS).length
  const healthy = Math.max(0, total - expired - expiringSoon)
  return {
    total,
    expired,
    expiringSoon,
    healthy,
  }
})

const fetchData = async () => {
  loading.value = true
  try {
    // We only need the latest check result for certificate ranking.
    const response = await fetch('/api/v1/endpoints/statuses?page=1&pageSize=1', {
      credentials: 'include'
    })
    if (response.status === 200) {
      const data = await response.json()
      endpointStatuses.value = Array.isArray(data) ? data : []
    } else {
      console.error('[CertificateMonitor][fetchData] Failed to fetch endpoint statuses:', await response.text())
      endpointStatuses.value = []
    }
  } catch (error) {
    console.error('[CertificateMonitor][fetchData] Error:', error)
    endpointStatuses.value = []
  } finally {
    loading.value = false
  }
}

const formatExpirationAt = (expirationTimestamp) => {
  return new Date(expirationTimestamp).toLocaleString(getCurrentLocale())
}

const formatRemaining = (remainingSeconds) => {
  const duration = formatDurationFromSeconds(remainingSeconds)
  if (remainingSeconds < 0) {
    return t('endpointDetails.expiredFor', { duration })
  }
  return t('endpointDetails.expiresIn', { duration })
}

const remainingClass = (remainingSeconds) => {
  if (remainingSeconds < 0) {
    return 'text-red-600'
  }
  if (remainingSeconds <= CERT_EXPIRING_SOON_SECONDS) {
    return 'text-amber-600'
  }
  return 'text-green-600'
}

const goToEndpoint = (key) => {
  router.push(`/endpoints/${key}`)
}

onMounted(() => {
  fetchData()
  nowTicker = window.setInterval(() => {
    nowTimestamp.value = Date.now()
  }, 60000)
})

onUnmounted(() => {
  if (nowTicker) {
    window.clearInterval(nowTicker)
  }
})
</script>
