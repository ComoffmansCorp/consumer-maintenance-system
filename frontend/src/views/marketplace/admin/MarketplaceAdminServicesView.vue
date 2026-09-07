<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { catalogApi } from '@/api/marketplace'
import { extractErrorMessage } from '@/api/client'
import { useToastStore } from '@/stores/toast'
import type { CategoryDTO, ServiceDTO } from '@/types'
import MarketplaceShell from '@/components/layout/MarketplaceShell.vue'

const toast = useToastStore()
const categories = ref<CategoryDTO[]>([])
const services = ref<ServiceDTO[]>([])
const loading = ref(true)
const error = ref('')
const saving = ref(false)

function flatCategories(cats: CategoryDTO[]): CategoryDTO[] {
  return cats.flatMap((c) => [c, ...(c.subcategories ?? [])])
}

async function load() {
  loading.value = true
  error.value = ''
  try {
    const [cats, svcs] = await Promise.all([catalogApi.listCategories(), catalogApi.listServices()])
    categories.value = flatCategories(cats)
    services.value = svcs
  } catch (e) {
    error.value = extractErrorMessage(e)
  } finally {
    loading.value = false
  }
}
onMounted(load)

function categoryName(id: number) {
  return categories.value.find((c) => c.id === id)?.name ?? `#${id}`
}

const showCreate = ref(false)
const form = ref({ categoryId: '', name: '', description: '', priceFrom: '', priceTo: '', unit: '' })

async function create() {
  if (!form.value.categoryId || !form.value.name.trim()) return
  saving.value = true
  try {
    await catalogApi.createService({
      categoryId: Number(form.value.categoryId),
      name: form.value.name.trim(),
      description: form.value.description.trim() || undefined,
      priceFrom: form.value.priceFrom ? Number(form.value.priceFrom) : undefined,
      priceTo: form.value.priceTo ? Number(form.value.priceTo) : undefined,
      unit: form.value.unit.trim() || undefined,
    })
    toast.success('Услуга создана')
    showCreate.value = false
    form.value = { categoryId: '', name: '', description: '', priceFrom: '', priceTo: '', unit: '' }
    await load()
  } catch (e) {
    toast.error(extractErrorMessage(e))
  } finally {
    saving.value = false
  }
}

async function toggleActive(svc: ServiceDTO) {
  try {
    await catalogApi.updateService(svc.id, {
      name: svc.name,
      description: svc.description,
      priceFrom: svc.priceFrom,
      priceTo: svc.priceTo,
      unit: svc.unit,
      active: !svc.active,
    })
    toast.success('Обновлено')
    await load()
  } catch (e) {
    toast.error(extractErrorMessage(e))
  }
}
</script>

