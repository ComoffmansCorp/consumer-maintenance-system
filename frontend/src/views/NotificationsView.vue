<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import { notificationsApi } from '@/api/notifications'
import type { NotificationDTO } from '@/types'
import { extractErrorMessage } from '@/api/client'
import { useToastStore } from '@/stores/toast'
import LoadingState from '@/components/ui/LoadingState.vue'
import ErrorState from '@/components/ui/ErrorState.vue'
import EmptyState from '@/components/ui/EmptyState.vue'
import AppButton from '@/components/ui/AppButton.vue'
import Pagination from '@/components/ui/Pagination.vue'

const toast = useToastStore()

const items = ref<NotificationDTO[]>([])
const loading = ref(true)
const error = ref('')
const page = ref(1)
const totalPages = ref(1)
const totalItems = ref(0)
const unreadOnly = ref(false)

async function load() {
  loading.value = true
  error.value = ''
  try {
    const result = await notificationsApi.list({ page: page.value, pageSize: 20, unread: unreadOnly.value })
    items.value = result.items
    totalPages.value = result.totalPages
    totalItems.value = result.totalItems
  } catch (e) {
    error.value = extractErrorMessage(e)
  } finally {
    loading.value = false
  }
}
onMounted(load)
watch(page, load)
watch(unreadOnly, () => {
  page.value = 1
  load()
})

async function markRead(n: NotificationDTO) {
  if (n.read) return
  try {
    await notificationsApi.markRead(n.id)
    n.read = true
  } catch (e) {
    toast.error(extractErrorMessage(e))
  }
}

async function markAllRead() {
  try {
    await notificationsApi.markAllRead()
    items.value.forEach((n) => (n.read = true))
    toast.success('Все уведомления отмечены прочитанными')
  } catch (e) {
    toast.error(extractErrorMessage(e))
  }
}

function formatDate(iso: string) {
  return new Date(iso).toLocaleString('ru-RU')
}
</script>

<template>
  <div class="flex flex-col gap-6">
    <div class="flex flex-wrap items-center justify-between gap-3">
      <h1 class="font-display text-2xl font-semibold text-ink dark:text-surface">Уведомления</h1>
      <div class="flex items-center gap-3">
        <label class="flex items-center gap-2 text-sm text-slate dark:text-mist">
          <input v-model="unreadOnly" type="checkbox" class="accent-primary" />
          Только непрочитанные
        </label>
        <AppButton variant="secondary" size="sm" @click="markAllRead">Прочитать все</AppButton>
      </div>
    </div>

    <LoadingState v-if="loading" />
    <ErrorState v-else-if="error" :message="error" @retry="load" />
    <EmptyState v-else-if="items.length === 0" title="Уведомлений нет" description="Здесь появятся события по вашим нарядам" />
    <div v-else class="card divide-y divide-line dark:divide-graphite">
      <button
        v-for="n in items"
        :key="n.id"
        class="flex w-full items-start gap-3 px-4 py-3 text-left hover:bg-surface dark:hover:bg-ink"
        @click="markRead(n)"
      >
        <span
          class="mt-1.5 h-2 w-2 shrink-0 rounded-full"
          :style="{ backgroundColor: n.read ? 'transparent' : 'var(--color-primary)' }"
        />
        <div class="flex-1">
          <p class="text-sm font-medium text-ink dark:text-surface">{{ n.title }}</p>
          <p class="text-sm text-slate dark:text-mist">{{ n.message }}</p>
          <p class="mt-1 text-xs text-mist">{{ formatDate(n.createdAt) }}</p>
        </div>
      </button>
      <div class="px-4">
        <Pagination v-model:page="page" :total-pages="totalPages" :total-items="totalItems" />
      </div>
    </div>
  </div>
</template>
