<script setup lang="ts">
import { onMounted, ref } from 'vue'

import { ApiError, api } from '@/api/client'
import { Button, Card, CardContent, CardHeader, CardTitle, Input } from '@/components/ui'
import type { Cinema, Movie, Screen, Showtime, ShowtimeStatus } from '@/types/catalog'

const movies = ref<Movie[]>([])
const screens = ref<Screen[]>([])
const cinemas = ref<Cinema[]>([])
const showtimes = ref<Showtime[]>([])
const loading = ref(true)
const errorMessage = ref('')
const editingId = ref<string | null>(null)

const filter = ref({
  cinemaId: '',
  movieId: '',
})

const form = ref({
  movieId: '',
  screenId: '',
  startsAtLocal: '',
  standard: 1200,
  vip: 1800,
  wheelchair: 1200,
  status: 'OPEN' as ShowtimeStatus,
})

function cinemaNameForScreen(screenId: string) {
  const screen = screens.value.find((s) => s.id === screenId)
  if (!screen) return ''
  return cinemas.value.find((c) => c.id === screen.cinemaId)?.name ?? ''
}

function movieTitle(id: string) {
  return movies.value.find((m) => m.id === id)?.title ?? id
}

function screenName(id: string) {
  return screens.value.find((s) => s.id === id)?.name ?? id
}

