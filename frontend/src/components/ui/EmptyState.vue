<script setup lang="ts">
withDefaults(defineProps<{ title: string; description?: string }>(), {})
</script>

<template>
  <div class="flex flex-col items-center justify-center gap-3 rounded-md border border-dashed border-line px-6 py-16 text-center dark:border-graphite">
    <div class="flex h-10 w-10 items-center justify-center rounded-full bg-surface text-mist dark:bg-graphite">
      <svg viewBox="0 0 24 24" fill="none" class="h-5 w-5">
        <path
          d="M4 7h16M4 12h16M4 17h10"
          stroke="currentColor"
          stroke-width="1.5"
          stroke-linecap="round"
        />
      </svg>
    </div>
    <p class="font-display text-base font-semibold text-ink dark:text-surface">{{ title }}</p>
    <p v-if="description" class="max-w-sm text-sm text-slate dark:text-mist">{{ description }}</p>
    <div v-if="$slots.action" class="mt-1">
      <slot name="action" />
    </div>
  </div>
</template>
