<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { extractErrorMessage } from '@/api/client'
import AppInput from '@/components/ui/AppInput.vue'
import AppSelect from '@/components/ui/AppSelect.vue'
import AppButton from '@/components/ui/AppButton.vue'

const auth = useAuthStore()
const router = useRouter()

const tenantName = ref('')
const tenantCode = ref('')
const plan = ref('FREE')
const username = ref('')
const password = ref('')
const fullName = ref('')
const loading = ref(false)
const error = ref('')

async function submit() {
  loading.value = true
  error.value = ''
  try {
    await auth.registerCompany({
      tenantName: tenantName.value,
      tenantCode: tenantCode.value,
      plan: plan.value,
      username: username.value,
      password: password.value,
      fullName: fullName.value,
    })
    router.push('/')
  } catch (e) {
    error.value = extractErrorMessage(e)
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="auth-shell flex min-h-screen items-center justify-center bg-surface px-4 py-10 dark:bg-ink">
    <div class="card w-full max-w-md p-8">
      <h1 class="font-display text-xl font-semibold text-ink dark:text-surface">Регистрация компании</h1>
      <p class="mt-1 text-sm text-slate dark:text-mist">
        Создаём тенант и первого администратора компании
      </p>

      <form class="mt-6 flex flex-col gap-4" @submit.prevent="submit">
        <AppInput v-model="tenantName" label="Название компании" required />
        <AppInput v-model="tenantCode" label="Код компании (для входа)" required />
        <AppSelect
          v-model="plan"
          label="Тариф"
          :options="[
            { value: 'FREE', label: 'Free — до 5 пользователей' },
            { value: 'BUSINESS', label: 'Business — до 50 пользователей' },
            { value: 'ENTERPRISE', label: 'Enterprise — без ограничений' },
          ]"
        />
        <AppInput v-model="fullName" label="Ваше имя" required />
        <AppInput v-model="username" label="Логин администратора" required />
        <AppInput v-model="password" type="password" label="Пароль" required />

        <p v-if="error" class="text-sm font-medium" style="color: var(--color-status-canceled)">
          {{ error }}
        </p>

        <AppButton type="submit" :loading="loading" class="w-full justify-center">
          Зарегистрировать
        </AppButton>
      </form>

      <RouterLink to="/login" class="mt-6 block text-center text-sm font-medium text-primary hover:underline">
        Уже есть аккаунт? Войти
      </RouterLink>
    </div>
  </div>
</template>
