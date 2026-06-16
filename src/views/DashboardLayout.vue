<script setup lang="ts">
import { ref, onMounted } from 'vue'

const showAnnouncement = ref(true)
const isDarkMode = ref(true)
const sidebarOpen = ref(true)

const applyTheme = (dark: boolean) => {
  const html = document.documentElement
  if (dark) {
    html.classList.add('dark')
    html.classList.remove('light')
  } else {
    html.classList.remove('dark')
    html.classList.add('light')
  }
}

const toggleDarkMode = () => {
  isDarkMode.value = !isDarkMode.value
  applyTheme(isDarkMode.value)
  localStorage.setItem('theme', isDarkMode.value ? 'dark' : 'light')
}

const toggleSidebar = () => {
  sidebarOpen.value = !sidebarOpen.value
}

const closeSidebar = () => {
  sidebarOpen.value = false
}

const navItems = [
  { name: 'Overview', path: '/dashboard', icon: 'M3 13.125C3 12.504 3.504 12 4.125 12h2.25c.621 0 1.125.504 1.125 1.125v6.75C7.5 20.496 6.996 21 6.375 21h-2.25A1.125 1.125 0 013 19.875v-6.75zM9.75 8.625c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125v11.25c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 01-1.125-1.125V8.625zM16.5 4.125c0-.621.504-1.125 1.125-1.125h2.25C20.496 3 21 3.504 21 4.125v15.75c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 01-1.125-1.125V4.125z' },
  { name: 'Routes', path: '/routes', icon: 'M8 7h12m0 0l-4-4m4 4l-4 4m0 6H4m0 0l4 4m-4-4l4-4' },
  { name: 'Block List', path: '/blocklist', icon: 'M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z' },
  { name: 'Logs', path: '/logs', icon: 'M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01' },
]

onMounted(() => {
  const saved = localStorage.getItem('theme')
  if (saved === 'light') {
    isDarkMode.value = false
    applyTheme(false)
  } else {
    isDarkMode.value = true
    applyTheme(true)
  }
})
</script>

