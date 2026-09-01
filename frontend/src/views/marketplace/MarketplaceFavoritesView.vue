<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { requestsApi } from '@/api/marketplace'
import { extractErrorMessage } from '@/api/client'
import { useToastStore } from '@/stores/toast'
import type { FavoriteDTO } from '@/types'
import MarketplaceShell from '@/components/layout/MarketplaceShell.vue'

const toast = useToastStore()
const items = ref<FavoriteDTO[]>([])
const loading = ref(true)
const error = ref('')

async function load() {
  loading.value = true
  error.value = ''
  try {
    items.value = await requestsApi.listFavorites()
  } catch (e) {
    error.value = extractErrorMessage(e)
  } finally {
    loading.value = false
  }
}
onMounted(load)

async function remove(masterId: number) {
  try {
    await requestsApi.removeFavorite(masterId)
    items.value = items.value.filter((f) => f.masterId !== masterId)
    toast.success('Удалено из избранного')
  } catch (e) {
    toast.error(extractErrorMessage(e))
  }
}
</script>

<template>
  <MarketplaceShell>
    <div class="mx-auto max-w-[820px]">
      <h1 class="mk-display text-2xl font-bold tracking-tight">Избранные мастера</h1>

      <div v-if="loading" class="py-16 text-center text-sm text-[#8D8A7E]">Загрузка…</div>
      <div v-else-if="error" class="mt-6 rounded-2xl border border-[#F3D3CE] bg-[#FBF0EE] px-6 py-10 text-center">
        <p class="text-sm font-medium text-[#B3261E]">{{ error }}</p>
      </div>
      <div v-else-if="items.length === 0" class="mt-6 rounded-2xl border border-dashed border-[#E2DED2] px-6 py-16 text-center">
        <p class="text-sm text-[#6E6B60]">Пока никого не добавили в избранное.</p>
      </div>
      <ul v-else class="mt-6 flex flex-col gap-3">
        <li
          v-for="f in items"
          :key="f.masterId"
          class="flex items-center justify-between rounded-2xl border border-[#E2DED2] bg-white p-5"
        >
          <span class="text-sm font-medium">Мастер #{{ f.masterId }}</span>
          <button
            type="button"
            class="rounded-[11px] border border-[#B3261E] px-3 py-2 text-sm text-[#B3261E] hover:bg-[#FBF0EE]"
            @click="remove(f.masterId)"
          >
            Убрать
          </button>
        </li>
      </ul>
    </div>
  </MarketplaceShell>
</template>
