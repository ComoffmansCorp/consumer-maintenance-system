<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { tasksApi } from '@/api/tasks'
import type { TaskStatus } from '@/types'
import LoadingState from '@/components/ui/LoadingState.vue'
import StatusChip from '@/components/ui/StatusChip.vue'
import OdometerNumber from '@/components/ui/OdometerNumber.vue'

const auth = useAuthStore()
const loading = ref(true)
const counts = ref<Record<TaskStatus, number>>({
  PENDING: 0,
  IN_PROGRESS: 0,
  COMPLETED: 0,
  CANCELED: 0,
})

const statuses: TaskStatus[] = ['PENDING', 'IN_PROGRESS', 'COMPLETED', 'CANCELED']

async function load() {
  if (!auth.role || auth.role === 'SUPER_ADMIN' || auth.role === 'ELECTRICIAN') {
    loading.value = false
    return
  }
  loading.value = true
  try {
    const results = await Promise.all(
      statuses.map((status) => tasksApi.list({ status, pageSize: 1 })),
    )
    results.forEach((page, i) => {
      counts.value[statuses[i]] = page.totalItems
    })
  } finally {
    loading.value = false
  }
}

onMounted(load)
</script>

<template>
  <div class="flex flex-col gap-6">
    <div>
      <h1 class="font-display text-2xl font-semibold text-ink dark:text-surface">
        Здравствуйте, {{ auth.fullName }}
      </h1>
      <p class="mt-1 text-sm text-slate dark:text-mist">{{ auth.tenantName || 'Платформа' }}</p>
    </div>

    <div v-if="auth.role === 'SUPER_ADMIN'" class="card p-6">
      <p class="text-sm text-slate dark:text-mist">
        Управление тенантами платформы доступно в разделе
        <RouterLink to="/tenants" class="font-medium text-primary hover:underline">Тенанты</RouterLink>.
      </p>
    </div>

    <div v-else-if="auth.role === 'ELECTRICIAN'" class="card p-6">
      <p class="text-sm text-slate dark:text-mist">
        Эта админка — рабочее место диспетчера. Наряды и акты выполняются в мобильном
        приложении инспектора.
      </p>
    </div>

    <template v-else>
      <LoadingState v-if="loading" />
      <div v-else class="grid grid-cols-2 gap-4 lg:grid-cols-4">
        <RouterLink
          v-for="status in statuses"
          :key="status"
          :to="{ path: '/tasks', query: { status } }"
          class="card flex flex-col gap-3 p-5 transition-colors hover:border-primary"
        >
          <StatusChip :status="status" />
          <OdometerNumber :value="counts[status]" size="lg" />
        </RouterLink>
      </div>

      <div class="card flex flex-wrap gap-3 p-6">
        <RouterLink
          to="/tasks"
          class="rounded-md border border-line px-4 py-2 text-sm font-medium text-ink hover:bg-surface dark:border-graphite dark:text-surface dark:hover:bg-ink"
        >
          Все наряды
        </RouterLink>
        <RouterLink
          to="/addresses"
          class="rounded-md border border-line px-4 py-2 text-sm font-medium text-ink hover:bg-surface dark:border-graphite dark:text-surface dark:hover:bg-ink"
        >
          Адреса
        </RouterLink>
        <RouterLink
          to="/consumers"
          class="rounded-md border border-line px-4 py-2 text-sm font-medium text-ink hover:bg-surface dark:border-graphite dark:text-surface dark:hover:bg-ink"
        >
          Потребители
        </RouterLink>
      </div>
    </template>
  </div>
</template>
