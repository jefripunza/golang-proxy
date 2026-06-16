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
}

const logs = ref<ProxyLog[]>([])
const loading = ref(true)
const limit = ref(100)

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

const getStatusClass = (code: number) => {
  if (code >= 200 && code < 300) return 'bg-green-950/40 border border-green-800/40 text-green-300'
  if (code >= 300 && code < 400) return 'bg-blue-950/40 border border-blue-800/40 text-blue-300'
  if (code >= 400 && code < 500) return 'bg-yellow-950/40 border border-yellow-800/40 text-yellow-300'
  return 'bg-red-950/40 border border-red-800/40 text-red-300'
}

onMounted(() => {
  fetchLogs()
  // Poll every 5 seconds for real-time monitoring
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
        <p class="text-body-sm text-ash mt-1 max-w-lg">Real-time request log stream. Monitor incoming traffic, status codes, and response latencies.</p>
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
          @click="fetchLogs"
          class="px-4 py-2 border border-graphite rounded-lg text-snow text-[13px] font-medium hover:bg-card-carbon/50 transition-colors cursor-pointer"
        >
          Refresh
        </button>
      </div>
    </div>

    <!-- Logs Table -->
    <div v-if="loading" class="text-center text-ash py-8">
      Loading traffic activity logs...
    </div>

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
          <tr v-for="log in logs" :key="log.id" class="hover:bg-deep-coal/50 transition-colors">
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
  </div>
</template>
