<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'

import { translateApiError } from '@/api/errors'
import { ApiError, api } from '@/api/client'
import BookingsTable from '@/components/admin/BookingsTable.vue'
import { Button, Card, CardContent, CardHeader, CardTitle, Input } from '@/components/ui'
import type { BookingSummary } from '@/types/admin'
import type { Movie } from '@/types/catalog'

const { t } = useI18n()
const route = useRoute()

const focusBookingId = computed(() => {
  const raw = route.query.bookingId
  return typeof raw === 'string' ? raw : undefined
})
const focusBookingRef = computed(() => {
  const raw = route.query.bookingRef
  return typeof raw === 'string' ? raw : undefined
})

const selectClass =
  'w-full rounded-lg border border-surface-border bg-surface px-3 py-2 text-sm text-copy-primary focus:outline-none focus:ring-2 focus:ring-accent-glow focus:border-brand/50'

const bookingRef = ref('')
const email = ref('')
const userId = ref('')
const showtimeId = ref('')
const movieId = ref('')
const locale = ref('')
const confirmedFrom = ref('')
const confirmedTo = ref('')
const movies = ref<Movie[]>([])
const bookings = ref<BookingSummary[]>([])
const loading = ref(false)
const errorMessage = ref('')
const searched = ref(false)
const page = ref(1)
const limit = 20
const total = ref(0)

const hasFilters = computed(
  () =>
    Boolean(
      bookingRef.value.trim() ||
        email.value.trim() ||
        userId.value.trim() ||
        showtimeId.value.trim() ||
        movieId.value.trim() ||
        locale.value.trim() ||
        confirmedFrom.value.trim() ||
        confirmedTo.value.trim(),
    ),
)

const totalPages = computed(() => Math.max(1, Math.ceil(total.value / limit)))
const canGoPrev = computed(() => page.value > 1)
const canGoNext = computed(() => page.value < totalPages.value)

const emptyMessage = computed(() => {
  if (!searched.value) return t('admin.bookings.empty')
  return hasFilters.value ? t('admin.bookings.noMatch') : t('admin.bookings.empty')
})

async function loadBookings() {
  loading.value = true
  errorMessage.value = ''
  searched.value = true
  try {
    const params = new URLSearchParams({
      page: String(page.value),
      limit: String(limit),
    })
    if (bookingRef.value.trim()) params.set('bookingRef', bookingRef.value.trim())
    if (email.value.trim()) params.set('email', email.value.trim())
    if (userId.value.trim()) params.set('userId', userId.value.trim())
    if (showtimeId.value.trim()) params.set('showtimeId', showtimeId.value.trim())
    if (movieId.value.trim()) params.set('movieId', movieId.value.trim())
    if (locale.value.trim()) params.set('locale', locale.value.trim())
    if (confirmedFrom.value.trim()) params.set('confirmedFrom', confirmedFrom.value.trim())
    if (confirmedTo.value.trim()) params.set('confirmedTo', confirmedTo.value.trim())

    const data = await api.get<{
      bookings: BookingSummary[]
      total: number
      page: number
      limit: number
    }>(`/admin/bookings?${params}`)
    bookings.value = data.bookings ?? []
    total.value = data.total ?? 0
    page.value = data.page ?? page.value
  } catch (error) {
    bookings.value = []
    total.value = 0
    errorMessage.value =
      error instanceof ApiError
        ? translateApiError(error.code, error.message)
        : t('admin.bookings.loadFailed')
  } finally {
    loading.value = false
  }
}

function search() {
  page.value = 1
  loadBookings()
}

function clearFilters() {
  bookingRef.value = ''
  email.value = ''
  userId.value = ''
  showtimeId.value = ''
  movieId.value = ''
  locale.value = ''
  confirmedFrom.value = ''
  confirmedTo.value = ''
  page.value = 1
  loadBookings()
}

async function loadMovies() {
  try {
    const data = await api.get<{ movies: Movie[] }>('/admin/movies')
    movies.value = (data.movies ?? []).slice().sort((a, b) => a.title.localeCompare(b.title))
  } catch {
    movies.value = []
  }
}

function goPrev() {
  if (!canGoPrev.value) return
  page.value -= 1
  loadBookings()
}

function goNext() {
  if (!canGoNext.value) return
  page.value += 1
  loadBookings()
}

onMounted(() => {
  void loadMovies()
  if (focusBookingRef.value) {
    bookingRef.value = focusBookingRef.value
  }
  void loadBookings()
})
</script>

