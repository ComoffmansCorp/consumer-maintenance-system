<script setup lang="ts">
import { useRouter } from 'vue-router'
import MarketplaceShell from '@/components/layout/MarketplaceShell.vue'

const router = useRouter()

// Plain string, not a static `src="..."` attribute -- see the same note in
// MarketplaceLandingView.vue: Vue's compiler rewrites static asset src
// attributes into build-time imports, which fails for paths served at
// runtime by the gateway/MinIO (docker-compose.yml minio/minio-init).
const heroAppScreenImg = '/media/hero/app-screen.svg'

const proStats = [
  { label: 'Специализация', value: 'Проверяется на сервере' },
  { label: 'Модель распределения', value: 'Self-claim' },
  { label: 'Комиссия платформы', value: '0 ₽' },
]

function goRegister() {
  router.push({ name: 'marketplace-register' })
}
</script>

<template>
  <MarketplaceShell>
    <div class="mx-auto max-w-[1280px] px-5 py-14 md:px-10">
      <div class="grid gap-10 rounded-[26px] bg-[#5B4BE0] p-8 text-white sm:grid-cols-[1fr_0.8fr] sm:items-center md:p-14">
        <div>
          <div class="mb-4.5 font-mono text-[11px] uppercase tracking-[0.14em] text-white/70">Для специалистов</div>
          <h1 class="mk-display mb-4 text-[28px] font-bold leading-[1.1] tracking-tight md:text-[40px]">
            Берите заявки по своей специализации — без обзвонов
          </h1>
          <p class="mb-7 max-w-[460px] text-[17px] leading-[1.55] text-white/82 md:text-[18px]">
            Заявки появляются в открытом пуле по мере того, как их оставляют клиенты. Смотрите
            подходящие вам и берите в работу — в вебе или в мобильном приложении.
          </p>
          <div class="flex flex-wrap gap-3">
            <button type="button" class="rounded-[13px] bg-white px-6.5 py-4 text-base font-medium text-[#4536BC]" @click="goRegister">
              Стать мастером
            </button>
          </div>
        </div>
        <div class="grid gap-3.5">
          <div v-for="stat in proStats" :key="stat.label" class="flex items-baseline justify-between gap-4 rounded-2xl bg-white/14 px-5.5 py-5">
            <span class="text-base text-white/80">{{ stat.label }}</span>
            <span class="mk-display text-right text-xl font-medium">{{ stat.value }}</span>
          </div>
        </div>
      </div>

      <!-- App block: the real Android app built for the MASTER role -->
      <div class="mt-10 grid items-center gap-10 overflow-hidden rounded-[26px] border border-[#E2DED2] bg-white p-8 sm:grid-cols-[1fr_0.9fr] md:p-11">
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
        <div class="grid h-[220px] place-items-center overflow-hidden rounded-[18px] bg-[#F1EDE3]">
          <img :src="heroAppScreenImg" alt="Экран «Доступные заявки»" class="h-full object-contain" />
        </div>
      </div>
    </div>
  </MarketplaceShell>
</template>
