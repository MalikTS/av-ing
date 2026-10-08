export default defineNuxtConfig({
  compatibilityDate: '2026-10-01',
  devtools: { enabled: false },
  css: ['~/assets/css/main.css'],
  app: {
    head: {
      htmlAttrs: { lang: 'ru' },
      title: 'av-ing — объявления Ингушетии',
      meta: [
        { name: 'viewport', content: 'width=device-width, initial-scale=1' },
        { name: 'description', content: 'Объявления о товарах и услугах в Республике Ингушетия.' },
        { name: 'theme-color', content: '#12354a' }
      ]
    }
  }
})
