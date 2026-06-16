<script setup lang="ts">
const systemSummary = {
  listenPort: '8080 (HTTP) / 443 (HTTPS when SSL active)',
  serverPort: '8000 (API & Dashboard Admin)',
  database: 'SQLite3 + GORM AutoMigrate Enabled',
  authMethod: 'Basic Authentication'
}
</script>

<template>
  <div class="space-y-8 max-w-2xl">
    <div>
      <span class="text-caption font-jetbrains-mono tracking-caption text-blue-cornflower uppercase font-medium">SYSTEM</span>
      <h2 class="text-heading-sm font-semibold text-snow mt-1 tracking-tight">Settings</h2>
      <p class="text-body-sm text-ash mt-1">Proxy server configuration and administrative session overview.</p>
    </div>

    <!-- Auth Verification Panel -->
    <div class="bg-card-carbon border border-steel-border rounded-lg p-6 space-y-3">
      <h3 class="text-[15px] font-semibold text-snow">Administrative Session</h3>
      <div class="flex items-center gap-2.5">
        <div class="w-2.5 h-2.5 rounded-full bg-green-500 animate-pulse" />
        <span class="text-[13px] text-green-400 font-semibold font-jetbrains-mono uppercase tracking-wider">Authenticated</span>
      </div>
      <p class="text-body-sm text-ash">Secured via Browser Basic Authentication.</p>
    </div>

    <!-- System Setup Overview -->
    <div class="bg-card-carbon border border-steel-border rounded-lg p-6 space-y-4">
      <h3 class="text-[15px] font-semibold text-snow">Proxy Server Configuration</h3>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
        <div class="border border-steel-border rounded-lg p-4 bg-deep-coal group">
          <span class="text-caption font-jetbrains-mono tracking-caption text-ash uppercase">PROXY LISTEN PORTS</span>
          <div class="text-[13px] font-medium text-snow mt-1.5 font-jetbrains-mono leading-relaxed">{{ systemSummary.listenPort }}</div>
        </div>
        <div class="border border-steel-border rounded-lg p-4 bg-deep-coal group">
          <span class="text-caption font-jetbrains-mono tracking-caption text-ash uppercase">ADMIN SERVER PORT</span>
          <div class="text-[13px] font-medium text-snow mt-1.5 font-jetbrains-mono leading-relaxed">{{ systemSummary.serverPort }}</div>
        </div>
        <div class="border border-steel-border rounded-lg p-4 bg-deep-coal group">
          <span class="text-caption font-jetbrains-mono tracking-caption text-ash uppercase">STORAGE ENGINE</span>
          <div class="text-[13px] font-medium text-snow mt-1.5 font-jetbrains-mono leading-relaxed">{{ systemSummary.database }}</div>
        </div>
        <div class="border border-steel-border rounded-lg p-4 bg-deep-coal group">
          <span class="text-caption font-jetbrains-mono tracking-caption text-ash uppercase">API SECURITY</span>
          <div class="text-[13px] font-medium text-snow mt-1.5 font-jetbrains-mono leading-relaxed">{{ systemSummary.authMethod }}</div>
        </div>
      </div>
    </div>
  </div>
</template>
