<script setup lang="ts">
import { onMounted, ref } from 'vue'
import OdometerNumber from '@/components/ui/OdometerNumber.vue'

withDefaults(defineProps<{ tag?: string }>(), { tag: 'ОБСЛУЖИВАНИЕ ПРИБОРОВ УЧЁТА' })

const serial = ref(0)
onMounted(() => {
  window.setTimeout(() => (serial.value = 240719), 150)
})
</script>

<template>
  <div class="relative hidden overflow-hidden bg-graphite md:flex md:w-[42%] md:min-w-[380px] md:flex-col md:justify-between">
    <!-- Oversized gauge face, cropped off-canvas like a close-up photo of the housing -->
    <svg
      viewBox="0 0 400 400"
      class="pointer-events-none absolute -bottom-24 -right-32 h-[560px] w-[560px] text-primary opacity-[0.14]"
      aria-hidden="true"
    >
      <circle cx="200" cy="200" r="150" fill="none" stroke="currentColor" stroke-width="10" />
      <circle cx="200" cy="200" r="96" fill="none" stroke="currentColor" stroke-width="3" />
      <circle cx="200" cy="200" r="24" fill="currentColor" />
      <g stroke="currentColor" stroke-width="6" stroke-linecap="round">
        <path d="M200,44 v34" />
        <path d="M200,322 v34" />
        <path d="M356,200 h-34" />
        <path d="M44,200 h34" />
        <path d="M310,90 l-24,24" />
        <path d="M90,310 l24,-24" />
        <path d="M310,310 l-24,-24" />
        <path d="M90,90 l24,24" />
      </g>
    </svg>

    <div class="relative z-10 flex items-center gap-2.5 px-10 pt-10">
      <svg viewBox="0 0 24 24" class="h-6 w-6 shrink-0 text-primary" aria-hidden="true">
        <circle cx="12" cy="12" r="9.25" fill="none" stroke="currentColor" stroke-width="1.5" />
        <circle cx="12" cy="12" r="2" fill="currentColor" />
        <path
          d="M12 4v2M12 18v2M20 12h-2M6 12H4M17.36 6.64l-1.42 1.42M8.06 15.94l-1.42 1.42M17.36 17.36l-1.42-1.42M8.06 8.06 6.64 6.64"
          stroke="currentColor"
          stroke-width="1.5"
          stroke-linecap="round"
        />
      </svg>
      <span class="font-display text-lg font-semibold text-surface">ПриборСервис</span>
    </div>

    <div class="relative z-10 px-10 pb-10">
      <p class="max-w-[22ch] font-display text-2xl leading-snug text-surface">
        Учёт нарядов и актов обслуживания приборов учёта
      </p>
      <p class="data-mono mt-6 text-xs uppercase tracking-widest text-mist">{{ tag }}</p>
      <div class="mt-3 flex items-center gap-2 on-dark-odometer">
        <span class="data-mono text-[10px] uppercase tracking-widest text-mist">Узел учёта №</span>
        <OdometerNumber :value="serial" :min-digits="6" />
      </div>
    </div>
  </div>
</template>

<style scoped>
.on-dark-odometer :deep(.odometer-digit) {
  background: var(--color-ink);
  color: var(--color-surface);
  border-right-color: var(--color-graphite);
}
</style>
