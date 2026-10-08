<script setup lang="ts">
import { cities } from '~/data/mock'

defineProps<{ open: boolean }>()
const emit = defineEmits<{ close: []; apply: []; reset: [] }>()
const city = defineModel<string>('city', { required: true })
const priceMin = defineModel<string>('priceMin', { required: true })
const priceMax = defineModel<string>('priceMax', { required: true })
</script>

<template>
  <Teleport to="body">
    <Transition name="fade"><button v-if="open" class="drawer-backdrop" aria-label="Закрыть фильтры" @click="emit('close')" /></Transition>
    <Transition name="drawer">
      <aside v-if="open" class="filter-drawer" aria-label="Фильтры">
        <div class="filter-drawer__header">
          <h2>Фильтры</h2>
          <button class="icon-button" type="button" aria-label="Закрыть" @click="emit('close')"><svg viewBox="0 0 24 24"><path d="M6 6l12 12M18 6 6 18" /></svg></button>
        </div>
        <div class="filter-group">
          <label for="filter-city">Город</label>
          <select id="filter-city" v-model="city" class="input-control"><option v-for="item in cities" :key="item">{{ item }}</option></select>
        </div>
        <div class="filter-group">
          <label>Цена, ₽</label>
          <div class="price-range">
            <input v-model="priceMin" class="input-control" inputmode="numeric" placeholder="От">
            <input v-model="priceMax" class="input-control" inputmode="numeric" placeholder="До">
          </div>
        </div>
        <div class="filter-group">
          <label>Состояние</label>
          <label class="check-row"><input type="checkbox"> Новое</label>
          <label class="check-row"><input type="checkbox"> Б/у</label>
        </div>
        <div class="filter-drawer__footer">
          <button class="button button--ghost" type="button" @click="emit('reset')">Сбросить</button>
          <button class="button button--primary" type="button" @click="emit('apply')">Показать</button>
        </div>
      </aside>
    </Transition>
  </Teleport>
</template>
