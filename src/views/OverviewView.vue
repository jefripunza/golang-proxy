<script setup lang="ts">
import { ref, onMounted } from 'vue'
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

const metrics = ref<Metrics | null>(null)
const loading = ref(true)

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

onMounted(() => {
  fetchMetrics()
  // Poll every 10 seconds for real-time dashboard updates
  const interval = setInterval(fetchMetrics, 10000)
  return () => clearInterval(interval)
})

const getVolumeChartOptions = (seriesData: { timestamp: number; value: number }[]): Highcharts.Options => ({
  chart: { type: 'area' },
  title: { text: 'Request Volume (Last 24 Hours)', align: 'left' },
  xAxis: { type: 'datetime' },
  yAxis: { title: { text: 'Requests' } },
  series: [{
    name: 'Requests',
    type: 'area',
    data: seriesData.map(d => [d.timestamp, d.value]),
    color: '#6798ff',
    fillColor: {
      linearGradient: { x1: 0, y1: 0, x2: 0, y2: 1 },
      stops: [
        [0, 'rgba(103, 152, 255, 0.4)'],
        [1, 'rgba(103, 152, 255, 0.05)']
      ]
    }
  }]
})

const getLatencyChartOptions = (seriesData: { timestamp: number; value: number }[]): Highcharts.Options => ({
  chart: { type: 'line' },
  title: { text: 'Response Latency (Last 24 Hours)', align: 'left' },
  xAxis: { type: 'datetime' },
  yAxis: { title: { text: 'Latency (ms)' } },
  series: [{
    name: 'Average Latency',
    type: 'line',
    data: seriesData.map(d => [d.timestamp, d.value]),
    color: '#a7a7a7'
  }]
})

const getStatusCodesChartOptions = (seriesData: { name: string; value: number }[]): Highcharts.Options => ({
  chart: { type: 'pie' },
  title: { text: 'HTTP Status Code Distribution', align: 'left' },
  tooltip: { pointFormat: '<b>{point.percentage:.1f}%</b> ({point.y} requests)' },
  plotOptions: {
    pie: {
      allowPointSelect: true,
      cursor: 'pointer',
      borderWidth: 1,
      borderColor: '#1e1e1e',
      dataLabels: {
        enabled: true,
        format: '{point.name}: {point.y}',
        style: { color: '#ffffff', textOutline: 'none' }
      }
    }
  },
  series: [{
    name: 'Status Codes',
    type: 'pie',
    data: seriesData.map(d => {
      let color = '#a7a7a7';
      if (d.name.startsWith('2')) color = '#22c55e'; // Success (green)
      else if (d.name.startsWith('3')) color = '#6798ff'; // Redirect (blue-cornflower)
      else if (d.name.startsWith('4')) color = '#eab308'; // Client Error (yellow)
      else if (d.name.startsWith('5')) color = '#ef4444'; // Server Error (red)
      else if (d.name.toLowerCase().includes('fail') || d.name.toLowerCase().includes('block')) color = '#ef4444'; // Blocked/Failed (red)
      return { name: d.name, y: d.value, color }
    })
  }]
})
</script>

<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <span class="text-caption font-jetbrains-mono tracking-caption text-blue-cornflower uppercase">DASHBOARD</span>
        <h2 class="text-3xl font-semibold text-snow mt-1 tracking-tight">Overview</h2>
      </div>
      <button
        @click="fetchMetrics"
        class="px-4 py-2 bg-transparent border border-graphite rounded-lg text-snow text-sm font-medium hover:bg-card-carbon transition-colors cursor-pointer"
      >
        Refresh Data
      </button>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="h-64 flex items-center justify-center text-ash">
      Loading analytics metrics...
    </div>

    <template v-else-if="metrics">
      <!-- Stats Grid -->
      <div class="grid grid-cols-1 md:grid-cols-4 gap-6">
        <!-- Total Requests -->
        <div class="bg-card-carbon border border-steel-border rounded-lg p-6">
          <span class="text-caption font-jetbrains-mono tracking-caption text-ash uppercase">TOTAL REQUESTS</span>
          <div class="text-4xl font-semibold text-snow mt-2 font-inter">{{ metrics.total_requests }}</div>
        </div>

        <!-- Success Rate -->
        <div class="bg-card-carbon border border-steel-border rounded-lg p-6">
          <span class="text-caption font-jetbrains-mono tracking-caption text-ash uppercase">SUCCESS RATE</span>
          <div class="text-4xl font-semibold text-snow mt-2 font-inter">
            {{ metrics.total_requests > 0 ? ((metrics.success_requests / metrics.total_requests) * 100).toFixed(1) : '100' }}%
          </div>
        </div>

        <!-- Error Requests -->
        <div class="bg-card-carbon border border-steel-border rounded-lg p-6">
          <span class="text-caption font-jetbrains-mono tracking-caption text-ash uppercase">ERRORS / BLOCKED</span>
          <div class="text-4xl font-semibold text-red-400 mt-2 font-inter">{{ metrics.error_requests }}</div>
        </div>

        <!-- Avg Latency -->
        <div class="bg-card-carbon border border-steel-border rounded-lg p-6">
          <span class="text-caption font-jetbrains-mono tracking-caption text-ash uppercase">AVG LATENCY</span>
          <div class="text-4xl font-semibold text-snow mt-2 font-inter">{{ metrics.average_latency_ms.toFixed(0) }} ms</div>
        </div>
      </div>

      <!-- Charts Grid -->
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div class="bg-card-carbon border border-steel-border rounded-lg p-6">
          <Chart :options="getVolumeChartOptions(metrics.volume_series)" />
        </div>
        <div class="bg-card-carbon border border-steel-border rounded-lg p-6">
          <Chart :options="getLatencyChartOptions(metrics.latency_series)" />
        </div>
        <div class="bg-card-carbon border border-steel-border rounded-lg p-6 lg:col-span-2">
          <Chart :options="getStatusCodesChartOptions(metrics.status_codes_series)" />
        </div>
      </div>
    </template>
  </div>
</template>
