<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { authApi } from '@/api/auth'
import { useAuthStore } from '@/stores/auth'
import { extractErrorMessage } from '@/api/client'
import AppInput from '@/components/ui/AppInput.vue'
import AppButton from '@/components/ui/AppButton.vue'

const auth = useAuthStore()
const router = useRouter()

const username = ref('')
const password = ref('')
const fullName = ref('')
const loading = ref(false)
const error = ref('')

async function submit() {
  loading.value = true
  error.value = ''
  try {
    const resp = await authApi.bootstrapSuperAdmin({
      username: username.value,
      password: password.value,
      fullName: fullName.value,
    })
    auth.applySession(resp)
    router.push('/')
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
      <h1 class="font-display text-xl font-semibold text-ink dark:text-surface">Первый запуск</h1>
      <p class="mt-1 text-sm text-slate dark:text-mist">Создание супер-администратора платформы</p>

      <form class="mt-6 flex flex-col gap-4" @submit.prevent="submit">
        <AppInput v-model="fullName" label="Полное имя" required />
        <AppInput v-model="username" label="Логин" required />
        <AppInput v-model="password" type="password" label="Пароль" required />

        <p v-if="error" class="text-sm font-medium" style="color: var(--color-status-canceled)">
          {{ error }}
        </p>

        <AppButton type="submit" :loading="loading" class="w-full justify-center">Создать</AppButton>
      </form>
    </div>
  </div>
</template>
