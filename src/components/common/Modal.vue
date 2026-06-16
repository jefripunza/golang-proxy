<script setup lang="ts">
defineProps<{
  show: boolean
  title: string
}>()

defineEmits<{
  (e: 'close'): void
}>()
</script>

<template>
  <div
    v-if="show"
    class="fixed inset-0 z-50 flex items-center justify-center bg-black/70 backdrop-blur-sm p-4"
  >
    <div
      class="w-full max-w-lg bg-card-carbon border border-steel-border rounded-lg overflow-hidden flex flex-col max-h-[90vh]"
    >
      <!-- Header -->
      <div class="px-6 py-4 border-b border-steel-border flex items-center justify-between">
        <h3 class="text-lg font-semibold text-snow font-inter">{{ title }}</h3>
        <button
          type="button"
          @click="$emit('close')"
          class="text-ash hover:text-snow transition-colors cursor-pointer"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            class="h-5 w-5"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M6 18L18 6M6 6l12 12"
            />
          </svg>
        </button>
      </div>

      <!-- Content -->
      <div class="p-6 overflow-y-auto flex-1 text-ash">
        <slot></slot>
      </div>

      <!-- Footer -->
      <div v-if="$slots.footer" class="px-6 py-4 border-t border-steel-border flex justify-end gap-3 bg-deep-coal">
        <slot name="footer"></slot>
      </div>
    </div>
  </div>
</template>
