<script setup lang="ts">
import { computed } from 'vue'
import type { RequestStatus } from '@/types'

const props = defineProps<{ status: RequestStatus }>()

const steps = [
  { key: 'OPEN', label: 'Открыта' },
  { key: 'ASSIGNED', label: 'В работе' },
  { key: 'COMPLETED', label: 'Выполнена' },
]

const activeIndex = computed(() => steps.findIndex((s) => s.key === props.status))
</script>

<template>
  <!-- CANCELED can happen from OPEN or ASSIGNED, so it doesn't fit a linear
       progress line -- shown as its own state instead of a 4th step. -->
  <div v-if="status === 'CANCELED'" class="flex items-center gap-2 text-sm font-medium text-[#B3261E]">
    <span class="h-2.5 w-2.5 rounded-full bg-[#B3261E]" />
    Заявка отменена
  </div>
  <div v-else class="flex items-center">
    <template v-for="(step, i) in steps" :key="step.key">
      <div class="flex flex-col items-center gap-1.5" :class="i > 0 ? 'flex-1' : ''">
        <div class="flex w-full items-center">
          <div v-if="i > 0" class="h-0.5 flex-1" :class="i <= activeIndex ? 'bg-[#5B4BE0]' : 'bg-[#E2DED2]'" />
          <div
            class="grid h-6 w-6 shrink-0 place-items-center rounded-full text-[11px] font-semibold"
            :class="i <= activeIndex ? 'bg-[#5B4BE0] text-white' : 'border border-[#E2DED2] bg-white text-[#A8A498]'"
          >
            <svg v-if="i < activeIndex" viewBox="0 0 16 16" class="h-3 w-3" fill="none">
              <path d="M3 8.5 6.2 12 13 4" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" />
            </svg>
            <span v-else>{{ i + 1 }}</span>
          </div>
        </div>
        <span class="text-xs" :class="i <= activeIndex ? 'font-medium text-[#17160F]' : 'text-[#A8A498]'">{{ step.label }}</span>
      </div>
    </template>
  </div>
</template>
