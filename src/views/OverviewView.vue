<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import api from '@/services/api'
import Chart from '@/components/common/Chart.vue'

interface Metrics {
  total_requests: number
  success_requests: number
  error_requests: number
  average_latency_ms: number
  volume_series: { timestamp: number; value: number }[]
  latency_series: { timestamp: number; value: number }[]
  status_codes_series: { name: string; value: number }[]
}

interface LiveStats {
  total_requests: number
  success_requests: number
  error_requests: number
  average_latency_ms: number
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
      }
    } catch {
      // ignore parse errors
    }
  }
  eventSource.onopen = () => {
    if (loading.value) loading.value = false
  }
  eventSource.onerror = () => {
    eventSource?.close()
    eventSource = null
    setTimeout(connectSSE, 3000)
  }
}

onMounted(() => {
  fetchMetrics()
  connectSSE()
  // Periodically refresh chart data
  const chartInterval = setInterval(fetchMetrics, 30000)
  return () => clearInterval(chartInterval)
})

onUnmounted(() => {
  eventSource?.close()
})

const getVolumeChartOptions = (seriesData: { timestamp: number; value: number }[]): Highcharts.Options => ({
  chart: { type: 'area', height: 280 },
  title: { text: 'Request Volume', align: 'left', margin: 0, y: 16 },
  xAxis: { type: 'datetime', lineWidth: 0, tickWidth: 0, labels: { y: 16 } },
  yAxis: {
    title: { text: undefined },
    gridLineDashStyle: 'Dash' as Highcharts.DashStyleValue,
    labels: { y: 8 }
  },
  plotOptions: {
    area: {
      marker: { radius: 2, symbol: 'circle', fillColor: '#6798ff', lineWidth: 0 },
      lineWidth: 1.5,
      states: { hover: { lineWidth: 2 } },
      threshold: null
    }
  },
  series: [{
    name: 'Requests',
    type: 'area',
    data: seriesData.map(d => [d.timestamp, d.value]),
    color: '#6798ff',
    fillColor: {
      linearGradient: { x1: 0, y1: 0, x2: 0, y2: 1 },
      stops: [
        [0, 'rgba(103, 152, 255, 0.35)'],
        [1, 'rgba(103, 152, 255, 0.02)']
      ]
    }
  }]
})

const getLatencyChartOptions = (seriesData: { timestamp: number; value: number }[]): Highcharts.Options => ({
  chart: { type: 'spline', height: 280 },
  title: { text: 'Response Latency', align: 'left', margin: 0, y: 16 },
  xAxis: { type: 'datetime', lineWidth: 0, tickWidth: 0, labels: { y: 16 } },
  yAxis: {
    title: { text: 'ms', align: 'high', rotation: 0, offset: 0, y: -12, x: 8 },
    gridLineDashStyle: 'Dash' as Highcharts.DashStyleValue,
    labels: { y: 8 }
  },
  plotOptions: {
    spline: {
      marker: { enabled: false },
      lineWidth: 1.5,
      states: { hover: { lineWidth: 2.5 } }
    }
  },
  series: [{
    name: 'Avg Latency',
    type: 'spline',
    data: seriesData.map(d => [d.timestamp, d.value]),
    color: '#a7a7a7'
  }]
})

