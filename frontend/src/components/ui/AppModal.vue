<script setup lang="ts">
defineProps<{ title: string }>()
const emit = defineEmits<{ close: [] }>()
</script>

<template>
  <Teleport to="body">
    <div class="fixed inset-0 z-40 flex items-center justify-center bg-ink/40 p-4" @click.self="emit('close')">
      <div class="card w-full max-w-lg p-6 shadow-lg">
        <div class="mb-4 flex items-center justify-between">
          <h2 class="font-display text-lg font-semibold text-ink dark:text-surface">{{ title }}</h2>
          <button
            class="rounded-md p-1 text-mist hover:bg-surface hover:text-slate dark:hover:bg-ink"
            aria-label="Закрыть"
            @click="emit('close')"
          >
            ✕
          </button>
        </div>
        <slot />
      </div>
    </div>
  </Teleport>
</template>
