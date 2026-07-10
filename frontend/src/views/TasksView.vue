<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { tasksApi } from '@/api/tasks'
import { addressesApi } from '@/api/addresses'
import { usersApi } from '@/api/auth'
import type { AddressDTO, TaskDTO, TaskStatus, TaskType, UserDTO } from '@/types'
import { extractErrorMessage } from '@/api/client'
import { useToastStore } from '@/stores/toast'
import { useAuthStore } from '@/stores/auth'
import LoadingState from '@/components/ui/LoadingState.vue'
import ErrorState from '@/components/ui/ErrorState.vue'
import EmptyState from '@/components/ui/EmptyState.vue'
import AppButton from '@/components/ui/AppButton.vue'
import AppModal from '@/components/ui/AppModal.vue'
import AppInput from '@/components/ui/AppInput.vue'
import AppSelect from '@/components/ui/AppSelect.vue'
import Pagination from '@/components/ui/Pagination.vue'
import StatusChip from '@/components/ui/StatusChip.vue'

const toast = useToastStore()
const auth = useAuthStore()
const route = useRoute()
const router = useRouter()

const typeLabels: Record<TaskType, string> = { INSPECTION: 'Осмотр', REPLACEMENT: 'Замена' }

const items = ref<TaskDTO[]>([])
const loading = ref(true)
const error = ref('')
const page = ref(1)
const totalPages = ref(1)
const totalItems = ref(0)
const statusFilter = ref<string>((route.query.status as string) || '')
const typeFilter = ref<string>('')

async function load() {
  loading.value = true
  error.value = ''
  try {
    const result = await tasksApi.list({
      page: page.value,
      pageSize: 20,
      status: (statusFilter.value || undefined) as TaskStatus | undefined,
      type: (typeFilter.value || undefined) as TaskType | undefined,
    })
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
watch([statusFilter, typeFilter], () => {
  page.value = 1
  load()
})

const canCreate = ref(auth.role === 'TENANT_ADMIN' || auth.role === 'DISPATCHER')
const addresses = ref<AddressDTO[]>([])
const electricians = ref<UserDTO[]>([])
const showCreate = ref(false)
const type = ref<TaskType>('INSPECTION')
const addressId = ref('')
const dueDate = ref('')
const assigneeId = ref('')
const saving = ref(false)

async function openCreate() {
  type.value = 'INSPECTION'
  addressId.value = ''
  dueDate.value = ''
  assigneeId.value = ''
  showCreate.value = true
  if (addresses.value.length === 0) {
    const result = await addressesApi.list({ pageSize: 100 })
    addresses.value = result.items
  }
  if (electricians.value.length === 0) {
    electricians.value = await usersApi.list('ELECTRICIAN')
  }
}

async function createTask() {
  saving.value = true
  try {
    await tasksApi.create({
      type: type.value,
      addressId: Number(addressId.value),
      dueDate: dueDate.value || null,
      assigneeId: assigneeId.value ? Number(assigneeId.value) : null,
    })
    toast.success('Наряд создан')
    showCreate.value = false
    await load()
  } catch (e) {
    toast.error(extractErrorMessage(e))
  } finally {
    saving.value = false
  }
}

function openTask(t: TaskDTO) {
  router.push({ name: 'task-detail', params: { id: t.id } })
}

function formatDate(iso?: string) {
  if (!iso) return '—'
  return new Date(iso).toLocaleDateString('ru-RU')
}
</script>

<template>
  <div class="flex flex-col gap-6">
    <div class="flex flex-wrap items-center justify-between gap-3">
      <h1 class="font-display text-2xl font-semibold text-ink dark:text-surface">Наряды</h1>
      <AppButton v-if="canCreate" @click="openCreate">Новый наряд</AppButton>
    </div>

    <div class="flex flex-wrap gap-3">
      <AppSelect
        v-model="statusFilter"
        placeholder="Все статусы"
        :options="[
          { value: 'PENDING', label: 'Ожидает' },
          { value: 'IN_PROGRESS', label: 'В работе' },
          { value: 'COMPLETED', label: 'Выполнен' },
          { value: 'CANCELED', label: 'Отменён' },
        ]"
      />
      <AppSelect
        v-model="typeFilter"
        placeholder="Все типы"
        :options="[
          { value: 'INSPECTION', label: 'Осмотр' },
          { value: 'REPLACEMENT', label: 'Замена' },
        ]"
      />
    </div>

    <LoadingState v-if="loading" />
    <ErrorState v-else-if="error" :message="error" @retry="load" />
    <EmptyState
      v-else-if="items.length === 0"
      title="Нарядов не найдено"
      description="Измените фильтры или создайте новый наряд"
    >
      <template #action>
        <AppButton v-if="canCreate" @click="openCreate">Новый наряд</AppButton>
      </template>
    </EmptyState>
    <div v-else class="card overflow-hidden">
      <table class="w-full text-sm">
        <thead>
          <tr class="border-b border-line text-left text-xs uppercase tracking-wide text-mist dark:border-graphite">
            <th class="px-4 py-3 font-medium">Адрес</th>
            <th class="px-4 py-3 font-medium">Тип</th>
            <th class="px-4 py-3 font-medium">Инспектор</th>
            <th class="px-4 py-3 font-medium">Срок</th>
            <th class="px-4 py-3 font-medium">Статус</th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="t in items"
            :key="t.id"
            class="cursor-pointer border-b border-line last:border-0 hover:bg-surface dark:border-graphite dark:hover:bg-ink"
            @click="openTask(t)"
          >
            <td class="px-4 py-3 font-medium text-ink dark:text-surface">{{ t.addressLabel }}</td>
            <td class="px-4 py-3 text-slate dark:text-mist">{{ typeLabels[t.type] }}</td>
            <td class="px-4 py-3 text-slate dark:text-mist">{{ t.assigneeName || '—' }}</td>
            <td class="data-mono px-4 py-3 text-slate dark:text-mist">{{ formatDate(t.dueDate) }}</td>
            <td class="px-4 py-3"><StatusChip :status="t.status" /></td>
          </tr>
        </tbody>
      </table>
      <div class="px-4">
        <Pagination v-model:page="page" :total-pages="totalPages" :total-items="totalItems" />
      </div>
    </div>

    <AppModal v-if="showCreate" title="Новый наряд" @close="showCreate = false">
      <form class="flex flex-col gap-4" @submit.prevent="createTask">
        <AppSelect
          v-model="type"
          label="Тип"
          :options="[
            { value: 'INSPECTION', label: 'Осмотр' },
            { value: 'REPLACEMENT', label: 'Замена' },
          ]"
        />
        <AppSelect
          v-model="addressId"
          label="Адрес"
          required
          placeholder="Выберите адрес"
          :options="addresses.map((a) => ({ value: a.id, label: `${a.street}, ${a.house}` }))"
        />
        <AppInput v-model="dueDate" type="date" label="Срок выполнения" />
        <AppSelect
          v-model="assigneeId"
          label="Инспектор"
          placeholder="Назначить позже"
          :options="electricians.map((e) => ({ value: e.id, label: e.fullName }))"
        />
        <div class="flex justify-end gap-2">
          <AppButton variant="secondary" type="button" @click="showCreate = false">Отмена</AppButton>
          <AppButton type="submit" :loading="saving">Создать</AppButton>
        </div>
      </form>
    </AppModal>
  </div>
</template>
