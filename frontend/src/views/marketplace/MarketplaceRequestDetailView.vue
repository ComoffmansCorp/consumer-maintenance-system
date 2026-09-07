<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'
import { requestsApi, paymentsApi, reviewsApi, chatApi } from '@/api/marketplace'
import { extractErrorMessage } from '@/api/client'
import { useToastStore } from '@/stores/toast'
import { useAuthStore } from '@/stores/auth'
import type { RequestDTO, OfferDTO, PaymentDTO, MessageDTO } from '@/types'
import MarketplaceShell from '@/components/layout/MarketplaceShell.vue'

const requestStatusLabels: Record<string, string> = {
  OPEN: 'Открыта',
  ASSIGNED: 'В работе',
  COMPLETED: 'Выполнена',
  CANCELED: 'Отменена',
}
const requestStatusStyles: Record<string, { bg: string; fg: string }> = {
  OPEN: { bg: '#EFEBE1', fg: '#55524A' },
  ASSIGNED: { bg: '#E7E3FC', fg: '#5B4BE0' },
  COMPLETED: { bg: '#DCEFE1', fg: '#2E7D4F' },
  CANCELED: { bg: '#FBE9E7', fg: '#B3261E' },
}
const paymentStatusLabels: Record<string, string> = {
  HELD: 'Удержана (эскроу)',
  RELEASED: 'Выплачена мастеру',
  REFUNDED: 'Возвращена клиенту',
}
const paymentStatusStyles: Record<string, { bg: string; fg: string }> = {
  HELD: { bg: '#EFEBE1', fg: '#55524A' },
  RELEASED: { bg: '#DCEFE1', fg: '#2E7D4F' },
  REFUNDED: { bg: '#FBE9E7', fg: '#B3261E' },
}

const props = defineProps<{ id: string }>()
const toast = useToastStore()
const auth = useAuthStore()

const request = ref<RequestDTO | null>(null)
const offers = ref<OfferDTO[]>([])
const payment = ref<PaymentDTO | null>(null)
const loading = ref(true)
const error = ref('')
const working = ref(false)

const messages = ref<MessageDTO[]>([])
const newMessage = ref('')
let chatTimer: ReturnType<typeof setInterval> | null = null

async function pollMessages() {
  if (!request.value) return
  try {
    messages.value = await chatApi.listMessages(request.value.id)
  } catch {
    // Chat polling failure is silent -- the rest of the page still works,
    // no need to interrupt the user with an error toast every few seconds.
  }
}

function startChatPolling() {
  pollMessages()
  chatTimer = setInterval(pollMessages, 5000)
}

async function sendMessage() {
  if (!request.value || !newMessage.value.trim()) return
  try {
    await chatApi.send(request.value.id, newMessage.value.trim())
    newMessage.value = ''
    await pollMessages()
  } catch (e) {
    toast.error(extractErrorMessage(e))
  }
}

onUnmounted(() => {
  if (chatTimer) clearInterval(chatTimer)
})

async function load() {
  loading.value = true
  error.value = ''
  try {
    request.value = await requestsApi.get(Number(props.id))
    if (request.value.status === 'OPEN') {
      offers.value = await requestsApi.listOffers(request.value.id)
    }
    if (request.value.status === 'ASSIGNED' || request.value.status === 'COMPLETED') {
      payment.value = await paymentsApi.getForRequest(request.value.id).catch(() => null)
    }
    if (request.value.masterId) {
      startChatPolling()
    }
  } catch (e) {
    error.value = extractErrorMessage(e)
  } finally {
    loading.value = false
  }
}
onMounted(load)

async function acceptOffer(offer: OfferDTO) {
  if (!request.value) return
  working.value = true
  try {
    request.value = await requestsApi.acceptOffer(request.value.id, offer.id)
    toast.success('Мастер выбран')
    offers.value = []
  } catch (e) {
    toast.error(extractErrorMessage(e))
  } finally {
    working.value = false
  }
}

