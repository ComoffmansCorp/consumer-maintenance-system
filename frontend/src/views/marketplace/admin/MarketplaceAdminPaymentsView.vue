<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import { adminApi } from '@/api/marketplace'
import { extractErrorMessage } from '@/api/client'
import type { PaymentDTO } from '@/types'
import MarketplaceShell from '@/components/layout/MarketplaceShell.vue'

const paymentStatusLabels: Record<string, string> = {
  HELD: 'Удержана (эскроу)',
  RELEASED: 'Выплачена',
  REFUNDED: 'Возвращена',
}
const paymentStatusStyles: Record<string, { bg: string; fg: string }> = {
  HELD: { bg: '#EFEBE1', fg: '#55524A' },
  RELEASED: { bg: '#DCEFE1', fg: '#2E7D4F' },
  REFUNDED: { bg: '#FBE9E7', fg: '#B3261E' },
}

const items = ref<PaymentDTO[]>([])
const loading = ref(true)
const error = ref('')
const page = ref(1)
const totalPages = ref(1)

async function load() {
  loading.value = true
  error.value = ''
  try {
    const result = await adminApi.listPayments(page.value, 20)
    items.value = result.items
    totalPages.value = result.totalPages
  } catch (e) {
    error.value = extractErrorMessage(e)
  } finally {
    loading.value = false
  }
}
onMounted(load)
watch(page, load)

function formatDate(iso: string) {
  return new Date(iso).toLocaleDateString('ru-RU')
}
</script>

<template>
  <MarketplaceShell>
    <div class="mx-auto max-w-[820px]">
      <h1 class="mk-display text-2xl font-bold tracking-tight">Платежи</h1>

      <div v-if="loading" class="py-16 text-center text-sm text-[#8D8A7E]">Загрузка…</div>
      <div v-else-if="error" class="mt-6 rounded-2xl border border-[#F3D3CE] bg-[#FBF0EE] px-6 py-10 text-center text-sm text-[#B3261E]">{{ error }}</div>
      <div v-else-if="items.length === 0" class="mt-6 rounded-2xl border border-dashed border-[#E2DED2] px-6 py-16 text-center text-sm text-[#6E6B60]">
        Платежей пока нет
      </div>
      <div v-else class="mt-6 flex flex-col gap-3">
        <div v-for="p in items" :key="p.id" class="flex items-center justify-between rounded-2xl border border-[#E2DED2] bg-white p-5">
          <div>
            <p class="text-sm font-medium">Заявка №{{ p.requestId }} — {{ p.amount }} ₽</p>
            <p class="mt-0.5 text-xs text-[#9B978A]">Комиссия {{ p.platformFee }} ₽ · {{ formatDate(p.createdAt) }}</p>
          </div>
          <span
            class="rounded-[3px] px-2.5 py-1 text-xs font-medium"
            :style="{ background: paymentStatusStyles[p.status].bg, color: paymentStatusStyles[p.status].fg }"
          >
            {{ paymentStatusLabels[p.status] }}
          </span>
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