<template>
  <div class="space-y-8">
    <div>
      <h1 class="text-2xl font-semibold text-copy-primary">{{ t('admin.bookings.title') }}</h1>
      <p class="mt-1 text-sm text-copy-secondary">{{ t('admin.bookings.subtitle') }}</p>
    </div>

    <Card>
      <CardHeader>
        <CardTitle>{{ t('admin.bookings.searchTitle') }}</CardTitle>
      </CardHeader>
      <CardContent class="space-y-4">
        <div class="grid gap-4 sm:grid-cols-2">
          <label class="block space-y-1.5">
            <span class="text-sm text-copy-secondary">{{ t('admin.bookings.bookingRef') }}</span>
            <Input v-model="bookingRef" :placeholder="t('admin.bookings.bookingRefPlaceholder')" autocomplete="off" />
          </label>
          <label class="block space-y-1.5">
            <span class="text-sm text-copy-secondary">{{ t('admin.bookings.customerEmail') }}</span>
            <Input v-model="email" type="email" :placeholder="t('admin.bookings.emailPlaceholder')" autocomplete="off" />
          </label>
          <label class="block space-y-1.5">
            <span class="text-sm text-copy-secondary">{{ t('admin.bookings.userId') }}</span>
            <Input v-model="userId" :placeholder="t('admin.bookings.objectIdPlaceholder')" autocomplete="off" />
          </label>
          <label class="block space-y-1.5">
            <span class="text-sm text-copy-secondary">{{ t('admin.bookings.showtimeId') }}</span>
            <Input v-model="showtimeId" :placeholder="t('admin.bookings.objectIdPlaceholder')" autocomplete="off" />
          </label>
          <label class="block space-y-1.5">
            <span class="text-sm text-copy-secondary">{{ t('admin.bookings.movie') }}</span>
            <select v-model="movieId" :class="selectClass">
              <option value="">{{ t('admin.bookings.movieAll') }}</option>
              <option v-for="movie in movies" :key="movie.id" :value="movie.id">
                {{ movie.title }}
              </option>
            </select>
          </label>
          <label class="block space-y-1.5">
            <span class="text-sm text-copy-secondary">{{ t('admin.bookings.locale') }}</span>
            <select v-model="locale" :class="selectClass">
              <option value="">{{ t('admin.bookings.localeAll') }}</option>
              <option value="en">{{ t('locale.en') }}</option>
              <option value="th">{{ t('locale.th') }}</option>
            </select>
          </label>
          <label class="block space-y-1.5">
            <span class="text-sm text-copy-secondary">{{ t('admin.bookings.confirmedFrom') }}</span>
            <Input v-model="confirmedFrom" type="date" autocomplete="off" />
          </label>
          <label class="block space-y-1.5">
            <span class="text-sm text-copy-secondary">{{ t('admin.bookings.confirmedTo') }}</span>
            <Input v-model="confirmedTo" type="date" autocomplete="off" />
          </label>
        </div>
        <div class="flex flex-wrap gap-3">
          <Button :disabled="loading" @click="search">{{ t('common.search') }}</Button>
          <Button variant="secondary" :disabled="loading" @click="clearFilters">{{ t('common.clear') }}</Button>
        </div>
        <p v-if="errorMessage" class="text-sm text-state-error" role="alert">{{ errorMessage }}</p>
      </CardContent>
    </Card>

    <Card>
      <CardHeader class="flex flex-row flex-wrap items-center justify-between gap-3">
        <CardTitle>{{ t('common.results') }}</CardTitle>
        <p v-if="searched && total > 0" class="text-sm text-copy-muted">
          {{ t('admin.bookings.pageInfo', { total, page, totalPages }) }}
        </p>
      </CardHeader>
      <CardContent class="space-y-4">
        <BookingsTable
          :bookings="bookings"
          :loading="loading"
          show-customer
          :empty-message="emptyMessage"
          :focus-booking-id="focusBookingId"
          :focus-booking-ref="focusBookingRef"
        />
        <div v-if="total > limit" class="flex flex-wrap items-center justify-end gap-3">
          <Button variant="secondary" :disabled="loading || !canGoPrev" @click="goPrev">
            {{ t('common.previous') }}
          </Button>
          <Button variant="secondary" :disabled="loading || !canGoNext" @click="goNext">
            {{ t('common.next') }}
          </Button>
        </div>
      </CardContent>
    </Card>
  </div>
</template>
