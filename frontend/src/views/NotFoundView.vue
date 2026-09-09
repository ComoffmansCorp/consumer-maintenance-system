<script setup lang="ts"></script>

<template>
  <div class="mk-home flex min-h-screen flex-col items-center justify-center gap-3 bg-[#FBFAF7] px-4 text-center">
    <span class="flex h-11 w-11 items-center justify-center rounded-[12px] bg-[#5B4BE0]">
      <svg viewBox="0 0 48 48" class="h-[26px] w-[26px]">
        <path fill-rule="evenodd" clip-rule="evenodd" fill="#FFFFFF"
          d="M24,9 L37.04,16.5 L37.04,31.5 L24,39 L10.96,31.5 L10.96,16.5 Z
             M30.19,24 a6.19,6.19 0 1,0 -12.38,0 a6.19,6.19 0 1,0 12.38,0 Z"/>
      </svg>
    </span>
    <p class="font-mono text-sm text-[#9B978A]">404</p>
    <h1 class="mk-display text-2xl font-bold tracking-tight text-[#17160F]">Страница не найдена</h1>
    <RouterLink to="/" class="text-sm font-medium text-[#5B4BE0] hover:underline">На главную</RouterLink>
  </div>
</template>
