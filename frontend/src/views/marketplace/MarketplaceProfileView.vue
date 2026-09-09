<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { catalogApi, masterApi, requestsApi } from '@/api/marketplace'
import { extractErrorMessage } from '@/api/client'
import { useAuthStore } from '@/stores/auth'
import { useToastStore } from '@/stores/toast'
import type { FavoriteDTO, ProfileDTO, RequestDTO, ServiceDTO } from '@/types'
import MarketplaceShell from '@/components/layout/MarketplaceShell.vue'
import StarRating from '@/components/StarRating.vue'

const auth = useAuthStore()
const toast = useToastStore()

const loading = ref(true)
const error = ref('')

// MASTER
const profile = ref<ProfileDTO | null>(null)
const services = ref<ServiceDTO[]>([])
const editing = ref(false)
const form = ref({ city: '', bio: '', specializationIds: [] as number[] })
const saving = ref(false)

// CLIENT
const myRequests = ref<RequestDTO[]>([])
const favorites = ref<FavoriteDTO[]>([])

const specializationNames = computed(() => {
  const byId = new Map(services.value.map((s) => [s.id, s.name]))
  return (profile.value?.specializationIds ?? []).map((id) => byId.get(id) ?? `#${id}`)
})

const activeRequestsCount = computed(
  () => myRequests.value.filter((r) => r.status === 'OPEN' || r.status === 'ASSIGNED').length,
)
const completedRequestsCount = computed(() => myRequests.value.filter((r) => r.status === 'COMPLETED').length)

async function load() {
  loading.value = true
  error.value = ''
  try {
    if (auth.role === 'MASTER') {
      const [p, s] = await Promise.all([masterApi.getProfile(), catalogApi.listServices()])
      profile.value = p
      services.value = s
      form.value = { city: p.city ?? '', bio: p.bio ?? '', specializationIds: [...p.specializationIds] }
    } else if (auth.role === 'CLIENT') {
      const [reqs, favs] = await Promise.all([requestsApi.listMine(1, 200), requestsApi.listFavorites()])
      myRequests.value = reqs.items
      favorites.value = favs
    }
  } catch (e) {
    error.value = extractErrorMessage(e)
  } finally {
    loading.value = false
  }
}
onMounted(load)

function startEdit() {
  if (!profile.value) return
  form.value = {
    city: profile.value.city ?? '',
    bio: profile.value.bio ?? '',
    specializationIds: [...profile.value.specializationIds],
  }
  editing.value = true
}

function toggleSpecialization(id: number) {
  const idx = form.value.specializationIds.indexOf(id)
  if (idx === -1) form.value.specializationIds.push(id)
  else form.value.specializationIds.splice(idx, 1)
}

async function save() {
  saving.value = true
  try {
    profile.value = await masterApi.updateProfile(form.value)
    editing.value = false
    toast.success('Профиль обновлён')
  } catch (e) {
    toast.error(extractErrorMessage(e))
  } finally {
    saving.value = false
  }
}

const roleLabel = computed(() => {
  switch (auth.role) {
    case 'MASTER':
      return 'Мастер'
    case 'CLIENT':
      return 'Клиент'
    case 'SUPER_ADMIN':
      return 'Администратор платформы'
    default:
      return ''
  }
})
</script>

