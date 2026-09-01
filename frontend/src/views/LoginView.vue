<script setup lang="ts">
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { extractErrorMessage } from '@/api/client'
import AppInput from '@/components/ui/AppInput.vue'
import AppButton from '@/components/ui/AppButton.vue'
import AuthSplitPanel from '@/components/layout/AuthSplitPanel.vue'

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
  <div class="flex min-h-screen bg-surface dark:bg-ink">
    <AuthSplitPanel />

    <div class="flex flex-1 items-center justify-center px-6 py-12">
      <div class="w-full max-w-sm">
        <p class="field-label text-primary">Вход в систему</p>
        <h1 class="mt-1 font-display text-3xl font-semibold text-ink dark:text-surface">
          Здравствуйте
        </h1>

        <form class="mt-8 flex flex-col gap-5" @submit.prevent="submit">
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

          <AppButton type="submit" :loading="loading" class="mt-2 w-full justify-center">Войти</AppButton>
        </form>

        <div class="mt-8 flex flex-col gap-1.5 border-t border-line pt-6 text-sm dark:border-graphite">
          <RouterLink to="/register-company" class="font-medium text-primary hover:underline">
            Зарегистрировать компанию
          </RouterLink>
          <RouterLink to="/bootstrap" class="text-xs text-mist hover:underline">
            Первый запуск: создать супер-админа
          </RouterLink>
        </div>
      </div>
    </div>
  </div>
</template>
