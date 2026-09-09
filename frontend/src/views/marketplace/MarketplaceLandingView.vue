<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { catalogApi } from '@/api/marketplace'
import type { CategoryDTO, ServiceDTO } from '@/types'
import MarketplaceShell from '@/components/layout/MarketplaceShell.vue'

const router = useRouter()

// Plain string, not a static `src="..."` attribute -- Vue's template
// compiler rewrites static asset src attributes into build-time module
// imports, which fails for paths served at runtime by the gateway/MinIO
// (see docker-compose.yml minio/minio-init) rather than bundled by Vite.
const heroMasterImg = '/media/hero/master-at-work.jpg'

const loading = ref(true)
const categories = ref<CategoryDTO[]>([])
const services = ref<ServiceDTO[]>([])
const query = ref('')

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
    // Public landing page: fail soft into the empty-catalog state rather
    // than showing a visitor a scary error banner on the first thing they see.
    categories.value = []
    services.value = []
  } finally {
    loading.value = false
  }
}
onMounted(load)

function serviceWord(n: number) {
  if (n % 100 >= 11 && n % 100 <= 19) return 'услуг'
  if (n % 10 === 1) return 'услуга'
  if (n % 10 >= 2 && n % 10 <= 4) return 'услуги'
  return 'услуг'
}
function categoryWord(n: number) {
  if (n % 100 >= 11 && n % 100 <= 19) return 'категорий'
  if (n % 10 === 1) return 'категория'
  if (n % 10 >= 2 && n % 10 <= 4) return 'категории'
  return 'категорий'
}

const suggestions = computed(() => services.value.slice(0, 5).map((s) => s.name))

function goToCatalog(params: { q?: string; categoryId?: number } = {}) {
  router.push({ name: 'marketplace-catalog', query: params })
}

function runSearch() {
  goToCatalog(query.value.trim() ? { q: query.value.trim() } : {})
}

function goRegister() {
  router.push({ name: 'marketplace-register' })
}
</script>

