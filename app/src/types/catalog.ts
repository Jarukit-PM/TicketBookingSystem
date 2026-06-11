export interface Cinema {
  id: string
  name: string
  address: string
  timezone: string
}

export type MovieStatus = 'NOW_SHOWING' | 'COMING_SOON' | 'ARCHIVED'

export interface Movie {
  id: string
  title: string
  posterUrl: string
  durationMin: number
  rating: string
  synopsis: string
  status: MovieStatus
}

export interface LayoutSeat {
  seatId: string
  row: number
  col: number
  type: string
}

export interface ScreenLayout {
  seats: LayoutSeat[]
}

export interface Screen {
  id: string
  cinemaId: string
  name: string
  layout: ScreenLayout
}

export type ShowtimeStatus = 'OPEN' | 'CANCELLED'

export interface PriceTiers {
  standard: number
  vip: number
  wheelchair: number
}

export interface Showtime {
  id: string
  movieId: string
  screenId: string
  screenName?: string
  startsAt: string
  priceTiers: PriceTiers
  status: ShowtimeStatus | string
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
