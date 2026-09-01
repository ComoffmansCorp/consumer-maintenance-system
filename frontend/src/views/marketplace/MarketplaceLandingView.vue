<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { catalogApi } from '@/api/marketplace'
import { useAuthStore } from '@/stores/auth'
import type { CategoryDTO, ServiceDTO } from '@/types'
import MarketplaceShell from '@/components/layout/MarketplaceShell.vue'

const auth = useAuthStore()
const router = useRouter()

const loading = ref(true)
const categories = ref<CategoryDTO[]>([])
const services = ref<ServiceDTO[]>([])
const selectedCategoryId = ref<number | null>(null)
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
function categoryWord(n: number) {
  if (n % 100 >= 11 && n % 100 <= 19) return 'категорий'
  if (n % 10 === 1) return 'категория'
  if (n % 10 >= 2 && n % 10 <= 4) return 'категории'
  return 'категорий'
}

const visibleServices = computed(() => {
  let items = services.value
  if (selectedCategoryId.value) items = items.filter((s) => s.categoryId === selectedCategoryId.value)
  const q = query.value.trim().toLowerCase()
  if (q) items = items.filter((s) => s.name.toLowerCase().includes(q) || s.description?.toLowerCase().includes(q))
  return items
})

const suggestions = computed(() => services.value.slice(0, 5).map((s) => s.name))

function pickCategory(id: number) {
  selectedCategoryId.value = selectedCategoryId.value === id ? null : id
  scrollToServices()
}

function pickSuggestion(label: string) {
  query.value = label
  scrollToServices()
}

function runSearch() {
  scrollToServices()
}

function scrollToServices() {
  document.getElementById('services')?.scrollIntoView({ behavior: 'smooth', block: 'start' })
}

function startRequest(serviceId?: number) {
  if (!auth.isAuthenticated || auth.role !== 'CLIENT') {
    router.push({ name: 'marketplace-register', query: serviceId ? { serviceId } : {} })
    return
  }
  router.push({ name: 'marketplace-new-request', query: serviceId ? { serviceId } : {} })
}

function goRegister() {
  router.push({ name: 'marketplace-register' })
}

// Everything below describes what the backend actually does (see
// internal/marketplace/service.go): a hard, server-enforced specialization
// check on claim, first-come self-claim (no bidding/offers), no payments,
// no ratings, no identity verification. Copy is written to match that —
// no feature is promised here that the app doesn't have.
const platformFacts = [
  {
    title: 'Специализация — жёсткое ограничение',
    text: 'Проверяется на сервере при каждом взятии заявки: мастер не может выйти за рамки своей специализации, даже в обход интерфейса.',
    tint: '#E7E3FC',
  },
  {
    title: 'Self-claim, без очереди офферов',
    text: 'Кто из подходящих мастеров откликнулся первым — тот и берёт заявку в работу. Никаких торгов и ожидания предложений.',
    tint: 'oklch(0.94 0.04 145)',
  },
  {
    title: 'Статус в реальном времени',
    text: 'Клиент видит смену статуса в личном кабинете: открыта → в работе → выполнена (или отменена).',
    tint: 'oklch(0.94 0.04 75)',
  },
  {
    title: 'Без комиссии',
    text: 'Платформа пока не берёт долю со сделок — оставить и взять заявку бесплатно.',
    tint: 'oklch(0.94 0.04 20)',
  },
]

const proStats = [
  { label: 'Специализация', value: 'Проверяется на сервере' },
  { label: 'Модель распределения', value: 'Self-claim' },
  { label: 'Комиссия платформы', value: '0 ₽' },
]

