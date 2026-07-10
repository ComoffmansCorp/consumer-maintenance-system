<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import { consumersApi } from '@/api/consumers'
import type { ConsumerDTO, ConsumerType } from '@/types'
import { extractErrorMessage } from '@/api/client'
import { useToastStore } from '@/stores/toast'
import LoadingState from '@/components/ui/LoadingState.vue'
import ErrorState from '@/components/ui/ErrorState.vue'
import EmptyState from '@/components/ui/EmptyState.vue'
import AppButton from '@/components/ui/AppButton.vue'
import AppModal from '@/components/ui/AppModal.vue'
import AppInput from '@/components/ui/AppInput.vue'
import AppSelect from '@/components/ui/AppSelect.vue'
import Pagination from '@/components/ui/Pagination.vue'

const toast = useToastStore()

const typeLabels: Record<ConsumerType, string> = {
  COMMERCIAL: 'Коммерческий',
  GOVERNMENT: 'Государственный',
  RESIDENTIAL: 'Жилой',
}

const items = ref<ConsumerDTO[]>([])
const loading = ref(true)
const error = ref('')
const page = ref(1)
const totalPages = ref(1)
const totalItems = ref(0)
const search = ref('')

async function load() {
  loading.value = true
  error.value = ''
  try {
    const result = await consumersApi.list({ page: page.value, pageSize: 20, search: search.value })
    items.value = result.items
    totalPages.value = result.totalPages
    totalItems.value = result.totalItems
  } catch (e) {
    error.value = extractErrorMessage(e)
  } finally {
    loading.value = false
  }
}
onMounted(load)
watch(page, load)
let searchTimer: ReturnType<typeof setTimeout>
watch(search, () => {
  clearTimeout(searchTimer)
  searchTimer = setTimeout(() => {
    page.value = 1
    load()
  }, 300)
})

const showModal = ref(false)
const editing = ref<ConsumerDTO | null>(null)
const name = ref('')
const type = ref<ConsumerType>('COMMERCIAL')
const description = ref('')
const saving = ref(false)

function openCreate() {
  editing.value = null
  name.value = ''
  type.value = 'COMMERCIAL'
  description.value = ''
  showModal.value = true
}

function openEdit(c: ConsumerDTO) {
  editing.value = c
  name.value = c.name
  type.value = c.type
  description.value = c.description ?? ''
  showModal.value = true
}

async function save() {
  saving.value = true
  try {
    if (editing.value) {
      await consumersApi.update(editing.value.id, { name: name.value, type: type.value, description: description.value })
      toast.success('Потребитель обновлён')
    } else {
      await consumersApi.create({ name: name.value, type: type.value, description: description.value })
      toast.success('Потребитель создан')
    }
    showModal.value = false
    await load()
  } catch (e) {
    toast.error(extractErrorMessage(e))
  } finally {
    saving.value = false
  }
}

async function deactivate(c: ConsumerDTO) {
  try {
    await consumersApi.deactivate(c.id)
    toast.success('Потребитель деактивирован')
    await load()
  } catch (e) {
    toast.error(extractErrorMessage(e))
  }
}
</script>

<template>
  <div class="flex flex-col gap-6">
    <div class="flex flex-wrap items-center justify-between gap-3">
      <h1 class="font-display text-2xl font-semibold text-ink dark:text-surface">Потребители</h1>
      <AppButton @click="openCreate">Новый потребитель</AppButton>
    </div>

    <AppInput v-model="search" placeholder="Поиск по названию…" class="max-w-xs" />

    <LoadingState v-if="loading" />
    <ErrorState v-else-if="error" :message="error" @retry="load" />
    <EmptyState
      v-else-if="items.length === 0"
      title="Потребителей пока нет"
      description="Добавьте компанию или жильца, чтобы привязывать к ним адреса"
    >
      <template #action>
        <AppButton @click="openCreate">Новый потребитель</AppButton>
      </template>
    </EmptyState>
    <div v-else class="card overflow-hidden">
      <table class="w-full text-sm">
        <thead>
          <tr class="border-b border-line text-left text-xs uppercase tracking-wide text-mist dark:border-graphite">
            <th class="px-4 py-3 font-medium">Название</th>
            <th class="px-4 py-3 font-medium">Тип</th>
            <th class="px-4 py-3 font-medium">Описание</th>
            <th class="px-4 py-3"></th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="c in items" :key="c.id" class="border-b border-line last:border-0 dark:border-graphite">
            <td class="px-4 py-3 font-medium text-ink dark:text-surface">{{ c.name }}</td>
            <td class="px-4 py-3 text-slate dark:text-mist">{{ typeLabels[c.type] }}</td>
            <td class="px-4 py-3 text-slate dark:text-mist">{{ c.description || '—' }}</td>
            <td class="px-4 py-3 text-right">
              <button class="text-sm font-medium text-primary hover:underline" @click="openEdit(c)">
                Изменить
              </button>
              <button
                class="ml-3 text-sm font-medium hover:underline"
                style="color: var(--color-status-canceled)"
                @click="deactivate(c)"
              >
                Деактивировать
              </button>
            </td>
          </tr>
        </tbody>
      </table>
      <div class="px-4">
        <Pagination v-model:page="page" :total-pages="totalPages" :total-items="totalItems" />
      </div>
    </div>

    <AppModal v-if="showModal" :title="editing ? 'Изменить потребителя' : 'Новый потребитель'" @close="showModal = false">
      <form class="flex flex-col gap-4" @submit.prevent="save">
        <AppInput v-model="name" label="Название" required />
        <AppSelect
          v-model="type"
          label="Тип"
          :options="[
            { value: 'COMMERCIAL', label: 'Коммерческий' },
            { value: 'GOVERNMENT', label: 'Государственный' },
            { value: 'RESIDENTIAL', label: 'Жилой' },
          ]"
        />
        <AppInput v-model="description" label="Описание" />
        <div class="flex justify-end gap-2">
          <AppButton variant="secondary" type="button" @click="showModal = false">Отмена</AppButton>
          <AppButton type="submit" :loading="saving">Сохранить</AppButton>
        </div>
      </form>
    </AppModal>
  </div>
</template>
