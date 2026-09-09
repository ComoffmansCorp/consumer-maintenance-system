<script setup lang="ts">
withDefaults(defineProps<{ rounded?: string }>(), { rounded: 'rounded-xl' })
</script>

<template>
  <div class="animate-pulse bg-[#EFEBE1]" :class="rounded" />
</template>