const getStatusCodesChartOptions = (seriesData: { name: string; value: number }[]): Highcharts.Options => ({
  chart: { type: 'pie', height: 280 },
  title: { text: 'Status Code Distribution', align: 'left', margin: 0, y: 16 },
  tooltip: { pointFormat: '<b>{point.percentage:.1f}%</b> ({point.y})' },
  plotOptions: {
    pie: {
      allowPointSelect: true,
      cursor: 'pointer',
      borderWidth: 1,
      borderColor: '#1e1e1e',
      innerSize: '55%',
      dataLabels: {
        enabled: true,
        distance: 8,
        format: '{point.name}',
        style: { color: '#a7a7a7', fontSize: '10px', textOutline: 'none', fontWeight: '500' }
      }
    }
  },
  series: [{
    name: 'Codes',
    type: 'pie',
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
</script>

<template>
  <div class="space-y-8">
    <div>
      <span class="text-caption font-jetbrains-mono tracking-caption text-blue-cornflower uppercase font-medium">COMMAND CENTER</span>
      <h2 class="text-heading-sm font-semibold text-snow mt-1 tracking-tight">Overview</h2>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="flex items-center justify-center py-20 text-ash text-sm">
      <div class="w-4 h-4 border-2 border-blue-cornflower/30 border-t-blue-cornflower rounded-full animate-spin mr-3" />
      Loading analytics metrics&hellip;
    </div>

    <template v-else-if="metrics">
      <!-- Stat Cards Row -->
      <div class="grid grid-cols-2 lg:grid-cols-4 gap-4">
        <div class="bg-card-carbon border border-steel-border rounded-lg p-5 group">
          <svg class="w-4 h-4 text-blue-cornflower mb-3" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M3 13.125C3 12.504 3.504 12 4.125 12h2.25c.621 0 1.125.504 1.125 1.125v6.75C7.5 20.496 6.996 21 6.375 21h-2.25A1.125 1.125 0 013 19.875v-6.75zM9.75 8.625c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125v11.25c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 01-1.125-1.125V8.625zM16.5 4.125c0-.621.504-1.125 1.125-1.125h2.25C20.496 3 21 3.504 21 4.125v15.75c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 01-1.125-1.125V4.125z" />
          </svg>
          <div class="text-[32px] font-semibold text-snow tracking-tight font-inter">{{ metrics.total_requests.toLocaleString() }}</div>
          <span class="text-caption font-jetbrains-mono tracking-caption text-ash uppercase mt-1 block">Total Requests</span>
        </div>

        <div class="bg-card-carbon border border-steel-border rounded-lg p-5 group">
          <svg class="w-4 h-4 text-green-400 mb-3" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M9 12.75L11.25 15 15 9.75M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <div class="text-[32px] font-semibold text-green-400 tracking-tight font-inter">
            {{ metrics.total_requests > 0 ? ((metrics.success_requests / metrics.total_requests) * 100).toFixed(1) : '100' }}%
          </div>
          <span class="text-caption font-jetbrains-mono tracking-caption text-ash uppercase mt-1 block">Success Rate</span>
        </div>

        <div class="bg-card-carbon border border-steel-border rounded-lg p-5 group">
          <svg class="w-4 h-4 text-red-400 mb-3" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v3.75m-9.303 3.376c-.866 1.5.217 3.374 1.948 3.374h14.71c1.73 0 2.813-1.874 1.948-3.374L13.949 3.378c-.866-1.5-3.032-1.5-3.898 0L2.697 16.126zM12 15.75h.007v.008H12v-.008z" />
          </svg>
          <div class="text-[32px] font-semibold text-red-400 tracking-tight font-inter">{{ metrics.error_requests.toLocaleString() }}</div>
          <span class="text-caption font-jetbrains-mono tracking-caption text-ash uppercase mt-1 block">Errors / Blocked</span>
        </div>

        <div class="bg-card-carbon border border-steel-border rounded-lg p-5 group">
          <svg class="w-4 h-4 text-ash mb-3" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 6v6h4.5m4.5 0a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <div class="text-[32px] font-semibold text-snow tracking-tight font-inter">{{ metrics.average_latency_ms.toFixed(0) }}<span class="text-base text-ash ml-1">ms</span></div>
          <span class="text-caption font-jetbrains-mono tracking-caption text-ash uppercase mt-1 block">Avg Latency</span>
        </div>
      </div>

      <!-- Charts Grid -->
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
        <div class="bg-card-carbon border border-steel-border rounded-lg p-5">
          <Chart :options="getVolumeChartOptions(metrics.volume_series)" />
        </div>
        <div class="bg-card-carbon border border-steel-border rounded-lg p-5">
          <Chart :options="getLatencyChartOptions(metrics.latency_series)" />
        </div>
        <div class="bg-card-carbon border border-steel-border rounded-lg p-5 lg:col-span-2">
          <Chart :options="getStatusCodesChartOptions(metrics.status_codes_series)" />
        </div>
      </div>
    </template>
  </div>
</template>
