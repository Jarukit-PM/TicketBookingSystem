<script setup lang="ts">
import { CalendarClock, Clapperboard, Film } from 'lucide-vue-next'
import { computed, onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { fetchCinemas, fetchMovies } from '@/api/catalog'
import AppHeader from '@/components/AppHeader.vue'
import CinemaPicker from '@/components/CinemaPicker.vue'
import MovieCard from '@/components/MovieCard.vue'
import MovieGridSkeleton from '@/components/skeletons/MovieGridSkeleton.vue'
import { EmptyState, ErrorAlert } from '@/components/ui'
import { useAuthStore } from '@/stores/auth'
import { useCatalogStore } from '@/stores/catalog'
import type { CatalogTab, Cinema, Movie } from '@/types/catalog'

const { t } = useI18n()
const catalog = useCatalogStore()

const cinemas = ref<Cinema[]>([])
const movies = ref<Movie[]>([])
const loading = ref(false)
const error = ref<string | null>(null)

const selectedCinema = computed({
  get: () => catalog.selectedCinemaId,
  set: (id: string | null) => {
    if (id) {
      catalog.setCinema(id)
    }
  },
})

const emptyMessage = computed(() => {
  if (!catalog.selectedCinemaId) {
    return t('catalog.selectCinema')
  }
  if (catalog.activeTab === 'now_showing') {
    return t('catalog.emptyNowShowing')
  }
  return t('catalog.emptyComingSoon')
})

const emptyDescription = computed(() => {
  if (catalog.selectedCinemaId && catalog.activeTab === 'now_showing') {
    return t('catalog.adminHint')
  }
  return undefined
})

async function loadMovies(): Promise<void> {
  if (!catalog.selectedCinemaId) {
    movies.value = []
    return
  }

  loading.value = true
  error.value = null
  try {
    movies.value = await fetchMovies(catalog.selectedCinemaId, catalog.activeTab)
  } catch {
    error.value = t('catalog.loadMoviesError')
    movies.value = []
  } finally {
    loading.value = false
  }
}

function switchTab(tab: CatalogTab): void {
  if (catalog.activeTab === tab) return
  catalog.setTab(tab)
}

onMounted(async () => {
  try {
    cinemas.value = await fetchCinemas()
    if (!catalog.selectedCinemaId && cinemas.value.length > 0) {
      catalog.setCinema(cinemas.value[0]!.id)
    }
  } catch {
    error.value = t('catalog.loadCinemasError')
  }
})

watch(
  () => [catalog.selectedCinemaId, catalog.activeTab] as const,
  () => void loadMovies(),
  { immediate: true },
)
</script>

<template>
  <div class="min-h-screen bg-base">
    <AppHeader />

    <main class="mx-auto max-w-6xl px-4 py-8 md:px-6 md:py-12">
      <div class="mb-8 flex flex-col gap-6 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <h1 class="text-3xl font-semibold tracking-tight text-copy-primary md:text-4xl">
            {{ t('catalog.browseMovies') }}
          </h1>
          <p class="mt-2 text-sm text-copy-secondary">
            {{ t('catalog.browseSubtitle') }}
          </p>
        </div>
        <div class="w-full max-w-xs">
          <CinemaPicker v-model="selectedCinema" :cinemas="cinemas" />
        </div>
      </div>

      <div class="mb-8 flex gap-2 rounded-lg border border-surface-border bg-surface p-1 sm:max-w-md">
        <button
          type="button"
          class="inline-flex flex-1 items-center justify-center gap-2 rounded-md px-4 py-2 text-sm font-medium transition-colors"
          :class="
            catalog.activeTab === 'now_showing'
              ? 'bg-gradient-brand text-white'
              : 'text-copy-secondary hover:text-copy-primary'
          "
          @click="switchTab('now_showing')"
        >
          <Clapperboard class="h-4 w-4" aria-hidden="true" />
          {{ t('catalog.nowShowing') }}
        </button>
        <button
          type="button"
          class="inline-flex flex-1 items-center justify-center gap-2 rounded-md px-4 py-2 text-sm font-medium transition-colors"
          :class="
            catalog.activeTab === 'coming_soon'
              ? 'bg-gradient-brand text-white'
              : 'text-copy-secondary hover:text-copy-primary'
          "
          @click="switchTab('coming_soon')"
        >
          <CalendarClock class="h-4 w-4" aria-hidden="true" />
          {{ t('catalog.comingSoon') }}
        </button>
      </div>

      <ErrorAlert v-if="error" :message="error" class="mb-6" />

      <MovieGridSkeleton v-if="loading" />

      <EmptyState
        v-else-if="movies.length === 0"
        :icon="Film"
        :title="emptyMessage"
        :description="emptyDescription"
      />

      <div v-else class="grid gap-6 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
        <MovieCard
          v-for="movie in movies"
          :key="movie.id"
          :movie="movie"
          :tab="catalog.activeTab"
        />
      </div>
    </main>
  </div>
</template>