<template>
  <MarketplaceShell>
    <div class="mx-auto max-w-[820px]">
      <h1 class="mk-display text-2xl font-bold tracking-tight">Личный кабинет</h1>

      <div v-if="loading" class="py-16 text-center text-sm text-[#8D8A7E]">Загрузка…</div>
      <div v-else-if="error" class="mt-6 rounded-2xl border border-[#F3D3CE] bg-[#FBF0EE] px-6 py-10 text-center">
        <p class="text-sm font-medium text-[#B3261E]">{{ error }}</p>
      </div>

      <template v-else>
        <!-- Identity card, same for every role -->
        <div class="mt-6 flex items-center gap-4 rounded-2xl border border-[#E2DED2] bg-white p-6">
          <img
            v-if="auth.role === 'MASTER' && profile?.avatarUrl"
            :src="profile.avatarUrl"
            alt=""
            class="h-16 w-16 rounded-full object-cover"
          />
          <div
            v-else
            class="grid h-16 w-16 place-items-center rounded-full bg-[#EFE9FF] text-xl font-semibold text-[#5B4BE0]"
          >
            {{ auth.fullName?.[0] ?? '?' }}
          </div>
          <div>
            <div class="mk-display text-lg font-medium">{{ auth.fullName }}</div>
            <div class="text-sm text-[#8D8A7E]">{{ roleLabel }}</div>
          </div>
          <div v-if="auth.role === 'MASTER' && profile" class="ml-auto">
            <StarRating :value="profile.ratingAvg" :count="profile.ratingCount" />
          </div>
        </div>

        <!-- MASTER: profile details + edit -->
        <div v-if="auth.role === 'MASTER' && profile" class="mt-4 rounded-2xl border border-[#E2DED2] bg-white p-6">
          <div class="flex items-center justify-between">
            <h2 class="mk-display text-lg font-medium">О мастере</h2>
            <button
              v-if="!editing"
              type="button"
              class="rounded-[11px] border border-[#DDD8CC] px-4 py-2 text-sm hover:border-[#17160F]"
              @click="startEdit"
            >
              Редактировать
            </button>
          </div>

          <template v-if="!editing">
            <p class="mt-3 text-sm text-[#6E6B60]">{{ profile.city || 'Город не указан' }}</p>
            <p class="mt-2 text-[15px] leading-[1.6] text-[#4C4A40]">{{ profile.bio || 'Описание не заполнено' }}</p>
            <div class="mt-4 flex flex-wrap gap-2">
              <span
                v-for="name in specializationNames"
                :key="name"
                class="rounded-full bg-[#EFE9FF] px-3 py-1.5 text-xs text-[#5B4BE0]"
              >
                {{ name }}
              </span>
              <span v-if="specializationNames.length === 0" class="text-xs text-[#8D8A7E]">Специализации не выбраны</span>
            </div>
          </template>

          <form v-else class="mt-4 flex flex-col gap-4" @submit.prevent="save">
            <label class="flex flex-col gap-1.5 text-sm">
              Город
              <input
                v-model="form.city"
                type="text"
                class="rounded-[11px] border border-[#DDD8CC] px-3.5 py-2.5 text-[15px] outline-none focus:border-[#5B4BE0]"
              />
            </label>
            <label class="flex flex-col gap-1.5 text-sm">
              О себе
              <textarea
                v-model="form.bio"
                rows="3"
                class="rounded-[11px] border border-[#DDD8CC] px-3.5 py-2.5 text-[15px] outline-none focus:border-[#5B4BE0]"
              />
            </label>
            <div class="flex flex-col gap-1.5 text-sm">
              Специализации
              <div class="flex flex-wrap gap-2">
                <button
                  v-for="s in services"
                  :key="s.id"
                  type="button"
                  class="rounded-full border px-3 py-1.5 text-xs"
                  :class="
                    form.specializationIds.includes(s.id)
                      ? 'border-[#5B4BE0] bg-[#EFE9FF] text-[#5B4BE0]'
                      : 'border-[#DDD8CC] text-[#6E6B60] hover:border-[#5B4BE0]'
                  "
                  @click="toggleSpecialization(s.id)"
                >
                  {{ s.name }}
                </button>
              </div>
            </div>
            <div class="flex gap-3">
              <button
                type="submit"
                :disabled="saving"
                class="rounded-[11px] bg-[#5B4BE0] px-5 py-2.5 text-sm font-medium text-white disabled:opacity-60"
              >
                {{ saving ? 'Сохранение…' : 'Сохранить' }}
              </button>
              <button
                type="button"
                class="rounded-[11px] border border-[#DDD8CC] px-5 py-2.5 text-sm"
                @click="editing = false"
              >
                Отмена
              </button>
            </div>
          </form>
        </div>

        <!-- CLIENT: account overview -->
        <div v-else-if="auth.role === 'CLIENT'" class="mt-4 grid grid-cols-2 gap-4">
          <RouterLink
            :to="{ name: 'marketplace-my-requests' }"
            class="rounded-2xl border border-[#E2DED2] bg-white p-6 hover:border-[#5B4BE0]"
          >
            <div class="mk-display text-2xl font-semibold text-[#5B4BE0]">{{ activeRequestsCount }}</div>
            <div class="mt-1 text-sm text-[#6E6B60]">активных заявок</div>
          </RouterLink>
          <RouterLink
            :to="{ name: 'marketplace-my-requests' }"
            class="rounded-2xl border border-[#E2DED2] bg-white p-6 hover:border-[#5B4BE0]"
          >
            <div class="mk-display text-2xl font-semibold">{{ completedRequestsCount }}</div>
            <div class="mt-1 text-sm text-[#6E6B60]">выполненных заявок</div>
          </RouterLink>
          <RouterLink
            :to="{ name: 'marketplace-favorites' }"
            class="col-span-2 rounded-2xl border border-[#E2DED2] bg-white p-6 hover:border-[#5B4BE0]"
          >
            <div class="mk-display text-2xl font-semibold">{{ favorites.length }}</div>
            <div class="mt-1 text-sm text-[#6E6B60]">мастеров в избранном</div>
          </RouterLink>
        </div>

        <!-- SUPER_ADMIN: quick links -->
        <div v-else-if="auth.role === 'SUPER_ADMIN'" class="mt-4 rounded-2xl border border-[#E2DED2] bg-white p-6">
          <p class="text-sm text-[#6E6B60]">
            Разделы администрирования — Категории, Услуги, Заявки, Мастера, Отзывы, Платежи — доступны в шапке сайта.
          </p>
        </div>
      </template>
    </div>
  </MarketplaceShell>
</template>
