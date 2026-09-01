<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { marketplaceApi, requestStatusLabels, requestStatusStyles } from '@/api/marketplace'
import { extractErrorMessage } from '@/api/client'
import { useToastStore } from '@/stores/toast'
import type { RequestDTO } from '@/types'
import MarketplaceShell from '@/components/layout/MarketplaceShell.vue'

const props = defineProps<{ id: string }>()
const toast = useToastStore()

const request = ref<RequestDTO | null>(null)
const loading = ref(true)
const error = ref('')
const working = ref(false)

async function load() {
  loading.value = true
  error.value = ''
  try {
    request.value = await marketplaceApi.getRequest(Number(props.id))
  } catch (e) {
    error.value = extractErrorMessage(e)
  } finally {
    loading.value = false
  }
}
onMounted(load)

const showCancel = ref(false)
const cancelReason = ref('')

async function cancel() {
  if (!request.value || !cancelReason.value.trim()) return
  working.value = true
  try {
    request.value = await marketplaceApi.cancelRequest(request.value.id, cancelReason.value)
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
            <p class="mt-1 text-sm">{{ request.masterName || 'Ещё не назначен' }}</p>
          </div>
          <div>
            <p class="font-mono text-[11px] uppercase tracking-[0.1em] text-[#9B978A]">Создана</p>
            <p class="mt-1 font-mono text-sm">{{ formatDateTime(request.createdAt) }}</p>
          </div>
          <div v-if="request.cancelReason" class="sm:col-span-2">
            <p class="font-mono text-[11px] uppercase tracking-[0.1em] text-[#9B978A]">Причина отмены</p>
            <p class="mt-1 text-sm">{{ request.cancelReason }}</p>
          </div>
        </div>

        <div
          v-if="request.status === 'OPEN' || request.status === 'IN_PROGRESS'"
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
