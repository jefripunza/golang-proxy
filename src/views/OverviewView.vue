<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import api from '@/services/api'
import Chart from '@/components/common/Chart.vue'

const route = useRoute()

interface Metrics {
  total_requests: number
  success_requests: number
  error_requests: number
  average_latency_ms: number
  max_latency_ms: number
  min_latency_ms: number
  volume_series: { timestamp: number; value: number }[]
  latency_series: { timestamp: number; value: number }[]
  status_codes_series: { name: string; value: number }[]
}

interface LiveStats {
  total_requests: number
  success_requests: number
  error_requests: number
  average_latency_ms: number
  max_latency_ms: number
  min_latency_ms: number
  volume_series: { timestamp: number; value: number }[]
  latency_series: { timestamp: number; value: number }[]
}

const metrics = ref<Metrics | null>(null)
const loading = ref(true)
let eventSource: EventSource | null = null

const fetchMetrics = async () => {
  try {
    const res = await api.get('/api/metrics')
    metrics.value = res.data
  } catch (err) {
    console.error('Failed to fetch metrics:', err)
  } finally {
    loading.value = false
  }
}

const clearMetrics = async () => {
  if (!confirm('Clear all proxy metrics data? This cannot be undone.')) return
  try {
    await api.delete('/api/metrics')
    fetchMetrics()
  } catch (err) {
    console.error('Failed to clear metrics:', err)
  }
}

const connectSSE = () => {
  if (eventSource) return
  eventSource = new EventSource('/api/metrics/stream', { withCredentials: true })
  eventSource.onmessage = (event) => {
    try {
      const data: LiveStats = JSON.parse(event.data)
      if (metrics.value) {
        metrics.value.total_requests = data.total_requests
        metrics.value.success_requests = data.success_requests
        metrics.value.error_requests = data.error_requests
        metrics.value.average_latency_ms = data.average_latency_ms
        metrics.value.max_latency_ms = data.max_latency_ms
        metrics.value.min_latency_ms = data.min_latency_ms
        if (data.volume_series) metrics.value.volume_series = data.volume_series
        if (data.latency_series) metrics.value.latency_series = data.latency_series
      }
    } catch { /* ignore */ }
  }
  eventSource.onopen = () => { if (loading.value) loading.value = false }
  eventSource.onerror = () => { eventSource?.close(); eventSource = null; setTimeout(connectSSE, 1000) }
}

const getCombinedChartOptions = (
  volumeData: { timestamp: number; value: number }[],
  latencyData: { timestamp: number; value: number }[]
): Highcharts.Options => ({
  chart: { type: 'line', height: 320, animation: false },
  title: { text: 'Request Volume & Response Latency', align: 'left', margin: 24, y: 24 },
  xAxis: { type: 'datetime', lineWidth: 0, tickWidth: 0, labels: { y: 16 } },
  yAxis: [
    { title: { text: 'Requests', style: { color: '#6798ff' } }, gridLineDashStyle: 'Dash' as Highcharts.DashStyleValue, labels: { y: 8, style: { color: '#6798ff' } } },
    { title: { text: 'ms', style: { color: '#a7a7a7' } }, opposite: true, gridLineWidth: 0, labels: { y: 8, style: { color: '#a7a7a7' } } }
  ],
  tooltip: { shared: true },
  plotOptions: { line: { marker: { enabled: false }, lineWidth: 2, animation: false } },
  series: [{
    name: 'Requests', type: 'line', yAxis: 0,
    data: volumeData.map(d => [d.timestamp, d.value]),
    color: '#6798ff'
  }, {
    name: 'Latency', type: 'spline', yAxis: 1,
    data: latencyData.map(d => [d.timestamp, d.value]),
    color: '#a7a7a7', dashStyle: 'Dash' as Highcharts.DashStyleValue
  }]
})

const getStatusCodesChartOptions = (seriesData: { name: string; value: number }[]): Highcharts.Options => ({
  chart: { type: 'pie', height: 200, animation: false },
  title: { text: 'Status Codes', align: 'left', margin: 20, y: 20 },
  tooltip: { pointFormat: '<b>{point.percentage:.1f}%</b> ({point.y})' },
  plotOptions: {
    pie: {
      allowPointSelect: true, cursor: 'pointer', borderWidth: 1, borderColor: '#1e1e1e',
      innerSize: '50%',
      dataLabels: { enabled: true, distance: 6, format: '{point.name}', style: { color: '#a7a7a7', fontSize: '9px', textOutline: 'none' } }
    }
  },
  series: [{
    name: 'Codes', type: 'pie',
    data: seriesData.map(d => {
      let color = '#454545'
      if (d.name.startsWith('2')) color = '#22c55e'
      else if (d.name.startsWith('3')) color = '#6798ff'
      else if (d.name.startsWith('4')) color = '#eab308'
      else if (d.name.startsWith('5')) color = '#ef4444'
      else color = '#ef4444'
      return { name: d.name, y: d.value, color }
    })
  }]
})

