<script setup lang="ts">
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { extractErrorMessage } from '@/api/client'
import MarketplaceShell from '@/components/layout/MarketplaceShell.vue'

const auth = useAuthStore()
const router = useRouter()
const route = useRoute()

const username = ref('')
const password = ref('')
const loading = ref(false)
const error = ref('')

async function submit() {
  loading.value = true
  error.value = ''
  try {
    await auth.login({ username: username.value, password: password.value })
    const redirect = (route.query.redirect as string) || '/my-requests'
    router.push(redirect)
  } catch (e) {
    error.value = extractErrorMessage(e)
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <MarketplaceShell>
    <div class="mx-auto flex min-h-[70vh] max-w-[420px] flex-col justify-center px-1 py-10">
      <p class="font-mono text-[11px] uppercase tracking-[0.14em] text-[#5B4BE0]">Вход</p>
      <h1 class="mk-display mt-2 text-[32px] font-bold tracking-tight">Здравствуйте</h1>
      <p class="mt-2 text-[15px] text-[#6E6B60]">Для клиентов и мастеров «Мастерской».</p>

      <form class="mt-8 flex flex-col gap-4" @submit.prevent="submit">
        <label class="flex flex-col gap-1.5">
          <span class="font-mono text-[11px] uppercase tracking-[0.1em] text-[#9B978A]">Логин</span>
          <input
            v-model="username"
            type="text"
            required
            class="rounded-xl border border-[#E2DED2] bg-white px-4 py-3 text-base outline-none focus-visible:border-[#5B4BE0]"
          />
        </label>
        <label class="flex flex-col gap-1.5">
          <span class="font-mono text-[11px] uppercase tracking-[0.1em] text-[#9B978A]">Пароль</span>
          <input
            v-model="password"
            type="password"
            required
            class="rounded-xl border border-[#E2DED2] bg-white px-4 py-3 text-base outline-none focus-visible:border-[#5B4BE0]"
          />
        </label>

        <p v-if="error" class="text-sm font-medium text-[#B3261E]">{{ error }}</p>

        <button
          type="submit"
          :disabled="loading"
          class="mt-2 rounded-[14px] bg-[#5B4BE0] px-6 py-3.5 text-base font-medium text-white hover:bg-[#4536BC] disabled:opacity-60"
        >
          {{ loading ? 'Входим…' : 'Войти' }}
        </button>
      </form>

      <div class="mt-8 flex flex-col gap-2 border-t border-[#E2DED2] pt-6 text-sm">
        <RouterLink :to="{ name: 'marketplace-register' }" class="font-medium text-[#5B4BE0] hover:underline">
          Нет аккаунта? Зарегистрироваться
        </RouterLink>
        <RouterLink :to="{ name: 'marketplace' }" class="text-[#9B978A] hover:underline">← Назад к услугам</RouterLink>
      </div>
    </div>
  </MarketplaceShell>
</template>
