<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import { adminApi } from '@/api/marketplace'
import { extractErrorMessage } from '@/api/client'
import type { ProfileDTO } from '@/types'
import MarketplaceShell from '@/components/layout/MarketplaceShell.vue'

const items = ref<ProfileDTO[]>([])
const loading = ref(true)
const error = ref('')
const page = ref(1)
const totalPages = ref(1)

async function load() {
  loading.value = true
  error.value = ''
  try {
    const result = await adminApi.listMasters(page.value, 20)
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
</script>

<template>
  <MarketplaceShell>
    <div class="mx-auto max-w-[820px]">
      <h1 class="mk-display text-2xl font-bold tracking-tight">Мастера</h1>

      <div v-if="loading" class="py-16 text-center text-sm text-[#8D8A7E]">Загрузка…</div>
      <div v-else-if="error" class="mt-6 rounded-2xl border border-[#F3D3CE] bg-[#FBF0EE] px-6 py-10 text-center text-sm text-[#B3261E]">{{ error }}</div>
      <div v-else-if="items.length === 0" class="mt-6 rounded-2xl border border-dashed border-[#E2DED2] px-6 py-16 text-center text-sm text-[#6E6B60]">
        Мастеров пока нет
      </div>
      <div v-else class="mt-6 flex flex-col gap-3">
        <div v-for="m in items" :key="m.userId" class="flex items-center justify-between rounded-2xl border border-[#E2DED2] bg-white p-5">
          <div class="flex items-center gap-3">
            <img
              v-if="m.avatarUrl"
              :src="m.avatarUrl"
              alt=""
              class="h-11 w-11 shrink-0 rounded-full object-cover"
            />
            <div v-else class="h-11 w-11 shrink-0 rounded-full bg-[#EFEBE1]" />
            <div>
              <p class="text-sm font-medium">Мастер #{{ m.userId }}<span v-if="m.city"> · {{ m.city }}</span></p>
              <p v-if="m.bio" class="mt-0.5 text-sm text-[#6E6B60]">{{ m.bio }}</p>
              <p class="mt-1 text-xs text-[#9B978A]">Специализаций: {{ m.specializationIds.length }}</p>
            </div>
          </div>
          <div class="text-right">
            <p class="text-sm font-medium">★ {{ m.ratingAvg.toFixed(1) }}</p>
            <p class="text-xs text-[#9B978A]">{{ m.ratingCount }} отзывов</p>
          </div>
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
