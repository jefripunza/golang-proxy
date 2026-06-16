<script setup lang="ts">
import { ref, onMounted } from 'vue'

const isDarkMode = ref(true)

const toggleDarkMode = () => {
  isDarkMode.value = !isDarkMode.value
  const html = document.documentElement
  if (isDarkMode.value) {
    html.classList.add('dark')
    html.classList.remove('light')
  } else {
    html.classList.remove('dark')
    html.classList.add('light')
  }
}

const navItems = [
  { name: 'Overview', path: '/dashboard', icon: 'M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z' },
  { name: 'Proxy Routes', path: '/routes', icon: 'M8 7h12m0 0l-4-4m4 4l-4 4m0 6H4m0 0l4 4m-4-4l4-4' },
  { name: 'IP Block List', path: '/blocklist', icon: 'M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z' },
  { name: 'Activity Logs', path: '/logs', icon: 'M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01' },
  { name: 'Settings', path: '/settings', icon: 'M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z' },
]

onMounted(() => {
  document.documentElement.classList.add('dark')
})
</script>

<template>
  <div class="min-h-screen bg-page-ink flex flex-col font-inter">
    <!-- Header -->
    <header class="h-16 border-b border-steel-border bg-deep-coal flex items-center justify-between px-8 sticky top-0 z-40 select-none">
      <div class="flex items-center gap-3">
        <div class="w-8 h-8 rounded-lg bg-blue-cornflower flex items-center justify-center font-bold text-page-ink text-base select-none">D</div>
        <div class="flex flex-col">
          <h1 class="text-base font-semibold text-snow leading-none tracking-tight">Dovetail</h1>
          <span class="text-[10px] font-jetbrains-mono tracking-wider text-ash uppercase mt-1">Dynamic Proxy Command Center</span>
        </div>
      </div>

      <div class="flex items-center gap-4">
        <!-- Dark Mode Toggle -->
        <button
          type="button"
          @click="toggleDarkMode"
          class="p-2 border border-graphite rounded-md text-ash hover:text-snow hover:border-steel-border transition-colors cursor-pointer flex items-center justify-center"
          title="Toggle Dark Mode"
        >
          <svg
            v-if="isDarkMode"
            xmlns="http://www.w3.org/2000/svg"
            width="18"
            height="18"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
            style="width: 18px; height: 18px; min-width: 18px; min-height: 18px;"
          >
            <circle cx="12" cy="12" r="5" />
            <line x1="12" y1="1" x2="12" y2="3" />
            <line x1="12" y1="21" x2="12" y2="23" />
            <line x1="4.22" y1="4.22" x2="5.64" y2="5.64" />
            <line x1="18.36" y1="18.36" x2="19.78" y2="19.78" />
            <line x1="1" y1="12" x2="3" y2="12" />
            <line x1="21" y1="12" x2="23" y2="12" />
            <line x1="4.22" y1="19.78" x2="5.64" y2="18.36" />
            <line x1="18.36" y1="5.64" x2="19.78" y2="4.22" />
          </svg>
          <svg
            v-else
            xmlns="http://www.w3.org/2000/svg"
            width="18"
            height="18"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
            style="width: 18px; height: 18px; min-width: 18px; min-height: 18px;"
          >
            <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z" />
          </svg>
        </button>
      </div>
    </header>

    <div class="flex flex-1 overflow-hidden">
      <!-- Sidebar -->
      <aside 
        class="border-r border-steel-border bg-deep-coal flex flex-col p-6 gap-2 shrink-0 select-none" 
        style="width: 256px; min-width: 256px; max-width: 256px; flex: 0 0 256px;"
      >
        <router-link
          v-for="item in navItems"
          :key="item.path"
          :to="item.path"
          class="flex items-center gap-3 px-4 py-3 rounded-lg text-ash hover:text-snow hover:bg-card-carbon transition-colors group"
          active-class="bg-card-carbon !text-blue-cornflower"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="20"
            height="20"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
            style="width: 20px; height: 20px; min-width: 20px; min-height: 20px; flex-shrink: 0;"
          >
            <path :d="item.icon" />
          </svg>
          <span class="text-sm font-medium whitespace-nowrap">{{ item.name }}</span>
        </router-link>
      </aside>

      <!-- Main Content -->
      <main class="flex-1 p-8 overflow-y-auto bg-page-ink">
        <router-view />
      </main>
    </div>
  </div>
</template>

