<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { RouterLink, useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { notificationsApi } from '@/api/notifications'

const auth = useAuthStore()
const router = useRouter()
const route = useRoute()

const isDark = ref(document.documentElement.classList.contains('dark'))
function toggleTheme() {
  isDark.value = !isDark.value
  document.documentElement.classList.toggle('dark', isDark.value)
  localStorage.setItem('cms.theme', isDark.value ? 'dark' : 'light')
}

const unreadCount = ref(0)
async function refreshUnread() {
  try {
    const { unreadCount: count } = await notificationsApi.unreadCount()
    unreadCount.value = count
  } catch {
    // notifications are non-critical; ignore errors here
  }
}
onMounted(() => {
  refreshUnread()
  window.setInterval(refreshUnread, 30000)
})

interface NavItem {
  label: string
  to: string
  roles?: string[]
}

const navItems: NavItem[] = [
  { label: 'Дашборд', to: '/' },
  { label: 'Тенанты', to: '/tenants', roles: ['SUPER_ADMIN'] },
  { label: 'Наряды', to: '/tasks', roles: ['TENANT_ADMIN', 'DISPATCHER'] },
  { label: 'Потребители', to: '/consumers', roles: ['TENANT_ADMIN', 'DISPATCHER'] },
  { label: 'Адреса', to: '/addresses', roles: ['TENANT_ADMIN', 'DISPATCHER'] },
  { label: 'Сотрудники', to: '/users', roles: ['TENANT_ADMIN', 'DISPATCHER'] },
]

function visible(item: NavItem) {
  return !item.roles || (auth.role && item.roles.includes(auth.role))
}

async function handleLogout() {
  await auth.logout()
  router.push({ name: 'login' })
}

const roleLabels: Record<string, string> = {
  SUPER_ADMIN: 'Супер-админ',
  TENANT_ADMIN: 'Администратор',
  DISPATCHER: 'Диспетчер',
  ELECTRICIAN: 'Инспектор',
}
</script>

<template>
  <div class="flex min-h-screen bg-surface dark:bg-ink">
    <aside class="hidden w-60 shrink-0 flex-col border-r border-line bg-white dark:border-graphite dark:bg-graphite md:flex">
      <div class="flex h-16 items-center gap-2.5 px-5">
        <svg viewBox="0 0 24 24" class="h-6 w-6 shrink-0 text-primary" aria-hidden="true">
          <circle cx="12" cy="12" r="9.25" fill="none" stroke="currentColor" stroke-width="1.5" />
          <circle cx="12" cy="12" r="2" fill="currentColor" />
          <path
            d="M12 4v2M12 18v2M20 12h-2M6 12H4M17.36 6.64l-1.42 1.42M8.06 15.94l-1.42 1.42M17.36 17.36l-1.42-1.42M8.06 8.06 6.64 6.64"
            stroke="currentColor"
            stroke-width="1.5"
            stroke-linecap="round"
          />
        </svg>
        <span class="font-display text-lg font-semibold text-ink dark:text-surface">ПриборСервис</span>
      </div>
      <nav class="flex-1 space-y-0.5 px-3">
        <RouterLink
          v-for="item in navItems.filter(visible)"
          :key="item.to"
          :to="item.to"
          class="block rounded-md border-l-[3px] border-transparent px-3 py-2 text-sm font-medium text-slate hover:bg-surface dark:text-mist dark:hover:bg-ink"
          :class="route.path === item.to && 'border-primary bg-[var(--color-primary-soft)] text-primary dark:bg-ink dark:text-primary'"
        >
          {{ item.label }}
        </RouterLink>
      </nav>
      <div class="border-t border-line p-4 text-xs text-mist dark:border-graphite">
        <p class="font-medium text-slate dark:text-mist">{{ auth.tenantName || 'Платформа' }}</p>
        <p>{{ auth.role ? roleLabels[auth.role] : '' }}</p>
      </div>
    </aside>

    <div class="flex min-h-screen flex-1 flex-col">
      <header class="flex h-16 items-center justify-between border-b border-line bg-white px-4 dark:border-graphite dark:bg-graphite md:px-6">
        <span class="font-display text-base font-semibold text-ink dark:text-surface md:hidden">ПриборСервис</span>
        <span class="hidden text-sm text-slate dark:text-mist md:block" />
        <div class="flex items-center gap-3">
          <button
            class="rounded-md p-2 text-slate hover:bg-surface dark:text-mist dark:hover:bg-ink"
            :aria-label="isDark ? 'Светлая тема' : 'Тёмная тема'"
            @click="toggleTheme"
          >
            <svg v-if="isDark" viewBox="0 0 24 24" class="h-4.5 w-4.5" aria-hidden="true">
              <circle cx="12" cy="12" r="4.5" fill="none" stroke="currentColor" stroke-width="1.5" />
              <path
                d="M12 3v2M12 19v2M21 12h-2M5 12H3M18.36 5.64l-1.42 1.42M7.06 16.94l-1.42 1.42M18.36 18.36l-1.42-1.42M7.06 7.06 5.64 5.64"
                stroke="currentColor"
                stroke-width="1.5"
                stroke-linecap="round"
              />
            </svg>
            <svg v-else viewBox="0 0 24 24" class="h-4.5 w-4.5" aria-hidden="true">
              <path
                d="M20 14.5A8.5 8.5 0 1 1 9.5 4a6.5 6.5 0 0 0 10.5 10.5Z"
                fill="none"
                stroke="currentColor"
                stroke-width="1.5"
                stroke-linejoin="round"
              />
            </svg>
          </button>
          <RouterLink
            to="/notifications"
            class="relative rounded-md p-2 text-slate hover:bg-surface dark:text-mist dark:hover:bg-ink"
            aria-label="Уведомления"
          >
            <svg viewBox="0 0 24 24" class="h-4.5 w-4.5" aria-hidden="true">
              <path
                d="M6 10a6 6 0 1 1 12 0c0 4 1.5 5.5 1.5 5.5h-15S6 14 6 10Z"
                fill="none"
                stroke="currentColor"
                stroke-width="1.5"
                stroke-linejoin="round"
              />
              <path d="M10 18.5a2 2 0 0 0 4 0" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" />
            </svg>
            <span
              v-if="unreadCount > 0"
              class="data-mono absolute -right-0.5 -top-0.5 flex h-4 min-w-4 items-center justify-center rounded-full bg-[var(--color-status-canceled)] px-1 text-[10px] font-semibold text-white"
            >
              {{ unreadCount > 9 ? '9+' : unreadCount }}
            </span>
          </RouterLink>
          <div class="hidden text-right sm:block">
            <p class="text-sm font-medium text-ink dark:text-surface">{{ auth.fullName }}</p>
          </div>
          <button
            class="rounded-md border border-line px-3 py-1.5 text-sm font-medium text-slate hover:bg-surface dark:border-graphite dark:text-mist dark:hover:bg-ink"
            @click="handleLogout"
          >
            Выйти
          </button>
        </div>
      </header>

      <main class="flex-1 px-4 py-6 md:px-8 md:py-8">
        <slot />
      </main>
    </div>
  </div>
</template>
