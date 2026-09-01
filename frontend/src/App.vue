<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import AppShell from '@/components/layout/AppShell.vue'
import ToastContainer from '@/components/ui/ToastContainer.vue'

const route = useRoute()

// The public marketplace never gets the tenant admin's sidebar shell, even
// on its authenticated (CLIENT/MASTER) routes -- `meta.public` alone can't
// decide this since marketplace's own auth'd pages aren't `public: true`.
const useAdminShell = computed(() => !route.meta.public && !route.path.startsWith('/marketplace'))
</script>

<template>
  <AppShell v-if="useAdminShell">
    <RouterView />
  </AppShell>
  <RouterView v-else />
  <ToastContainer />
</template>
