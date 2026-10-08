<script setup lang="ts">
import { cities } from '~/data/mock'

const search = defineModel<string>('search', { required: true })
const city = defineModel<string>('city', { required: true })
const emit = defineEmits<{ filters: [] }>()
</script>

<template>
  <section class="hero">
    <div class="shell hero__inner">
      <div class="hero__copy">
        <p class="eyebrow">Всё нужное — рядом</p>
        <h1>Объявления<br><span>Ингушетии</span></h1>
        <p class="hero__description">Покупайте, продавайте и находите услуги у людей поблизости.</p>
      </div>
      <form class="search-panel" @submit.prevent>
        <label class="search-field search-field--main">
          <span class="sr-only">Поиск объявлений</span>
          <svg viewBox="0 0 24 24" aria-hidden="true"><circle cx="11" cy="11" r="7" /><path d="m20 20-4-4" /></svg>
          <input v-model="search" type="search" placeholder="Что вы ищете?" autocomplete="off">
        </label>
        <label class="search-field search-field--city">
          <span class="sr-only">Город</span>
          <svg viewBox="0 0 24 24" aria-hidden="true"><path d="M12 21s7-6.2 7-12a7 7 0 1 0-14 0c0 5.8 7 12 7 12Z" /><circle cx="12" cy="9" r="2.3" /></svg>
          <select v-model="city"><option v-for="item in cities" :key="item" :value="item">{{ item }}</option></select>
        </label>
        <button class="button button--accent search-submit" type="submit">Найти</button>
        <button class="filter-trigger" type="button" aria-label="Открыть фильтры" @click="emit('filters')">
          <svg viewBox="0 0 24 24"><path d="M4 7h16M7 12h10m-7 5h4" /></svg>
        </button>
      </form>
      <div class="hero__facts" aria-label="Преимущества сервиса">
        <span><b>2 584</b> активных объявления</span>
        <span><b>9</b> категорий</span>
        <span><b>6</b> городов</span>
      </div>
    </div>
  </section>
</template>
