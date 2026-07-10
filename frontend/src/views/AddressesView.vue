<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import { addressesApi } from '@/api/addresses'
import { consumersApi } from '@/api/consumers'
import type { AddressDTO, ConsumerDTO } from '@/types'
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

const items = ref<AddressDTO[]>([])
const consumers = ref<ConsumerDTO[]>([])
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
    const result = await addressesApi.list({ page: page.value, pageSize: 20, search: search.value })
    items.value = result.items
    totalPages.value = result.totalPages
    totalItems.value = result.totalItems
  } catch (e) {
    error.value = extractErrorMessage(e)
  } finally {
    loading.value = false
  }
}

async function loadConsumers() {
  const result = await consumersApi.list({ pageSize: 100 })
  consumers.value = result.items
}

onMounted(() => {
  load()
  loadConsumers()
})
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
const editing = ref<AddressDTO | null>(null)
const street = ref('')
const house = ref('')
const building = ref('')
const apartment = ref('')
const consumerId = ref<string>('')
const saving = ref(false)

function openCreate() {
  editing.value = null
  street.value = ''
  house.value = ''
  building.value = ''
  apartment.value = ''
  consumerId.value = ''
  showModal.value = true
}

function openEdit(a: AddressDTO) {
  editing.value = a
  street.value = a.street
  house.value = a.house
  building.value = a.building ?? ''
  apartment.value = a.apartment ?? ''
  consumerId.value = a.consumerId ? String(a.consumerId) : ''
  showModal.value = true
}

async function save() {
  saving.value = true
  const payload = {
    street: street.value,
    house: house.value,
    building: building.value,
    apartment: apartment.value,
    consumerId: consumerId.value ? Number(consumerId.value) : null,
  }
  try {
    if (editing.value) {
      await addressesApi.update(editing.value.id, payload)
      toast.success('Адрес обновлён')
    } else {
      await addressesApi.create(payload)
      toast.success('Адрес добавлен')
    }
    showModal.value = false
    await load()
  } catch (e) {
    toast.error(extractErrorMessage(e))
  } finally {
    saving.value = false
  }
}
</script>

<template>
  <div class="flex flex-col gap-6">
    <div class="flex flex-wrap items-center justify-between gap-3">
      <h1 class="font-display text-2xl font-semibold text-ink dark:text-surface">Адреса</h1>
      <AppButton @click="openCreate">Новый адрес</AppButton>
    </div>

    <AppInput v-model="search" placeholder="Поиск по улице или дому…" class="max-w-xs" />

    <LoadingState v-if="loading" />
    <ErrorState v-else-if="error" :message="error" @retry="load" />
    <EmptyState
      v-else-if="items.length === 0"
      title="Адресов пока нет"
      description="Добавьте адрес, чтобы создавать по нему наряды"
    >
      <template #action>
        <AppButton @click="openCreate">Новый адрес</AppButton>
      </template>
    </EmptyState>
    <div v-else class="card overflow-hidden">
      <table class="w-full text-sm">
        <thead>
          <tr class="border-b border-line text-left text-xs uppercase tracking-wide text-mist dark:border-graphite">
            <th class="px-4 py-3 font-medium">Адрес</th>
            <th class="px-4 py-3 font-medium">Потребитель</th>
            <th class="px-4 py-3"></th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="a in items" :key="a.id" class="border-b border-line last:border-0 dark:border-graphite">
            <td class="px-4 py-3 font-medium text-ink dark:text-surface">
              {{ a.street }}, {{ a.house }}<span v-if="a.building"> корп. {{ a.building }}</span
              ><span v-if="a.apartment">, кв. {{ a.apartment }}</span>
            </td>
            <td class="px-4 py-3 text-slate dark:text-mist">{{ a.consumerName || '—' }}</td>
            <td class="px-4 py-3 text-right">
              <button class="text-sm font-medium text-primary hover:underline" @click="openEdit(a)">
                Изменить
              </button>
            </td>
          </tr>
        </tbody>
      </table>
      <div class="px-4">
        <Pagination v-model:page="page" :total-pages="totalPages" :total-items="totalItems" />
      </div>
    </div>

    <AppModal v-if="showModal" :title="editing ? 'Изменить адрес' : 'Новый адрес'" @close="showModal = false">
      <form class="flex flex-col gap-4" @submit.prevent="save">
        <AppInput v-model="street" label="Улица" required />
        <div class="grid grid-cols-2 gap-4">
          <AppInput v-model="house" label="Дом" required />
          <AppInput v-model="building" label="Корпус" />
        </div>
        <AppInput v-model="apartment" label="Квартира / офис" />
        <AppSelect
          v-model="consumerId"
          label="Потребитель"
          placeholder="Не выбран"
          :options="consumers.map((c) => ({ value: c.id, label: c.name }))"
        />
        <div class="flex justify-end gap-2">
          <AppButton variant="secondary" type="button" @click="showModal = false">Отмена</AppButton>
          <AppButton type="submit" :loading="saving">Сохранить</AppButton>
        </div>
      </form>
    </AppModal>
  </div>
</template>