<template>
  <MarketplaceShell>
    <div class="mx-auto max-w-[900px]">
      <div class="flex items-center justify-between gap-3">
        <h1 class="mk-display text-2xl font-bold tracking-tight">Услуги</h1>
        <button type="button" class="rounded-[11px] bg-[#5B4BE0] px-4 py-2.5 text-sm font-medium text-white hover:opacity-90" @click="showCreate = true">
          Добавить услугу
        </button>
      </div>

      <div v-if="loading" class="py-16 text-center text-sm text-[#8D8A7E]">Загрузка…</div>
      <div v-else-if="error" class="mt-6 rounded-2xl border border-[#F3D3CE] bg-[#FBF0EE] px-6 py-10 text-center text-sm text-[#B3261E]">
        {{ error }}
      </div>
      <div v-else class="mt-6 flex flex-col gap-3">
        <div v-for="svc in services" :key="svc.id" class="flex items-center justify-between rounded-2xl border border-[#E2DED2] bg-white p-5">
          <div class="flex items-center gap-3">
            <img
              v-if="svc.imageUrl"
              :src="svc.imageUrl"
              alt=""
              class="h-11 w-11 shrink-0 rounded-lg object-cover"
            />
            <div v-else class="h-11 w-11 shrink-0 rounded-lg bg-[#EFEBE1]" />
            <div>
              <p class="text-sm font-medium">{{ svc.name }}</p>
              <p class="mt-0.5 text-xs text-[#9B978A]">{{ categoryName(svc.categoryId) }} · {{ svc.priceFrom ?? '—' }}–{{ svc.priceTo ?? '—' }} ₽{{ svc.unit ? ` / ${svc.unit}` : '' }}</p>
            </div>
          </div>
          <button type="button" class="rounded-[11px] border border-[#E2DED2] px-3 py-2 text-sm hover:border-[#17160F]" @click="toggleActive(svc)">
            {{ svc.active ? 'Скрыть' : 'Включить' }}
          </button>
        </div>
      </div>
    </div>

    <Teleport to="body">
      <div v-if="showCreate" class="fixed inset-0 z-40 flex items-center justify-center bg-[#17160F]/40 p-4" @click.self="showCreate = false">
        <div class="w-full max-w-md rounded-2xl border border-[#E2DED2] bg-white p-6 shadow-lg">
          <h2 class="mk-display text-lg font-semibold">Новая услуга</h2>
          <form class="mt-4 flex flex-col gap-4" @submit.prevent="create">
            <label class="flex flex-col gap-1.5">
              <span class="font-mono text-[11px] uppercase tracking-[0.1em] text-[#9B978A]">Категория</span>
              <select v-model="form.categoryId" required class="rounded-xl border border-[#E2DED2] px-4 py-2.5 text-base outline-none focus-visible:border-[#5B4BE0]">
                <option value="" disabled>Выберите категорию</option>
                <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
              </select>
            </label>
            <label class="flex flex-col gap-1.5">
              <span class="font-mono text-[11px] uppercase tracking-[0.1em] text-[#9B978A]">Название</span>
              <input v-model="form.name" required class="rounded-xl border border-[#E2DED2] px-4 py-2.5 text-base outline-none focus-visible:border-[#5B4BE0]" />
            </label>
            <label class="flex flex-col gap-1.5">
              <span class="font-mono text-[11px] uppercase tracking-[0.1em] text-[#9B978A]">Описание</span>
              <textarea v-model="form.description" rows="2" class="rounded-xl border border-[#E2DED2] px-4 py-2.5 text-base outline-none focus-visible:border-[#5B4BE0]" />
            </label>
            <div class="grid grid-cols-3 gap-3">
              <label class="flex flex-col gap-1.5">
                <span class="font-mono text-[11px] uppercase tracking-[0.1em] text-[#9B978A]">От, ₽</span>
                <input v-model="form.priceFrom" type="number" class="rounded-xl border border-[#E2DED2] px-3 py-2.5 text-base outline-none focus-visible:border-[#5B4BE0]" />
              </label>
              <label class="flex flex-col gap-1.5">
                <span class="font-mono text-[11px] uppercase tracking-[0.1em] text-[#9B978A]">До, ₽</span>
                <input v-model="form.priceTo" type="number" class="rounded-xl border border-[#E2DED2] px-3 py-2.5 text-base outline-none focus-visible:border-[#5B4BE0]" />
              </label>
              <label class="flex flex-col gap-1.5">
                <span class="font-mono text-[11px] uppercase tracking-[0.1em] text-[#9B978A]">Ед.</span>
                <input v-model="form.unit" placeholder="услуга" class="rounded-xl border border-[#E2DED2] px-3 py-2.5 text-base outline-none focus-visible:border-[#5B4BE0]" />
              </label>
            </div>
            <div class="flex justify-end gap-2">
              <button type="button" class="rounded-[11px] border border-[#E2DED2] px-4 py-2.5 text-sm hover:border-[#17160F]" @click="showCreate = false">Отмена</button>
              <button type="submit" :disabled="saving" class="rounded-[11px] bg-[#5B4BE0] px-4 py-2.5 text-sm font-medium text-white hover:opacity-90 disabled:opacity-60">
                {{ saving ? 'Сохраняем…' : 'Создать' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </Teleport>
  </MarketplaceShell>
</template>
