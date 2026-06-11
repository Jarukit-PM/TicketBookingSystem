import { apiGet } from '@/api/client'
import type { CatalogTab, Cinema, Movie, MovieDetail, Showtime } from '@/types/catalog'

export function fetchCinemas(): Promise<Cinema[]> {
  return apiGet<Cinema[]>('/cinemas')
}

export function fetchMovies(cinemaId: string, tab: CatalogTab): Promise<Movie[]> {
  return apiGet<Movie[]>('/movies', { cinemaId, tab })
}

export function fetchMovieDetail(movieId: string, cinemaId: string): Promise<MovieDetail> {
  return apiGet<MovieDetail>(`/movies/${movieId}`, { cinemaId })
}

export function fetchShowtimes(
  cinemaId: string,
  movieId: string,
  date?: string,
): Promise<Showtime[]> {
  const params: Record<string, string> = { cinemaId, movieId }
  if (date) {
    params.date = date
  }
  return apiGet<Showtime[]>('/showtimes', params)
}
