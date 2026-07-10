<script setup lang="ts">
const props = defineProps<{ page: number; totalPages: number; totalItems: number }>()
const emit = defineEmits<{ 'update:page': [value: number] }>()

function go(page: number) {
  if (page < 1 || page > props.totalPages || page === props.page) return
  emit('update:page', page)
}
</script>

<template>
  <div v-if="totalPages > 1" class="flex items-center justify-between border-t border-line px-1 py-3 text-sm dark:border-graphite">
    <span class="text-slate dark:text-mist">Всего: {{ totalItems }}</span>
    <div class="flex items-center gap-1">
      <button
        class="rounded-md px-2.5 py-1 text-slate hover:bg-surface disabled:opacity-40 dark:text-mist dark:hover:bg-graphite"
        :disabled="page <= 1"
        @click="go(page - 1)"
      >
        Назад
      </button>
      <span class="data-mono px-2 text-ink dark:text-surface">{{ page }} / {{ totalPages }}</span>
      <button
        class="rounded-md px-2.5 py-1 text-slate hover:bg-surface disabled:opacity-40 dark:text-mist dark:hover:bg-graphite"
        :disabled="page >= totalPages"
        @click="go(page + 1)"
      >
        Далее
      </button>
    </div>
  </div>
</template>
