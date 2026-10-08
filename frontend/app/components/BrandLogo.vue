<script setup lang="ts">
const appName = inject<string>('appName', 'av-ing')
</script>

<template>
  <NuxtLink to="/" class="brand" aria-label="На главную">
    <span class="brand__mark" aria-hidden="true">
      <svg viewBox="0 0 36 36"><path d="M18 3.8a10.7 10.7 0 0 0-10.7 10.7c0 7.6 10.7 17.7 10.7 17.7s10.7-10.1 10.7-17.7A10.7 10.7 0 0 0 18 3.8Z" /><path class="brand__mark-dot" d="M13.1 14.5a4.9 4.9 0 1 0 9.8 0 4.9 4.9 0 0 0-9.8 0Z" /></svg>
    </span>
    <span>{{ appName }}</span>
  </NuxtLink>
</template>
