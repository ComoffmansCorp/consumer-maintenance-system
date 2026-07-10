<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { tenantsApi } from '@/api/tenants'
import type { TenantDTO, TenantPlan } from '@/types'
import { extractErrorMessage } from '@/api/client'
import { useToastStore } from '@/stores/toast'
import LoadingState from '@/components/ui/LoadingState.vue'
import ErrorState from '@/components/ui/ErrorState.vue'
import EmptyState from '@/components/ui/EmptyState.vue'
import AppButton from '@/components/ui/AppButton.vue'
import AppModal from '@/components/ui/AppModal.vue'
import AppInput from '@/components/ui/AppInput.vue'
import AppSelect from '@/components/ui/AppSelect.vue'

const toast = useToastStore()

const tenants = ref<TenantDTO[]>([])
const loading = ref(true)
const error = ref('')

async function load() {
  loading.value = true
  error.value = ''
  try {
    tenants.value = await tenantsApi.list()
  } catch (e) {
    error.value = extractErrorMessage(e)
  } finally {
    loading.value = false
  }
}
onMounted(load)

const showCreate = ref(false)
const name = ref('')
const code = ref('')
const plan = ref<TenantPlan>('FREE')
const saving = ref(false)

async function createTenant() {
  saving.value = true
  try {
    await tenantsApi.create({ name: name.value, code: code.value, plan: plan.value })
    toast.success('Тенант создан')
    showCreate.value = false
    name.value = ''
    code.value = ''
    plan.value = 'FREE'
    await load()
  } catch (e) {
    toast.error(extractErrorMessage(e))
  } finally {
    saving.value = false
  }
}

async function changePlan(tenant: TenantDTO, plan: string) {
  try {
    await tenantsApi.updatePlan(tenant.id, plan as TenantPlan)
    toast.success('Тариф обновлён')
    await load()
  } catch (e) {
    toast.error(extractErrorMessage(e))
  }
}

async function deactivate(tenant: TenantDTO) {
  try {
    await tenantsApi.deactivate(tenant.id)
    toast.success('Тенант деактивирован')
    await load()
  } catch (e) {
    toast.error(extractErrorMessage(e))
  }
}
</script>

<template>
  <div class="flex flex-col gap-6">
    <div class="flex items-center justify-between">
      <h1 class="font-display text-2xl font-semibold text-ink dark:text-surface">Тенанты платформы</h1>
      <AppButton @click="showCreate = true">Новый тенант</AppButton>
    </div>

    <LoadingState v-if="loading" />
    <ErrorState v-else-if="error" :message="error" @retry="load" />
    <EmptyState v-else-if="tenants.length === 0" title="Тенантов пока нет" description="Создайте первый, чтобы начать" />
    <div v-else class="card overflow-hidden">
      <table class="w-full text-sm">
        <thead>
          <tr class="border-b border-line text-left text-xs uppercase tracking-wide text-mist dark:border-graphite">
            <th class="px-4 py-3 font-medium">Название</th>
            <th class="px-4 py-3 font-medium">Код</th>
            <th class="px-4 py-3 font-medium">Тариф</th>
            <th class="px-4 py-3 font-medium">Статус</th>
            <th class="px-4 py-3"></th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="t in tenants" :key="t.id" class="border-b border-line last:border-0 dark:border-graphite">
            <td class="px-4 py-3 font-medium text-ink dark:text-surface">{{ t.name }}</td>
            <td class="data-mono px-4 py-3 text-slate dark:text-mist">{{ t.code }}</td>
            <td class="px-4 py-3">
              <select
                :value="t.plan"
                class="rounded-md border border-line bg-white px-2 py-1 text-sm dark:border-graphite dark:bg-ink dark:text-surface"
                @change="changePlan(t, ($event.target as HTMLSelectElement).value)"
              >
                <option value="FREE">Free</option>
                <option value="BUSINESS">Business</option>
                <option value="ENTERPRISE">Enterprise</option>
              </select>
            </td>
            <td class="px-4 py-3">
              <span
                class="inline-flex items-center gap-1.5 text-xs font-medium"
                :style="{ color: t.active ? 'var(--color-status-completed)' : 'var(--color-mist)' }"
              >
                <span class="h-1.5 w-1.5 rounded-full" :style="{ backgroundColor: t.active ? 'var(--color-status-completed)' : 'var(--color-mist)' }" />
                {{ t.active ? 'Активен' : 'Деактивирован' }}
              </span>
            </td>
            <td class="px-4 py-3 text-right">
              <button
                v-if="t.active"
                class="text-sm font-medium hover:underline"
                style="color: var(--color-status-canceled)"
                @click="deactivate(t)"
              >
                Деактивировать
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <AppModal v-if="showCreate" title="Новый тенант" @close="showCreate = false">
      <form class="flex flex-col gap-4" @submit.prevent="createTenant">
        <AppInput v-model="name" label="Название компании" required />
        <AppInput v-model="code" label="Код" required />
        <AppSelect
          v-model="plan"
          label="Тариф"
          :options="[
            { value: 'FREE', label: 'Free' },
            { value: 'BUSINESS', label: 'Business' },
            { value: 'ENTERPRISE', label: 'Enterprise' },
          ]"
        />
        <div class="flex justify-end gap-2">
          <AppButton variant="secondary" type="button" @click="showCreate = false">Отмена</AppButton>
          <AppButton type="submit" :loading="saving">Создать</AppButton>
        </div>
      </form>
    </AppModal>
  </div>
</template>
