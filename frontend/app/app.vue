<script setup lang="ts">
const appName = 'av-ing'
provide('appName', appName)
</script>

<template>
  <NuxtRouteAnnouncer />
  <NuxtPage />
</template>
