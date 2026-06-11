<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'

import { translateApiError } from '@/api/errors'
import { ApiError, api } from '@/api/client'
import { Button, Card, CardContent, CardHeader, CardTitle, Input } from '@/components/ui'
import type { Movie, MovieStatus } from '@/types/catalog'

const { t } = useI18n()

const movies = ref<Movie[]>([])
const loading = ref(true)
const errorMessage = ref('')
const editingId = ref<string | null>(null)

const form = ref({
  title: '',
  posterUrl: '',
  durationMin: 120,
  rating: 'PG',
  synopsis: '',
  status: 'NOW_SHOWING' as MovieStatus,
})

const statusOptions: MovieStatus[] = ['NOW_SHOWING', 'COMING_SOON', 'ARCHIVED']

function resetForm() {
  editingId.value = null
  form.value = {
    title: '',
    posterUrl: '',
    durationMin: 120,
    rating: 'PG',
    synopsis: '',
    status: 'NOW_SHOWING',
  }
}

async function loadMovies() {
  loading.value = true
  errorMessage.value = ''
  try {
    const response = await api.get<{ movies: Movie[] }>('/admin/movies')
    movies.value = response.movies ?? []
  } catch (error) {
    errorMessage.value =
      error instanceof ApiError
        ? translateApiError(error.code, error.message)
        : t('admin.movies.loadFailed')
  } finally {
    loading.value = false
  }
}

function startEdit(movie: Movie) {
  editingId.value = movie.id
  form.value = {
    title: movie.title,
    posterUrl: movie.posterUrl,
    durationMin: movie.durationMin,
    rating: movie.rating,
    synopsis: movie.synopsis,
    status: movie.status,
  }
}

async function onSubmit() {
  errorMessage.value = ''
  try {
    if (editingId.value) {
      await api.put(`/admin/movies/${editingId.value}`, form.value)
    } else {
      await api.post('/admin/movies', form.value)
    }
    resetForm()
    await loadMovies()
  } catch (error) {
    errorMessage.value =
      error instanceof ApiError
        ? translateApiError(error.code, error.message)
        : t('admin.movies.saveFailed')
  }
}

async function onDelete(id: string) {
  if (!confirm(t('admin.movies.confirmDelete'))) return
  try {
    await api.delete(`/admin/movies/${id}`)
    if (editingId.value === id) resetForm()
    await loadMovies()
  } catch (error) {
    errorMessage.value =
      error instanceof ApiError
        ? translateApiError(error.code, error.message)
        : t('admin.movies.deleteFailed')
  }
}

onMounted(loadMovies)
</script>

<template>
  <div class="space-y-8">
    <div>
      <h1 class="text-2xl font-semibold text-copy-primary">{{ t('admin.movies.title') }}</h1>
      <p class="mt-1 text-sm text-copy-secondary">{{ t('admin.movies.subtitle') }}</p>
    </div>

    <Card>
      <CardHeader>
        <CardTitle>{{ editingId ? t('admin.movies.editTitle') : t('admin.movies.addTitle') }}</CardTitle>
      </CardHeader>
      <CardContent>
        <form class="grid gap-4 md:grid-cols-2" @submit.prevent="onSubmit">
          <div class="space-y-2 md:col-span-2">
            <label class="text-sm text-copy-secondary" for="title">{{ t('common.title') }}</label>
            <Input id="title" v-model="form.title" required />
          </div>
          <div class="space-y-2 md:col-span-2">
            <label class="text-sm text-copy-secondary" for="posterUrl">{{ t('admin.movies.posterUrl') }}</label>
            <Input id="posterUrl" v-model="form.posterUrl" :placeholder="t('admin.movies.posterPlaceholder')" />
          </div>
          <div class="space-y-2">
            <label class="text-sm text-copy-secondary" for="durationMin">{{ t('common.durationMin') }}</label>
            <Input id="durationMin" :model-value="String(form.durationMin)" @update:model-value="form.durationMin = Number($event) || 0" type="number" min="1" required />
          </div>
          <div class="space-y-2">
            <label class="text-sm text-copy-secondary" for="rating">{{ t('admin.movies.rating') }}</label>
            <Input id="rating" v-model="form.rating" />
          </div>
          <div class="space-y-2 md:col-span-2">
            <label class="text-sm text-copy-secondary" for="status">{{ t('common.status') }}</label>
            <select
              id="status"
              v-model="form.status"
              class="w-full rounded-lg border border-surface-border bg-surface px-3 py-2 text-sm text-copy-primary"
            >
              <option v-for="status in statusOptions" :key="status" :value="status">
                {{ t(`admin.movies.status.${status}`) }}
              </option>
            </select>
          </div>
          <div class="space-y-2 md:col-span-2">
            <label class="text-sm text-copy-secondary" for="synopsis">{{ t('admin.movies.synopsis') }}</label>
            <textarea
              id="synopsis"
              v-model="form.synopsis"
              rows="3"
              class="w-full rounded-lg border border-surface-border bg-surface px-3 py-2 text-sm text-copy-primary"
            />
          </div>
          <div class="flex gap-2 md:col-span-2">
            <Button type="submit">{{ editingId ? t('common.update') : t('common.create') }}</Button>
            <Button v-if="editingId" type="button" variant="secondary" @click="resetForm">
              {{ t('common.cancel') }}
            </Button>
          </div>
        </form>
      </CardContent>
    </Card>

    <p v-if="errorMessage" class="text-sm text-state-error" role="alert">{{ errorMessage }}</p>

    <Card>
      <CardHeader>
        <CardTitle>{{ t('admin.movies.allTitle') }}</CardTitle>
      </CardHeader>
      <CardContent>
        <p v-if="loading" class="text-sm text-copy-muted">{{ t('common.loading') }}</p>
        <div v-else class="overflow-x-auto">
          <table class="w-full text-left text-sm">
            <thead class="text-copy-muted">
              <tr>
                <th class="pb-3 pr-4 font-medium">{{ t('common.title') }}</th>
                <th class="pb-3 pr-4 font-medium">{{ t('common.status') }}</th>
                <th class="pb-3 pr-4 font-medium">{{ t('common.duration') }}</th>
                <th class="pb-3 font-medium">{{ t('common.actions') }}</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="movie in movies"
                :key="movie.id"
                class="border-t border-surface-border"
              >
                <td class="py-3 pr-4 text-copy-primary">{{ movie.title }}</td>
                <td class="py-3 pr-4 text-copy-secondary">{{ t(`admin.movies.status.${movie.status}`) }}</td>
                <td class="py-3 pr-4 text-copy-secondary">{{ t('common.minutesShort', { count: movie.durationMin }) }}</td>
                <td class="py-3">
                  <div class="flex gap-2">
                    <Button variant="ghost" @click="startEdit(movie)">{{ t('common.edit') }}</Button>
                    <Button variant="destructive" @click="onDelete(movie.id)">{{ t('common.delete') }}</Button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
          <p v-if="!movies.length" class="text-sm text-copy-muted">{{ t('admin.movies.empty') }}</p>
        </div>
      </CardContent>
    </Card>
  </div>
</template>
