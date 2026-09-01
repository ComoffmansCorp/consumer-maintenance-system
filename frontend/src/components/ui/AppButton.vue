<script setup lang="ts">
withDefaults(
  defineProps<{
    variant?: 'primary' | 'secondary' | 'ghost' | 'danger'
    size?: 'sm' | 'md'
    type?: 'button' | 'submit'
    disabled?: boolean
    loading?: boolean
  }>(),
  { variant: 'primary', size: 'md', type: 'button', disabled: false, loading: false },
)
</script>

<template>
  <button
    :type="type"
    :disabled="disabled || loading"
    class="inline-flex items-center justify-center gap-2 font-medium transition-colors focus-visible:outline-2 disabled:cursor-not-allowed disabled:opacity-50"
    :class="[
      size === 'sm' ? 'px-3 py-1.5 text-sm' : 'px-4 py-2 text-sm',
      variant === 'primary' &&
        'notch-corner bg-primary font-semibold tracking-wide text-white hover:bg-[var(--color-primary-hover)] disabled:hover:bg-primary',
      variant === 'secondary' &&
        'rounded-md border border-line bg-white text-ink hover:bg-surface dark:border-graphite dark:bg-graphite dark:text-surface dark:hover:bg-ink',
      variant === 'ghost' && 'rounded-md text-slate hover:bg-surface dark:text-mist dark:hover:bg-graphite',
      variant === 'danger' && 'notch-corner bg-[var(--color-status-canceled)] text-white hover:opacity-90',
    ]"
  >
    <svg
      v-if="loading"
      class="h-4 w-4 animate-spin motion-reduce:animate-none"
      viewBox="0 0 24 24"
      fill="none"
    >
      <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
      <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v4a4 4 0 00-4 4H4z" />
    </svg>
    <slot />
  </button>
</template>
