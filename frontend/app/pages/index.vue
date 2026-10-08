<script setup lang="ts">
import { ads, categories } from '~/data/mock'

const search = ref('')
const city = ref('Все города')
const selectedCategory = ref('')
const priceMin = ref('')
const priceMax = ref('')
const sort = ref('newest')
const filtersOpen = ref(false)
const favorites = ref(new Set<string>(['3']))

const filteredAds = computed(() => {
  const query = search.value.trim().toLocaleLowerCase('ru')
  const min = Number(priceMin.value) || 0
  const max = Number(priceMax.value) || Number.POSITIVE_INFINITY
  const result = ads.filter((ad) => {
    const matchesSearch = !query || ad.title.toLocaleLowerCase('ru').includes(query)
    const matchesCity = city.value === 'Все города' || ad.city === city.value
    const matchesCategory = !selectedCategory.value || ad.categoryId === selectedCategory.value
    return matchesSearch && matchesCity && matchesCategory && ad.price >= min && ad.price <= max
  })
  return [...result].sort((a, b) => {
    if (sort.value === 'price-low') return a.price - b.price
    if (sort.value === 'price-high') return b.price - a.price
    return Number(b.promoted ?? false) - Number(a.promoted ?? false)
  })
})

function toggleFavorite(id: string) {
  const next = new Set(favorites.value)
  next.has(id) ? next.delete(id) : next.add(id)
  favorites.value = next
}

function resetFilters() {
  city.value = 'Все города'
  priceMin.value = ''
  priceMax.value = ''
  selectedCategory.value = ''
}

useSeoMeta({ title: 'av-ing — объявления Ингушетии', description: 'Товары, недвижимость, транспорт и услуги в городах Республики Ингушетия.' })
</script>

<template>
  <div class="page-shell">
    <AppHeader @location="filtersOpen = true" />
    <main>
      <HeroSearch v-model:search="search" v-model:city="city" @filters="filtersOpen = true" />

      <section id="catalog" class="section shell">
        <div class="section-heading">
          <div><p class="eyebrow eyebrow--dark">Найдите нужное</p><h2>Категории</h2></div>
          <button v-if="selectedCategory" class="text-button" type="button" @click="selectedCategory = ''">Показать все</button>
        </div>
        <CategoryGrid :categories="categories" :selected="selectedCategory" @select="selectedCategory = $event" />
      </section>

      <section class="section section--listings">
        <div class="shell">
          <div class="listings-toolbar">
            <div><p class="eyebrow eyebrow--dark">Рядом с вами</p><h2>{{ selectedCategory ? categories.find(item => item.id === selectedCategory)?.name : 'Свежие объявления' }}</h2></div>
            <div class="toolbar-actions">
              <button class="button button--outline filter-button" type="button" @click="filtersOpen = true"><svg viewBox="0 0 24 24"><path d="M4 7h16M7 12h10m-7 5h4" /></svg>Фильтры</button>
              <label class="sort-select"><span class="sr-only">Сортировка</span><select v-model="sort"><option value="newest">Сначала новые</option><option value="price-low">Сначала дешевле</option><option value="price-high">Сначала дороже</option></select></label>
            </div>
          </div>

          <div v-if="filteredAds.length" class="ad-grid">
            <AdCard v-for="ad in filteredAds" :key="ad.id" :ad="ad" :favorite="favorites.has(ad.id)" @favorite="toggleFavorite" />
          </div>
          <div v-else class="empty-state">
            <span class="empty-state__icon"><CategoryIcon name="briefcase" /></span>
            <h3>Объявления не найдены</h3><p>Попробуйте изменить запрос или сбросить фильтры.</p>
            <button class="button button--outline" type="button" @click="resetFilters(); search = ''">Сбросить фильтры</button>
          </div>
          <button v-if="filteredAds.length >= 6" class="button button--load-more" type="button">Показать ещё</button>
        </div>
      </section>

      <section class="trust-section">
        <div class="shell trust-section__inner">
          <div><p class="eyebrow">Просто и по-соседски</p><h2>Покупайте и продавайте<br>в своём регионе</h2></div>
          <div class="trust-points"><div><b>01</b><span>Разместите объявление за пару минут</span></div><div><b>02</b><span>Общайтесь напрямую без посредников</span></div><div><b>03</b><span>Оценивайте продавцов и делитесь опытом</span></div></div>
        </div>
      </section>
    </main>

    <footer class="site-footer">
      <div class="shell site-footer__inner"><BrandLogo /><p>Объявления Республики Ингушетия</p><nav><a href="#">Помощь</a><a href="#">Правила</a><a href="#">Конфиденциальность</a></nav><span>© 2026</span></div>
    </footer>

    <FilterDrawer v-model:city="city" v-model:price-min="priceMin" v-model:price-max="priceMax" :open="filtersOpen" @close="filtersOpen = false" @apply="filtersOpen = false" @reset="resetFilters" />
    <MobileNav />
  </div>
</template>
