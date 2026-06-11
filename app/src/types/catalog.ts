export interface Cinema {
  id: string
  name: string
  address: string
  timezone: string
}

export interface Movie {
  id: string
  title: string
  posterUrl: string
  durationMin: number
  rating: string
  synopsis: string
  status: 'NOW_SHOWING' | 'COMING_SOON' | 'ARCHIVED'
}

export interface PriceTiers {
  standard: number
  vip: number
  wheelchair: number
}

export interface Showtime {
  id: string
  movieId: string
  screenId: string
  screenName: string
  startsAt: string
  priceTiers: PriceTiers
  status: string
}

export interface MovieDetail extends Movie {
  showtimes: Showtime[]
}

export type CatalogTab = 'now_showing' | 'coming_soon'

export interface ApiErrorBody {
  error: {
    code: string
    message: string
  }
}
