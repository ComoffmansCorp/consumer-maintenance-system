<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'
import { tasksApi } from '@/api/tasks'
import { actsApi } from '@/api/acts'
import { usersApi } from '@/api/auth'
import type { TaskDTO, UserDTO } from '@/types'
import { extractErrorMessage } from '@/api/client'
import { useToastStore } from '@/stores/toast'
import { useAuthStore } from '@/stores/auth'
import LoadingState from '@/components/ui/LoadingState.vue'
import ErrorState from '@/components/ui/ErrorState.vue'
import AppButton from '@/components/ui/AppButton.vue'
import AppModal from '@/components/ui/AppModal.vue'
import AppSelect from '@/components/ui/AppSelect.vue'
import AppInput from '@/components/ui/AppInput.vue'
import StatusChip from '@/components/ui/StatusChip.vue'

const props = defineProps<{ id: string }>()
const toast = useToastStore()
const auth = useAuthStore()
const router = useRouter()

const task = ref<TaskDTO | null>(null)
const loading = ref(true)
const error = ref('')
const actId = ref<number | null>(null)
const actExists = ref(false)
const working = ref(false)

const typeLabels = { INSPECTION: 'Осмотр', REPLACEMENT: 'Замена' }

async function load() {
  loading.value = true
  error.value = ''
  try {
    task.value = await tasksApi.get(Number(props.id))
    await checkAct()
  } catch (e) {
    error.value = extractErrorMessage(e)
  } finally {
    loading.value = false
  }
}

async function checkAct() {
  if (!task.value) return
  try {
    if (task.value.type === 'INSPECTION') {
      const act = await actsApi.getInspectionByTask(task.value.id)
      actId.value = act.id
      actExists.value = true
    } else {
      const act = await actsApi.getReplacementByTask(task.value.id)
      actId.value = act.id
      actExists.value = true
    }
  } catch (e) {
    if (axios.isAxiosError(e) && e.response?.status === 404) {
      actExists.value = false
    }
  }
}

onMounted(load)

const isDispatcher = computed(() => auth.role === 'TENANT_ADMIN' || auth.role === 'DISPATCHER')

function openAct() {
  if (!actId.value || !task.value) return
  if (task.value.type === 'INSPECTION') router.push({ name: 'inspection-act', params: { id: actId.value } })
  else router.push({ name: 'replacement-act', params: { id: actId.value } })
}

const showAssign = ref(false)
const electricians = ref<UserDTO[]>([])
const assigneeId = ref('')

async function openAssign() {
  if (electricians.value.length === 0) electricians.value = await usersApi.list('ELECTRICIAN')
  assigneeId.value = task.value?.assigneeId ? String(task.value.assigneeId) : ''
  showAssign.value = true
}

async function assign() {
  if (!task.value || !assigneeId.value) return
  working.value = true
  try {
    task.value = await tasksApi.assign(task.value.id, Number(assigneeId.value))
    toast.success('Инспектор назначен')
    showAssign.value = false
  } catch (e) {
    toast.error(extractErrorMessage(e))
  } finally {
    working.value = false
  }
}

const showCancel = ref(false)
const cancelReason = ref('')

async function cancel() {
  if (!task.value || !cancelReason.value.trim()) return
  working.value = true
  try {
    task.value = await tasksApi.cancel(task.value.id, cancelReason.value)
    toast.success('Наряд отменён')
    showCancel.value = false
    cancelReason.value = ''
  } catch (e) {
    toast.error(extractErrorMessage(e))
  } finally {
    working.value = false
  }
}

function formatDate(iso?: string) {
  if (!iso) return '—'
  return new Date(iso).toLocaleDateString('ru-RU')
}
</script>

<template>
  <div class="flex flex-col gap-6">
    <LoadingState v-if="loading" />
    <ErrorState v-else-if="error" :message="error" @retry="load" />

    <template v-else-if="task">
      <div class="flex flex-wrap items-start justify-between gap-3">
        <div>
          <p class="text-sm text-slate dark:text-mist">{{ typeLabels[task.type] }} · наряд №{{ task.id }}</p>
          <h1 class="font-display text-2xl font-semibold text-ink dark:text-surface">{{ task.addressLabel }}</h1>
        </div>
        <StatusChip :status="task.status" />
      </div>

      <div class="card grid grid-cols-1 gap-4 p-6 sm:grid-cols-2">
        <div>
          <p class="text-xs uppercase tracking-wide text-mist">Инспектор</p>
          <p class="mt-1 text-sm text-ink dark:text-surface">{{ task.assigneeName || 'Не назначен' }}</p>
        </div>
        <div>
          <p class="text-xs uppercase tracking-wide text-mist">Срок</p>
          <p class="data-mono mt-1 text-sm text-ink dark:text-surface">{{ formatDate(task.dueDate) }}</p>
        </div>
        <div v-if="task.cancelReason">
          <p class="text-xs uppercase tracking-wide text-mist">Причина отмены</p>
          <p class="mt-1 text-sm text-ink dark:text-surface">{{ task.cancelReason }}</p>
        </div>
      </div>

      <div class="card flex flex-wrap items-center gap-3 p-6">
        <AppButton v-if="isDispatcher && task.status === 'PENDING'" variant="secondary" @click="openAssign">
          {{ task.assigneeId ? 'Переназначить' : 'Назначить инспектора' }}
        </AppButton>
        <AppButton
          v-if="isDispatcher && (task.status === 'PENDING' || task.status === 'IN_PROGRESS')"
          variant="danger"
          @click="showCancel = true"
        >
          Отменить наряд
        </AppButton>

        <AppButton v-if="actExists" variant="secondary" @click="openAct">Открыть акт</AppButton>
        <span
          v-else-if="task.status === 'IN_PROGRESS'"
          class="text-sm text-slate dark:text-mist"
        >
          Акт ещё не заполнен — инспектор оформит его в мобильном приложении.
        </span>
        <span v-else-if="task.status === 'PENDING'" class="text-sm text-slate dark:text-mist">
          Акт появится, когда инспектор возьмёт наряд в работу.
        </span>
      </div>
    </template>

    <AppModal v-if="showAssign" title="Назначить инспектора" @close="showAssign = false">
      <form class="flex flex-col gap-4" @submit.prevent="assign">
        <AppSelect
          v-model="assigneeId"
          label="Инспектор"
          required
          placeholder="Выберите инспектора"
          :options="electricians.map((e) => ({ value: e.id, label: e.fullName }))"
        />
        <div class="flex justify-end gap-2">
          <AppButton variant="secondary" type="button" @click="showAssign = false">Отмена</AppButton>
          <AppButton type="submit" :loading="working">Назначить</AppButton>
        </div>
      </form>
    </AppModal>

    <AppModal v-if="showCancel" title="Отменить наряд" @close="showCancel = false">
      <form class="flex flex-col gap-4" @submit.prevent="cancel">
        <AppInput v-model="cancelReason" label="Причина отмены" required />
        <div class="flex justify-end gap-2">
          <AppButton variant="secondary" type="button" @click="showCancel = false">Назад</AppButton>
          <AppButton variant="danger" type="submit" :loading="working">Отменить наряд</AppButton>
        </div>
      </form>
    </AppModal>
  </div>
</template>
