<script setup lang="ts">
import { ref, onMounted } from 'vue'

const showAnnouncement = ref(true)
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
  { name: 'Product', disabled: true, path: '/dashboard' },
  { name: 'Overview', disabled: false, path: '/dashboard' },
  { name: 'Routes', disabled: false, path: '/routes' },
  { name: 'Block List', disabled: false, path: '/blocklist' },
  { name: 'Logs', disabled: false, path: '/logs' },
  { name: 'Settings', disabled: false, path: '/settings' },
]

const trustLogos = ['Shopify', 'AWS', 'Atlassian', 'Notion', 'Canva', 'Lovable']

onMounted(() => {
  document.documentElement.classList.add('dark')
})
</script>

<template>
  <div class="min-h-screen bg-page-ink flex flex-col font-inter relative">
    <!-- Announcement Bar -->
    <div
      v-if="showAnnouncement"
      class="relative w-full bg-blue-cornflower flex items-center justify-center gap-3 py-2 px-4 z-50 select-none"
    >
      <div class="flex items-center gap-2">
        <div class="flex -space-x-1.5">
          <div class="w-4 h-4 rounded-full bg-white/30 border border-white/20" />
          <div class="w-4 h-4 rounded-full bg-white/30 border border-white/20" />
          <div class="w-4 h-4 rounded-full bg-white/30 border border-white/20" />
        </div>
        <span class="text-[11px] font-inter font-medium uppercase tracking-wider text-white">
          The best never guess &mdash; proxy command center live
        </span>
      </div>
      <button
        type="button"
        @click="showAnnouncement = false"
        class="absolute right-3 w-5 h-5 flex items-center justify-center rounded-[4px] text-white/70 hover:text-white hover:bg-white/10 transition-colors cursor-pointer"
      >
        <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round">
          <line x1="18" y1="6" x2="6" y2="18" /><line x1="6" y1="6" x2="18" y2="18" />
        </svg>
      </button>
    </div>

    <!-- Header / Nav -->
    <header class="h-16 border-b border-steel-border/60 bg-page-ink/80 backdrop-blur-sm flex items-center justify-between px-8 sticky top-0 z-40 select-none">
      <div class="flex items-center gap-8">
        <!-- Logo -->
        <div class="flex items-center gap-2.5">
          <svg width="24" height="24" viewBox="0 0 32 32" fill="none" class="shrink-0">
            <rect width="32" height="32" rx="8" fill="#6798ff" />
            <path d="M8 22V10l8 6-8 6z" fill="#0a0a0a" />
            <path d="M18 22V10l8 6-8 6z" fill="#0a0a0a" opacity="0.6"/>
          </svg>
          <span class="text-[15px] font-semibold text-snow tracking-tight">Dovetail</span>
        </div>

        <!-- Pill Nav - Center -->
        <nav class="hidden md:flex items-center gap-1">
          <template v-for="item in navItems" :key="item.name">
            <button
              v-if="item.disabled"
              type="button"
              disabled
              class="px-3.5 py-1.5 rounded-full text-[13px] font-medium text-fog cursor-not-allowed select-none"
            >
              {{ item.name }}
            </button>
            <router-link
              v-else
              :to="item.path"
              class="px-3.5 py-1.5 rounded-full text-[13px] font-medium transition-all duration-150 select-none"
              :class="
                $route.path === item.path
                  ? 'bg-card-carbon text-blue-cornflower'
                  : 'text-ash hover:text-snow hover:bg-card-carbon/50'
              "
            >
              {{ item.name }}
            </router-link>
          </template>
        </nav>
      </div>

      <!-- Right Actions -->
      <div class="flex items-center gap-3">
        <button
          type="button"
          @click="toggleDarkMode"
          class="p-1.5 rounded-lg text-ash hover:text-snow hover:bg-card-carbon/50 transition-colors cursor-pointer"
          title="Toggle theme"
        >
          <svg v-if="isDarkMode" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round">
            <circle cx="12" cy="12" r="5" />
            <line x1="12" y1="1" x2="12" y2="3" /><line x1="12" y1="21" x2="12" y2="23" />
            <line x1="4.22" y1="4.22" x2="5.64" y2="5.64" /><line x1="18.36" y1="18.36" x2="19.78" y2="19.78" />
            <line x1="1" y1="12" x2="3" y2="12" /><line x1="21" y1="12" x2="23" y2="12" />
          </svg>
          <svg v-else width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round">
            <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z" />
          </svg>
        </button>
        <span class="text-[13px] text-ash font-medium hidden sm:inline select-none">Admin</span>
        <div class="w-7 h-7 rounded-full bg-card-carbon border border-steel-border flex items-center justify-center select-none">
          <span class="text-[11px] font-semibold text-snow">A</span>
        </div>
      </div>
    </header>

    <!-- Main Content -->
    <main class="flex-1 max-w-[1200px] w-full mx-auto px-6 sm:px-8 py-10 relative z-10">
      <router-view />
    </main>

    <!-- Footer Trust Bar -->
    <footer class="border-t border-steel-border/40 py-10 px-8 select-none relative z-10">
      <div class="max-w-[1200px] mx-auto flex flex-col lg:flex-row items-center justify-between gap-6">
        <div class="flex flex-col items-center lg:items-start gap-3">
          <span class="text-caption font-jetbrains-mono tracking-caption text-ash uppercase">
            Connecting the world&rsquo;s leading companies to their customers
          </span>
          <div class="flex items-center gap-6 opacity-50">
            <span
              v-for="logo in trustLogos"
              :key="logo"
              class="text-[12px] font-semibold text-ash tracking-wider uppercase select-none"
            >
              {{ logo }}
            </span>
          </div>
        </div>
        <div class="flex items-center gap-8">
          <div class="flex flex-col items-end gap-0.5">
            <div class="flex items-center gap-0.5">
              <svg v-for="i in 4" :key="'g2-'+i" width="12" height="12" viewBox="0 0 24 24" fill="#ffffff">
                <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z" />
              </svg>
              <svg width="12" height="12" viewBox="0 0 24 24" class="opacity-50">
                <defs><linearGradient id="half-star"><stop offset="50%" stop-color="#ffffff"/><stop offset="50%" stop-color="transparent"/></linearGradient></defs>
                <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z" fill="url(#half-star)" stroke="#ffffff" stroke-width="0.5"/>
              </svg>
            </div>
            <span class="text-[10px] text-ash font-inter">4.5/5 G2</span>
          </div>
          <div class="flex flex-col items-end gap-0.5">
            <div class="flex items-center gap-0.5">
              <svg v-for="i in 4" :key="'cap-'+i" width="12" height="12" viewBox="0 0 24 24" fill="#ffffff">
                <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z" />
              </svg>
              <svg width="12" height="12" viewBox="0 0 24 24" class="opacity-50">
                <defs><linearGradient id="half-star-cap"><stop offset="50%" stop-color="#ffffff"/><stop offset="50%" stop-color="transparent"/></linearGradient></defs>
                <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z" fill="url(#half-star-cap)" stroke="#ffffff" stroke-width="0.5"/>
              </svg>
            </div>
            <span class="text-[10px] text-ash font-inter">4.6/5 Capterra</span>
          </div>
        </div>
      </div>
    </footer>
  </div>
</template>
