<script setup lang="ts">
withDefaults(
  defineProps<{
    modelValue: string | number | null
    label?: string
    options: { value: string | number; label: string }[]
    placeholder?: string
    required?: boolean
  }>(),
  { required: false },
)
defineEmits<{ 'update:modelValue': [value: string] }>()
</script>

<template>
  <label class="flex flex-col gap-1.5">
    <span v-if="label" class="field-label">{{ label }}</span>
    <select
      :value="modelValue ?? ''"
      :required="required"
      class="field-underline"
      @change="$emit('update:modelValue', ($event.target as HTMLSelectElement).value)"
    >
      <option v-if="placeholder" value="" disabled>{{ placeholder }}</option>
      <option v-for="opt in options" :key="opt.value" :value="opt.value">{{ opt.label }}</option>
    </select>
  </label>
</template>
