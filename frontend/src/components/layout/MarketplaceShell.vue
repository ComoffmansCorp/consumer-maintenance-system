<script setup lang="ts">
import { onUnmounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { chatApi } from '@/api/marketplace'

const auth = useAuthStore()
const router = useRouter()

async function handleLogout() {
  await auth.logout()
  router.push({ name: 'marketplace' })
}

// Polled, not pushed -- there's no websocket layer here (chat itself is
// poll-based too, see MarketplaceRequestDetailView), so "live" means "at
// most 20s stale", which is fine for a badge count.
const unreadCount = ref(0)
let unreadTimer: ReturnType<typeof setInterval> | null = null

async function pollUnread() {
  if (!auth.isAuthenticated) {
    unreadCount.value = 0
    return
  }
  try {
    unreadCount.value = await chatApi.unreadCount()
  } catch {
    // Silent, same reasoning as chat message polling -- a badge that
    // occasionally misses a refresh isn't worth an error toast.
  }
}

function startUnreadPolling() {
  stopUnreadPolling()
  pollUnread()
  unreadTimer = setInterval(pollUnread, 20000)
}

function stopUnreadPolling() {
  if (unreadTimer) {
    clearInterval(unreadTimer)
    unreadTimer = null
  }
}

watch(
  () => auth.isAuthenticated,
  (isAuthed) => (isAuthed ? startUnreadPolling() : (stopUnreadPolling(), (unreadCount.value = 0))),
  { immediate: true },
)
onUnmounted(stopUnreadPolling)
</script>

<template>
  <div class="mk-home flex min-h-screen flex-col bg-[#FBFAF7] text-[#17160F]">
    <header class="sticky top-0 z-20 border-b border-[#E7E3D9] bg-[#FBFAF7]/92 backdrop-blur-md">
      <div class="mx-auto flex max-w-[1280px] items-center gap-8 px-5 py-4 md:px-10">
        <RouterLink :to="{ name: 'marketplace' }" class="flex items-center gap-2.5">
          <span class="flex h-[30px] w-[30px] items-center justify-center rounded-[9px] bg-[#5B4BE0]">
            <svg viewBox="0 0 48 48" class="h-[18px] w-[18px]">
              <path fill-rule="evenodd" clip-rule="evenodd" fill="#FFFFFF"
                d="M24,9 L37.04,16.5 L37.04,31.5 L24,39 L10.96,31.5 L10.96,16.5 Z
                   M30.19,24 a6.19,6.19 0 1,0 -12.38,0 a6.19,6.19 0 1,0 12.38,0 Z"/>
            </svg>
          </span>
          <span class="mk-display text-[17px] font-medium tracking-tight">Мастерская</span>
        </RouterLink>
        <nav class="hidden items-center gap-6 text-[15px] text-[#4C4A40] md:flex">
          <RouterLink :to="{ name: 'marketplace-catalog' }" class="hover:text-[#5B4BE0]">Каталог</RouterLink>
          <RouterLink :to="{ name: 'marketplace-masters' }" class="hover:text-[#5B4BE0]">Мастера</RouterLink>
          <RouterLink :to="{ name: 'marketplace-for-masters' }" class="hover:text-[#5B4BE0]">Мастерам</RouterLink>
          <RouterLink :to="{ name: 'marketplace-about' }" class="hover:text-[#5B4BE0]">О платформе</RouterLink>
          <RouterLink
            v-if="auth.isAuthenticated && auth.role === 'CLIENT'"
            :to="{ name: 'marketplace-my-requests' }"
            class="relative hover:text-[#5B4BE0]"
          >
            Мои заявки
            <span
              v-if="unreadCount > 0"
              class="ml-1 inline-flex h-4 min-w-4 items-center justify-center rounded-full bg-[#B3261E] px-1 text-[10px] font-semibold leading-none text-white"
            >
              {{ unreadCount > 9 ? '9+' : unreadCount }}
            </span>
          </RouterLink>
          <RouterLink
            v-if="auth.isAuthenticated && auth.role === 'CLIENT'"
            :to="{ name: 'marketplace-favorites' }"
            class="hover:text-[#5B4BE0]"
          >
            Избранное
          </RouterLink>
          <RouterLink
            v-if="auth.isAuthenticated"
            :to="{ name: 'marketplace-profile' }"
            class="hover:text-[#5B4BE0]"
          >
            Личный кабинет
          </RouterLink>
          <RouterLink
            v-if="auth.isAuthenticated && auth.role === 'SUPER_ADMIN'"
            :to="{ name: 'admin-categories' }"
            class="hover:text-[#5B4BE0]"
          >
            Админ: Категории
          </RouterLink>
          <RouterLink
            v-if="auth.isAuthenticated && auth.role === 'SUPER_ADMIN'"
            :to="{ name: 'admin-services' }"
            class="hover:text-[#5B4BE0]"
          >
            Админ: Услуги
          </RouterLink>
          <RouterLink
            v-if="auth.isAuthenticated && auth.role === 'SUPER_ADMIN'"
            :to="{ name: 'admin-requests' }"
            class="hover:text-[#5B4BE0]"
          >
            Админ: Заявки
          </RouterLink>
          <RouterLink
            v-if="auth.isAuthenticated && auth.role === 'SUPER_ADMIN'"
            :to="{ name: 'admin-masters' }"
            class="hover:text-[#5B4BE0]"
          >
            Админ: Мастера
          </RouterLink>
          <RouterLink
            v-if="auth.isAuthenticated && auth.role === 'SUPER_ADMIN'"
            :to="{ name: 'admin-reviews' }"
            class="hover:text-[#5B4BE0]"
          >
            Админ: Отзывы
          </RouterLink>
          <RouterLink
            v-if="auth.isAuthenticated && auth.role === 'SUPER_ADMIN'"
            :to="{ name: 'admin-payments' }"
            class="hover:text-[#5B4BE0]"
          >
            Админ: Платежи
          </RouterLink>
        </nav>
        <div class="ml-auto flex items-center gap-3">
          <template v-if="auth.isAuthenticated">
            <span class="hidden text-[15px] text-[#55524A] sm:block">{{ auth.fullName }}</span>
            <button
              type="button"
              class="rounded-[11px] border border-[#DDD8CC] px-4 py-2.5 text-[15px] hover:border-[#17160F]"
              @click="handleLogout"
            >
              Выйти
            </button>
          </template>
          <template v-else>
            <RouterLink
              :to="{ name: 'marketplace-login' }"
              class="rounded-[11px] border border-[#DDD8CC] px-4 py-2.5 text-[15px] hover:border-[#17160F]"
            >
              Войти
            </RouterLink>
            <RouterLink
              :to="{ name: 'marketplace-register' }"
              class="rounded-[11px] bg-[#17160F] px-4.5 py-2.5 text-[15px] text-[#FBFAF7] hover:bg-[#5B4BE0]"
            >
              Стать мастером
            </RouterLink>
          </template>
        </div>
      </div>
    </header>

    <!-- No width cap here on purpose: some pages (landing) need full-bleed -->
    <!-- bands outside the 1280px column. Each page wraps its own content. -->
    <main class="flex-1">
      <slot />
    </main>

    <footer class="border-t border-[#E7E3D9] bg-[#FBFAF7] py-6">
      <p class="mx-auto max-w-[1280px] px-5 text-xs text-[#9B978A] md:px-10">
        Мастерская — дипломный проект, учебный маркетплейс мастеров и заявок на бытовые услуги.
      </p>
    </footer>
  </div>
</template>

<style>
/* Same "Мастерская" identity as MarketplaceLandingView.vue — see the note
   there. Kept as a global (non-scoped) block here too so it applies
   consistently across every marketplace page that uses this shell. */
.mk-home .mk-display {
  font-family: 'Playfair Display', Georgia, serif;
}
/* Highlights the current page in the header nav -- vue-router adds this
   class to a RouterLink automatically, no per-link wiring needed. */
.mk-home nav a.router-link-active {
  color: #5b4be0;
  font-weight: 500;
}
.mk-home {
  font-family: 'Golos Text', 'IBM Plex Sans', ui-sans-serif, system-ui, sans-serif;
}
.mk-home .font-mono {
  font-family: 'JetBrains Mono', ui-monospace, SFMono-Regular, monospace;
}
</style>
