<script setup lang="ts">
import type { Ad } from '~/types'

const props = defineProps<{ ad: Ad; favorite: boolean }>()
const emit = defineEmits<{ favorite: [id: string] }>()
const formattedPrice = computed(() => props.ad.price === 0 ? 'Бесплатно' : `${new Intl.NumberFormat('ru-RU').format(props.ad.price)} ₽`)
</script>

<template>
  <article class="ad-card">
    <NuxtLink :to="`/ads/${ad.id}`" class="ad-card__visual" :style="{ background: ad.accent }" :aria-label="ad.title">
      <span v-if="ad.promoted" class="ad-card__badge">Продвигается</span>
      <span class="ad-card__placeholder" aria-hidden="true"><CategoryIcon :name="ad.icon" /></span>
      <span class="ad-card__image-count"><svg viewBox="0 0 24 24"><rect x="3" y="5" width="18" height="14" rx="2" /><circle cx="8.5" cy="10" r="1.5" /><path d="m21 15-4.5-4.5L8 19" /></svg>{{ Number(ad.id) % 5 + 2 }}</span>
    </NuxtLink>
    <button class="favorite-button" :class="{ 'favorite-button--active': favorite }" type="button" :aria-label="favorite ? 'Удалить из избранного' : 'Добавить в избранное'" @click="emit('favorite', ad.id)">
      <svg viewBox="0 0 24 24"><path d="M20.8 4.6a5.5 5.5 0 0 0-7.8 0L12 5.7l-1.1-1.1a5.5 5.5 0 0 0-7.8 7.8l1.1 1.1L12 21l7.8-7.5 1.1-1.1a5.5 5.5 0 0 0-.1-7.8Z" /></svg>
    </button>
    <div class="ad-card__body">
      <NuxtLink :to="`/ads/${ad.id}`" class="ad-card__title">{{ ad.title }}</NuxtLink>
      <p class="ad-card__price">{{ formattedPrice }}</p>
      <div class="ad-card__meta"><span>{{ ad.city }}</span><span>{{ ad.publishedAt }}</span></div>
    </div>
  </article>
</template>
