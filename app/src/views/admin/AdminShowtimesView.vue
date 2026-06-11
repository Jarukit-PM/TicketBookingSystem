<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'

import { translateApiError } from '@/api/errors'
import { ApiError, api } from '@/api/client'
import TableSkeleton from '@/components/skeletons/TableSkeleton.vue'
import { Button, Card, CardContent, CardHeader, CardTitle, EmptyState, ErrorAlert, Input } from '@/components/ui'
import { Calendar } from 'lucide-vue-next'
import { useLocaleFormat } from '@/composables/useLocaleFormat'
import type { Cinema, Movie, Screen, Showtime, ShowtimeStatus } from '@/types/catalog'

const { t } = useI18n()
const { formatDateTime } = useLocaleFormat()

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

const SATANG_PER_BAHT = 100

function satangToBaht(satang: number): number {
  return satang / SATANG_PER_BAHT
}

function bahtToSatang(baht: number): number {
  return Math.round(baht * SATANG_PER_BAHT)
}

const form = ref({
  movieId: '',
  screenId: '',
  startsAtLocal: '',
  standard: 220,
  vip: 320,
  wheelchair: 220,
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
    standard: 220,
    vip: 320,
    wheelchair: 220,
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
  movies.value = movieRes.movies ?? []
  screens.value = screenRes.screens ?? []
  cinemas.value = cinemaRes.cinemas ?? []
}

async function loadShowtimes() {
  loading.value = true
  errorMessage.value = ''
  try {
    const response = await api.get<{ showtimes: Showtime[] }>(
      `/admin/showtimes${buildQuery()}`,
    )
    showtimes.value = response.showtimes ?? []
  } catch (error) {
    errorMessage.value =
      error instanceof ApiError
        ? translateApiError(error.code, error.message)
        : t('admin.showtimes.loadFailed')
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
    errorMessage.value =
      error instanceof ApiError
        ? translateApiError(error.code, error.message)
        : t('admin.showtimes.loadDataFailed')
    loading.value = false
  }
}

function startEdit(showtime: Showtime) {
  editingId.value = showtime.id
  form.value = {
    movieId: showtime.movieId,
    screenId: showtime.screenId,
    startsAtLocal: toLocalInputValue(showtime.startsAt),
    standard: satangToBaht(showtime.priceTiers.standard),
    vip: satangToBaht(showtime.priceTiers.vip),
    wheelchair: satangToBaht(showtime.priceTiers.wheelchair),
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
        standard: bahtToSatang(Number(form.value.standard)),
        vip: bahtToSatang(Number(form.value.vip)),
        wheelchair: bahtToSatang(Number(form.value.wheelchair)),
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
    errorMessage.value =
      error instanceof ApiError
        ? translateApiError(error.code, error.message)
        : t('admin.showtimes.saveFailed')
  }
}

async function onDelete(id: string) {
  if (!confirm(t('admin.showtimes.confirmDelete'))) return
  try {
    await api.delete(`/admin/showtimes/${id}`)
    if (editingId.value === id) resetForm()
    await loadShowtimes()
  } catch (error) {
    errorMessage.value =
      error instanceof ApiError
        ? translateApiError(error.code, error.message)
        : t('admin.showtimes.deleteFailed')
  }
}

onMounted(loadAll)
</script>

<template>
  <div class="space-y-8">
    <div>
      <h1 class="text-2xl font-semibold text-copy-primary">{{ t('admin.showtimes.title') }}</h1>
      <p class="mt-1 text-sm text-copy-secondary">{{ t('admin.showtimes.subtitle') }}</p>
    </div>

    <Card>
      <CardHeader>
        <CardTitle>{{ editingId ? t('admin.showtimes.editTitle') : t('admin.showtimes.addTitle') }}</CardTitle>
      </CardHeader>
      <CardContent>
        <form class="grid gap-4 md:grid-cols-2" @submit.prevent="onSubmit">
          <div class="space-y-2">
            <label class="text-sm text-copy-secondary" for="movieId">{{ t('common.movie') }}</label>
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
            <label class="text-sm text-copy-secondary" for="screenId">{{ t('common.screen') }}</label>
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
            <label class="text-sm text-copy-secondary" for="startsAt">{{ t('admin.showtimes.startsAt') }}</label>
            <Input id="startsAt" v-model="form.startsAtLocal" type="datetime-local" required />
          </div>
          <div class="space-y-2">
            <label class="text-sm text-copy-secondary" for="showtimeStatus">{{ t('common.status') }}</label>
            <select
              id="showtimeStatus"
              v-model="form.status"
              class="w-full rounded-lg border border-surface-border bg-surface px-3 py-2 text-sm text-copy-primary"
            >
              <option value="OPEN">{{ t('admin.showtimes.status.OPEN') }}</option>
              <option value="CANCELLED">{{ t('admin.showtimes.status.CANCELLED') }}</option>
            </select>
          </div>
          <div class="space-y-2">
            <label class="text-sm text-copy-secondary" for="standard">{{ t('admin.showtimes.standardTHB') }}</label>
            <Input id="standard" :model-value="String(form.standard)" @update:model-value="form.standard = Number($event) || 0" type="number" min="0" step="1" required />
          </div>
          <div class="space-y-2">
            <label class="text-sm text-copy-secondary" for="vip">{{ t('admin.showtimes.vipTHB') }}</label>
            <Input id="vip" :model-value="String(form.vip)" @update:model-value="form.vip = Number($event) || 0" type="number" min="0" step="1" required />
          </div>
          <div class="space-y-2 md:col-span-2">
            <label class="text-sm text-copy-secondary" for="wheelchair">{{ t('admin.showtimes.wheelchairTHB') }}</label>
            <Input id="wheelchair" :model-value="String(form.wheelchair)" @update:model-value="form.wheelchair = Number($event) || 0" type="number" min="0" step="1" required />
          </div>
          <div class="flex gap-2 md:col-span-2">
            <Button type="submit" :disabled="movies.length === 0 || screens.length === 0">
              {{ editingId ? t('common.update') : t('common.create') }}
            </Button>
            <Button v-if="editingId" type="button" variant="secondary" @click="resetForm">
              {{ t('common.cancel') }}
            </Button>
          </div>
        </form>
      </CardContent>
    </Card>

    <ErrorAlert v-if="errorMessage" :message="errorMessage" />

    <Card>
      <CardHeader>
        <div class="flex flex-wrap items-end gap-3">
          <CardTitle>{{ t('admin.showtimes.allTitle') }}</CardTitle>
          <select
            v-model="filter.cinemaId"
            class="rounded-lg border border-surface-border bg-surface px-3 py-2 text-sm text-copy-primary"
            @change="loadShowtimes"
          >
            <option value="">{{ t('common.allCinemas') }}</option>
            <option v-for="cinema in cinemas" :key="cinema.id" :value="cinema.id">
              {{ cinema.name }}
            </option>
          </select>
          <select
            v-model="filter.movieId"
            class="rounded-lg border border-surface-border bg-surface px-3 py-2 text-sm text-copy-primary"
            @change="loadShowtimes"
          >
            <option value="">{{ t('common.allMovies') }}</option>
            <option v-for="movie in movies" :key="movie.id" :value="movie.id">
              {{ movie.title }}
            </option>
          </select>
        </div>
      </CardHeader>
      <CardContent>
        <TableSkeleton v-if="loading" :columns="5" :rows="6" />
        <EmptyState
          v-else-if="!showtimes.length"
          :icon="Calendar"
          :title="t('admin.showtimes.empty')"
          class="py-10"
        />
        <div v-else class="overflow-x-auto">
          <table class="w-full text-left text-sm">
            <thead class="text-copy-muted">
              <tr>
                <th class="pb-3 pr-4 font-medium">{{ t('common.movie') }}</th>
                <th class="pb-3 pr-4 font-medium">{{ t('common.screen') }}</th>
                <th class="pb-3 pr-4 font-medium">{{ t('admin.showtimes.startsAt') }}</th>
                <th class="pb-3 pr-4 font-medium">{{ t('common.status') }}</th>
                <th class="pb-3 font-medium">{{ t('common.actions') }}</th>
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
                  {{ formatDateTime(showtime.startsAt) }}
                </td>
                <td class="py-3 pr-4 text-copy-secondary">
                  {{ t(`admin.showtimes.status.${showtime.status}`) }}
                </td>
                <td class="py-3">
                  <div class="flex gap-2">
                    <Button variant="ghost" @click="startEdit(showtime)">{{ t('common.edit') }}</Button>
                    <Button variant="destructive" @click="onDelete(showtime.id)">{{ t('common.delete') }}</Button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </CardContent>
    </Card>
  </div>
</template>
