<script setup lang="ts">
import { ref } from 'vue'
import { masterApi, adminApi } from '@/api/marketplace'
import { extractErrorMessage } from '@/api/client'
import { useToastStore } from '@/stores/toast'
import type { ReviewDTO } from '@/types'
import MarketplaceShell from '@/components/layout/MarketplaceShell.vue'

// No "list all reviews" endpoint exists on the backend (only per-master
// public listing + hide-by-id) -- admin moderation works by looking up a
// specific master's reviews, same data the public page shows.
const toast = useToastStore()
const masterIdInput = ref('')
const items = ref<ReviewDTO[]>([])
const loading = ref(false)
const error = ref('')
const searched = ref(false)

async function search() {
  const id = Number(masterIdInput.value)
  if (!id) return
  loading.value = true
  error.value = ''
  searched.value = true
  try {
    const result = await masterApi.listReviews(id, 1, 50)
    items.value = result.items
  } catch (e) {
    error.value = extractErrorMessage(e)
  } finally {
    loading.value = false
  }
}

async function hide(review: ReviewDTO) {
  try {
    await adminApi.hideReview(review.id)
    items.value = items.value.filter((r) => r.id !== review.id)
    toast.success('Отзыв скрыт')
  } catch (e) {
    toast.error(extractErrorMessage(e))
  }
}
</script>

<template>
  <MarketplaceShell>
    <div class="mx-auto max-w-[820px]">
      <h1 class="mk-display text-2xl font-bold tracking-tight">Модерация отзывов</h1>
      <p class="mt-1 text-sm text-[#6E6B60]">Введите ID мастера, чтобы посмотреть его отзывы.</p>

      <form class="mt-4 flex gap-2" @submit.prevent="search">
        <input v-model="masterIdInput" type="number" placeholder="ID мастера" class="w-40 rounded-xl border border-[#E2DED2] px-4 py-2.5 text-base outline-none focus-visible:border-[#5B4BE0]" />
        <button type="submit" class="rounded-[11px] bg-[#5B4BE0] px-4 py-2.5 text-sm font-medium text-white hover:opacity-90">Найти</button>
      </form>

      <div v-if="loading" class="py-16 text-center text-sm text-[#8D8A7E]">Загрузка…</div>
      <div v-else-if="error" class="mt-6 rounded-2xl border border-[#F3D3CE] bg-[#FBF0EE] px-6 py-10 text-center text-sm text-[#B3261E]">{{ error }}</div>
      <div v-else-if="searched && items.length === 0" class="mt-6 rounded-2xl border border-dashed border-[#E2DED2] px-6 py-16 text-center text-sm text-[#6E6B60]">
        Отзывов нет
      </div>
      <div v-else class="mt-6 flex flex-col gap-3">
        <div v-for="r in items" :key="r.id" class="flex items-center justify-between rounded-2xl border border-[#E2DED2] bg-white p-5">
          <div>
            <p class="text-sm font-medium">★ {{ r.rating }} — заявка №{{ r.requestId }}</p>
            <p v-if="r.comment" class="mt-0.5 text-sm text-[#6E6B60]">{{ r.comment }}</p>
          </div>
          <button type="button" class="rounded-[11px] border border-[#B3261E] px-3 py-2 text-sm text-[#B3261E] hover:bg-[#FBF0EE]" @click="hide(r)">
            Скрыть
          </button>
        </div>
      </div>
    </div>
  </MarketplaceShell>
</template>
