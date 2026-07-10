<script setup lang="ts">
withDefaults(defineProps<{ message: string; retryable?: boolean }>(), { retryable: true })
defineEmits<{ retry: [] }>()
</script>

<template>
  <div
    class="flex flex-col items-center gap-3 rounded-md border border-[color-mix(in_srgb,var(--color-status-canceled)_35%,transparent)] bg-[color-mix(in_srgb,var(--color-status-canceled)_6%,transparent)] px-6 py-10 text-center"
  >
    <p class="text-sm font-medium" style="color: var(--color-status-canceled)">{{ message }}</p>
    <button
      v-if="retryable"
      class="text-sm font-medium text-primary hover:underline"
      @click="$emit('retry')"
    >
      Повторить
    </button>
  </div>
</template>
