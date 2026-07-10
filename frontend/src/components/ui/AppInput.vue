<script setup lang="ts">
withDefaults(
  defineProps<{
    modelValue: string | number | null
    type?: string
    label?: string
    placeholder?: string
    required?: boolean
    mono?: boolean
    error?: string
  }>(),
  { type: 'text', required: false, mono: false },
)
defineEmits<{ 'update:modelValue': [value: string] }>()
</script>

<template>
  <label class="flex flex-col gap-1.5">
    <span v-if="label" class="text-sm font-medium text-slate dark:text-mist">{{ label }}</span>
    <input
      :type="type"
      :value="modelValue ?? ''"
      :placeholder="placeholder"
      :required="required"
      class="rounded-md border border-line bg-white px-3 py-2 text-sm text-ink placeholder:text-mist focus-visible:border-primary dark:border-graphite dark:bg-ink dark:text-surface"
      :class="mono && 'data-mono'"
      @input="$emit('update:modelValue', ($event.target as HTMLInputElement).value)"
    />
    <span v-if="error" class="text-xs font-medium" style="color: var(--color-status-canceled)">{{
      error
    }}</span>
  </label>
</template>
