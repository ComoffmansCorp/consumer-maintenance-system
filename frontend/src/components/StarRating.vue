<script setup lang="ts">
const props = defineProps<{ value: number; count?: number; size?: 'sm' | 'md' }>()

const stars = [1, 2, 3, 4, 5]
function fillFor(star: number) {
  const frac = props.value - (star - 1)
  if (frac >= 1) return 1
  if (frac <= 0) return 0
  return frac
}
</script>

<template>
  <span class="inline-flex items-center gap-1.5" :class="size === 'sm' ? 'text-sm' : 'text-base'">
    <span class="relative inline-flex" :class="size === 'sm' ? 'gap-0.5' : 'gap-1'">
      <span v-for="star in stars" :key="star" class="relative inline-block" :class="size === 'sm' ? 'h-3.5 w-3.5' : 'h-4.5 w-4.5'">
        <svg viewBox="0 0 20 20" class="absolute inset-0 h-full w-full text-[#E2DED2]" fill="currentColor">
          <path d="M10 1.5l2.6 5.6 6.1.6-4.6 4.1 1.3 6-5.4-3.2-5.4 3.2 1.3-6-4.6-4.1 6.1-.6z" />
        </svg>
        <span class="absolute inset-0 overflow-hidden" :style="{ width: `${fillFor(star) * 100}%` }">
          <svg viewBox="0 0 20 20" class="h-full w-full text-[#5B4BE0]" fill="currentColor" style="width: 20px; max-width: none">
            <path d="M10 1.5l2.6 5.6 6.1.6-4.6 4.1 1.3 6-5.4-3.2-5.4 3.2 1.3-6-4.6-4.1 6.1-.6z" />
          </svg>
        </span>
      </span>
    </span>
    <span v-if="value > 0" class="font-medium text-[#4C4A40]">{{ value.toFixed(1) }}</span>
    <span v-else class="text-[#9B978A]">Нет оценок</span>
    <span v-if="count !== undefined" class="text-[#8D8A7E]">({{ count }})</span>
  </span>
</template>
