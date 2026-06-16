<script setup lang="ts">
import { ref, watch, onUnmounted, nextTick } from 'vue'
import { Terminal } from '@xterm/xterm'
import '@xterm/xterm/css/xterm.css'

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

const props = defineProps<{
  show: boolean
  routeId: number
  domain: string
}>()

const emit = defineEmits<{
  (e: 'close'): void
}>()

const terminalRef = ref<HTMLElement | null>(null)
let terminal: Terminal | null = null
let eventSource: EventSource | null = null

const statusColor = (code: number): string => {
  if (code >= 200 && code < 300) return '\x1b[32m'
  if (code >= 400 && code < 500) return '\x1b[33m'
  if (code >= 500 || code === 0) return '\x1b[31m'
  return '\x1b[36m'
}

const methodColor = (method: string): string => {
  switch (method) {
    case 'GET': return '\x1b[34m'
    case 'POST': return '\x1b[35m'
    case 'PUT': return '\x1b[33m'
    case 'DELETE': return '\x1b[31m'
    default: return '\x1b[37m'
  }
}

const writeLog = (log: ProxyLog) => {
  if (!terminal) return
  const time = new Date(log.timestamp).toLocaleTimeString('en-US', { hour12: false })
  const reset = '\x1b[0m'
  const dim = '\x1b[2m'
  terminal.writeln(
    `${dim}${time}${reset} ${dim}${log.source_ip.padEnd(15)}${reset} ` +
    `${methodColor(log.method)}${log.method.padEnd(7)}${reset} ` +
    `${statusColor(log.status_code)}${String(log.status_code || '---').padStart(3)}${reset} ` +
    `${dim}${String(log.response_time_ms).padStart(5)}ms${reset} ` +
    `${log.path}${log.error_message ? ` ${dim}(${log.error_message})${reset}` : ''}`
  )
}

const connectStream = () => {
  if (eventSource) return
  eventSource = new EventSource(`/api/routes/${props.routeId}/stream`)
  eventSource.onmessage = (event) => {
    try {
      const log: ProxyLog = JSON.parse(event.data)
      writeLog(log)
    } catch {
      // ignore
    }
  }
  eventSource.onerror = () => {
    eventSource?.close()
    eventSource = null
    setTimeout(connectStream, 3000)
  }
}

const disconnectStream = () => {
  eventSource?.close()
  eventSource = null
}

const initTerminal = () => {
  if (!terminalRef.value || terminal) return
  terminal = new Terminal({
    cursorBlink: false,
    disableStdin: true,
    fontSize: 12,
    fontFamily: "'JetBrains Mono', monospace",
    theme: {
      background: '#0a0a0a',
      foreground: '#a7a7a7',
      cursor: '#6798ff',
      selectionBackground: '#6798ff33',
      black: '#1e1e1e',
      red: '#ef4444',
      green: '#22c55e',
      yellow: '#eab308',
      blue: '#6798ff',
      magenta: '#a855f7',
      cyan: '#06b6d4',
      white: '#a7a7a7',
      brightBlack: '#454545',
      brightRed: '#f87171',
      brightGreen: '#4ade80',
      brightYellow: '#facc15',
      brightBlue: '#93c5fd',
      brightMagenta: '#c084fc',
      brightCyan: '#67e8f9',
      brightWhite: '#ffffff',
    },
    cols: 120,
    rows: 20,
  })
  terminal.open(terminalRef.value)
  terminal.writeln('\x1b[1;37m══════════════════════════════════════════════════════════════\x1b[0m')
  terminal.writeln(`\x1b[1;36m  LIVE TRAFFIC MONITOR — ${props.domain}\x1b[0m`)
  terminal.writeln('\x1b[1;37m══════════════════════════════════════════════════════════════\x1b[0m')
  terminal.writeln('')
}

watch(
  () => props.show,
  async (val) => {
    if (val) {
      await nextTick()
      initTerminal()
      connectStream()
    } else {
      disconnectStream()
    }
  }
)

onUnmounted(() => {
  disconnectStream()
  terminal?.dispose()
})
</script>

<template>
  <div
    v-if="show"
    class="fixed inset-0 z-50 flex items-center justify-center bg-black/70 backdrop-blur-sm p-4"
  >
    <div class="w-full max-w-4xl bg-card-carbon border border-steel-border rounded-lg overflow-hidden flex flex-col max-h-[85vh]">
      <div class="px-5 py-3 border-b border-steel-border flex items-center justify-between shrink-0">
        <div class="flex items-center gap-3">
          <div class="flex gap-1.5">
            <div
              class="w-3 h-3 rounded-full bg-red-400/60 cursor-pointer hover:bg-red-400"
              @click="emit('close')"
            />
            <div class="w-3 h-3 rounded-full bg-yellow-400/60" />
            <div class="w-3 h-3 rounded-full bg-green-400/60" />
          </div>
          <h3 class="text-[13px] font-semibold text-snow font-inter ml-2">traffic — {{ domain }}</h3>
          <span class="px-1.5 py-0.5 rounded-[3px] bg-blue-cornflower/10 border border-blue-cornflower/20 text-[9px] font-jetbrains-mono text-blue-cornflower uppercase tracking-wider">LIVE</span>
        </div>
        <div class="flex items-center gap-2 text-[10px] font-jetbrains-mono text-ash uppercase tracking-wider">
          <span>Request stream</span>
          <span class="w-1.5 h-1.5 rounded-full bg-green-500 animate-pulse" />
        </div>
      </div>
      <div ref="terminalRef" class="flex-1 overflow-hidden p-2" />
    </div>
  </div>
</template>
