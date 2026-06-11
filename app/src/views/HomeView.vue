<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { RouterLink } from 'vue-router'
import { fetchCinemas, fetchMovies } from '@/api/catalog'
import CinemaPicker from '@/components/CinemaPicker.vue'
import LocaleSwitcher from '@/components/LocaleSwitcher.vue'
import MovieCard from '@/components/MovieCard.vue'
import { Button } from '@/components/ui'
import { useAuthStore } from '@/stores/auth'
import { useCatalogStore } from '@/stores/catalog'
import type { CatalogTab, Cinema, Movie } from '@/types/catalog'

const { t } = useI18n()
const auth = useAuthStore()
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
    <header
      class="sticky top-0 z-10 flex h-16 items-center border-b border-surface-border bg-base/80 px-4 backdrop-blur-md md:px-6"
    >
      <span
        class="bg-gradient-brand bg-clip-text text-xl font-semibold tracking-tight text-transparent"
      >
        {{ t('common.appName') }}
      </span>
      <div class="ml-auto flex items-center gap-3">
        <LocaleSwitcher />
        <template v-if="auth.isAuthenticated">
          <RouterLink
            to="/my-bookings"
            class="text-sm text-copy-secondary transition-colors hover:text-copy-primary"
          >
            {{ t('nav.myBookings') }}
          </RouterLink>
          <RouterLink
            v-if="auth.isAdmin"
            to="/admin"
            class="text-sm text-copy-secondary transition-colors hover:text-copy-primary"
          >
            {{ t('nav.admin') }}
          </RouterLink>
          <span class="hidden text-sm text-copy-muted sm:inline">{{ auth.user?.email }}</span>
          <Button variant="secondary" @click="auth.logout()">{{ t('nav.signOut') }}</Button>
        </template>
        <template v-else>
          <RouterLink
            to="/login"
            class="text-sm text-copy-secondary transition-colors hover:text-copy-primary"
          >
            {{ t('nav.signIn') }}
          </RouterLink>
          <RouterLink to="/register">
            <Button variant="secondary">{{ t('nav.register') }}</Button>
          </RouterLink>
        </template>
      </div>
    </header>

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
          class="flex-1 rounded-md px-4 py-2 text-sm font-medium transition-colors"
          :class="
            catalog.activeTab === 'now_showing'
              ? 'bg-gradient-brand text-white'
              : 'text-copy-secondary hover:text-copy-primary'
          "
          @click="switchTab('now_showing')"
        >
          {{ t('catalog.nowShowing') }}
        </button>
        <button
          type="button"
          class="flex-1 rounded-md px-4 py-2 text-sm font-medium transition-colors"
          :class="
            catalog.activeTab === 'coming_soon'
              ? 'bg-gradient-brand text-white'
              : 'text-copy-secondary hover:text-copy-primary'
          "
          @click="switchTab('coming_soon')"
        >
          {{ t('catalog.comingSoon') }}
        </button>
      </div>

      <p
        v-if="error"
        class="mb-6 rounded-lg border border-state-error/30 bg-state-error-dim px-4 py-3 text-sm text-state-error"
      >
        {{ error }}
      </p>

      <p v-if="loading" class="py-16 text-center text-copy-secondary">{{ t('catalog.loadingMovies') }}</p>

      <div
        v-else-if="movies.length === 0"
        class="rounded-xl border border-surface-border bg-surface px-6 py-16 text-center"
      >
        <p class="text-copy-primary">{{ emptyMessage }}</p>
        <p
          v-if="catalog.selectedCinemaId && catalog.activeTab === 'now_showing'"
          class="mt-2 text-sm text-copy-secondary"
        >
          {{ t('catalog.adminHint') }}
        </p>
      </div>

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
