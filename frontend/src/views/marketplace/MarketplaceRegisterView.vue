<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { extractErrorMessage } from '@/api/client'
import MarketplaceShell from '@/components/layout/MarketplaceShell.vue'

const auth = useAuthStore()
const router = useRouter()

const accountType = ref<'CLIENT' | 'MASTER'>('CLIENT')
const fullName = ref('')
const username = ref('')
const password = ref('')
const loading = ref(false)
const error = ref('')

async function submit() {
  loading.value = true
  error.value = ''
  try {
    const req = { username: username.value, password: password.value, fullName: fullName.value }
    if (accountType.value === 'CLIENT') {
      await auth.registerClient(req)
      router.push('/marketplace/my-requests')
    } else {
      await auth.registerMaster(req)
      router.push('/marketplace')
    }
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
      <p class="font-mono text-[11px] uppercase tracking-[0.14em] text-[#5B4BE0]">Новый аккаунт</p>
      <h1 class="mk-display mt-2 text-[32px] font-bold tracking-tight">Регистрация</h1>

      <div class="mt-6 grid grid-cols-2 gap-2">
        <button
          type="button"
          class="rounded-xl border px-4 py-2.5 text-sm font-medium transition-colors"
          :class="accountType === 'CLIENT' ? 'border-[#5B4BE0] bg-[#E7E3FC] text-[#5B4BE0]' : 'border-[#E2DED2] text-[#55524A] hover:border-[#17160F]'"
          @click="accountType = 'CLIENT'"
        >
          Я заказчик
        </button>
        <button
          type="button"
          class="rounded-xl border px-4 py-2.5 text-sm font-medium transition-colors"
          :class="accountType === 'MASTER' ? 'border-[#5B4BE0] bg-[#E7E3FC] text-[#5B4BE0]' : 'border-[#E2DED2] text-[#55524A] hover:border-[#17160F]'"
          @click="accountType = 'MASTER'"
        >
          Я мастер
        </button>
      </div>
      <p class="mt-2 text-xs text-[#9B978A]">
        {{
          accountType === 'CLIENT'
            ? 'Оставляйте заявки на услуги и следите за их статусом.'
            : 'Берите заявки из открытого пула по своей специализации.'
        }}
      </p>

      <form class="mt-6 flex flex-col gap-4" @submit.prevent="submit">
        <label class="flex flex-col gap-1.5">
          <span class="font-mono text-[11px] uppercase tracking-[0.1em] text-[#9B978A]">Ваше имя</span>
          <input
            v-model="fullName"
            type="text"
            required
            class="rounded-xl border border-[#E2DED2] bg-white px-4 py-3 text-base outline-none focus-visible:border-[#5B4BE0]"
          />
        </label>
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
          {{ loading ? 'Регистрируем…' : 'Зарегистрироваться' }}
        </button>
      </form>

      <div class="mt-8 border-t border-[#E2DED2] pt-6">
        <RouterLink :to="{ name: 'marketplace-login' }" class="text-sm font-medium text-[#5B4BE0] hover:underline">
          Уже есть аккаунт? Войти
        </RouterLink>
      </div>
    </div>
  </MarketplaceShell>
</template>
