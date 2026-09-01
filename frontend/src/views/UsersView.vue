<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { usersApi } from '@/api/auth'
import type { Role, UserDTO } from '@/types'
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

const toast = useToastStore()
const auth = useAuthStore()

const roleLabels: Record<Role, string> = {
  SUPER_ADMIN: 'Супер-админ',
  TENANT_ADMIN: 'Администратор',
  DISPATCHER: 'Диспетчер',
  ELECTRICIAN: 'Инспектор',
  // MASTER/CLIENT never appear here -- tenant staff management only ever
  // deals with DISPATCHER/ELECTRICIAN (enforced server-side too), but the
  // Role union is shared with the marketplace so this map must stay total.
  MASTER: 'Мастер',
  CLIENT: 'Клиент',
}

const items = ref<UserDTO[]>([])
const loading = ref(true)
const error = ref('')

async function load() {
  loading.value = true
  error.value = ''
  try {
    items.value = await usersApi.list()
  } catch (e) {
    error.value = extractErrorMessage(e)
  } finally {
    loading.value = false
  }
}
onMounted(load)

const showModal = ref(false)
const username = ref('')
const password = ref('')
const fullName = ref('')
const role = ref<Role>('ELECTRICIAN')
const saving = ref(false)

function openCreate() {
  username.value = ''
  password.value = ''
  fullName.value = ''
  role.value = 'ELECTRICIAN'
  showModal.value = true
}

async function save() {
  saving.value = true
  try {
    await usersApi.create({
      username: username.value,
      password: password.value,
      fullName: fullName.value,
      role: role.value,
    })
    toast.success('Сотрудник добавлен')
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
      <h1 class="font-display text-2xl font-semibold text-ink dark:text-surface">Сотрудники</h1>
      <AppButton v-if="auth.role === 'TENANT_ADMIN'" @click="openCreate">Новый сотрудник</AppButton>
    </div>

    <LoadingState v-if="loading" />
    <ErrorState v-else-if="error" :message="error" @retry="load" />
    <EmptyState v-else-if="items.length === 0" title="Сотрудников пока нет" />
    <div v-else class="card overflow-hidden">
      <table class="w-full text-sm">
        <thead>
          <tr class="border-b border-line text-left text-xs uppercase tracking-wide text-mist dark:border-graphite">
            <th class="px-4 py-3 font-medium">Имя</th>
            <th class="px-4 py-3 font-medium">Логин</th>
            <th class="px-4 py-3 font-medium">Роль</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="u in items" :key="u.id" class="border-b border-line last:border-0 dark:border-graphite">
            <td class="px-4 py-3 font-medium text-ink dark:text-surface">{{ u.fullName }}</td>
            <td class="data-mono px-4 py-3 text-slate dark:text-mist">{{ u.username }}</td>
            <td class="px-4 py-3 text-slate dark:text-mist">{{ roleLabels[u.role] }}</td>
          </tr>
        </tbody>
      </table>
    </div>

    <AppModal v-if="showModal" title="Новый сотрудник" @close="showModal = false">
      <form class="flex flex-col gap-4" @submit.prevent="save">
        <AppInput v-model="fullName" label="Полное имя" required />
        <AppInput v-model="username" label="Логин" required />
        <AppInput v-model="password" type="password" label="Пароль" required />
        <AppSelect
          v-model="role"
          label="Роль"
          :options="[
            { value: 'DISPATCHER', label: 'Диспетчер' },
            { value: 'ELECTRICIAN', label: 'Инспектор' },
          ]"
        />
        <div class="flex justify-end gap-2">
          <AppButton variant="secondary" type="button" @click="showModal = false">Отмена</AppButton>
          <AppButton type="submit" :loading="saving">Создать</AppButton>
        </div>
      </form>
    </AppModal>
  </div>
</template>
