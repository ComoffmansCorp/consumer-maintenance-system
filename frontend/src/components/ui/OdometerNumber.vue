<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(
  defineProps<{
    value: number
    minDigits?: number
    size?: 'md' | 'lg'
  }>(),
  { minDigits: 2, size: 'md' },
)

const digits = computed(() => {
  const str = String(Math.max(0, Math.trunc(props.value)))
  return str.padStart(props.minDigits, '0').split('')
})
</script>

<template>
  <span class="odometer" :aria-label="String(value)">
    <span
      v-for="(d, i) in digits"
      :key="i"
      class="odometer-digit"
      :class="size === 'lg' ? 'h-11 w-7 text-2xl font-semibold' : 'h-7 w-5 text-sm font-semibold'"
    >
      <Transition name="odometer-roll" mode="out-in">
        <span :key="d">{{ d }}</span>
      </Transition>
    </span>
  </span>
</template>

<style scoped>
.odometer-roll-enter-active,
.odometer-roll-leave-active {
  transition:
    transform 0.18s ease,
    opacity 0.18s ease;
}
.odometer-roll-enter-from {
  transform: translateY(40%);
  opacity: 0;
}
.odometer-roll-leave-to {
  transform: translateY(-40%);
  opacity: 0;
}
</style>
