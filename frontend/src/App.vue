<script setup lang="ts">
import { useRoute } from 'vue-router'
import AppShell from '@/components/layout/AppShell.vue'
import ToastContainer from '@/components/ui/ToastContainer.vue'

const route = useRoute()
</script>

<template>
  <template v-if="route.meta.public">
    <RouterView />
  </template>
  <AppShell v-else>
    <RouterView />
  </AppShell>
  <ToastContainer />
</template>
