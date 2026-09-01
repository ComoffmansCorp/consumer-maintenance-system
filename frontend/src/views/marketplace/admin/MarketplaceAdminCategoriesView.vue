<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { catalogApi } from '@/api/marketplace'
import { extractErrorMessage } from '@/api/client'
import { useToastStore } from '@/stores/toast'
import type { CategoryDTO } from '@/types'
import MarketplaceShell from '@/components/layout/MarketplaceShell.vue'

const toast = useToastStore()
const categories = ref<CategoryDTO[]>([])
const loading = ref(true)
const error = ref('')
const saving = ref(false)

async function load() {
  loading.value = true
  error.value = ''
  try {
    categories.value = await catalogApi.listCategories()
  } catch (e) {
    error.value = extractErrorMessage(e)
  } finally {
    loading.value = false
  }
}
onMounted(load)

const showCreate = ref(false)
const newName = ref('')
const newParentId = ref<number | null>(null)

async function create() {
  if (!newName.value.trim()) return
  saving.value = true
  try {
    await catalogApi.createCategory({
      name: newName.value.trim(),
      parentCategoryId: newParentId.value ?? undefined,
    })
    toast.success('Категория создана')
    showCreate.value = false
    newName.value = ''
    newParentId.value = null
    await load()
  } catch (e) {
    toast.error(extractErrorMessage(e))
  } finally {
    saving.value = false
  }
}

async function toggleActive(cat: CategoryDTO) {
  try {
    await catalogApi.updateCategory(cat.id, { name: cat.name, active: !cat.active })
    toast.success('Обновлено')
    await load()
  } catch (e) {
    toast.error(extractErrorMessage(e))
  }
}
</script>

<template>
  <MarketplaceShell>
    <div class="mx-auto max-w-[820px]">
      <div class="flex items-center justify-between gap-3">
        <h1 class="mk-display text-2xl font-bold tracking-tight">Категории каталога</h1>
        <button type="button" class="rounded-[11px] bg-[#5B4BE0] px-4 py-2.5 text-sm font-medium text-white hover:opacity-90" @click="showCreate = true">
          Добавить категорию
        </button>
      </div>

      <div v-if="loading" class="py-16 text-center text-sm text-[#8D8A7E]">Загрузка…</div>
      <div v-else-if="error" class="mt-6 rounded-2xl border border-[#F3D3CE] bg-[#FBF0EE] px-6 py-10 text-center text-sm text-[#B3261E]">
        {{ error }}
      </div>
      <div v-else class="mt-6 flex flex-col gap-3">
        <template v-for="cat in categories" :key="cat.id">
          <div class="flex items-center justify-between rounded-2xl border border-[#E2DED2] bg-white p-5">
            <div>
              <p class="text-sm font-medium">{{ cat.name }}</p>
              <p class="mt-0.5 text-xs text-[#9B978A]">{{ cat.active ? 'Активна' : 'Скрыта' }}</p>
            </div>
            <button type="button" class="rounded-[11px] border border-[#E2DED2] px-3 py-2 text-sm hover:border-[#17160F]" @click="toggleActive(cat)">
              {{ cat.active ? 'Скрыть' : 'Включить' }}
            </button>
          </div>
          <div v-for="sub in cat.subcategories ?? []" :key="sub.id" class="ml-6 flex items-center justify-between rounded-2xl border border-[#E2DED2] bg-white p-4">
            <div>
              <p class="text-sm">{{ sub.name }}</p>
              <p class="mt-0.5 text-xs text-[#9B978A]">{{ sub.active ? 'Активна' : 'Скрыта' }}</p>
            </div>
            <button type="button" class="rounded-[11px] border border-[#E2DED2] px-3 py-2 text-sm hover:border-[#17160F]" @click="toggleActive(sub)">
              {{ sub.active ? 'Скрыть' : 'Включить' }}
            </button>
          </div>
        </template>
      </div>
    </div>

    <Teleport to="body">
      <div v-if="showCreate" class="fixed inset-0 z-40 flex items-center justify-center bg-[#17160F]/40 p-4" @click.self="showCreate = false">
        <div class="w-full max-w-md rounded-2xl border border-[#E2DED2] bg-white p-6 shadow-lg">
          <h2 class="mk-display text-lg font-semibold">Новая категория</h2>
          <form class="mt-4 flex flex-col gap-4" @submit.prevent="create">
            <label class="flex flex-col gap-1.5">
              <span class="font-mono text-[11px] uppercase tracking-[0.1em] text-[#9B978A]">Название</span>
              <input v-model="newName" required class="rounded-xl border border-[#E2DED2] px-4 py-2.5 text-base outline-none focus-visible:border-[#5B4BE0]" />
            </label>
            <label class="flex flex-col gap-1.5">
              <span class="font-mono text-[11px] uppercase tracking-[0.1em] text-[#9B978A]">Родительская категория (опционально)</span>
              <select v-model.number="newParentId" class="rounded-xl border border-[#E2DED2] px-4 py-2.5 text-base outline-none focus-visible:border-[#5B4BE0]">
                <option :value="null">— верхний уровень —</option>
                <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
              </select>
            </label>
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
