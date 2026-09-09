<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { catalogApi, masterApi } from '@/api/marketplace'
import { extractErrorMessage } from '@/api/client'
import type { ProfileDTO, ReviewDTO, ServiceDTO } from '@/types'
import MarketplaceShell from '@/components/layout/MarketplaceShell.vue'
import StarRating from '@/components/StarRating.vue'
import Skeleton from '@/components/Skeleton.vue'

const route = useRoute()
const masterId = computed(() => Number(route.params.id))

const loading = ref(true)
const error = ref('')
const profile = ref<ProfileDTO | null>(null)
const services = ref<ServiceDTO[]>([])
const reviews = ref<ReviewDTO[]>([])

const specializationNames = computed(() => {
  const byId = new Map(services.value.map((s) => [s.id, s.name]))
  return (profile.value?.specializationIds ?? []).map((id) => byId.get(id) ?? `#${id}`)
})

async function load() {
  loading.value = true
  error.value = ''
  try {
    const [p, s, r] = await Promise.all([
      masterApi.getPublicProfile(masterId.value),
      catalogApi.listServices(),
      masterApi.listReviews(masterId.value, 1, 50),
    ])
    profile.value = p
    services.value = s
    reviews.value = r.items
  } catch (e) {
    error.value = extractErrorMessage(e)
  } finally {
    loading.value = false
  }
}
onMounted(load)
</script>

<template>
  <MarketplaceShell>
    <div class="mx-auto max-w-[820px] px-5 py-10 md:px-0">
      <div v-if="loading" class="flex items-center gap-4 rounded-2xl border border-[#E2DED2] bg-white p-6">
        <Skeleton class="h-20 w-20" rounded="rounded-full" />
        <div class="flex-1">
          <Skeleton class="h-5 w-1/2" />
          <Skeleton class="mt-2 h-4 w-1/3" />
        </div>
      </div>
      <div v-else-if="error" class="rounded-2xl border border-[#F3D3CE] bg-[#FBF0EE] px-6 py-10 text-center">
        <p class="text-sm font-medium text-[#B3261E]">{{ error }}</p>
      </div>

      <template v-else-if="profile">
        <div class="flex flex-col items-start gap-4 rounded-2xl border border-[#E2DED2] bg-white p-6 sm:flex-row sm:items-center">
          <img v-if="profile.avatarUrl" :src="profile.avatarUrl" alt="" class="h-20 w-20 rounded-full object-cover" />
          <div v-else class="grid h-20 w-20 place-items-center rounded-full bg-[#EFE9FF] text-2xl font-semibold text-[#5B4BE0]">
            {{ (profile.fullName ?? '?')[0] }}
          </div>
          <div>
            <h1 class="mk-display text-2xl font-bold tracking-tight">{{ profile.fullName || `Мастер #${profile.userId}` }}</h1>
            <p class="mt-1 text-sm text-[#8D8A7E]">{{ profile.city || 'Город не указан' }}</p>
            <StarRating :value="profile.ratingAvg" :count="profile.ratingCount" class="mt-2" />
          </div>
        </div>

        <p v-if="profile.bio" class="mt-6 text-[15px] leading-[1.6] text-[#4C4A40]">{{ profile.bio }}</p>

        <div class="mt-6 flex flex-wrap gap-2">
          <span v-for="name in specializationNames" :key="name" class="rounded-full bg-[#EFE9FF] px-3 py-1.5 text-xs text-[#5B4BE0]">
            {{ name }}
          </span>
        </div>

        <h2 class="mk-display mt-10 text-xl font-semibold tracking-tight">Отзывы</h2>
        <p v-if="reviews.length === 0" class="mt-3 text-sm text-[#6E6B60]">Пока нет отзывов.</p>
        <ul v-else class="mt-4 flex flex-col gap-3">
          <li v-for="rv in reviews" :key="rv.id" class="rounded-2xl border border-[#E2DED2] bg-white p-5">
            <StarRating :value="rv.rating" size="sm" />
            <p v-if="rv.comment" class="mt-2 text-sm leading-[1.55] text-[#4C4A40]">{{ rv.comment }}</p>
          </li>
        </ul>
      </template>
    </div>
  </MarketplaceShell>
</template>
