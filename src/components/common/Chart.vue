<script setup lang="ts">
import { onMounted, onUnmounted, ref, watch } from 'vue'
import Highcharts from 'highcharts'

const props = defineProps<{
  options: Highcharts.Options
}>()

const container = ref<HTMLElement | null>(null)
let chartInstance: Highcharts.Chart | null = null

const initChart = () => {
  if (container.value) {
    // Apply Dovetail style defaults
    const darkThemeOptions: Highcharts.Options = {
      chart: {
        backgroundColor: 'transparent',
        style: {
          fontFamily: 'Inter, sans-serif'
        }
      },
      title: {
        style: {
          color: '#ffffff',
          fontSize: '14px',
          fontWeight: '600'
        }
      },
      xAxis: {
        gridLineColor: 'rgba(255, 255, 255, 0.05)',
        lineColor: '#313131',
        tickColor: '#313131',
        labels: {
          style: {
            color: '#a7a7a7',
            fontSize: '11px'
          }
        }
      },
      yAxis: {
        gridLineColor: 'rgba(255, 255, 255, 0.05)',
        lineColor: '#313131',
        tickColor: '#313131',
        labels: {
          style: {
            color: '#a7a7a7',
            fontSize: '11px'
          }
        },
        title: {
          style: {
            color: '#a7a7a7',
            fontSize: '11px'
          }
        }
      },
      legend: {
        itemStyle: {
          color: '#ffffff'
        },
        itemHoverStyle: {
          color: '#6798ff'
        }
      },
      plotOptions: {
        series: {
          borderWidth: 0
        }
      },
      credits: {
        enabled: false
      }
    }

    const mergedOptions = Highcharts.merge(darkThemeOptions, props.options)
    chartInstance = Highcharts.chart(container.value, mergedOptions)
  }
}

watch(
  () => props.options,
  (newOptions) => {
    if (chartInstance) {
      chartInstance.update(newOptions, true, true)
    } else {
      initChart()
    }
  },
  { deep: true }
)

onMounted(() => {
  initChart()
})

onUnmounted(() => {
  if (chartInstance) {
    chartInstance.destroy()
  }
})
</script>

<template>
  <div ref="container" class="w-full h-full min-h-[300px]"></div>
</template>
