<script setup lang="ts">
import { computed } from 'vue'

// One hand-drawn glyph per seeded category, matched by exact name -- falls
// back to a generic tool icon for anything an admin adds later that isn't
// in this list yet. Single-tone (currentColor), matching the icon language
// already used elsewhere (header gauge glyph, seal marks) rather than
// introducing a clashing full-color illustration style.
const props = defineProps<{ name: string }>()

const kind = computed(() => {
  const known: Record<string, string> = {
    'Приборы учёта': 'meter',
    Сантехника: 'plumbing',
    Электрика: 'electric',
    Клининг: 'cleaning',
    'Бытовая техника': 'appliance',
    Мебель: 'furniture',
    'Компьютеры и телефоны': 'devices',
    'Ремонт и строительство': 'construction',
    Авто: 'auto',
    Репетиторство: 'tutoring',
    'Красота и здоровье': 'beauty',
    'Переезды и грузоперевозки': 'moving',
    'Няни и уход': 'care',
  }
  return known[props.name] ?? 'tool'
})
</script>

<template>
  <svg viewBox="0 0 24 24" class="h-6 w-6" aria-hidden="true">
    <template v-if="kind === 'meter'">
      <circle cx="12" cy="12" r="8.5" fill="none" stroke="currentColor" stroke-width="1.6" />
      <circle cx="12" cy="12" r="1.8" fill="currentColor" />
      <path d="M12 6v2M12 16v2M18 12h-2M8 12H6" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" />
    </template>
    <template v-else-if="kind === 'plumbing'">
      <path
        d="M6 5v4a4 4 0 0 0 4 4h1v3a3 3 0 1 0 3-3V9a4 4 0 0 0-4-4H6Z"
        fill="none"
        stroke="currentColor"
        stroke-width="1.6"
        stroke-linejoin="round"
      />
    </template>
    <template v-else-if="kind === 'electric'">
      <path d="M13 3 5 13h5l-1 8 8-10h-5l1-8Z" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linejoin="round" />
    </template>
    <template v-else-if="kind === 'cleaning'">
      <path d="M9 3v6l-3 9a3 3 0 0 0 3 3h6a3 3 0 0 0 3-3l-3-9V3" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linejoin="round" />
      <path d="M7 3h10" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" />
    </template>
    <template v-else-if="kind === 'appliance'">
      <rect x="5" y="3" width="14" height="18" rx="1.5" fill="none" stroke="currentColor" stroke-width="1.6" />
      <circle cx="12" cy="13" r="4" fill="none" stroke="currentColor" stroke-width="1.6" />
      <path d="M8 6.5h.01M11 6.5h.01" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" />
    </template>
    <template v-else-if="kind === 'furniture'">
      <path d="M4 10V6a2 2 0 0 1 2-2h12a2 2 0 0 1 2 2v4" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" />
      <path d="M4 10h16v4H4z" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linejoin="round" />
      <path d="M5 14v6M19 14v6" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" />
    </template>
    <template v-else-if="kind === 'devices'">
      <rect x="3" y="4" width="13" height="9" rx="1" fill="none" stroke="currentColor" stroke-width="1.6" />
      <path d="M3 16h13M8 20h3" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" />
      <rect x="17.5" y="7" width="4.5" height="12" rx="1" fill="none" stroke="currentColor" stroke-width="1.5" />
    </template>
    <template v-else-if="kind === 'construction'">
      <path d="M4 20 14 10" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" />
      <path
        d="m13 6 5 5-1.5 1.5a3.5 3.5 0 0 1-5-5L13 6Z"
        fill="none"
        stroke="currentColor"
        stroke-width="1.6"
        stroke-linejoin="round"
      />
    </template>
    <template v-else-if="kind === 'auto'">
      <path
        d="M4 15v-2l2-4h12l2 4v2"
        fill="none"
        stroke="currentColor"
        stroke-width="1.6"
        stroke-linejoin="round"
      />
      <rect x="3" y="15" width="18" height="4" rx="1" fill="none" stroke="currentColor" stroke-width="1.6" />
      <circle cx="7" cy="19" r="1.3" fill="currentColor" />
      <circle cx="17" cy="19" r="1.3" fill="currentColor" />
    </template>
    <template v-else-if="kind === 'tutoring'">
      <path d="M3 8 12 4l9 4-9 4-9-4Z" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linejoin="round" />
      <path d="M7 10.5V15c0 1.5 2.2 3 5 3s5-1.5 5-3v-4.5" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" />
    </template>
    <template v-else-if="kind === 'beauty'">
      <path
        d="M12 3c2 3 4 5.5 4 8.5a4 4 0 1 1-8 0C8 8.5 10 6 12 3Z"
        fill="none"
        stroke="currentColor"
        stroke-width="1.6"
        stroke-linejoin="round"
      />
      <path d="M12 15.5V21" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" />
    </template>
    <template v-else-if="kind === 'moving'">
      <rect x="3" y="9" width="11" height="9" rx="1" fill="none" stroke="currentColor" stroke-width="1.6" />
      <path d="M14 12h4l3 3v3h-7" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linejoin="round" />
      <circle cx="7.5" cy="19.5" r="1.3" fill="currentColor" />
      <circle cx="17" cy="19.5" r="1.3" fill="currentColor" />
    </template>
    <template v-else-if="kind === 'care'">
      <path
        d="M12 20s-7-4.4-7-9.5A4 4 0 0 1 12 8a4 4 0 0 1 7 2.5C19 15.6 12 20 12 20Z"
        fill="none"
        stroke="currentColor"
        stroke-width="1.6"
        stroke-linejoin="round"
      />
    </template>
    <template v-else>
      <path
        d="M14.7 6.3a4 4 0 0 1-5.4 5.4L4 17l3 3 5.3-5.3a4 4 0 0 1 5.4-5.4L14.7 6.3Z"
        fill="none"
        stroke="currentColor"
        stroke-width="1.6"
        stroke-linejoin="round"
      />
    </template>
  </svg>
</template>
