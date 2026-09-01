<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { requestsApi } from '@/api/marketplace'
import type { RequestDTO } from '@/types'
import { extractErrorMessage } from '@/api/client'
import MarketplaceShell from '@/components/layout/MarketplaceShell.vue'

const requestStatusLabels: Record<string, string> = {
  OPEN: 'Открыта',
  ASSIGNED: 'В работе',
  COMPLETED: 'Выполнена',
  CANCELED: 'Отменена',
}
const requestStatusStyles: Record<string, { bg: string; fg: string }> = {
  OPEN: { bg: '#EFEBE1', fg: '#55524A' },
  ASSIGNED: { bg: '#E7E3FC', fg: '#5B4BE0' },
  COMPLETED: { bg: '#DCEFE1', fg: '#2E7D4F' },
  CANCELED: { bg: '#FBE9E7', fg: '#B3261E' },
}

const router = useRouter()

const items = ref<RequestDTO[]>([])
const loading = ref(true)
const error = ref('')
const page = ref(1)
const totalPages = ref(1)
const totalItems = ref(0)

async function load() {
  loading.value = true
  error.value = ''
  try {
    const result = await requestsApi.listMine(page.value, 20)
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

function openRequest(r: RequestDTO) {
  router.push({ name: 'marketplace-request-detail', params: { id: r.id } })
}

function goNew() {
  router.push({ name: 'marketplace-new-request' })
}

function formatDate(iso: string) {
  return new Date(iso).toLocaleDateString('ru-RU')
}
</script>

<template>
  <MarketplaceShell>
    <div class="mx-auto max-w-[820px]">
      <div class="flex flex-wrap items-center justify-between gap-3">
        <h1 class="mk-display text-2xl font-bold tracking-tight">Мои заявки</h1>
        <button type="button" class="rounded-[11px] bg-[#5B4BE0] px-4.5 py-2.5 text-sm font-medium text-white hover:bg-[#4536BC]" @click="goNew">
          Новая заявка
        </button>
      </div>

      <div v-if="loading" class="flex items-center justify-center gap-2 py-16 text-sm text-[#8D8A7E]">
        <svg class="h-4 w-4 animate-spin motion-reduce:animate-none" viewBox="0 0 24 24" fill="none">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v4a4 4 0 00-4 4H4z" />
        </svg>
        Загрузка…
      </div>

      <div v-else-if="error" class="mt-6 flex flex-col items-center gap-3 rounded-2xl border border-[#F3D3CE] bg-[#FBF0EE] px-6 py-10 text-center">
        <p class="text-sm font-medium text-[#B3261E]">{{ error }}</p>
        <button type="button" class="text-sm font-medium text-[#5B4BE0] hover:underline" @click="load">Повторить</button>
      </div>

      <div v-else-if="items.length === 0" class="mt-6 flex flex-col items-center gap-3 rounded-2xl border border-dashed border-[#E2DED2] px-6 py-16 text-center">
        <p class="mk-display text-base font-semibold">Заявок пока нет</p>
        <p class="max-w-sm text-sm text-[#6E6B60]">Оставьте первую заявку — мастер откликнётся, как только возьмёт её в работу.</p>
        <button type="button" class="mt-1 rounded-[11px] bg-[#5B4BE0] px-4.5 py-2.5 text-sm font-medium text-white hover:bg-[#4536BC]" @click="goNew">
          Оставить заявку
        </button>
      </div>

      <div v-else class="mt-6 flex flex-col gap-3">
        <button
          v-for="r in items"
          :key="r.id"
          type="button"
          class="flex flex-col gap-2 rounded-2xl border border-[#E2DED2] bg-white p-5 text-left transition-colors hover:border-[#17160F] sm:flex-row sm:items-center sm:justify-between"
          @click="openRequest(r)"
        >
          <div>
            <p class="mk-display text-base font-semibold">{{ r.serviceName || 'Услуга' }}</p>
            <p class="mt-0.5 text-sm text-[#6E6B60]">{{ r.addressText }}</p>
            <p class="mt-1 font-mono text-xs text-[#9B978A]">{{ formatDate(r.createdAt) }}</p>
          </div>
          <span
            class="inline-flex w-fit items-center gap-1.5 rounded-[3px] border-l-[3px] px-2.5 py-1 text-xs font-medium tracking-wide"
            :style="{ borderLeftColor: requestStatusStyles[r.status].fg, background: requestStatusStyles[r.status].bg, color: requestStatusStyles[r.status].fg }"
          >
            {{ requestStatusLabels[r.status] }}
          </span>
        </button>

        <div v-if="totalPages > 1" class="flex items-center justify-between border-t border-[#E2DED2] px-1 py-3 text-sm">
          <span class="text-[#6E6B60]">Всего: {{ totalItems }}</span>
          <div class="flex items-center gap-1">
            <button type="button" class="rounded-md px-2.5 py-1 text-[#6E6B60] hover:bg-[#F1EDE3] disabled:opacity-40" :disabled="page <= 1" @click="page--">
              Назад
            </button>
            <span class="font-mono px-2">{{ page }} / {{ totalPages }}</span>
            <button type="button" class="rounded-md px-2.5 py-1 text-[#6E6B60] hover:bg-[#F1EDE3] disabled:opacity-40" :disabled="page >= totalPages" @click="page++">
              Далее
            </button>
          </div>
        </div>
      </div>
    </div>
  </MarketplaceShell>
</template>
