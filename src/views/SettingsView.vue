<script setup lang="ts">
const systemSummary = {
  listenPort: '8080 (HTTP) / 443 (HTTPS when SSL active)',
  serverPort: '8000 (API & Dashboard Admin)',
  database: 'SQLite3 + GORM AutoMigrate Enabled',
  authMethod: 'Basic Authentication'
}
</script>

<template>
  <div class="space-y-6 max-w-3xl">
    <div>
      <span class="text-caption font-jetbrains-mono tracking-caption text-blue-cornflower uppercase font-medium">SYSTEM</span>
      <h2 class="text-3xl font-semibold text-snow mt-1 tracking-tight">Settings</h2>
    </div>

    <!-- Auth Verification Panel -->
    <div class="bg-card-carbon border border-steel-border rounded-lg p-6 space-y-4">
      <h3 class="text-lg font-semibold text-snow">Administrative Session</h3>
      <div class="flex items-center gap-3">
        <div
          class="w-3.5 h-3.5 rounded-full bg-green-500 animate-pulse"
        ></div>
        <span class="text-sm text-green-400 font-semibold">
          Session Authenticated
        </span>
      </div>
      <p class="text-sm text-ash">
        Logged in via Browser Basic Authentication.
      </p>
    </div>

    <!-- System Setup Overview -->
    <div class="bg-card-carbon border border-steel-border rounded-lg p-6 space-y-4">
      <h3 class="text-lg font-semibold text-snow">Proxy Server Configuration</h3>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div class="border border-steel-border rounded-lg p-4 bg-deep-coal">
          <span class="text-caption font-jetbrains-mono tracking-caption text-ash uppercase">PROXY LISTEN PORTS</span>
          <div class="text-sm font-medium text-snow mt-1.5 font-jetbrains-mono">{{ systemSummary.listenPort }}</div>
        </div>
        <div class="border border-steel-border rounded-lg p-4 bg-deep-coal">
          <span class="text-caption font-jetbrains-mono tracking-caption text-ash uppercase">ADMIN SERVER PORT</span>
          <div class="text-sm font-medium text-snow mt-1.5 font-jetbrains-mono">{{ systemSummary.serverPort }}</div>
        </div>
        <div class="border border-steel-border rounded-lg p-4 bg-deep-coal">
          <span class="text-caption font-jetbrains-mono tracking-caption text-ash uppercase">STORAGE ENGINE</span>
          <div class="text-sm font-medium text-snow mt-1.5 font-jetbrains-mono">{{ systemSummary.database }}</div>
        </div>
        <div class="border border-steel-border rounded-lg p-4 bg-deep-coal">
          <span class="text-caption font-jetbrains-mono tracking-caption text-ash uppercase">API SECURITY</span>
          <div class="text-sm font-medium text-snow mt-1.5 font-jetbrains-mono">{{ systemSummary.authMethod }}</div>
        </div>
      </div>
    </div>
  </div>
</template>