onMounted(() => { fetchMetrics(); connectSSE() })
onUnmounted(() => { eventSource?.close() })

// Re-fetch and show skeleton when navigating back to this page
watch(() => route.path, (path) => {
  if (path === '/dashboard') {
    loading.value = true
    fetchMetrics()
  }
})
</script>

<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <span class="text-caption font-jetbrains-mono tracking-caption text-blue-cornflower uppercase font-medium">COMMAND CENTER</span>
        <h2 class="text-heading-sm font-semibold text-snow mt-1 tracking-tight">Overview</h2>
      </div>
      <button type="button" @click="clearMetrics" class="px-3 py-1.5 border border-red-900/40 bg-red-950/10 rounded-lg text-red-400 text-xs font-medium hover:text-red-300 hover:bg-red-950/20 transition-colors cursor-pointer">Clear Metrics</button>
    </div>

    <!-- Skeleton Loading -->
    <div v-if="loading" class="space-y-6">
      <div class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-5 gap-3">
        <div v-for="i in 5" :key="i" class="bg-card-carbon border border-steel-border rounded-lg p-4 animate-pulse">
          <div class="h-2.5 bg-deep-coal rounded w-16 mb-3" />
          <div class="h-6 bg-deep-coal rounded w-20" />
        </div>
      </div>
      <div class="grid grid-cols-1 lg:grid-cols-12 gap-4">
        <div class="lg:col-span-9 bg-card-carbon border border-steel-border rounded-lg p-4 animate-pulse">
          <div class="h-4 bg-deep-coal rounded w-48 mb-4" />
          <div class="h-[280px] bg-deep-coal rounded" />
        </div>
        <div class="lg:col-span-3 bg-card-carbon border border-steel-border rounded-lg p-4 animate-pulse">
          <div class="h-4 bg-deep-coal rounded w-24 mb-4" />
          <div class="h-[180px] bg-deep-coal rounded" />
        </div>
      </div>
    </div>

    <template v-else-if="metrics">
      <!-- Stat Cards: Success Rate, Errors, Avg/Max/Min Latency -->
      <div class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-5 gap-3">
        <div class="bg-card-carbon border border-steel-border rounded-lg p-4">
          <span class="text-[9px] font-jetbrains-mono tracking-wider text-ash uppercase">Success Rate</span>
          <div class="text-xl font-semibold text-green-400 mt-1 font-inter">
            {{ metrics.total_requests > 0 ? ((metrics.success_requests / metrics.total_requests) * 100).toFixed(2) : '100.00' }}%
          </div>
        </div>
        <div class="bg-card-carbon border border-steel-border rounded-lg p-4">
          <span class="text-[9px] font-jetbrains-mono tracking-wider text-ash uppercase">Errors</span>
          <div class="text-xl font-semibold text-red-400 mt-1 font-inter">{{ metrics.error_requests }}</div>
        </div>
        <div class="bg-card-carbon border border-steel-border rounded-lg p-4">
          <span class="text-[9px] font-jetbrains-mono tracking-wider text-ash uppercase">Avg Latency</span>
          <div class="text-xl font-semibold text-snow mt-1 font-inter">{{ metrics.average_latency_ms.toFixed(2) }}ms</div>
        </div>
        <div class="bg-card-carbon border border-steel-border rounded-lg p-4">
          <span class="text-[9px] font-jetbrains-mono tracking-wider text-ash uppercase">Max Latency</span>
          <div class="text-xl font-semibold text-snow mt-1 font-inter">{{ metrics.max_latency_ms.toFixed(2) }}ms</div>
        </div>
        <div class="bg-card-carbon border border-steel-border rounded-lg p-4">
          <span class="text-[9px] font-jetbrains-mono tracking-wider text-ash uppercase">Min Latency</span>
          <div class="text-xl font-semibold text-snow mt-1 font-inter">{{ metrics.min_latency_ms.toFixed(2) }}ms</div>
        </div>
      </div>

      <!-- Charts: 9-col combined + 3-col donut -->
      <div class="grid grid-cols-1 lg:grid-cols-12 gap-4">
        <div class="lg:col-span-9 bg-card-carbon border border-steel-border rounded-lg p-4">
          <Chart :options="getCombinedChartOptions(metrics.volume_series, metrics.latency_series)" />
        </div>
        <div class="lg:col-span-3 bg-card-carbon border border-steel-border rounded-lg p-4">
          <Chart :options="getStatusCodesChartOptions(metrics.status_codes_series)" />
        </div>
      </div>
    </template>
  </div>
</template>