const faqItems = [
  {
    q: 'Сколько стоит оставить заявку?',
    a: 'Ничего. Платформа сейчас не берёт комиссию ни с клиента, ни с мастера — оставить и взять заявку бесплатно.',
  },
  {
    q: 'Как проверяется специализация мастера?',
    a: 'Мастер указывает специализации в своём профиле. При попытке взять заявку сервер повторно сверяет категорию услуги со специализациями мастера — несовпадение отклоняется, даже если запрос пришёл напрямую в API, в обход приложения.',
  },
  {
    q: 'Что происходит после того, как мастер взял заявку?',
    a: 'Статус меняется на «В работе», клиент сразу видит это в личном кабинете. Когда мастер заканчивает, он закрывает заявку — статус становится «Выполнена».',
  },
  {
    q: 'Можно ли отменить заявку?',
    a: 'Да, с указанием причины — пока она не завершена. Отменить может и клиент, и мастер, которому она назначена.',
  },
]
const openFaq = ref(0)
function toggleFaq(i: number) {
  openFaq.value = openFaq.value === i ? -1 : i
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
            @click="pickSuggestion(s)"
          >
            {{ s }}
          </button>
        </div>
      </div>

      <div class="grid gap-4">
        <div
          class="relative grid aspect-[4/3.4] place-items-center overflow-hidden rounded-[22px] border border-[#E2DED2]"
          style="background: repeating-linear-gradient(135deg, #f1ede3 0 10px, #e9e4d8 10px 20px)"
        >
          <span class="rounded-lg bg-[#FBFAF7]/90 px-2.5 py-1.5 font-mono text-xs text-[#8D8A7E]">фото мастера за работой</span>
        </div>
        <div class="grid grid-cols-2 gap-4">
          <div class="rounded-2xl bg-[#17160F] p-5 text-[#F5F2EA]">
            <div class="mk-display text-[28px] font-medium tracking-tight sm:text-[30px]">{{ categories.length }}</div>
            <div class="mt-1.5 text-sm leading-tight text-[#F5F2EA]/65">{{ categoryWord(categories.length) }} услуг в каталоге</div>
          </div>
          <div class="rounded-2xl border border-[#E2DED2] bg-white p-5">
            <div class="mk-display text-[28px] font-medium tracking-tight text-[#5B4BE0] sm:text-[30px]">{{ services.length }}</div>
            <div class="mt-1.5 text-sm leading-tight text-[#6E6B60]">{{ serviceWord(services.length) }} доступно мастерам</div>
          </div>
        </div>
      </div>
    </section>

    <!-- Categories (real catalog) -->
    <section id="cats" class="mx-auto max-w-[1280px] px-5 pb-16 pt-8 md:px-10">
      <div class="mb-6 flex items-end justify-between gap-5">
        <h2 class="mk-display text-[28px] font-bold tracking-tight md:text-[36px]">Категории услуг</h2>
        <span class="font-mono text-sm text-[#A8A498]">{{ services.length }} {{ serviceWord(services.length) }}</span>
      </div>

      <p v-if="!loading && categories.length === 0" class="text-sm text-[#6E6B60]">
        Каталог пока пуст — загляните позже.
      </p>

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
            <span class="text-[17px] font-semibold leading-tight">{{ c.name }}</span>
            <span class="whitespace-nowrap font-mono text-[11px] text-[#A8A498]">{{ servicesInCategory(c.id).length }}</span>
          </div>
          <div class="mt-auto text-sm leading-snug text-[#6E6B60]">{{ categoryExamples(c.id) }}</div>
        </button>
      </div>
    </section>

    <!-- How it works -->
    <section id="how" class="bg-[#17160F] text-[#F5F2EA]">
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
          <a href="#trust" class="rounded-[13px] border border-[#F5F2EA]/30 px-6.5 py-4 text-base text-[#F5F2EA]">
            Как проверяется специализация
          </a>
        </div>
      </div>
    </section>

    <!-- Platform facts (real, server-enforced) -->
    <section id="trust" class="mx-auto max-w-[1280px] px-5 py-16 md:px-10 md:py-[84px]">
      <div class="grid overflow-hidden rounded-[22px] border border-[#E2DED2] bg-white sm:grid-cols-2 lg:grid-cols-4">
        <div
          v-for="(g, i) in platformFacts"
          :key="g.title"
          class="border-[#E7E3D9] p-7"
          :class="i < platformFacts.length - 1 ? 'border-b lg:border-b-0 lg:border-r' : ''"
        >
          <div class="mb-4.5 h-8.5 w-8.5 rounded-[10px]" :style="{ background: g.tint }" />
          <h3 class="mb-2 text-lg font-semibold">{{ g.title }}</h3>
          <p class="text-[15px] leading-normal text-[#6E6B60]">{{ g.text }}</p>
        </div>
      </div>
    </section>

    <!-- For masters CTA -->
    <section id="pro" class="mx-auto mb-16 max-w-[1280px] px-5 md:mb-[84px] md:px-10">
      <div class="grid gap-10 rounded-[26px] bg-[#5B4BE0] p-8 text-white sm:grid-cols-[1fr_0.8fr] sm:items-center md:p-14">
        <div>
          <div class="mb-4.5 font-mono text-[11px] uppercase tracking-[0.14em] text-white/70">Для специалистов</div>
          <h2 class="mk-display mb-4 text-[28px] font-bold leading-[1.1] tracking-tight md:text-[40px]">
            Берите заявки по своей специализации — без обзвонов
          </h2>
          <p class="mb-7 max-w-[460px] text-[17px] leading-[1.55] text-white/82 md:text-[18px]">
            Заявки появляются в открытом пуле по мере того, как их оставляют клиенты. Смотрите
            подходящие вам и берите в работу — в вебе или в мобильном приложении.
          </p>
          <div class="flex flex-wrap gap-3">
            <button type="button" class="rounded-[13px] bg-white px-6.5 py-4 text-base font-medium text-[#4536BC]" @click="goRegister">
              Стать мастером
            </button>
            <a href="#how" class="rounded-[13px] border border-white/45 px-6.5 py-4 text-base text-white">
              Как это работает
            </a>
          </div>
        </div>
        <div class="grid gap-3.5">
          <div v-for="stat in proStats" :key="stat.label" class="flex items-baseline justify-between gap-4 rounded-2xl bg-white/14 px-5.5 py-5">
            <span class="text-base text-white/80">{{ stat.label }}</span>
            <span class="mk-display text-right text-xl font-medium">{{ stat.value }}</span>
          </div>
        </div>
      </div>
    </section>

    <!-- FAQ -->
    <section id="faq" class="mx-auto max-w-[1280px] px-5 pb-16 md:px-10 md:pb-[84px]">
      <div class="grid gap-10 md:grid-cols-[0.7fr_1.3fr]">
        <h2 class="mk-display text-[28px] font-bold tracking-tight md:text-[36px]">Частые вопросы</h2>
        <div class="border-t border-[#E2DED2]">
          <div v-for="(f, i) in faqItems" :key="f.q" class="border-b border-[#E7E3D9]">
            <button
              type="button"
              class="flex w-full items-center justify-between gap-5 py-5.5 text-left text-lg font-medium"
              @click="toggleFaq(i)"
            >
              {{ f.q }}
              <span class="font-mono text-base text-[#A8A498]">{{ openFaq === i ? '−' : '+' }}</span>
            </button>
            <p v-if="openFaq === i" class="max-w-[720px] pb-6 pr-15 text-base leading-[1.6] text-[#6E6B60]">
              {{ f.a }}
            </p>
          </div>
        </div>
      </div>
    </section>

    <!-- Real-catalog service results (functional, tied to search/category filter above) -->
    <section id="services" class="mx-auto max-w-[1280px] px-5 pb-16 md:px-10 md:pb-[84px]">
      <div class="mb-5 flex items-center justify-between">
        <h2 class="mk-display text-2xl font-bold tracking-tight">
          {{ selectedCategoryId ? categories.find((c) => c.id === selectedCategoryId)?.name : 'Все услуги в каталоге' }}
        </h2>
        <button v-if="selectedCategoryId" type="button" class="text-sm font-medium text-[#5B4BE0] hover:underline" @click="selectedCategoryId = null">
          Сбросить фильтр
        </button>
      </div>

      <p v-if="!loading && visibleServices.length === 0" class="text-sm text-[#6E6B60]">
        Ничего не нашлось — попробуйте изменить запрос или выбрать другую категорию.
      </p>

      <div v-else class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        <button
          v-for="s in visibleServices"
          :key="s.id"
          type="button"
          class="flex flex-col items-start gap-2 rounded-[18px] border border-l-[3px] border-[#E2DED2] border-l-[#5B4BE0] bg-white p-5 text-left transition-colors hover:border-[#17160F]"
          @click="startRequest(s.id)"
        >
          <span class="text-base font-semibold">{{ s.name }}</span>
          <span v-if="s.description" class="text-sm text-[#6E6B60]">{{ s.description }}</span>
          <span class="mt-auto pt-2 text-sm font-medium text-[#5B4BE0]">Оставить заявку →</span>
        </button>
      </div>

      <p class="mt-8 text-center text-xs text-[#9B978A]">
        Заявка автоматически попадает к мастерам с подходящей специализацией — случайный
        исполнитель не сможет её взять.
      </p>
    </section>

    <!-- App block: the real Android app built for the MASTER role -->
    <section class="mx-auto max-w-[1280px] px-5 pb-16 md:px-10 md:pb-[84px]">
      <div class="grid items-center gap-10 overflow-hidden rounded-[26px] border border-[#E2DED2] bg-white p-8 sm:grid-cols-[1fr_0.9fr] md:p-11">
        <div>
          <h2 class="mk-display mb-3 text-2xl font-bold tracking-tight sm:text-[32px]">Приложение для мастеров</h2>
          <p class="mb-6.5 max-w-[420px] text-[17px] leading-[1.55] text-[#6E6B60]">
            Экран «Доступные заявки» показывает открытый пул по вашей специализации. Взяли заявку —
            она переходит в «Мои заявки», там же отмечаете, что работа выполнена.
          </p>
          <div class="flex flex-wrap gap-3">
            <span class="rounded-xl border border-[#DDD8CC] px-5.5 py-3.5 text-sm">Android</span>
          </div>
        </div>
        <div
          class="grid h-[220px] place-items-center rounded-[18px]"
          style="background: repeating-linear-gradient(135deg, #f1ede3 0 10px, #e9e4d8 10px 20px)"
        >
          <span class="rounded-lg bg-white/90 px-2.5 py-1.5 font-mono text-xs text-[#8D8A7E]">экран «Доступные заявки»</span>
        </div>
      </div>
    </section>

    <!-- Extra footer band with in-page anchors (the shell's own footer is a -->
    <!-- one-liner shared across every marketplace page) -->
    <footer class="bg-[#17160F] text-[#F5F2EA]">
      <div class="mx-auto max-w-[1280px] px-5 pb-10 pt-16 md:px-10">
        <div class="grid gap-10 border-b border-[#F5F2EA]/16 pb-11 sm:grid-cols-2 lg:grid-cols-[1.3fr_1fr_1fr_1fr]">
          <div>
            <div class="mb-3.5 flex items-center gap-2.5">
              <span class="mk-display flex h-7 w-7 items-center justify-center rounded-lg bg-[#F5F2EA] text-sm font-semibold text-[#17160F]">М</span>
              <span class="mk-display text-base font-medium">Мастерская</span>
            </div>
            <p class="max-w-[280px] text-[15px] leading-[1.55] text-[#F5F2EA]/55">
              Маркетплейс мастеров и заявок на бытовые услуги.
            </p>
          </div>
          <div class="grid gap-2.5 text-[15px]">
            <div class="mb-1 font-mono text-[11px] uppercase tracking-[0.12em] text-[#F5F2EA]/40">Клиентам</div>
            <RouterLink :to="{ name: 'marketplace-register' }" class="text-[#F5F2EA]/80 hover:text-[#F5F2EA]">Оставить заявку</RouterLink>
            <a href="#cats" class="text-[#F5F2EA]/80 hover:text-[#F5F2EA]">Все услуги</a>
            <a href="#faq" class="text-[#F5F2EA]/80 hover:text-[#F5F2EA]">Частые вопросы</a>
          </div>
          <div class="grid gap-2.5 text-[15px]">
            <div class="mb-1 font-mono text-[11px] uppercase tracking-[0.12em] text-[#F5F2EA]/40">Мастерам</div>
            <RouterLink :to="{ name: 'marketplace-register' }" class="text-[#F5F2EA]/80 hover:text-[#F5F2EA]">Начать работать</RouterLink>
            <a href="#trust" class="text-[#F5F2EA]/80 hover:text-[#F5F2EA]">Как устроена платформа</a>
          </div>
          <div class="grid gap-2.5 text-[15px]">
            <div class="mb-1 font-mono text-[11px] uppercase tracking-[0.12em] text-[#F5F2EA]/40">Компания</div>
            <RouterLink :to="{ name: 'marketplace-login' }" class="text-[#F5F2EA]/80 hover:text-[#F5F2EA]">Войти</RouterLink>
            <span class="text-[#F5F2EA]/80">Дипломный проект</span>
          </div>
        </div>
        <div class="flex flex-wrap items-center justify-between gap-5 border-t border-[#F5F2EA]/16 pt-5 text-sm text-[#F5F2EA]/40">
          <span>© 2026 Мастерская</span>
          <span>Учебный проект — не коммерческий сервис</span>
        </div>
      </div>
    </footer>
  </MarketplaceShell>
</template>
