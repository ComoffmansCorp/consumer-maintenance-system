<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()
const router = useRouter()

async function handleLogout() {
  await auth.logout()
  router.push({ name: 'marketplace' })
}
</script>

<template>
  <div class="mk-home flex min-h-screen flex-col bg-[#FBFAF7] text-[#17160F]">
    <header class="sticky top-0 z-20 border-b border-[#E7E3D9] bg-[#FBFAF7]/92 backdrop-blur-md">
      <div class="mx-auto flex max-w-[1280px] items-center gap-8 px-5 py-4 md:px-10">
        <RouterLink :to="{ name: 'marketplace' }" class="flex items-center gap-2.5">
          <span class="mk-display flex h-[30px] w-[30px] items-center justify-center rounded-[9px] bg-[#5B4BE0] text-[15px] font-semibold text-white">М</span>
          <span class="mk-display text-[17px] font-medium tracking-tight">Мастерская</span>
        </RouterLink>
        <nav class="hidden items-center gap-6 text-[15px] text-[#4C4A40] md:flex">
          <RouterLink :to="{ name: 'marketplace' }" class="hover:text-[#5B4BE0]">Все услуги</RouterLink>
          <RouterLink
            v-if="auth.isAuthenticated && auth.role === 'CLIENT'"
            :to="{ name: 'marketplace-my-requests' }"
            class="hover:text-[#5B4BE0]"
          >
            Мои заявки
          </RouterLink>
          <RouterLink
            v-if="auth.isAuthenticated && auth.role === 'CLIENT'"
            :to="{ name: 'marketplace-favorites' }"
            class="hover:text-[#5B4BE0]"
          >
            Избранное
          </RouterLink>
          <RouterLink
            v-if="auth.isAuthenticated && auth.role === 'MASTER'"
            :to="{ name: 'marketplace' }"
            class="hover:text-[#5B4BE0]"
          >
            Профиль мастера
          </RouterLink>
          <RouterLink
            v-if="auth.isAuthenticated && auth.role === 'SUPER_ADMIN'"
            :to="{ name: 'admin-categories' }"
            class="hover:text-[#5B4BE0]"
          >
            Категории
          </RouterLink>
          <RouterLink
            v-if="auth.isAuthenticated && auth.role === 'SUPER_ADMIN'"
            :to="{ name: 'admin-services' }"
            class="hover:text-[#5B4BE0]"
          >
            Услуги
          </RouterLink>
          <RouterLink
            v-if="auth.isAuthenticated && auth.role === 'SUPER_ADMIN'"
            :to="{ name: 'admin-requests' }"
            class="hover:text-[#5B4BE0]"
          >
            Заявки
          </RouterLink>
          <RouterLink
            v-if="auth.isAuthenticated && auth.role === 'SUPER_ADMIN'"
            :to="{ name: 'admin-masters' }"
            class="hover:text-[#5B4BE0]"
          >
            Мастера
          </RouterLink>
          <RouterLink
            v-if="auth.isAuthenticated && auth.role === 'SUPER_ADMIN'"
            :to="{ name: 'admin-reviews' }"
            class="hover:text-[#5B4BE0]"
          >
            Отзывы
          </RouterLink>
          <RouterLink
            v-if="auth.isAuthenticated && auth.role === 'SUPER_ADMIN'"
            :to="{ name: 'admin-payments' }"
            class="hover:text-[#5B4BE0]"
          >
            Платежи
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
.mk-home {
  font-family: 'Golos Text', 'IBM Plex Sans', ui-sans-serif, system-ui, sans-serif;
}
.mk-home .font-mono {
  font-family: 'JetBrains Mono', ui-monospace, SFMono-Regular, monospace;
}
</style>
