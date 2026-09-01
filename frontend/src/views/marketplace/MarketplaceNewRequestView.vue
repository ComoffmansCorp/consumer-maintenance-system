<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { marketplaceApi } from '@/api/marketplace'
import { extractErrorMessage } from '@/api/client'
import { useToastStore } from '@/stores/toast'
import { useYandexSuggest } from '@/composables/useYandexSuggest'
import type { CategoryDTO, ServiceDTO } from '@/types'
import MarketplaceShell from '@/components/layout/MarketplaceShell.vue'

const route = useRoute()
const router = useRouter()
const toast = useToastStore()
const { available: suggestAvailable, suggestions, onInput: onAddressInput, clear: clearSuggestions, resolve } = useYandexSuggest()

const loading = ref(true)
const submitting = ref(false)
const error = ref('')

const categories = ref<CategoryDTO[]>([])
const services = ref<ServiceDTO[]>([])
const categoryId = ref('')
const serviceId = ref('')
const description = ref('')
const addressText = ref('')
const latitude = ref<number | undefined>()
const longitude = ref<number | undefined>()
const showSuggestions = ref(false)

async function load() {
  loading.value = true
  try {
    const [categoriesResult, servicesResult] = await Promise.all([
      marketplaceApi.listCategories(),
      marketplaceApi.listServices(),
    ])
    categories.value = categoriesResult
    services.value = servicesResult

    const preselected = route.query.serviceId ? Number(route.query.serviceId) : null
    if (preselected) {
      const service = services.value.find((s) => s.id === preselected)
      if (service) {
        serviceId.value = String(service.id)
        categoryId.value = String(service.categoryId)
      }
    }
  } finally {
    loading.value = false
  }
}

onMounted(load)

function onAddressTyped() {
  // Manual edits invalidate any coordinates captured from a previous pick --
  // the text no longer necessarily matches the geocoded point.
  latitude.value = undefined
  longitude.value = undefined
  showSuggestions.value = true
  onAddressInput(addressText.value)
}

async function pickSuggestion(suggestion: { text: string; uri?: string }) {
  const resolved = await resolve(suggestion)
  addressText.value = resolved.value
  latitude.value = resolved.latitude
  longitude.value = resolved.longitude
  showSuggestions.value = false
  clearSuggestions()
}

function servicesForCategory(catId: string) {
  if (!catId) return services.value
  return services.value.filter((s) => String(s.categoryId) === catId)
}

async function submit() {
  if (!serviceId.value) {
    error.value = 'Выберите услугу'
    return
  }
  if (!description.value.trim()) {
    error.value = 'Опишите, что нужно сделать'
    return
  }
  if (!addressText.value.trim()) {
    error.value = 'Укажите адрес'
    return
  }
  submitting.value = true
  error.value = ''
  try {
    await marketplaceApi.createRequest({
      serviceId: Number(serviceId.value),
      description: description.value.trim(),
      addressText: addressText.value.trim(),
      latitude: latitude.value,
      longitude: longitude.value,
    })
    toast.success('Заявка создана')
    router.push({ name: 'marketplace-my-requests' })
  } catch (e) {
    error.value = extractErrorMessage(e)
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <MarketplaceShell>
    <div class="mx-auto max-w-lg">
      <p class="font-mono text-[11px] uppercase tracking-[0.14em] text-[#5B4BE0]">Новая заявка</p>
      <h1 class="mk-display mt-1 text-2xl font-bold tracking-tight">Опишите, что нужно сделать</h1>

      <div v-if="loading" class="flex items-center justify-center gap-2 py-16 text-sm text-[#8D8A7E]">
        <svg class="h-4 w-4 animate-spin motion-reduce:animate-none" viewBox="0 0 24 24" fill="none">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v4a4 4 0 00-4 4H4z" />
        </svg>
        Загрузка…
      </div>

      <form v-else class="mt-8 flex flex-col gap-5" @submit.prevent="submit">
        <label class="flex flex-col gap-1.5">
          <span class="font-mono text-[11px] uppercase tracking-[0.1em] text-[#9B978A]">Категория</span>
          <select
            v-model="categoryId"
            class="rounded-xl border border-[#E2DED2] bg-white px-4 py-3 text-base outline-none focus-visible:border-[#5B4BE0]"
            @change="serviceId = ''"
          >
            <option value="" disabled>Выберите категорию</option>
            <option v-for="c in categories" :key="c.id" :value="String(c.id)">{{ c.name }}</option>
          </select>
        </label>

        <label class="flex flex-col gap-1.5">
          <span class="font-mono text-[11px] uppercase tracking-[0.1em] text-[#9B978A]">Услуга</span>
          <select
            v-model="serviceId"
            required
            class="rounded-xl border border-[#E2DED2] bg-white px-4 py-3 text-base outline-none focus-visible:border-[#5B4BE0]"
          >
            <option value="" disabled>Выберите услугу</option>
            <option v-for="s in servicesForCategory(categoryId)" :key="s.id" :value="String(s.id)">{{ s.name }}</option>
          </select>
        </label>

        <label class="flex flex-col gap-1.5">
          <span class="font-mono text-[11px] uppercase tracking-[0.1em] text-[#9B978A]">Описание задачи</span>
          <textarea
            v-model="description"
            required
            rows="4"
            placeholder="Что случилось, что нужно сделать"
            class="resize-none rounded-xl border border-[#E2DED2] bg-white px-4 py-3 text-base outline-none focus-visible:border-[#5B4BE0]"
          />
        </label>

        <label class="relative flex flex-col gap-1.5">
          <span class="font-mono text-[11px] uppercase tracking-[0.1em] text-[#9B978A]">Адрес</span>
          <input
            v-model="addressText"
            type="text"
            required
            placeholder="Город, улица, дом"
            class="rounded-xl border border-[#E2DED2] bg-white px-4 py-3 text-base outline-none focus-visible:border-[#5B4BE0]"
            autocomplete="off"
            @input="onAddressTyped"
            @focus="showSuggestions = true"
            @blur="showSuggestions = false"
          />
          <span v-if="!suggestAvailable" class="text-xs text-[#9B978A]">Введите адрес вручную</span>
          <span v-else-if="latitude" class="text-xs text-[#9B978A]">Адрес распознан на карте</span>

          <ul
            v-if="showSuggestions && suggestions.length > 0"
            class="absolute top-full z-10 mt-1 w-full overflow-hidden rounded-xl border border-[#E2DED2] bg-white py-1 shadow-lg"
          >
            <li v-for="s in suggestions" :key="s.text">
              <button type="button" class="w-full px-3 py-2 text-left text-sm hover:bg-[#F1EDE3]" @mousedown.prevent="pickSuggestion(s)">
                {{ s.text }}
              </button>
            </li>
          </ul>
        </label>

        <p v-if="error" class="text-sm font-medium text-[#B3261E]">{{ error }}</p>

        <button
          type="submit"
          :disabled="submitting"
          class="mt-2 rounded-[14px] bg-[#5B4BE0] px-6 py-3.5 text-base font-medium text-white hover:bg-[#4536BC] disabled:opacity-60"
        >
          {{ submitting ? 'Отправляем…' : 'Отправить заявку' }}
        </button>
      </form>
    </div>
  </MarketplaceShell>
</template>
