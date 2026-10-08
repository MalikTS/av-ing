<script setup lang="ts">
import type { Category } from '~/types'

defineProps<{ categories: Category[]; selected: string }>()
const emit = defineEmits<{ select: [id: string] }>()
</script>

<template>
  <div class="category-grid">
    <button v-for="category in categories" :key="category.id" class="category-card" :class="{ 'category-card--active': selected === category.id }" type="button" @click="emit('select', selected === category.id ? '' : category.id)">
      <span class="category-card__icon"><CategoryIcon :name="category.icon" /></span>
      <span class="category-card__name">{{ category.name }}</span>
      <span class="category-card__count">{{ category.count }} объявлений</span>
    </button>
  </div>
</template>
