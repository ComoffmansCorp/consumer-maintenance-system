<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import { adminApi } from '@/api/marketplace'
import { extractErrorMessage } from '@/api/client'
import type { RequestDTO } from '@/types'
import MarketplaceShell from '@/components/layout/MarketplaceShell.vue'

const requestStatusLabels: Record<string, string> = {
  OPEN: 'Открыта',
  ASSIGNED: 'В работе',
  COMPLETED: 'Выполнена',
  CANCELED: 'Отменена',
}

const items = ref<RequestDTO[]>([])
const loading = ref(true)
const error = ref('')
const page = ref(1)
const totalPages = ref(1)
const status = ref('')

async function load() {
  loading.value = true
  error.value = ''
  try {
    const result = await adminApi.listRequests(status.value || undefined, page.value, 20)
    items.value = result.items
    totalPages.value = result.totalPages
  } catch (e) {
    error.value = extractErrorMessage(e)
  } finally {
    loading.value = false
  }
}
onMounted(load)
watch([page, status], load)

function formatDate(iso: string) {
  return new Date(iso).toLocaleDateString('ru-RU')
}
</script>

<template>
  <MarketplaceShell>
    <div class="mx-auto max-w-[900px]">
      <div class="flex flex-wrap items-center justify-between gap-3">
        <h1 class="mk-display text-2xl font-bold tracking-tight">Все заявки</h1>
        <select v-model="status" class="rounded-xl border border-[#E2DED2] px-3 py-2 text-sm outline-none focus-visible:border-[#5B4BE0]">
          <option value="">Все статусы</option>
          <option value="OPEN">Открыта</option>
          <option value="ASSIGNED">В работе</option>
          <option value="COMPLETED">Выполнена</option>
          <option value="CANCELED">Отменена</option>
        </select>
      </div>

      <div v-if="loading" class="py-16 text-center text-sm text-[#8D8A7E]">Загрузка…</div>
      <div v-else-if="error" class="mt-6 rounded-2xl border border-[#F3D3CE] bg-[#FBF0EE] px-6 py-10 text-center text-sm text-[#B3261E]">{{ error }}</div>
      <div v-else-if="items.length === 0" class="mt-6 rounded-2xl border border-dashed border-[#E2DED2] px-6 py-16 text-center text-sm text-[#6E6B60]">
        Заявок нет
      </div>
      <div v-else class="mt-6 flex flex-col gap-3">
        <div v-for="r in items" :key="r.id" class="flex items-center justify-between rounded-2xl border border-[#E2DED2] bg-white p-5">
          <div>
            <p class="text-sm font-medium">№{{ r.id }} — {{ r.serviceName || 'Услуга' }}</p>
            <p class="mt-0.5 text-xs text-[#9B978A]">{{ r.addressText }} · {{ formatDate(r.createdAt) }} · клиент #{{ r.clientId }}<span v-if="r.masterId"> · мастер #{{ r.masterId }}</span></p>
          </div>
          <span class="rounded-[3px] px-2.5 py-1 text-xs font-medium" style="background:#EFEBE1">{{ requestStatusLabels[r.status] }}</span>
        </div>

        <div v-if="totalPages > 1" class="flex items-center justify-center gap-3 border-t border-[#E2DED2] px-1 py-3 text-sm">
          <button type="button" class="rounded-md px-2.5 py-1 text-[#6E6B60] hover:bg-[#F1EDE3] disabled:opacity-40" :disabled="page <= 1" @click="page--">Назад</button>
          <span class="font-mono">{{ page }} / {{ totalPages }}</span>
          <button type="button" class="rounded-md px-2.5 py-1 text-[#6E6B60] hover:bg-[#F1EDE3] disabled:opacity-40" :disabled="page >= totalPages" @click="page++">Далее</button>
        </div>
      </div>
    </div>
  </MarketplaceShell>
</template>
