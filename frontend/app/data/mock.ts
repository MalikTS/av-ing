import type { Ad, Category } from '~/types'

export const cities = ['Все города', 'Магас', 'Назрань', 'Сунжа', 'Карабулак', 'Малгобек', 'Другой']

export const categories: Category[] = [
  { id: 'transport', name: 'Транспорт', icon: 'car', count: 428 },
  { id: 'real-estate', name: 'Недвижимость', icon: 'home', count: 312 },
  { id: 'electronics', name: 'Электроника', icon: 'phone', count: 597 },
  { id: 'jobs', name: 'Работа', icon: 'briefcase', count: 184 },
  { id: 'services', name: 'Услуги', icon: 'tools', count: 263 },
  { id: 'home', name: 'Дом и сад', icon: 'sofa', count: 346 },
  { id: 'fashion', name: 'Одежда', icon: 'shirt', count: 221 },
  { id: 'hobby', name: 'Хобби', icon: 'game', count: 137 },
  { id: 'animals', name: 'Животные', icon: 'paw', count: 96 }
]

export const ads: Ad[] = [
  { id: '1', title: 'Toyota Camry, 2019', price: 2450000, categoryId: 'transport', city: 'Назрань', publishedAt: 'Сегодня, 10:24', views: 184, promoted: true, accent: 'linear-gradient(145deg, #d9e3e8, #8fa5af)', icon: 'car' },
  { id: '2', title: '2-комнатная квартира, 64 м²', price: 6300000, categoryId: 'real-estate', city: 'Магас', publishedAt: 'Сегодня, 09:10', views: 96, accent: 'linear-gradient(145deg, #e8e1d4, #baa98d)', icon: 'home' },
  { id: '3', title: 'iPhone 15 Pro, 256 ГБ', price: 92000, categoryId: 'electronics', city: 'Сунжа', publishedAt: 'Вчера, 21:45', views: 241, promoted: true, accent: 'linear-gradient(145deg, #dedde3, #9d9ba7)', icon: 'phone' },
  { id: '4', title: 'Диван в отличном состоянии', price: 38000, categoryId: 'home', city: 'Карабулак', publishedAt: 'Вчера, 18:30', views: 57, accent: 'linear-gradient(145deg, #d9e3dd, #8da398)', icon: 'sofa' },
  { id: '5', title: 'Ремонт квартир под ключ', price: 2500, categoryId: 'services', city: 'Назрань', publishedAt: '7 октября', views: 122, accent: 'linear-gradient(145deg, #e9e0d4, #b89b78)', icon: 'tools' },
  { id: '6', title: 'Горный велосипед', price: 34500, categoryId: 'hobby', city: 'Малгобек', publishedAt: '7 октября', views: 73, accent: 'linear-gradient(145deg, #d8e2e7, #7f9eac)', icon: 'game' },
  { id: '7', title: 'Мастер по ремонту телефонов', price: 1000, categoryId: 'services', city: 'Магас', publishedAt: '6 октября', views: 88, accent: 'linear-gradient(145deg, #e1e5e7, #a1adb2)', icon: 'tools' },
  { id: '8', title: 'Котята в добрые руки', price: 0, categoryId: 'animals', city: 'Сунжа', publishedAt: '6 октября', views: 205, accent: 'linear-gradient(145deg, #e8e1dc, #bfa998)', icon: 'paw' }
]
