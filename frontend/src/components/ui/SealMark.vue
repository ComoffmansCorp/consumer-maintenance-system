<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{ state: 'INTACT' | 'BROKEN' | 'MISSING' }>()

const color = computed(
  () =>
    ({
      INTACT: 'var(--color-status-completed)',
      BROKEN: 'var(--color-status-canceled)',
      MISSING: 'var(--color-mist)',
    })[props.state],
)
</script>

<template>
  <svg viewBox="0 0 20 20" class="h-4 w-4 shrink-0" :style="{ color }" aria-hidden="true">
    <!-- Wire-seal tag: intact = closed loop with pressed tab; broken = loop
         severed; missing = dashed outline, nothing to press. -->
    <template v-if="state === 'INTACT'">
      <circle cx="10" cy="11" r="6.5" fill="none" stroke="currentColor" stroke-width="1.6" />
      <path d="M10 4.5v2.2" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" />
      <circle cx="10" cy="11" r="2" fill="currentColor" />
    </template>
    <template v-else-if="state === 'BROKEN'">
      <path
        d="M4.6 8a6.5 6.5 0 0 1 10.8 0M15.4 14a6.5 6.5 0 0 1-10.8 0"
        fill="none"
        stroke="currentColor"
        stroke-width="1.6"
      />
      <path d="M10 4.5v2.2" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" />
      <path d="M7.2 9.2 12.8 12.8M12.8 9.2 7.2 12.8" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" />
    </template>
    <template v-else>
      <circle cx="10" cy="11" r="6.5" fill="none" stroke="currentColor" stroke-width="1.4" stroke-dasharray="2.2 2.2" />
    </template>
  </svg>
</template>