<template>
  <MarketplaceShell>
    <!-- Hero -->
    <section class="mx-auto grid max-w-[1280px] gap-12 px-5 pb-14 pt-16 md:grid-cols-[1.15fr_0.85fr] md:gap-16 md:px-10 md:pb-14 md:pt-[72px]">
      <div>
        <div class="mb-6 inline-flex items-center gap-2.5 rounded-full border border-[#E2DED2] bg-white py-1.5 pl-2.5 pr-3.5 text-[13px] text-[#55524A]">
          <span class="font-mono text-[11px] text-[#5B4BE0]">{{ categories.length }}</span>
          {{ categories.length === 1 ? 'категория услуг' : 'категории услуг' }} на платформе
        </div>
        <h1 class="mk-display mb-5 text-[42px] font-bold leading-[1.02] tracking-tight sm:text-[56px] md:text-[68px]">
          Найдите мастера<br />за пять минут
        </h1>
        <p class="mb-8 max-w-[480px] text-[17px] leading-[1.55] text-[#55524A] md:text-[19px]">
          Опишите задачу — она сразу становится доступна мастерам с нужной специализацией.
          Кто первый откликнется, тот и берёт заявку в работу.
        </p>

        <form
          class="flex flex-col gap-2 rounded-2xl border border-[#E2DED2] bg-white p-2.5 shadow-[0_12px_34px_-22px_rgba(23,22,15,0.35)] sm:flex-row sm:items-center"
          @submit.prevent="runSearch"
        >
          <div class="flex flex-1 items-center gap-2.5 px-3.5 py-1">
            <span class="font-mono text-xs text-[#A8A498]">Q</span>
            <input
              v-model="query"
              type="text"
              placeholder="Какая услуга нужна? Например, замена счётчика"
              class="w-full border-0 bg-transparent py-3.5 text-base outline-none placeholder:text-[#9B978A]"
            />
          </div>
          <button
            type="submit"
            class="rounded-[14px] bg-[#5B4BE0] px-7 py-4 text-base font-medium text-white hover:bg-[#4536BC]"
          >
            Найти
          </button>
        </form>

        <div class="mt-4 flex flex-wrap items-center gap-2">
          <span class="px-0.5 py-2 text-sm text-[#8D8A7E]">Часто ищут:</span>
          <button
            v-for="s in suggestions"
            :key="s"
            type="button"
            class="rounded-full border border-[#E2DED2] bg-white px-4 py-2 text-sm text-[#4C4A40] hover:border-[#5B4BE0] hover:text-[#5B4BE0]"
            @click="goToCatalog({ q: s })"
          >
            {{ s }}
          </button>
        </div>
      </div>

      <div class="grid gap-4">
        <div class="aspect-[4/3.4] overflow-hidden rounded-[22px] border border-[#E2DED2]">
          <img :src="heroMasterImg" alt="Мастер за работой" class="h-full w-full object-cover" />
        </div>
        <div class="grid grid-cols-2 gap-4">
          <button
            type="button"
            class="rounded-2xl bg-[#17160F] p-5 text-left text-[#F5F2EA] hover:bg-[#26241A]"
            @click="goToCatalog()"
          >
            <div class="mk-display text-[28px] font-medium tracking-tight sm:text-[30px]">{{ categories.length }}</div>
            <div class="mt-1.5 text-sm leading-tight text-[#F5F2EA]/65">{{ categoryWord(categories.length) }} услуг в каталоге</div>
          </button>
          <button
            type="button"
            class="rounded-2xl border border-[#E2DED2] bg-white p-5 text-left hover:border-[#17160F]"
            @click="goToCatalog()"
          >
            <div class="mk-display text-[28px] font-medium tracking-tight text-[#5B4BE0] sm:text-[30px]">{{ services.length }}</div>
            <div class="mt-1.5 text-sm leading-tight text-[#6E6B60]">{{ serviceWord(services.length) }} доступно мастерам</div>
          </button>
        </div>
      </div>
    </section>

    <!-- Popular categories teaser -- full catalog lives on its own page -->
    <section class="mx-auto max-w-[1280px] px-5 pb-16 pt-8 md:px-10">
      <div class="mb-6 flex items-end justify-between gap-5">
        <h2 class="mk-display text-[28px] font-bold tracking-tight md:text-[36px]">Категории услуг</h2>
        <RouterLink :to="{ name: 'marketplace-catalog' }" class="text-sm font-medium text-[#5B4BE0] hover:underline">
          Весь каталог →
        </RouterLink>
      </div>

      <p v-if="!loading && categories.length === 0" class="text-sm text-[#6E6B60]">Каталог пока пуст — загляните позже.</p>

      <div v-else class="flex flex-wrap gap-3">
        <button
          v-for="c in categories"
          :key="c.id"
          type="button"
          class="rounded-full border border-[#E2DED2] bg-white px-5 py-3 text-[15px] font-medium hover:border-[#17160F]"
          @click="goToCatalog({ categoryId: c.id })"
        >
          {{ c.name }}
        </button>
      </div>
    </section>

    <!-- How it works -->
    <section class="bg-[#17160F] text-[#F5F2EA]">
      <div class="mx-auto max-w-[1280px] px-5 py-16 md:px-10 md:py-[84px]">
        <h2 class="mk-display mb-2 text-[28px] font-bold tracking-tight md:text-[36px]">Как это работает</h2>
        <p class="mb-10 text-[17px] text-[#F5F2EA]/60">Три шага — без торгов и ожидания предложений.</p>
        <div class="grid gap-7 sm:grid-cols-3">
          <div class="border-t border-[#F5F2EA]/18 pt-5">
            <div class="mb-4 font-mono text-xs" style="color: oklch(0.78 0.14 75)">01</div>
            <h3 class="mb-2.5 text-[21px] font-semibold">Опишите задачу</h3>
            <p class="text-base leading-[1.55] text-[#F5F2EA]/65">Выберите услугу, укажите адрес и что нужно сделать. Без звонков и объявлений.</p>
          </div>
          <div class="border-t border-[#F5F2EA]/18 pt-5">
            <div class="mb-4 font-mono text-xs" style="color: oklch(0.78 0.14 75)">02</div>
            <h3 class="mb-2.5 text-[21px] font-semibold">Заявка уходит в пул</h3>
            <p class="text-base leading-[1.55] text-[#F5F2EA]/65">Её видят только мастера с подходящей специализацией. Кто откликнулся первым — берёт в работу.</p>
          </div>
          <div class="border-t border-[#F5F2EA]/18 pt-5">
            <div class="mb-4 font-mono text-xs" style="color: oklch(0.78 0.14 75)">03</div>
            <h3 class="mb-2.5 text-[21px] font-semibold">Мастер выполняет и закрывает</h3>
            <p class="text-base leading-[1.55] text-[#F5F2EA]/65">Статус меняется на глазах у клиента — открыта → в работе → выполнена, прямо в личном кабинете.</p>
          </div>
        </div>
        <div class="mt-11 flex flex-wrap gap-3.5">
          <button type="button" class="rounded-[13px] bg-[#F5F2EA] px-6.5 py-4 text-base font-medium text-[#17160F]" @click="goRegister">
            Оставить заявку
          </button>
          <RouterLink :to="{ name: 'marketplace-about' }" class="rounded-[13px] border border-[#F5F2EA]/30 px-6.5 py-4 text-base text-[#F5F2EA]">
            Как проверяется специализация
          </RouterLink>
        </div>
      </div>
    </section>

    <!-- For masters teaser -- full page at /for-masters -->
    <section class="mx-auto my-16 max-w-[1280px] px-5 md:my-[84px] md:px-10">
      <div class="flex flex-col items-start gap-6 rounded-[26px] bg-[#5B4BE0] p-8 text-white sm:flex-row sm:items-center sm:justify-between md:p-14">
        <div>
          <div class="mb-3 font-mono text-[11px] uppercase tracking-[0.14em] text-white/70">Для специалистов</div>
          <h2 class="mk-display text-[24px] font-bold leading-[1.15] tracking-tight md:text-[32px]">
            Берите заявки по своей специализации — без обзвонов
          </h2>
        </div>
        <RouterLink
          :to="{ name: 'marketplace-for-masters' }"
          class="whitespace-nowrap rounded-[13px] bg-white px-6.5 py-4 text-base font-medium text-[#4536BC]"
        >
          Узнать больше →
        </RouterLink>
      </div>
    </section>
  </MarketplaceShell>
</template>
