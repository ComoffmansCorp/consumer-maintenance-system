<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { authApi } from '@/api/auth'
import { useAuthStore } from '@/stores/auth'
import { extractErrorMessage } from '@/api/client'
import AppInput from '@/components/ui/AppInput.vue'
import AppButton from '@/components/ui/AppButton.vue'
import AuthSplitPanel from '@/components/layout/AuthSplitPanel.vue'

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
  <div class="flex min-h-screen bg-surface dark:bg-ink">
    <AuthSplitPanel tag="ПЕРВЫЙ ЗАПУСК СИСТЕМЫ" />

    <div class="flex flex-1 items-center justify-center px-6 py-12">
      <div class="w-full max-w-sm">
        <p class="field-label text-primary">Первый запуск</p>
        <h1 class="mt-1 font-display text-3xl font-semibold text-ink dark:text-surface">
          Супер-админ
        </h1>
        <p class="mt-2 text-sm text-slate dark:text-mist">Создание супер-администратора платформы</p>

        <form class="mt-8 flex flex-col gap-5" @submit.prevent="submit">
          <AppInput v-model="fullName" label="Полное имя" required />
          <AppInput v-model="username" label="Логин" required />
          <AppInput v-model="password" type="password" label="Пароль" required />

          <p v-if="error" class="text-sm font-medium" style="color: var(--color-status-canceled)">
            {{ error }}
          </p>

          <AppButton type="submit" :loading="loading" class="mt-2 w-full justify-center">Создать</AppButton>
        </form>
      </div>
    </div>
  </div>
</template>
