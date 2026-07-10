<script setup lang="ts"></script>

<template>
  <div class="flex min-h-screen flex-col items-center justify-center gap-3 bg-surface px-4 text-center dark:bg-ink">
    <p class="data-mono text-sm text-mist">404</p>
    <h1 class="font-display text-xl font-semibold text-ink dark:text-surface">Страница не найдена</h1>
    <RouterLink to="/" class="text-sm font-medium text-primary hover:underline">На главную</RouterLink>
  </div>
</template>
