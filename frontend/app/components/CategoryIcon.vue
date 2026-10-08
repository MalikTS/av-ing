<script setup lang="ts">
import type { CategoryIconName } from '~/types'
defineProps<{ name: CategoryIconName }>()
</script>

<template>
  <svg v-if="name === 'car'" viewBox="0 0 24 24" aria-hidden="true"><path d="M5 17h14l1.2-4.2a2 2 0 0 0-1.9-2.6H5.7a2 2 0 0 0-1.9 2.6L5 17Zm2-6.8 1.2-3h7.6l1.2 3M7 17v2m10-2v2M7.5 14h.01m8.99 0h.01" /></svg>
  <svg v-else-if="name === 'home'" viewBox="0 0 24 24" aria-hidden="true"><path d="m3.5 11 8.5-7 8.5 7M5.5 9.5V20h13V9.5M9.5 20v-6h5v6" /></svg>
  <svg v-else-if="name === 'phone'" viewBox="0 0 24 24" aria-hidden="true"><rect x="7" y="2.5" width="10" height="19" rx="2" /><path d="M10 5h4m-3 13.5h2" /></svg>
  <svg v-else-if="name === 'briefcase'" viewBox="0 0 24 24" aria-hidden="true"><rect x="3" y="7" width="18" height="13" rx="2" /><path d="M9 7V4h6v3m-12 5.5c5.4 2 12.6 2 18 0M10 13h4" /></svg>
  <svg v-else-if="name === 'tools'" viewBox="0 0 24 24" aria-hidden="true"><path d="M14.2 6.1a4.5 4.5 0 0 0-5.7 5.7L3 17.3 6.7 21l5.5-5.5a4.5 4.5 0 0 0 5.7-5.7l-3 3-2.7-2.7 3-3Z" /></svg>
  <svg v-else-if="name === 'sofa'" viewBox="0 0 24 24" aria-hidden="true"><path d="M5 12V8.5A2.5 2.5 0 0 1 7.5 6h9A2.5 2.5 0 0 1 19 8.5V12M5 11a2 2 0 0 0-2 2v5h18v-5a2 2 0 0 0-2-2 2 2 0 0 0-2 2v1H7v-1a2 2 0 0 0-2-2Zm0 7v2m14-2v2" /></svg>
  <svg v-else-if="name === 'shirt'" viewBox="0 0 24 24" aria-hidden="true"><path d="m8 4-5 3 2.5 4L8 9.5V20h8V9.5l2.5 1.5L21 7l-5-3c-.7 1.3-2 2-4 2S8.7 5.3 8 4Z" /></svg>
  <svg v-else-if="name === 'game'" viewBox="0 0 24 24" aria-hidden="true"><path d="M8.5 8h7a5.5 5.5 0 0 1 5.2 7.3l-.8 2.3a2.2 2.2 0 0 1-3.7.8L14.8 17H9.2l-1.4 1.4a2.2 2.2 0 0 1-3.7-.8l-.8-2.3A5.5 5.5 0 0 1 8.5 8ZM7 11v4m-2-2h4m7-1h.01m2 2h.01M10 8V5h4" /></svg>
  <svg v-else viewBox="0 0 24 24" aria-hidden="true"><path d="M8.2 10.5c-1.4-1.9-3.4-2.5-4.5-1.2s-.3 3.4 1.6 4.3m10.5-3.1c1.4-1.9 3.4-2.5 4.5-1.2s.3 3.4-1.6 4.3M8.5 7C8 4.8 8.9 3 10.3 3s2.2 1.8 1.7 4m3.5 0c.5-2.2-.4-4-1.8-4S11.5 4.8 12 7m0 3c-4 0-7 3-7 6.2C5 19 7.2 21 9.8 20c1.4-.5 3-.5 4.4 0 2.6 1 4.8-1 4.8-3.8C19 13 16 10 12 10Z" /></svg>
</template>

<style scoped>
svg { width: 1.65rem; height: 1.65rem; fill: none; stroke: currentColor; stroke-width: 1.65; stroke-linecap: round; stroke-linejoin: round; }
</style>
