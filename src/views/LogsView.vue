<script setup lang="ts">
import { ref, onMounted } from 'vue'
import api from '@/services/api'

interface ProxyLog {
  id: number
  timestamp: string
  domain: string
  path: string
  method: string
  status_code: number
  response_time_ms: number
  source_ip: string
  error_message: string
  request_headers: string
  response_headers: string
  request_body: string
  response_body: string
}

const logs = ref<ProxyLog[]>([])
const loading = ref(true)
const limit = ref(100)
const selectedLog = ref<ProxyLog | null>(null)
const autoClear = ref(0) // 0 = never

const fetchLogs = async () => {
  try {
    const res = await api.get(`/api/logs?limit=${limit.value}`)
    logs.value = res.data
  } catch (err) {
    console.error('Failed to fetch logs:', err)
  } finally {
    loading.value = false
  }
}

const fetchAutoClear = async () => {
  try {
    const res = await api.get('/api/logs/autoclear')
    autoClear.value = res.data.interval || 0
  } catch { /* ignore */ }
}

const setAutoClear = async (val: number) => {
  autoClear.value = val
  try {
    await api.put('/api/logs/autoclear', { interval: val })
  } catch { /* ignore */ }
}

const getStatusClass = (code: number) => {
  if (code >= 200 && code < 300) return 'bg-green-950/40 border border-green-800/40 text-green-300'
  if (code >= 300 && code < 400) return 'bg-blue-950/40 border border-blue-800/40 text-blue-300'
  if (code >= 400 && code < 500) return 'bg-yellow-950/40 border border-yellow-800/40 text-yellow-300'
  return 'bg-red-950/40 border border-red-800/40 text-red-300'
}

const parseHeaders = (raw: string): Record<string, string> => {
  if (!raw) return {}
  try {
    return JSON.parse(raw)
  } catch {
    return {}
  }
}

const openDetail = (log: ProxyLog) => {
  selectedLog.value = log
}

const clearLogs = async () => {
  if (!confirm('Delete all proxy logs? This cannot be undone.')) return
  try {
    await api.delete('/api/logs')
    fetchLogs()
  } catch (err) {
    console.error('Failed to clear logs:', err)
  }
}

onMounted(() => {
  fetchLogs()
  fetchAutoClear()
  const interval = setInterval(fetchLogs, 5000)
  return () => clearInterval(interval)
})
</script>

