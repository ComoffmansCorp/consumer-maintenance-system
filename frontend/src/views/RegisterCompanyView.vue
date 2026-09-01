<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { extractErrorMessage } from '@/api/client'
import AppInput from '@/components/ui/AppInput.vue'
import AppSelect from '@/components/ui/AppSelect.vue'
import AppButton from '@/components/ui/AppButton.vue'
import AuthSplitPanel from '@/components/layout/AuthSplitPanel.vue'

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
  <div class="flex min-h-screen bg-surface dark:bg-ink">
    <AuthSplitPanel tag="РЕГИСТРАЦИЯ КОМПАНИИ" />

    <div class="flex flex-1 items-center justify-center px-6 py-12">
      <div class="w-full max-w-sm">
        <p class="field-label text-primary">Новая компания</p>
        <h1 class="mt-1 font-display text-3xl font-semibold text-ink dark:text-surface">
          Регистрация
        </h1>
        <p class="mt-2 text-sm text-slate dark:text-mist">
          Создаём тенант и первого администратора компании
        </p>

        <form class="mt-8 flex flex-col gap-5" @submit.prevent="submit">
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

          <AppButton type="submit" :loading="loading" class="mt-2 w-full justify-center">
            Зарегистрировать
          </AppButton>
        </form>

        <div class="mt-8 border-t border-line pt-6 dark:border-graphite">
          <RouterLink to="/login" class="text-sm font-medium text-primary hover:underline">
            Уже есть аккаунт? Войти
          </RouterLink>
        </div>
      </div>
    </div>
  </div>
</template>
