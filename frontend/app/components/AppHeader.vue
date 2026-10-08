<script setup lang="ts">
const emit = defineEmits<{ location: [] }>()
</script>

<template>
  <header class="site-header">
    <div class="shell site-header__inner">
      <BrandLogo />
      <button class="location-button desktop-only" type="button" @click="emit('location')">
        <svg viewBox="0 0 24 24" aria-hidden="true"><path d="M12 21s7-6.2 7-12a7 7 0 1 0-14 0c0 5.8 7 12 7 12Z" /><circle cx="12" cy="9" r="2.3" /></svg>
        Республика Ингушетия
      </button>
      <nav class="header-actions" aria-label="Основная навигация">
        <button class="icon-button desktop-only" type="button" aria-label="Уведомления">
          <svg viewBox="0 0 24 24"><path d="M18 8a6 6 0 0 0-12 0c0 7-3 7-3 9h18c0-2-3-2-3-9M10 21h4" /></svg>
          <span class="notification-dot" />
        </button>
        <NuxtLink to="#favorites" class="header-link desktop-only">Избранное</NuxtLink>
        <NuxtLink to="#profile" class="header-link desktop-only">Войти</NuxtLink>
        <button class="button button--primary publish-button" type="button">
          <svg viewBox="0 0 24 24" aria-hidden="true"><path d="M12 5v14M5 12h14" /></svg>
          <span class="publish-button__full">Разместить объявление</span>
          <span class="publish-button__short">Разместить</span>
        </button>
      </nav>
    </div>
  </header>
</template>
