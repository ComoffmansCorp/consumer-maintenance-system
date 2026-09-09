<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { catalogApi } from '@/api/marketplace'
import { useAuthStore } from '@/stores/auth'
import type { CategoryDTO, ServiceDTO } from '@/types'
import MarketplaceShell from '@/components/layout/MarketplaceShell.vue'
import Skeleton from '@/components/Skeleton.vue'
import CategoryIcon from '@/components/marketplace/CategoryIcon.vue'

const auth = useAuthStore()
const router = useRouter()
const route = useRoute()

const loading = ref(true)
const categories = ref<CategoryDTO[]>([])
const services = ref<ServiceDTO[]>([])
const selectedCategoryId = ref<number | null>(
  route.query.categoryId ? Number(route.query.categoryId) : null,
)
const query = ref((route.query.q as string) ?? '')

async function load() {
  loading.value = true
  try {
    const [categoriesResult, servicesResult] = await Promise.all([
      catalogApi.listCategories(),
      catalogApi.listServices(),
    ])
    categories.value = categoriesResult
    services.value = servicesResult
  } catch {
    categories.value = []
    services.value = []
  } finally {
    loading.value = false
  }
}
onMounted(load)

function servicesInCategory(categoryId: number) {
  return services.value.filter((s) => s.categoryId === categoryId)
}

function categoryExamples(categoryId: number) {
  const names = servicesInCategory(categoryId).map((s) => s.name)
  return names.length ? names.slice(0, 3).join(', ') : 'Скоро появятся услуги'
}

function serviceWord(n: number) {
  if (n % 100 >= 11 && n % 100 <= 19) return 'услуг'
  if (n % 10 === 1) return 'услуга'
  if (n % 10 >= 2 && n % 10 <= 4) return 'услуги'
  return 'услуг'
}

const visibleServices = computed(() => {
  let items = services.value
  if (selectedCategoryId.value) items = items.filter((s) => s.categoryId === selectedCategoryId.value)
  const q = query.value.trim().toLowerCase()
  if (q) items = items.filter((s) => s.name.toLowerCase().includes(q) || s.description?.toLowerCase().includes(q))
  return items
})

function pickCategory(id: number) {
  selectedCategoryId.value = selectedCategoryId.value === id ? null : id
}

function startRequest(serviceId?: number) {
  if (!auth.isAuthenticated || auth.role !== 'CLIENT') {
    router.push({ name: 'marketplace-register', query: serviceId ? { serviceId } : {} })
    return
  }
  router.push({ name: 'marketplace-new-request', query: serviceId ? { serviceId } : {} })
}
</script>

<template>
  <MarketplaceShell>
    <div class="mx-auto max-w-[1280px] px-5 py-10 md:px-10">
      <h1 class="mk-display text-[32px] font-bold tracking-tight md:text-[40px]">Каталог услуг</h1>
      <p class="mt-2 text-[15px] text-[#6E6B60]">
        Выберите категорию или найдите услугу поиском — заявка сразу уйдёт мастерам с подходящей специализацией.
      </p>

      <div class="mt-7 flex items-center gap-2.5 rounded-2xl border border-[#E2DED2] bg-white px-4 py-1">
        <span class="font-mono text-xs text-[#A8A498]">Q</span>
        <input
          v-model="query"
          type="text"
          placeholder="Какая услуга нужна? Например, замена счётчика"
          class="w-full border-0 bg-transparent py-3.5 text-base outline-none placeholder:text-[#9B978A]"
        />
      </div>

      <!-- Categories -->
      <div class="mt-9">
        <div class="mb-4 flex items-end justify-between gap-5">
          <h2 class="mk-display text-xl font-semibold tracking-tight">Категории</h2>
          <span class="font-mono text-sm text-[#A8A498]">{{ services.length }} {{ serviceWord(services.length) }}</span>
        </div>

        <div v-if="loading" class="grid grid-cols-2 gap-4 sm:grid-cols-3 md:grid-cols-4">
          <Skeleton v-for="i in 8" :key="i" class="h-[168px]" rounded="rounded-[18px]" />
        </div>
        <p v-else-if="categories.length === 0" class="text-sm text-[#6E6B60]">Каталог пока пуст — загляните позже.</p>

        <div v-else class="grid grid-cols-2 gap-4 sm:grid-cols-3 md:grid-cols-4">
          <button
            v-for="c in categories"
            :key="c.id"
            type="button"
            class="flex min-h-[168px] flex-col gap-8 rounded-[18px] border bg-white p-5 text-left transition-colors"
            :class="selectedCategoryId === c.id ? 'border-[#17160F]' : 'border-[#E2DED2] hover:border-[#17160F]'"
            @click="pickCategory(c.id)"
          >
            <div class="flex items-start justify-between gap-3">
              <span class="grid h-10 w-10 shrink-0 place-items-center rounded-full bg-[#EFE9FF] text-[#5B4BE0]">
                <CategoryIcon :name="c.name" />
              </span>
              <span class="whitespace-nowrap font-mono text-[11px] text-[#A8A498]">{{ servicesInCategory(c.id).length }}</span>
            </div>
            <div>
              <div class="text-[17px] font-semibold leading-tight">{{ c.name }}</div>
              <div class="mt-1 text-sm leading-snug text-[#6E6B60]">{{ categoryExamples(c.id) }}</div>
            </div>
          </button>
        </div>
      </div>

      <!-- Services -->
      <div class="mt-12">
        <div class="mb-5 flex items-center justify-between">
          <h2 class="mk-display text-xl font-semibold tracking-tight">
            {{ selectedCategoryId ? categories.find((c) => c.id === selectedCategoryId)?.name : 'Все услуги' }}
          </h2>
          <button v-if="selectedCategoryId" type="button" class="text-sm font-medium text-[#5B4BE0] hover:underline" @click="selectedCategoryId = null">
            Сбросить фильтр
          </button>
        </div>

        <div v-if="loading" class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          <Skeleton v-for="i in 6" :key="i" class="h-[172px]" rounded="rounded-[18px]" />
        </div>
        <p v-else-if="visibleServices.length === 0" class="text-sm text-[#6E6B60]">
          Ничего не нашлось — попробуйте изменить запрос или выбрать другую категорию.
        </p>

        <div v-else class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          <button
            v-for="s in visibleServices"
            :key="s.id"
            type="button"
            class="flex flex-col items-start gap-3 overflow-hidden rounded-[18px] border border-l-[3px] border-[#E2DED2] border-l-[#5B4BE0] bg-white text-left transition-colors hover:border-[#17160F]"
            @click="startRequest(s.id)"
          >
            <img v-if="s.imageUrl" :src="s.imageUrl" :alt="s.name" class="h-32 w-full object-cover" loading="lazy" />
            <div class="flex flex-1 flex-col items-start gap-2 px-5 pb-5" :class="s.imageUrl ? 'pt-4' : 'pt-5'">
              <span class="text-base font-semibold">{{ s.name }}</span>
              <span v-if="s.description" class="text-sm text-[#6E6B60]">{{ s.description }}</span>
              <span class="mt-auto pt-2 text-sm font-medium text-[#5B4BE0]">Оставить заявку →</span>
            </div>
          </button>
        </div>

        <p class="mt-8 text-center text-xs text-[#9B978A]">
          Заявка автоматически попадает к мастерам с подходящей специализацией — случайный исполнитель не сможет её взять.
        </p>
      </div>
    </div>
  </MarketplaceShell>
</template>
