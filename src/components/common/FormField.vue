<script setup lang="ts">
defineProps<{
  label: string
  error?: string
  id: string
  required?: boolean
}>()
</script>

<template>
  <div class="flex flex-col gap-2 mb-4 w-full">
    <label :for="id" class="text-caption font-jetbrains-mono tracking-caption text-ash uppercase flex items-center gap-1">
      {{ label }}
      <span v-if="required" class="text-red-500">*</span>
    </label>
    <div>
      <slot></slot>
    </div>
    <span v-if="error" class="text-xs text-red-400 font-inter mt-1">{{ error }}</span>
  </div>
</template>
