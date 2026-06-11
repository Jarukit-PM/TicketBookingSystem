<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import { fetchCinemas, fetchMovieDetail } from '@/api/catalog'
import CinemaPicker from '@/components/CinemaPicker.vue'
import LocaleSwitcher from '@/components/LocaleSwitcher.vue'
import ShowtimeDateFilter from '@/components/ShowtimeDateFilter.vue'
import { Button, Card, CardContent } from '@/components/ui'
import { useLocaleFormat } from '@/composables/useLocaleFormat'
import { useShowtimeDates } from '@/composables/useShowtimeDates'
import { formatDuration, lowestTierPrice } from '@/lib/format'
import { useCatalogStore } from '@/stores/catalog'
import type { Cinema, MovieDetail } from '@/types/catalog'

const { t } = useI18n()
const { formatTime, formatTHB } = useLocaleFormat()
const route = useRoute()
const router = useRouter()
const catalog = useCatalogStore()

const cinemas = ref<Cinema[]>([])
const detail = ref<MovieDetail | null>(null)
const loading = ref(true)
const error = ref<string | null>(null)

const movieId = computed(() => route.params.id as string)

const selectedCinema = computed({
  get: () => catalog.selectedCinemaId,
  set: (id: string | null) => {
    if (id) {
      catalog.setCinema(id)
    }
  },
})

const selectedCinemaMeta = computed(() =>
  cinemas.value.find((c) => c.id === catalog.selectedCinemaId),
)

const posterStyle = computed(() => ({
  backgroundImage: detail.value?.posterUrl ? `url(${detail.value.posterUrl})` : undefined,
}))

const posterAria = computed(() =>
  detail.value ? t('catalog.posterAria', { title: detail.value.title }) : '',
)

const cinemaTimezone = computed(() => selectedCinemaMeta.value?.timezone)

const showtimes = computed(() => detail.value?.showtimes ?? [])

const { selectedDateKey, dateOptions, filteredShowtimes } = useShowtimeDates(
  showtimes,
  cinemaTimezone,
)

async function loadDetail(): Promise<void> {
  if (!catalog.selectedCinemaId) {
    detail.value = null
    return
  }

  detail.value = await fetchMovieDetail(movieId.value, catalog.selectedCinemaId)
}

async function refresh(): Promise<void> {
  loading.value = true
  error.value = null
  try {
    if (cinemas.value.length === 0) {
      cinemas.value = await fetchCinemas()
      if (!catalog.selectedCinemaId && cinemas.value.length > 0) {
        catalog.setCinema(cinemas.value[0]!.id)
      }
    }
    await loadDetail()
  } catch {
    error.value = t('movie.loadError')
    detail.value = null
  } finally {
    loading.value = false
  }
}

function goBack(): void {
  router.push({ name: 'home' })
}

function bookShowtime(showtimeId: string): void {
  router.push({ name: 'book', params: { showtimeId } })
}

onMounted(refresh)

watch(
  () => [movieId.value, catalog.selectedCinemaId] as const,
  () => {
    void refresh()
  },
)
</script>

<template>
  <div class="min-h-screen bg-base">
    <header
      class="sticky top-0 z-10 flex h-16 items-center gap-4 border-b border-surface-border bg-base/80 px-4 backdrop-blur-md md:px-6"
    >
      <button
        type="button"
        class="text-sm text-copy-secondary transition-colors hover:text-copy-primary"
        @click="goBack"
      >
        ← {{ t('common.back') }}
      </button>
      <span
        class="bg-gradient-brand bg-clip-text text-xl font-semibold tracking-tight text-transparent"
      >
        {{ t('common.appName') }}
      </span>
      <div class="ml-auto">
        <LocaleSwitcher />
      </div>
    </header>

    <main class="mx-auto max-w-6xl px-4 py-8 md:px-6">
      <div class="mb-8 max-w-xs">
        <CinemaPicker v-model="selectedCinema" :cinemas="cinemas" />
      </div>

      <p
        v-if="error"
        class="mb-6 rounded-lg border border-state-error/30 bg-state-error-dim px-4 py-3 text-sm text-state-error"
      >
        {{ error }}
      </p>

      <div v-if="loading" class="py-16 text-center text-copy-secondary">{{ t('movie.loading') }}</div>

      <div
        v-else-if="!catalog.selectedCinemaId"
        class="py-16 text-center text-copy-secondary"
      >
        {{ t('movie.selectCinemaForShowtimes') }}
      </div>

      <div v-else-if="detail" class="grid gap-8 lg:grid-cols-[280px_1fr]">
        <div
          class="aspect-[2/3] w-full max-w-xs rounded-xl bg-subtle bg-cover bg-center ring-1 ring-white/10"
          :style="posterStyle"
          role="img"
          :aria-label="posterAria"
        />

        <div class="flex flex-col gap-6">
          <div>
            <h1 class="text-3xl font-semibold tracking-tight text-copy-primary md:text-4xl">
              {{ detail.title }}
            </h1>
            <p class="mt-2 text-sm text-copy-secondary">
              {{ detail.rating }} · {{ formatDuration(detail.durationMin) }}
            </p>
            <p class="mt-4 max-w-2xl text-copy-secondary">{{ detail.synopsis }}</p>
          </div>

          <section>
            <h2 class="mb-4 text-lg font-medium text-copy-primary">{{ t('movie.showtimes') }}</h2>

            <div
              v-if="detail.showtimes.length === 0"
              class="rounded-lg border border-surface-border bg-subtle px-4 py-8 text-center"
            >
              <p class="text-copy-primary">{{ t('movie.showtimesNotAnnounced') }}</p>
              <p class="mt-1 text-sm text-copy-secondary">
                {{ t('movie.showtimesCheckBack') }}
              </p>
            </div>

            <template v-else>
              <ShowtimeDateFilter
                v-model="selectedDateKey"
                :dates="dateOptions"
                :time-zone="cinemaTimezone"
              />

              <div class="grid gap-3 sm:grid-cols-2">
                <Card
                  v-for="showtime in filteredShowtimes"
                  :key="showtime.id"
                  class="overflow-hidden transition-shadow hover:shadow-glow-brand/10"
                >
                  <CardContent class="flex items-center justify-between gap-4 p-4">
                    <div>
                      <p class="font-medium text-copy-primary">
                        {{ formatTime(showtime.startsAt, cinemaTimezone) }}
                      </p>
                      <p class="text-sm text-copy-secondary">{{ showtime.screenName }}</p>
                      <p class="mt-1 text-xs text-copy-muted">
                        {{ t('common.from') }}
                        {{ formatTHB(lowestTierPrice(showtime.priceTiers)) }}
                      </p>
                    </div>
                    <Button variant="primary" @click="bookShowtime(showtime.id)">
                      {{ t('movie.book') }}
                    </Button>
                  </CardContent>
                </Card>
              </div>
            </template>
          </section>
        </div>
      </div>
    </main>
  </div>
</template>