function toLocalInputValue(iso: string) {
  const date = new Date(iso)
  const pad = (n: number) => String(n).padStart(2, '0')
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())}T${pad(date.getHours())}:${pad(date.getMinutes())}`
}

function resetForm() {
  editingId.value = null
  const tomorrow = new Date()
  tomorrow.setDate(tomorrow.getDate() + 1)
  tomorrow.setMinutes(0, 0, 0)
  form.value = {
    movieId: movies.value[0]?.id || '',
    screenId: screens.value[0]?.id || '',
    startsAtLocal: toLocalInputValue(tomorrow.toISOString()),
    standard: 1200,
    vip: 1800,
    wheelchair: 1200,
    status: 'OPEN',
  }
}

function buildQuery() {
  const params = new URLSearchParams()
  if (filter.value.cinemaId) params.set('cinemaId', filter.value.cinemaId)
  if (filter.value.movieId) params.set('movieId', filter.value.movieId)
  const query = params.toString()
  return query ? `?${query}` : ''
}

async function loadReferenceData() {
  const [movieRes, screenRes, cinemaRes] = await Promise.all([
    api.get<{ movies: Movie[] }>('/admin/movies'),
    api.get<{ screens: Screen[] }>('/admin/screens'),
    api.get<{ cinemas: Cinema[] }>('/admin/cinemas'),
  ])
  movies.value = movieRes.movies
  screens.value = screenRes.screens
  cinemas.value = cinemaRes.cinemas
}

async function loadShowtimes() {
  loading.value = true
  errorMessage.value = ''
  try {
    const response = await api.get<{ showtimes: Showtime[] }>(
      `/admin/showtimes${buildQuery()}`,
    )
    showtimes.value = response.showtimes
  } catch (error) {
    errorMessage.value = error instanceof ApiError ? error.message : 'Failed to load showtimes'
  } finally {
    loading.value = false
  }
}

async function loadAll() {
  loading.value = true
  try {
    await loadReferenceData()
    resetForm()
    await loadShowtimes()
  } catch (error) {
    errorMessage.value = error instanceof ApiError ? error.message : 'Failed to load data'
    loading.value = false
  }
}

function startEdit(showtime: Showtime) {
  editingId.value = showtime.id
  form.value = {
    movieId: showtime.movieId,
    screenId: showtime.screenId,
    startsAtLocal: toLocalInputValue(showtime.startsAt),
    standard: showtime.priceTiers.standard,
    vip: showtime.priceTiers.vip,
    wheelchair: showtime.priceTiers.wheelchair,
    status: showtime.status as ShowtimeStatus,
  }
}

async function onSubmit() {
  errorMessage.value = ''
  try {
    const payload = {
      movieId: form.value.movieId,
      screenId: form.value.screenId,
      startsAt: new Date(form.value.startsAtLocal).toISOString(),
      priceTiers: {
        standard: Number(form.value.standard),
        vip: Number(form.value.vip),
        wheelchair: Number(form.value.wheelchair),
      },
      status: form.value.status,
    }
    if (editingId.value) {
      await api.put(`/admin/showtimes/${editingId.value}`, payload)
    } else {
      await api.post('/admin/showtimes', payload)
    }
    resetForm()
    await loadShowtimes()
  } catch (error) {
    errorMessage.value = error instanceof ApiError ? error.message : 'Save failed'
  }
}

async function onDelete(id: string) {
  if (!confirm('Delete this showtime?')) return
  try {
    await api.delete(`/admin/showtimes/${id}`)
    if (editingId.value === id) resetForm()
    await loadShowtimes()
  } catch (error) {
    errorMessage.value = error instanceof ApiError ? error.message : 'Delete failed'
  }
}

onMounted(loadAll)
</script>

<template>
  <div class="space-y-8">
    <div>
      <h1 class="text-2xl font-semibold text-copy-primary">Showtimes</h1>
      <p class="mt-1 text-sm text-copy-secondary">
        Schedule screenings by linking a movie to a screen.
      </p>
    </div>

    <Card>
      <CardHeader>
        <CardTitle>{{ editingId ? 'Edit showtime' : 'Add showtime' }}</CardTitle>
      </CardHeader>
      <CardContent>
        <form class="grid gap-4 md:grid-cols-2" @submit.prevent="onSubmit">
          <div class="space-y-2">
            <label class="text-sm text-copy-secondary" for="movieId">Movie</label>
            <select
              id="movieId"
              v-model="form.movieId"
              class="w-full rounded-lg border border-surface-border bg-surface px-3 py-2 text-sm text-copy-primary"
              required
            >
              <option v-for="movie in movies" :key="movie.id" :value="movie.id">
                {{ movie.title }}
              </option>
            </select>
          </div>
          <div class="space-y-2">
            <label class="text-sm text-copy-secondary" for="screenId">Screen</label>
            <select
              id="screenId"
              v-model="form.screenId"
              class="w-full rounded-lg border border-surface-border bg-surface px-3 py-2 text-sm text-copy-primary"
              required
            >
              <option v-for="screen in screens" :key="screen.id" :value="screen.id">
                {{ screen.name }} ({{ cinemaNameForScreen(screen.id) }})
              </option>
            </select>
          </div>
          <div class="space-y-2">
            <label class="text-sm text-copy-secondary" for="startsAt">Starts at</label>
            <Input id="startsAt" v-model="form.startsAtLocal" type="datetime-local" required />
          </div>
          <div class="space-y-2">
            <label class="text-sm text-copy-secondary" for="showtimeStatus">Status</label>
            <select
              id="showtimeStatus"
              v-model="form.status"
              class="w-full rounded-lg border border-surface-border bg-surface px-3 py-2 text-sm text-copy-primary"
            >
              <option value="OPEN">OPEN</option>
              <option value="CANCELLED">CANCELLED</option>
            </select>
          </div>
          <div class="space-y-2">
            <label class="text-sm text-copy-secondary" for="standard">Standard (cents)</label>
            <Input id="standard" :model-value="String(form.standard)" @update:model-value="form.standard = Number($event) || 0" type="number" min="0" required />
          </div>
          <div class="space-y-2">
            <label class="text-sm text-copy-secondary" for="vip">VIP (cents)</label>
            <Input id="vip" :model-value="String(form.vip)" @update:model-value="form.vip = Number($event) || 0" type="number" min="0" required />
          </div>
          <div class="space-y-2 md:col-span-2">
            <label class="text-sm text-copy-secondary" for="wheelchair">Wheelchair (cents)</label>
            <Input id="wheelchair" :model-value="String(form.wheelchair)" @update:model-value="form.wheelchair = Number($event) || 0" type="number" min="0" required />
          </div>
          <div class="flex gap-2 md:col-span-2">
            <Button type="submit" :disabled="!movies.length || !screens.length">
              {{ editingId ? 'Update' : 'Create' }}
            </Button>
            <Button v-if="editingId" type="button" variant="secondary" @click="resetForm">
              Cancel
            </Button>
          </div>
        </form>
      </CardContent>
    </Card>

    <p v-if="errorMessage" class="text-sm text-state-error" role="alert">{{ errorMessage }}</p>

    <Card>
      <CardHeader>
        <div class="flex flex-wrap items-end gap-3">
          <CardTitle>All showtimes</CardTitle>
          <select
            v-model="filter.cinemaId"
            class="rounded-lg border border-surface-border bg-surface px-3 py-2 text-sm text-copy-primary"
            @change="loadShowtimes"
          >
            <option value="">All cinemas</option>
            <option v-for="cinema in cinemas" :key="cinema.id" :value="cinema.id">
              {{ cinema.name }}
            </option>
          </select>
          <select
            v-model="filter.movieId"
            class="rounded-lg border border-surface-border bg-surface px-3 py-2 text-sm text-copy-primary"
            @change="loadShowtimes"
          >
            <option value="">All movies</option>
            <option v-for="movie in movies" :key="movie.id" :value="movie.id">
              {{ movie.title }}
            </option>
          </select>
        </div>
      </CardHeader>
      <CardContent>
        <p v-if="loading" class="text-sm text-copy-muted">Loading…</p>
        <div v-else class="overflow-x-auto">
          <table class="w-full text-left text-sm">
            <thead class="text-copy-muted">
              <tr>
                <th class="pb-3 pr-4 font-medium">Movie</th>
                <th class="pb-3 pr-4 font-medium">Screen</th>
                <th class="pb-3 pr-4 font-medium">Starts</th>
                <th class="pb-3 pr-4 font-medium">Status</th>
                <th class="pb-3 font-medium">Actions</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="showtime in showtimes"
                :key="showtime.id"
                class="border-t border-surface-border"
              >
                <td class="py-3 pr-4 text-copy-primary">{{ movieTitle(showtime.movieId) }}</td>
                <td class="py-3 pr-4 text-copy-secondary">{{ screenName(showtime.screenId) }}</td>
                <td class="py-3 pr-4 text-copy-secondary">
                  {{ new Date(showtime.startsAt).toLocaleString() }}
                </td>
                <td class="py-3 pr-4 text-copy-secondary">{{ showtime.status }}</td>
                <td class="py-3">
                  <div class="flex gap-2">
                    <Button variant="ghost" @click="startEdit(showtime)">Edit</Button>
                    <Button variant="destructive" @click="onDelete(showtime.id)">Delete</Button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
          <p v-if="!showtimes.length" class="text-sm text-copy-muted">No showtimes yet.</p>
        </div>
      </CardContent>
    </Card>
  </div>
</template>