<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <span class="text-caption font-jetbrains-mono tracking-caption text-blue-cornflower uppercase font-medium">TRAFFIC MONITOR</span>
        <h2 class="text-heading-sm font-semibold text-snow mt-1 tracking-tight">Logs</h2>
        <p class="text-body-sm text-ash mt-1 max-w-lg">Real-time request log stream. Click any row to inspect request and response details.</p>
      </div>
      <div class="flex items-center gap-3">
        <div class="flex items-center gap-2 select-none">
          <span class="text-[10px] text-ash font-jetbrains-mono uppercase tracking-wider">Limit</span>
          <select
            v-model="limit"
            @change="fetchLogs"
            class="bg-deep-coal border border-graphite rounded-lg px-2.5 py-1.5 text-snow text-[12px] font-medium focus:outline-none focus:border-blue-cornflower transition-colors cursor-pointer"
          >
            <option :value="50">50</option>
            <option :value="100">100</option>
            <option :value="200">200</option>
            <option :value="500">500</option>
          </select>
        </div>
        <button
          type="button"
          @click="fetchLogs"
          class="px-4 py-2 border border-graphite rounded-lg text-snow text-[13px] font-medium hover:bg-card-carbon/50 transition-colors cursor-pointer"
        >
          Refresh
        </button>
        <button
          type="button"
          @click="clearLogs"
          class="px-3 py-1.5 border border-red-900/40 bg-red-950/10 rounded-lg text-red-400 text-xs font-medium hover:text-red-300 hover:bg-red-950/20 transition-colors cursor-pointer"
        >
          Clear
        </button>
        <div class="flex items-center gap-1.5 select-none border-l border-steel-border pl-3 ml-1">
          <span class="text-[9px] text-ash font-jetbrains-mono uppercase tracking-wider">Auto Clear</span>
          <select
            :value="autoClear"
            @change="setAutoClear(Number(($event.target as HTMLSelectElement).value))"
            class="bg-deep-coal border border-graphite rounded-lg px-2 py-1.5 text-snow text-[11px] font-medium focus:outline-none focus:border-blue-cornflower transition-colors cursor-pointer"
          >
            <option :value="0">Never</option>
            <option :value="1">1 Hour</option>
            <option :value="6">6 Hours</option>
            <option :value="12">12 Hours</option>
            <option :value="24">24 Hours</option>
          </select>
        </div>
      </div>
    </div>

    <div v-if="loading" class="text-center text-ash py-8">Loading traffic activity logs...</div>

    <div v-else-if="logs.length === 0" class="border border-dashed border-steel-border rounded-lg p-12 text-center text-ash">
      No activity logs found. Ensure requests are hitting the proxy engine.
    </div>

    <div v-else class="overflow-x-auto bg-card-carbon border border-steel-border rounded-lg">
      <table class="w-full text-left border-collapse">
        <thead>
          <tr class="border-b border-steel-border bg-deep-coal text-caption font-jetbrains-mono tracking-caption text-ash uppercase">
            <th class="px-6 py-3 font-medium">Timestamp</th>
            <th class="px-6 py-3 font-medium">Source IP</th>
            <th class="px-6 py-3 font-medium">Method</th>
            <th class="px-6 py-3 font-medium">Host / Path</th>
            <th class="px-6 py-3 font-medium">Status</th>
            <th class="px-6 py-3 font-medium">Latency</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-steel-border text-sm text-snow font-jetbrains-mono">
          <tr
            v-for="log in logs"
            :key="log.id"
            class="hover:bg-deep-coal/50 transition-colors cursor-pointer"
            @click="openDetail(log)"
          >
            <td class="px-6 py-4 text-ash text-xs">
              {{ new Date(log.timestamp).toLocaleString() }}
            </td>
            <td class="px-6 py-4 text-xs text-ash">{{ log.source_ip }}</td>
            <td class="px-6 py-4">
              <span class="font-bold text-xs uppercase" :class="log.method === 'GET' ? 'text-blue-300' : 'text-purple-300'">
                {{ log.method }}
              </span>
            </td>
            <td class="px-6 py-4 max-w-sm truncate text-xs">
              <span class="text-blue-cornflower">{{ log.domain }}</span><span class="text-ash">{{ log.path }}</span>
              <div v-if="log.error_message" class="text-[10px] text-red-400 mt-0.5 truncate">{{ log.error_message }}</div>
            </td>
            <td class="px-6 py-4">
              <span
                class="px-2 py-0.5 rounded-[4px] text-[10px] font-medium uppercase tracking-wider"
                :class="getStatusClass(log.status_code)"
              >
                {{ log.status_code || 'Blocked' }}
              </span>
            </td>
            <td class="px-6 py-4 text-xs text-ash">{{ log.response_time_ms }} ms</td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Detail Modal -->
    <div
      v-if="selectedLog"
      class="fixed inset-0 z-[60] flex items-center justify-center bg-black/70 backdrop-blur-sm p-4"
      @click.self="selectedLog = null"
    >
      <div class="w-full max-w-2xl max-h-[80vh] bg-card-carbon border border-steel-border rounded-lg overflow-hidden flex flex-col">
        <div class="px-5 py-3 border-b border-steel-border flex items-center justify-between shrink-0">
          <div class="flex items-center gap-3">
            <span class="text-[13px] font-semibold text-snow font-inter">Request Detail</span>
            <span
              class="px-2 py-0.5 rounded-[4px] text-[10px] font-jetbrains-mono font-medium uppercase tracking-wider"
              :class="getStatusClass(selectedLog.status_code)"
            >
              {{ selectedLog.status_code || 'Blocked' }}
            </span>
          </div>
          <button
            type="button"
            @click="selectedLog = null"
            class="text-ash hover:text-snow transition-colors cursor-pointer"
          >
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round">
              <line x1="18" y1="6" x2="6" y2="18" /><line x1="6" y1="6" x2="18" y2="18" />
            </svg>
          </button>
        </div>
        <div class="p-5 overflow-y-auto flex-1 space-y-4 text-sm">
          <!-- Summary -->
          <div class="grid grid-cols-2 gap-3">
            <div class="bg-deep-coal rounded-lg p-3">
              <span class="text-[10px] font-jetbrains-mono tracking-wider text-ash uppercase">Method</span>
              <div class="text-snow mt-0.5 font-semibold">{{ selectedLog.method }}</div>
            </div>
            <div class="bg-deep-coal rounded-lg p-3">
              <span class="text-[10px] font-jetbrains-mono tracking-wider text-ash uppercase">Status</span>
              <div class="text-snow mt-0.5 font-semibold">{{ selectedLog.status_code || '—' }}</div>
            </div>
            <div class="bg-deep-coal rounded-lg p-3">
              <span class="text-[10px] font-jetbrains-mono tracking-wider text-ash uppercase">Source IP</span>
              <div class="text-snow mt-0.5 font-medium font-jetbrains-mono">{{ selectedLog.source_ip }}</div>
            </div>
            <div class="bg-deep-coal rounded-lg p-3">
              <span class="text-[10px] font-jetbrains-mono tracking-wider text-ash uppercase">Latency</span>
              <div class="text-snow mt-0.5 font-medium">{{ selectedLog.response_time_ms }} ms</div>
            </div>
            <div class="col-span-2 bg-deep-coal rounded-lg p-3">
              <span class="text-[10px] font-jetbrains-mono tracking-wider text-ash uppercase">URL</span>
              <div class="text-blue-cornflower mt-0.5 font-medium font-jetbrains-mono text-xs break-all">{{ selectedLog.domain }}{{ selectedLog.path }}</div>
            </div>
          </div>

          <!-- Request Headers -->
          <div>
            <span class="text-caption font-jetbrains-mono tracking-caption text-ash uppercase">Request Headers</span>
            <div class="mt-1.5 bg-deep-coal rounded-lg p-3 font-jetbrains-mono text-xs max-h-40 overflow-y-auto">
              <div v-if="Object.keys(parseHeaders(selectedLog.request_headers)).length === 0" class="text-ash">No headers captured</div>
              <div v-else v-for="(val, key) in parseHeaders(selectedLog.request_headers)" :key="key" class="flex gap-2 py-0.5">
                <span class="text-blue-cornflower shrink-0">{{ key }}:</span>
                <span class="text-ash break-all">{{ val }}</span>
              </div>
            </div>
          </div>

          <!-- Response Headers -->
          <div>
            <span class="text-caption font-jetbrains-mono tracking-caption text-ash uppercase">Response Headers</span>
            <div class="mt-1.5 bg-deep-coal rounded-lg p-3 font-jetbrains-mono text-xs max-h-40 overflow-y-auto">
              <div v-if="Object.keys(parseHeaders(selectedLog.response_headers)).length === 0" class="text-ash">No headers captured</div>
              <div v-else v-for="(val, key) in parseHeaders(selectedLog.response_headers)" :key="key" class="flex gap-2 py-0.5">
                <span class="text-green-400 shrink-0">{{ key }}:</span>
                <span class="text-ash break-all">{{ val }}</span>
              </div>
            </div>
          </div>

          <!-- Request Body -->
          <div>
            <span class="text-caption font-jetbrains-mono tracking-caption text-ash uppercase">Request Body</span>
            <div class="mt-1.5 bg-deep-coal rounded-lg p-3 font-jetbrains-mono text-xs max-h-40 overflow-y-auto">
              <div v-if="!selectedLog.request_body" class="text-ash">No body</div>
              <div v-else-if="selectedLog.request_body.startsWith('[file')" class="text-yellow-400">{{ selectedLog.request_body }}</div>
              <pre v-else class="text-ash whitespace-pre-wrap break-all">{{ selectedLog.request_body }}</pre>
            </div>
          </div>

          <!-- Response Body -->
          <div>
            <span class="text-caption font-jetbrains-mono tracking-caption text-ash uppercase">Response Body</span>
            <div class="mt-1.5 bg-deep-coal rounded-lg p-3 font-jetbrains-mono text-xs max-h-40 overflow-y-auto">
              <div v-if="!selectedLog.response_body" class="text-ash">No body</div>
              <div v-else-if="selectedLog.response_body.startsWith('[file')" class="text-yellow-400">{{ selectedLog.response_body }}</div>
              <pre v-else class="text-ash whitespace-pre-wrap break-all">{{ selectedLog.response_body }}</pre>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
