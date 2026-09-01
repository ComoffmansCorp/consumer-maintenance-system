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
    <span v-if="label" class="field-label">{{ label }}</span>
    <input
      :type="type"
      :value="modelValue ?? ''"
      :placeholder="placeholder"
      :required="required"
      class="field-underline"
      :class="mono && 'data-mono'"
      @input="$emit('update:modelValue', ($event.target as HTMLInputElement).value)"
    />
    <span v-if="error" class="text-xs font-medium" style="color: var(--color-status-canceled)">{{
      error
    }}</span>
  </label>
</template>
