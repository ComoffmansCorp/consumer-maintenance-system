<script setup lang="ts">
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { extractErrorMessage } from '@/api/client'
import AppInput from '@/components/ui/AppInput.vue'
import AppButton from '@/components/ui/AppButton.vue'

const auth = useAuthStore()
const router = useRouter()
const route = useRoute()

const tenantCode = ref('')
const username = ref('')
const password = ref('')
const loading = ref(false)
const error = ref('')

async function submit() {
  loading.value = true
  error.value = ''
  try {
    await auth.login({ tenantCode: tenantCode.value, username: username.value, password: password.value })
    const redirect = (route.query.redirect as string) || '/'
    router.push(redirect)
  } catch (e) {
    error.value = extractErrorMessage(e)
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="auth-shell flex min-h-screen items-center justify-center bg-surface px-4 dark:bg-ink">
    <div class="card w-full max-w-sm p-8">
      <h1 class="font-display text-xl font-semibold text-ink dark:text-surface">Вход в систему</h1>
      <p class="mt-1 text-sm text-slate dark:text-mist">Обслуживание приборов учёта</p>

      <form class="mt-6 flex flex-col gap-4" @submit.prevent="submit">
        <AppInput
          v-model="tenantCode"
          label="Код компании"
          placeholder="Оставьте пустым для входа как супер-админ"
        />
        <AppInput v-model="username" label="Логин" required />
        <AppInput v-model="password" type="password" label="Пароль" required />

        <p v-if="error" class="text-sm font-medium" style="color: var(--color-status-canceled)">
          {{ error }}
        </p>

        <AppButton type="submit" :loading="loading" class="w-full justify-center">Войти</AppButton>
      </form>

      <div class="mt-6 flex flex-col gap-1 text-center text-sm text-slate dark:text-mist">
        <RouterLink to="/register-company" class="font-medium text-primary hover:underline">
          Зарегистрировать компанию
        </RouterLink>
        <RouterLink to="/bootstrap" class="text-xs text-mist hover:underline">
          Первый запуск: создать супер-админа
        </RouterLink>
      </div>
    </div>
  </div>
</template>
