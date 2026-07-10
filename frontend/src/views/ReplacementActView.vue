<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref } from 'vue'
import { actsApi } from '@/api/acts'
import { photosApi } from '@/api/photos'
import { fetchBlobUrl, extractErrorMessage } from '@/api/client'
import type { PhotoDTO, ReplacementActDTO } from '@/types'
import { useToastStore } from '@/stores/toast'
import LoadingState from '@/components/ui/LoadingState.vue'
import ErrorState from '@/components/ui/ErrorState.vue'
import AppButton from '@/components/ui/AppButton.vue'
import OdometerNumber from '@/components/ui/OdometerNumber.vue'

const props = defineProps<{ id: string }>()
const toast = useToastStore()

const act = ref<ReplacementActDTO | null>(null)
const photos = ref<PhotoDTO[]>([])
const photoUrls = ref<Record<number, string>>({})
const loading = ref(true)
const error = ref('')
const pdfLoading = ref(false)

async function load() {
  loading.value = true
  error.value = ''
  try {
    const id = Number(props.id)
    act.value = await actsApi.getReplacement(id)
    photos.value = await photosApi.listReplacement(id)
    for (const p of photos.value) {
      fetchBlobUrl(photosApi.downloadPath(p.id))
        .then((url) => (photoUrls.value[p.id] = url))
        .catch(() => {})
    }
  } catch (e) {
    error.value = extractErrorMessage(e)
  } finally {
    loading.value = false
  }
}
onMounted(load)

onBeforeUnmount(() => {
  Object.values(photoUrls.value).forEach((url) => URL.revokeObjectURL(url))
})

async function downloadPdf() {
  if (!act.value) return
  pdfLoading.value = true
  try {
    const url = await fetchBlobUrl(actsApi.replacementPdfPath(act.value.id))
    window.open(url, '_blank')
  } catch (e) {
    toast.error(extractErrorMessage(e))
  } finally {
    pdfLoading.value = false
  }
}

function formatDate(iso?: string) {
  if (!iso) return '—'
  return new Date(iso).toLocaleDateString('ru-RU')
}
</script>

<template>
  <div class="flex flex-col gap-6">
    <LoadingState v-if="loading" />
    <ErrorState v-else-if="error" :message="error" @retry="load" />

    <template v-else-if="act">
      <div class="flex flex-wrap items-start justify-between gap-3">
        <div>
          <p class="text-sm text-slate dark:text-mist">Акт замены №{{ act.id }}</p>
          <h1 class="font-display text-2xl font-semibold text-ink dark:text-surface">{{ act.addressLabel }}</h1>
        </div>
        <AppButton :loading="pdfLoading" @click="downloadPdf">Скачать PDF</AppButton>
      </div>

      <div class="card grid grid-cols-1 gap-4 p-6 sm:grid-cols-2">
        <div>
          <p class="text-xs uppercase tracking-wide text-mist">Лицевой счёт</p>
          <p class="data-mono mt-1 text-sm text-ink dark:text-surface">{{ act.accountNumber }}</p>
        </div>
        <div>
          <p class="text-xs uppercase tracking-wide text-mist">Дата замены</p>
          <p class="data-mono mt-1 text-sm text-ink dark:text-surface">{{ formatDate(act.installationDate) }}</p>
        </div>
      </div>

      <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
        <div class="card p-6">
          <h2 class="mb-4 font-display text-base font-semibold text-ink dark:text-surface">Снятый прибор</h2>
          <dl class="space-y-3 text-sm">
            <div>
              <dt class="text-xs uppercase tracking-wide text-mist">Марка</dt>
              <dd class="mt-0.5 text-ink dark:text-surface">{{ act.oldBrand || '—' }}</dd>
            </div>
            <div>
              <dt class="text-xs uppercase tracking-wide text-mist">Серийный номер</dt>
              <dd class="data-mono mt-0.5 text-ink dark:text-surface">{{ act.oldSerialNumber || '—' }}</dd>
            </div>
            <div>
              <dt class="text-xs uppercase tracking-wide text-mist">Показания</dt>
              <dd class="mt-1">
                <OdometerNumber v-if="act.oldReadings != null" :value="act.oldReadings" :min-digits="5" />
                <span v-else class="text-sm text-ink dark:text-surface">—</span>
              </dd>
            </div>
          </dl>
        </div>
        <div class="card p-6">
          <h2 class="mb-4 font-display text-base font-semibold text-ink dark:text-surface">Установленный прибор</h2>
          <dl class="space-y-3 text-sm">
            <div>
              <dt class="text-xs uppercase tracking-wide text-mist">Марка</dt>
              <dd class="mt-0.5 text-ink dark:text-surface">{{ act.newBrand || '—' }}</dd>
            </div>
            <div>
              <dt class="text-xs uppercase tracking-wide text-mist">Серийный номер</dt>
              <dd class="data-mono mt-0.5 text-ink dark:text-surface">{{ act.newSerialNumber || '—' }}</dd>
            </div>
            <div>
              <dt class="text-xs uppercase tracking-wide text-mist">Показания</dt>
              <dd class="mt-1">
                <OdometerNumber v-if="act.newReadings != null" :value="act.newReadings" :min-digits="5" />
                <span v-else class="text-sm text-ink dark:text-surface">—</span>
              </dd>
            </div>
          </dl>
        </div>
      </div>

      <div class="card p-6">
        <h2 class="mb-4 font-display text-base font-semibold text-ink dark:text-surface">
          Фотофиксация ({{ photos.length }})
        </h2>
        <p v-if="photos.length === 0" class="text-sm text-slate dark:text-mist">Фотографий нет.</p>
        <div v-else class="grid grid-cols-2 gap-3 sm:grid-cols-4">
          <a
            v-for="p in photos"
            :key="p.id"
            :href="photoUrls[p.id]"
            target="_blank"
            rel="noopener"
            class="hairline block overflow-hidden rounded-md bg-surface dark:bg-ink"
          >
            <img
              v-if="photoUrls[p.id]"
              :src="photoUrls[p.id]"
              :alt="p.note || p.originalFilename"
              class="aspect-square w-full object-cover"
            />
            <div v-else class="aspect-square w-full animate-pulse motion-reduce:animate-none" />
          </a>
        </div>
      </div>
    </template>
  </div>
</template>