<template>
  <div class="min-h-screen bg-page-ink flex flex-col font-inter relative">
    <!-- Announcement Bar -->
    <div
      v-if="showAnnouncement"
      class="relative w-full bg-blue-cornflower flex items-center justify-center gap-2.5 py-1.5 px-4 z-50 select-none"
    >
      <div class="flex items-center gap-2.5">
        <svg width="18" height="12" viewBox="0 0 18 12" fill="none" class="shrink-0">
          <ellipse cx="5" cy="6" rx="4" ry="5" fill="white" opacity="0.5" />
          <ellipse cx="12" cy="6" rx="4" ry="5" fill="white" opacity="0.3" />
        </svg>
        <span class="text-[11px] font-inter font-semibold uppercase tracking-wider text-white">
          The best never guess &mdash; proxy command center live
        </span>
      </div>
      <button
        type="button"
        @click="showAnnouncement = false"
        class="absolute right-3 w-5 h-5 flex items-center justify-center rounded-[4px] text-white/70 hover:text-white hover:bg-white/10 transition-colors cursor-pointer"
      >
        <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round">
          <line x1="18" y1="6" x2="6" y2="18" /><line x1="6" y1="6" x2="18" y2="18" />
        </svg>
      </button>
    </div>

    <div class="flex flex-1 overflow-hidden relative">
      <!-- Mobile overlay backdrop -->
      <div
        v-if="sidebarOpen"
        class="fixed inset-0 z-30 bg-black/50 lg:hidden"
        @click="closeSidebar"
      />

      <!-- Sidebar -->
      <aside
        class="border-r border-steel-border/60 bg-deep-coal flex flex-col shrink-0 select-none transition-all duration-200 overflow-hidden z-40
          fixed inset-y-0 left-0 lg:relative lg:inset-auto"
        :class="sidebarOpen ? '' : '-translate-x-full lg:translate-x-0'"
        :style="sidebarOpen
          ? { width: '232px', minWidth: '232px', maxWidth: '232px', flex: '0 0 232px' }
          : { width: '60px', minWidth: '60px', maxWidth: '60px', flex: '0 0 60px' }">
        <!-- App Identity Row -->
        <div class="flex items-center gap-2.5 px-4 pt-5 pb-3 shrink-0">
          <svg width="24" height="24" viewBox="0 0 32 32" fill="none" class="shrink-0">
            <rect width="32" height="32" rx="8" fill="#6798ff" />
            <path d="M8 22V10l8 6-8 6z" fill="#0a0a0a" />
            <path d="M18 22V10l8 6-8 6z" fill="#0a0a0a" opacity="0.6" />
          </svg>
          <span
            class="text-[15px] font-semibold text-snow tracking-tight whitespace-nowrap transition-opacity duration-200"
            :class="sidebarOpen ? 'opacity-100' : 'opacity-0'"
          >Golang Proxy</span>
          <button
            type="button"
            @click="toggleSidebar"
            class="absolute right-0 translate-x-1/2 z-20 w-6 h-6 flex items-center justify-center rounded-full bg-card-carbon border border-steel-border text-ash hover:text-snow hover:border-blue-cornflower/40 transition-colors cursor-pointer hidden lg:flex"
            :style="{ top: '20px' }"
            title="Toggle sidebar"
          >
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round">
              <line x1="3" y1="6" x2="21" y2="6" />
              <line x1="3" y1="12" x2="21" y2="12" />
              <line x1="3" y1="18" x2="21" y2="18" />
            </svg>
          </button>
        </div>

        <div class="mx-3 mb-3 border-t border-steel-border/40" />
        <nav class="flex flex-col gap-0.5 px-3 pb-5">
          <router-link
            v-for="item in navItems" :key="item.path" :to="item.path"
            @click="closeSidebar"
            class="flex items-center gap-3 px-3 py-2.5 rounded-full text-[13px] font-medium transition-all duration-150 select-none border whitespace-nowrap"
            :class="$route.path === item.path ? 'border-blue-cornflower/40 bg-blue-cornflower/5 text-blue-cornflower' : 'border-transparent text-ash hover:text-snow hover:bg-card-carbon/40'"
          >
            <svg width="17" height="17" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" style="width:17px;height:17px;min-width:17px;min-height:17px;flex-shrink:0;">
              <path :d="item.icon" />
            </svg>
            <span class="transition-opacity duration-200" :class="sidebarOpen ? 'opacity-100' : 'opacity-0'">{{ item.name }}</span>
          </router-link>
        </nav>
      </aside>

      <!-- Main Content -->
      <main class="flex-1 overflow-y-auto bg-page-ink relative z-10">
        <div class="sticky top-0 z-30 flex items-center justify-between px-4 py-2.5 bg-page-ink/80 backdrop-blur-sm border-b border-steel-border/60 lg:justify-end">
          <button type="button" @click="toggleSidebar" class="p-1.5 rounded-lg text-ash hover:text-snow hover:bg-card-carbon/50 transition-colors cursor-pointer lg:hidden" title="Menu">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round">
              <line x1="3" y1="6" x2="21" y2="6" /><line x1="3" y1="12" x2="21" y2="12" /><line x1="3" y1="18" x2="21" y2="18" />
            </svg>
          </button>
          <button type="button" @click="toggleDarkMode" class="p-1.5 rounded-lg text-ash hover:text-snow hover:bg-card-carbon/50 transition-colors cursor-pointer" title="Toggle theme">
            <svg v-if="isDarkMode" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round">
              <circle cx="12" cy="12" r="5" /><line x1="12" y1="1" x2="12" y2="3" /><line x1="12" y1="21" x2="12" y2="23" /><line x1="4.22" y1="4.22" x2="5.64" y2="5.64" /><line x1="18.36" y1="18.36" x2="19.78" y2="19.78" /><line x1="1" y1="12" x2="3" y2="12" /><line x1="21" y1="12" x2="23" y2="12" />
            </svg>
            <svg v-else width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round">
              <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z" />
            </svg>
          </button>
        </div>
        <div class="max-w-[1200px] mx-auto px-4 sm:px-6 py-6">
          <router-view />
        </div>
      </main>
    </div>
  </div>
</template>
