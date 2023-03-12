type Directory = {
  id: string
  name: string
  loved: boolean
  thumbnail: string
  lastModified: Date
}

type Paginated<T> = {
  list: T[]
  page: number
  pageSize: number
  pages: number
  totalElements: number
}

type GalleryResponse = {
  avif: string[]
  avifAvailable: boolean
  webp: string[]
  webpAvailable: boolean
  source: string[]
  cached: boolean
}