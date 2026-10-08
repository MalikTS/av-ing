export interface Category {
  id: string
  name: string
  icon: CategoryIconName
  count: number
}

export type CategoryIconName = 'car' | 'home' | 'phone' | 'briefcase' | 'tools' | 'sofa' | 'shirt' | 'game' | 'paw'

export interface Ad {
  id: string
  title: string
  price: number
  categoryId: string
  city: string
  publishedAt: string
  views: number
  promoted?: boolean
  accent: string
  icon: CategoryIconName
}