const reviewRating = ref(5)
const reviewComment = ref('')
const reviewSubmitted = ref(false)

async function submitReview() {
  if (!request.value) return
  working.value = true
  try {
    await reviewsApi.create({
      requestId: request.value.id,
      rating: reviewRating.value,
      comment: reviewComment.value,
    })
    reviewSubmitted.value = true
    toast.success('Спасибо за отзыв')
  } catch (e) {
    toast.error(extractErrorMessage(e))
  } finally {
    working.value = false
  }
}

const showCancel = ref(false)
const cancelReason = ref('')

async function cancel() {
  if (!request.value || !cancelReason.value.trim()) return
  working.value = true
  try {
    request.value = await requestsApi.cancel(request.value.id, cancelReason.value)
    toast.success('Заявка отменена')
    showCancel.value = false
    cancelReason.value = ''
  } catch (e) {
    toast.error(extractErrorMessage(e))
  } finally {
    working.value = false
  }
}

function formatDateTime(iso?: string) {
  if (!iso) return '—'
  return new Date(iso).toLocaleString('ru-RU')
}
</script>

<template>
  <MarketplaceShell>
    <div class="mx-auto max-w-lg">
      <div v-if="loading" class="flex items-center justify-center gap-2 py-16 text-sm text-[#8D8A7E]">
        <svg class="h-4 w-4 animate-spin motion-reduce:animate-none" viewBox="0 0 24 24" fill="none">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v4a4 4 0 00-4 4H4z" />
        </svg>
        Загрузка…
      </div>

      <div v-else-if="error" class="flex flex-col items-center gap-3 rounded-2xl border border-[#F3D3CE] bg-[#FBF0EE] px-6 py-10 text-center">
        <p class="text-sm font-medium text-[#B3261E]">{{ error }}</p>
        <button type="button" class="text-sm font-medium text-[#5B4BE0] hover:underline" @click="load">Повторить</button>
      </div>

      <template v-else-if="request">
        <div class="flex flex-wrap items-start justify-between gap-3">
          <div>
            <p class="text-sm text-[#6E6B60]">Заявка №{{ request.id }}</p>
            <h1 class="mk-display text-2xl font-bold tracking-tight">{{ request.serviceName || 'Услуга' }}</h1>
          </div>
          <span
            class="inline-flex w-fit items-center gap-1.5 rounded-[3px] border-l-[3px] px-2.5 py-1 text-xs font-medium tracking-wide"
            :style="{ borderLeftColor: requestStatusStyles[request.status].fg, background: requestStatusStyles[request.status].bg, color: requestStatusStyles[request.status].fg }"
          >
            {{ requestStatusLabels[request.status] }}
          </span>
        </div>

        <div class="mt-6 grid grid-cols-1 gap-4 rounded-2xl border border-[#E2DED2] bg-white p-6 sm:grid-cols-2">
          <div class="sm:col-span-2">
            <p class="font-mono text-[11px] uppercase tracking-[0.1em] text-[#9B978A]">Описание</p>
            <p class="mt-1 text-sm">{{ request.description }}</p>
          </div>
          <div class="sm:col-span-2">
            <p class="font-mono text-[11px] uppercase tracking-[0.1em] text-[#9B978A]">Адрес</p>
            <p class="mt-1 text-sm">{{ request.addressText }}</p>
          </div>
          <div>
            <p class="font-mono text-[11px] uppercase tracking-[0.1em] text-[#9B978A]">Мастер</p>
            <p class="mt-1 text-sm">{{ request.masterId ? `Мастер #${request.masterId}` : 'Ещё не назначен' }}</p>
          </div>
          <div>
            <p class="font-mono text-[11px] uppercase tracking-[0.1em] text-[#9B978A]">Создана</p>
            <p class="mt-1 font-mono text-sm">{{ formatDateTime(request.createdAt) }}</p>
          </div>
          <div v-if="request.agreedPrice" class="sm:col-span-2">
            <p class="font-mono text-[11px] uppercase tracking-[0.1em] text-[#9B978A]">Согласованная цена</p>
            <p class="mt-1 text-sm">{{ request.agreedPrice }} ₽</p>
          </div>
          <div v-if="payment" class="sm:col-span-2">
            <p class="font-mono text-[11px] uppercase tracking-[0.1em] text-[#9B978A]">Оплата</p>
            <span
              class="mt-1 inline-flex w-fit items-center gap-1.5 rounded-[3px] border-l-[3px] px-2.5 py-1 text-xs font-medium"
              :style="{ borderLeftColor: paymentStatusStyles[payment.status].fg, background: paymentStatusStyles[payment.status].bg, color: paymentStatusStyles[payment.status].fg }"
            >
              {{ paymentStatusLabels[payment.status] }}
            </span>
          </div>
          <div v-if="request.cancelReason" class="sm:col-span-2">
            <p class="font-mono text-[11px] uppercase tracking-[0.1em] text-[#9B978A]">Причина отмены</p>
            <p class="mt-1 text-sm">{{ request.cancelReason }}</p>
          </div>
        </div>

        <div v-if="request.status === 'OPEN' && offers.length" class="mt-4 rounded-2xl border border-[#E2DED2] bg-white p-6">
          <p class="mk-display text-base font-semibold">Отклики мастеров ({{ offers.length }})</p>
          <ul class="mt-3 flex flex-col gap-2">
            <li
              v-for="offer in offers"
              :key="offer.id"
              class="flex flex-wrap items-center justify-between gap-3 rounded-xl border border-[#E2DED2] px-4 py-3"
            >
              <div class="flex items-center gap-3">
                <img
                  v-if="offer.masterAvatarUrl"
                  :src="offer.masterAvatarUrl"
                  alt=""
                  class="h-9 w-9 shrink-0 rounded-full object-cover"
                />
                <div v-else class="h-9 w-9 shrink-0 rounded-full bg-[#EFEBE1]" />
                <div>
                  <p class="text-sm font-medium">{{ offer.price }} ₽ — мастер #{{ offer.masterId }}</p>
                  <p v-if="offer.comment" class="mt-0.5 text-sm text-[#6E6B60]">{{ offer.comment }}</p>
                </div>
              </div>
              <button
                type="button"
                :disabled="working"
                class="rounded-[11px] bg-[#5B4BE0] px-3 py-2 text-sm font-medium text-white hover:opacity-90 disabled:opacity-60"
                @click="acceptOffer(offer)"
              >
                Выбрать
              </button>
            </li>
          </ul>
        </div>
        <p v-else-if="request.status === 'OPEN'" class="mt-4 text-sm text-[#6E6B60]">
          Пока никто из мастеров не откликнулся.
        </p>

        <div
          v-if="request.status === 'OPEN' || request.status === 'ASSIGNED'"
          class="mt-4 flex flex-wrap items-center gap-3 rounded-2xl border border-[#E2DED2] bg-white p-6"
        >
          <button
            type="button"
            class="rounded-[11px] border border-[#B3261E] px-4 py-2.5 text-sm font-medium text-[#B3261E] hover:bg-[#FBF0EE]"
            @click="showCancel = true"
          >
            Отменить заявку
          </button>
        </div>

        <div v-if="request.status === 'COMPLETED'" class="mt-4 rounded-2xl border border-[#E2DED2] bg-white p-6">
          <template v-if="reviewSubmitted">
            <p class="text-sm text-[#2E7D4F]">Отзыв отправлен, спасибо!</p>
          </template>
          <template v-else>
            <p class="mk-display text-base font-semibold">Оставить отзыв мастеру</p>
            <form class="mt-3 flex flex-col gap-3" @submit.prevent="submitReview">
              <label class="flex flex-col gap-1.5">
                <span class="font-mono text-[11px] uppercase tracking-[0.1em] text-[#9B978A]">Оценка (1–5)</span>
                <input v-model.number="reviewRating" type="number" min="1" max="5" required class="w-24 rounded-xl border border-[#E2DED2] px-4 py-2.5 text-base outline-none focus-visible:border-[#5B4BE0]" />
              </label>
              <label class="flex flex-col gap-1.5">
                <span class="font-mono text-[11px] uppercase tracking-[0.1em] text-[#9B978A]">Комментарий</span>
                <textarea v-model="reviewComment" rows="3" class="rounded-xl border border-[#E2DED2] px-4 py-2.5 text-base outline-none focus-visible:border-[#5B4BE0]" />
              </label>
              <button type="submit" :disabled="working" class="w-fit rounded-[11px] bg-[#5B4BE0] px-4 py-2.5 text-sm font-medium text-white hover:opacity-90 disabled:opacity-60">
                Отправить отзыв
              </button>
            </form>
          </template>
        </div>

        <div v-if="request.masterId" class="mt-4 flex flex-col rounded-2xl border border-[#E2DED2] bg-white p-6">
          <p class="mk-display text-base font-semibold">Чат с мастером</p>
          <div class="mt-3 flex max-h-80 flex-col gap-2 overflow-y-auto">
            <p v-if="messages.length === 0" class="text-sm text-[#9B978A]">Сообщений пока нет.</p>
            <div
              v-for="m in messages"
              :key="m.id"
              class="max-w-[80%] rounded-xl px-3 py-2 text-sm"
              :class="m.senderId === auth.userId ? 'self-end bg-[#5B4BE0] text-white' : 'self-start bg-[#F1EDE3] text-[#17160F]'"
            >
              {{ m.text }}
            </div>
          </div>
          <form class="mt-3 flex gap-2" @submit.prevent="sendMessage">
            <input
              v-model="newMessage"
              type="text"
              placeholder="Сообщение…"
              class="flex-1 rounded-xl border border-[#E2DED2] px-4 py-2.5 text-base outline-none focus-visible:border-[#5B4BE0]"
            />
            <button type="submit" class="rounded-[11px] bg-[#5B4BE0] px-4 py-2.5 text-sm font-medium text-white hover:opacity-90">
              Отправить
            </button>
          </form>
        </div>
      </template>
    </div>

    <Teleport to="body">
      <div v-if="showCancel" class="fixed inset-0 z-40 flex items-center justify-center bg-[#17160F]/40 p-4" @click.self="showCancel = false">
        <div class="w-full max-w-lg rounded-2xl border border-[#E2DED2] bg-white p-6 shadow-lg">
          <div class="mb-4 flex items-center justify-between">
            <h2 class="mk-display text-lg font-semibold">Отменить заявку</h2>
            <button type="button" class="rounded-md p-1 text-[#9B978A] hover:bg-[#F1EDE3]" aria-label="Закрыть" @click="showCancel = false">✕</button>
          </div>
          <form class="flex flex-col gap-4" @submit.prevent="cancel">
            <label class="flex flex-col gap-1.5">
              <span class="font-mono text-[11px] uppercase tracking-[0.1em] text-[#9B978A]">Причина отмены</span>
              <input
                v-model="cancelReason"
                type="text"
                required
                class="rounded-xl border border-[#E2DED2] bg-white px-4 py-3 text-base outline-none focus-visible:border-[#5B4BE0]"
              />
            </label>
            <div class="flex justify-end gap-2">
              <button type="button" class="rounded-[11px] border border-[#E2DED2] px-4 py-2.5 text-sm font-medium hover:border-[#17160F]" @click="showCancel = false">
                Назад
              </button>
              <button type="submit" :disabled="working" class="rounded-[11px] bg-[#B3261E] px-4 py-2.5 text-sm font-medium text-white hover:opacity-90 disabled:opacity-60">
                {{ working ? 'Отменяем…' : 'Отменить заявку' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </Teleport>
  </MarketplaceShell>
</template>
