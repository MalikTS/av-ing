<template>
  <nav class="mobile-nav" aria-label="Мобильная навигация">
    <NuxtLink to="/" class="mobile-nav__item mobile-nav__item--active"><svg viewBox="0 0 24 24"><path d="m3 11 9-8 9 8v10h-6v-6H9v6H3V11Z" /></svg><span>Главная</span></NuxtLink>
    <NuxtLink to="#catalog" class="mobile-nav__item"><svg viewBox="0 0 24 24"><rect x="3" y="3" width="7" height="7" rx="1" /><rect x="14" y="3" width="7" height="7" rx="1" /><rect x="3" y="14" width="7" height="7" rx="1" /><rect x="14" y="14" width="7" height="7" rx="1" /></svg><span>Каталог</span></NuxtLink>
    <button class="mobile-nav__publish" type="button" aria-label="Разместить объявление"><svg viewBox="0 0 24 24"><path d="M12 5v14M5 12h14" /></svg></button>
    <NuxtLink to="#favorites" class="mobile-nav__item"><svg viewBox="0 0 24 24"><path d="M20.8 4.6a5.5 5.5 0 0 0-7.8 0L12 5.7l-1.1-1.1a5.5 5.5 0 0 0-7.8 7.8l1.1 1.1L12 21l7.8-7.5 1.1-1.1a5.5 5.5 0 0 0-.1-7.8Z" /></svg><span>Избранное</span></NuxtLink>
    <NuxtLink to="#profile" class="mobile-nav__item"><svg viewBox="0 0 24 24"><circle cx="12" cy="8" r="4" /><path d="M4 21a8 8 0 0 1 16 0" /></svg><span>Профиль</span></NuxtLink>
  </nav>
</template>
