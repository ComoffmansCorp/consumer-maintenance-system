<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref } from 'vue'
import { actsApi } from '@/api/acts'
import { metersApi } from '@/api/meters'
import { photosApi } from '@/api/photos'
import { fetchBlobUrl, extractErrorMessage } from '@/api/client'
import type { InspectionActDTO, MeterDTO, PhotoDTO } from '@/types'
import { useToastStore } from '@/stores/toast'
import LoadingState from '@/components/ui/LoadingState.vue'
import ErrorState from '@/components/ui/ErrorState.vue'
import AppButton from '@/components/ui/AppButton.vue'
import SealMark from '@/components/ui/SealMark.vue'

const props = defineProps<{ id: string }>()
const toast = useToastStore()

const act = ref<InspectionActDTO | null>(null)
const meters = ref<MeterDTO[]>([])
const photos = ref<PhotoDTO[]>([])
const photoUrls = ref<Record<number, string>>({})
const loading = ref(true)
const error = ref('')
const pdfLoading = ref(false)

const typeLabels = { SCHEDULED: 'Плановый', UNSCHEDULED: 'Внеплановый' }
const meterTypeLabels: Record<string, string> = {
  SINGLE_PHASE: 'Однофазный',
  THREE_PHASE_DIRECT: 'Трёхфазный прямого включения',
  THREE_PHASE_TRANSFORMER: 'Трёхфазный трансформаторный',
}
const sealLabels: Record<string, string> = { INTACT: 'Не нарушена', BROKEN: 'Нарушена', MISSING: 'Отсутствует' }

async function load() {
  loading.value = true
  error.value = ''
  try {
    const id = Number(props.id)
    act.value = await actsApi.getInspection(id)
    ;[meters.value, photos.value] = await Promise.all([metersApi.list(id), photosApi.listInspection(id)])
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
    const url = await fetchBlobUrl(actsApi.inspectionPdfPath(act.value.id))
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
          <p class="text-sm text-slate dark:text-mist">Акт осмотра №{{ act.id }}</p>
          <h1 class="font-display text-2xl font-semibold text-ink dark:text-surface">{{ act.addressLabel }}</h1>
        </div>
        <AppButton :loading="pdfLoading" @click="downloadPdf">Скачать PDF</AppButton>
      </div>

      <div class="card grid grid-cols-1 gap-4 p-6 sm:grid-cols-2">
        <div>
          <p class="text-xs uppercase tracking-wide text-mist">Дата осмотра</p>
          <p class="data-mono mt-1 text-sm text-ink dark:text-surface">{{ formatDate(act.inspectionDate) }}</p>
        </div>
        <div>
          <p class="text-xs uppercase tracking-wide text-mist">Вид осмотра</p>
          <p class="mt-1 text-sm text-ink dark:text-surface">{{ typeLabels[act.inspectionType] }}</p>
        </div>
        <div v-if="act.consumerName">
          <p class="text-xs uppercase tracking-wide text-mist">Потребитель</p>
          <p class="mt-1 text-sm text-ink dark:text-surface">{{ act.consumerName }}</p>
        </div>
        <div v-if="act.notes" class="sm:col-span-2">
          <p class="text-xs uppercase tracking-wide text-mist">Примечания</p>
          <p class="mt-1 text-sm text-ink dark:text-surface">{{ act.notes }}</p>
        </div>
      </div>

      <div class="card overflow-hidden">
        <h2 class="border-b border-line px-6 py-4 font-display text-base font-semibold text-ink dark:border-graphite dark:text-surface">
          Приборы учёта
        </h2>
        <p v-if="meters.length === 0" class="px-6 py-6 text-sm text-slate dark:text-mist">
          Приборы не зафиксированы.
        </p>
        <table v-else class="w-full text-sm">
          <thead>
            <tr class="border-b border-line text-left text-xs uppercase tracking-wide text-mist dark:border-graphite">
              <th class="px-6 py-3 font-medium">Тип</th>
              <th class="px-6 py-3 font-medium">Серийный номер</th>
              <th class="px-6 py-3 font-medium">Поверка</th>
              <th class="px-6 py-3 font-medium">Пломба</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="m in meters" :key="m.id" class="border-b border-line last:border-0 dark:border-graphite">
              <td class="px-6 py-3 text-ink dark:text-surface">{{ meterTypeLabels[m.type] }}</td>
              <td class="data-mono px-6 py-3 text-ink dark:text-surface">{{ m.serialNumber }}</td>
              <td class="data-mono px-6 py-3 text-slate dark:text-mist">{{ formatDate(m.verificationDate) }}</td>
              <td class="px-6 py-3 text-slate dark:text-mist">
                <span v-if="m.sealState" class="inline-flex items-center gap-1.5">
                  <SealMark :state="m.sealState" />
                  {{ sealLabels[m.sealState] }}
                </span>
                <template v-else>—</template>
              </td>
            </tr>
          </tbody>
        </table>
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
