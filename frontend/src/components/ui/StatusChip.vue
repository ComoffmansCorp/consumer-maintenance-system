<script setup lang="ts">
import { computed } from 'vue'
import type { TaskStatus } from '@/types'

const props = defineProps<{ status: TaskStatus }>()

const labels: Record<TaskStatus, string> = {
  PENDING: 'Ожидает',
  IN_PROGRESS: 'В работе',
  COMPLETED: 'Выполнен',
  CANCELED: 'Отменён',
}

const colors: Record<TaskStatus, string> = {
  PENDING: 'var(--color-status-created)',
  IN_PROGRESS: 'var(--color-status-in-progress)',
  COMPLETED: 'var(--color-status-completed)',
  CANCELED: 'var(--color-status-canceled)',
}

const color = computed(() => colors[props.status])
const label = computed(() => labels[props.status])
</script>

<template>
  <span
    class="seal-tag text-ink dark:text-surface"
    :style="{ borderLeftColor: color }"
  >
    {{ label }}
  </span>
</template>
