<script setup lang="ts">
import { ref } from 'vue'
import MarketplaceShell from '@/components/layout/MarketplaceShell.vue'

// Everything below describes what the backend actually does (see
// internal/marketplace/service.go): a hard, server-enforced specialization
// check on claim, first-come self-claim (no bidding/offers), no payments,
// no ratings, no identity verification. Copy is written to match that --
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
    <div class="mx-auto max-w-[1280px] px-5 py-14 md:px-10">
      <h1 class="mk-display text-[32px] font-bold tracking-tight md:text-[40px]">Как устроена платформа</h1>

      <div class="mt-9 grid overflow-hidden rounded-[22px] border border-[#E2DED2] bg-white sm:grid-cols-2 lg:grid-cols-4">
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

      <div class="mt-16 grid gap-10 md:grid-cols-[0.7fr_1.3fr]">
        <h2 class="mk-display text-[28px] font-bold tracking-tight">Частые вопросы</h2>
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
    </div>
  </MarketplaceShell>
</template>
