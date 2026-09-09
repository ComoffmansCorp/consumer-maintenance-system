<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { masterApi } from '@/api/marketplace'
import { extractErrorMessage } from '@/api/client'
import type { ProfileDTO } from '@/types'
import MarketplaceShell from '@/components/layout/MarketplaceShell.vue'
import StarRating from '@/components/StarRating.vue'
import Skeleton from '@/components/Skeleton.vue'

const loading = ref(true)
const error = ref('')
const masters = ref<ProfileDTO[]>([])
const cityFilter = ref('')

async function load() {
  loading.value = true
  error.value = ''
  try {
    const page = await masterApi.listPublic(1, 50)
    masters.value = page.items
  } catch (e) {
    error.value = extractErrorMessage(e)
  } finally {
    loading.value = false
  }
}
onMounted(load)

const cities = computed(() => [...new Set(masters.value.map((m) => m.city).filter(Boolean))] as string[])
const visibleMasters = computed(() =>
  cityFilter.value ? masters.value.filter((m) => m.city === cityFilter.value) : masters.value,
)
</script>

<template>
  <MarketplaceShell>
    <div class="mx-auto max-w-[1280px] px-5 py-10 md:px-10">
      <h1 class="mk-display text-[32px] font-bold tracking-tight md:text-[40px]">Мастера</h1>
      <p class="mt-2 text-[15px] text-[#6E6B60]">Отсортировано по рейтингу — сначала лучшие исполнители.</p>

      <div v-if="!loading && cities.length > 1" class="mt-6 flex flex-wrap gap-2">
        <button
          type="button"
          class="rounded-full border px-4 py-2 text-sm"
          :class="!cityFilter ? 'border-[#17160F] bg-[#17160F] text-white' : 'border-[#E2DED2] hover:border-[#17160F]'"
          @click="cityFilter = ''"
        >
          Все города
        </button>
        <button
          v-for="c in cities"
          :key="c"
          type="button"
          class="rounded-full border px-4 py-2 text-sm"
          :class="cityFilter === c ? 'border-[#17160F] bg-[#17160F] text-white' : 'border-[#E2DED2] hover:border-[#17160F]'"
          @click="cityFilter = c"
        >
          {{ c }}
        </button>
      </div>

      <div v-if="loading" class="mt-8 grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        <div v-for="i in 6" :key="i" class="flex items-center gap-4 rounded-2xl border border-[#E2DED2] bg-white p-5">
          <Skeleton class="h-14 w-14 shrink-0" rounded="rounded-full" />
          <div class="flex-1">
            <Skeleton class="h-4 w-3/4" />
            <Skeleton class="mt-2 h-3 w-1/2" />
          </div>
        </div>
      </div>
      <div v-else-if="error" class="mt-8 rounded-2xl border border-[#F3D3CE] bg-[#FBF0EE] px-6 py-10 text-center">
        <p class="text-sm font-medium text-[#B3261E]">{{ error }}</p>
      </div>
      <div v-else-if="visibleMasters.length === 0" class="mt-8 rounded-2xl border border-dashed border-[#E2DED2] px-6 py-16 text-center">
        <p class="text-sm text-[#6E6B60]">Мастеров пока нет.</p>
      </div>
      <div v-else class="mt-8 grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        <RouterLink
          v-for="m in visibleMasters"
          :key="m.userId"
          :to="{ name: 'marketplace-master-profile', params: { id: m.userId } }"
          class="flex items-center gap-4 rounded-2xl border border-[#E2DED2] bg-white p-5 transition-colors hover:border-[#17160F]"
        >
          <img v-if="m.avatarUrl" :src="m.avatarUrl" alt="" class="h-14 w-14 shrink-0 rounded-full object-cover" />
          <div v-else class="grid h-14 w-14 shrink-0 place-items-center rounded-full bg-[#EFE9FF] text-lg font-semibold text-[#5B4BE0]">
            {{ (m.fullName ?? '?')[0] }}
          </div>
          <div class="min-w-0">
            <div class="truncate font-semibold">{{ m.fullName || `Мастер #${m.userId}` }}</div>
            <div class="text-sm text-[#8D8A7E]">{{ m.city || 'Город не указан' }}</div>
            <StarRating :value="m.ratingAvg" :count="m.ratingCount" size="sm" class="mt-1" />
          </div>
        </RouterLink>
      </div>
    </div>
  </MarketplaceShell>
</template>
