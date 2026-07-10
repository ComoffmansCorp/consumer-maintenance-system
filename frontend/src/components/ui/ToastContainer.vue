<script setup lang="ts">
import { useToastStore } from '@/stores/toast'

const toast = useToastStore()
</script>

<template>
  <div class="pointer-events-none fixed bottom-4 right-4 z-50 flex flex-col gap-2">
    <TransitionGroup name="toast">
      <div
        v-for="t in toast.toasts"
        :key="t.id"
        class="pointer-events-auto flex min-w-64 max-w-sm items-center gap-2 rounded-md border px-4 py-3 text-sm shadow-sm"
        :class="
          t.variant === 'success'
            ? 'border-line bg-white text-ink dark:border-graphite dark:bg-graphite dark:text-surface'
            : 'border-[color-mix(in_srgb,var(--color-status-canceled)_35%,transparent)] bg-white text-ink dark:bg-graphite dark:text-surface'
        "
      >
        <span
          class="h-1.5 w-1.5 shrink-0 rounded-full"
          :style="{
            backgroundColor:
              t.variant === 'success' ? 'var(--color-status-completed)' : 'var(--color-status-canceled)',
          }"
        />
        <span>{{ t.message }}</span>
        <button class="ml-auto text-mist hover:text-slate" @click="toast.dismiss(t.id)">✕</button>
      </div>
    </TransitionGroup>
  </div>
</template>

<style scoped>
.toast-enter-active,
.toast-leave-active {
  transition:
    opacity 0.15s ease,
    transform 0.15s ease;
}
.toast-enter-from,
.toast-leave-to {
  opacity: 0;
  transform: translateY(8px);
}
</style>
