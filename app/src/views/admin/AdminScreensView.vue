<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'

import { translateApiError } from '@/api/errors'
import { ApiError, api } from '@/api/client'
import { Button, Card, CardContent, CardHeader, CardTitle, Input } from '@/components/ui'
import type { Cinema, Screen, ScreenLayout } from '@/types/catalog'

const { t } = useI18n()

const cinemas = ref<Cinema[]>([])
const screens = ref<Screen[]>([])
const loading = ref(true)
const errorMessage = ref('')
const editingId = ref<string | null>(null)
const filterCinemaId = ref('')

const defaultLayout: ScreenLayout = {
  seats: [
    { seatId: 'A-1', row: 1, col: 1, type: 'standard' },
    { seatId: 'A-2', row: 1, col: 2, type: 'standard' },
    { seatId: 'A-3', row: 1, col: 3, type: 'vip' },
  ],
}

const form = ref({
  cinemaId: '',
  name: '',
  layoutJson: JSON.stringify(defaultLayout, null, 2),
})

const filteredScreens = computed(() => {
  if (!filterCinemaId.value) return screens.value
  return screens.value.filter((s) => s.cinemaId === filterCinemaId.value)
})

function cinemaName(id: string) {
  return cinemas.value.find((c) => c.id === id)?.name ?? id
}

function resetForm() {
  editingId.value = null
  form.value = {
    cinemaId: filterCinemaId.value || cinemas.value[0]?.id || '',
    name: '',
    layoutJson: JSON.stringify(defaultLayout, null, 2),
  }
}

function parseLayout(): ScreenLayout {
  const parsed = JSON.parse(form.value.layoutJson) as ScreenLayout
  if (!parsed?.seats?.length) {
    throw new Error(t('admin.screens.layoutMustIncludeSeats'))
  }
  return parsed
}

async function loadData() {
  loading.value = true
  errorMessage.value = ''
  try {
    const [cinemaRes, screenRes] = await Promise.all([
      api.get<{ cinemas: Cinema[] }>('/admin/cinemas'),
      api.get<{ screens: Screen[] }>('/admin/screens'),
    ])
    cinemas.value = cinemaRes.cinemas ?? []
    screens.value = screenRes.screens ?? []
    if (!form.value.cinemaId && cinemas.value[0]) {
      form.value.cinemaId = cinemas.value[0].id
    }
  } catch (error) {
    errorMessage.value =
      error instanceof ApiError
        ? translateApiError(error.code, error.message)
        : t('admin.screens.loadFailed')
  } finally {
    loading.value = false
  }
}

async function loadScreens() {
  const query = filterCinemaId.value ? `?cinemaId=${filterCinemaId.value}` : ''
  const response = await api.get<{ screens: Screen[] }>(`/admin/screens${query}`)
  screens.value = response.screens ?? []
}

function startEdit(screen: Screen) {
  editingId.value = screen.id
  form.value = {
    cinemaId: screen.cinemaId,
    name: screen.name,
    layoutJson: JSON.stringify(screen.layout, null, 2),
  }
}

async function onSubmit() {
  errorMessage.value = ''
  try {
    const layout = parseLayout()
    const payload = {
      cinemaId: form.value.cinemaId,
      name: form.value.name,
      layout,
    }
    if (editingId.value) {
      await api.put(`/admin/screens/${editingId.value}`, payload)
    } else {
      await api.post('/admin/screens', payload)
    }
    resetForm()
    await loadScreens()
  } catch (error) {
    if (error instanceof ApiError) {
      errorMessage.value = translateApiError(error.code, error.message)
    } else if (error instanceof Error) {
      errorMessage.value = error.message
    } else {
      errorMessage.value = t('admin.screens.saveFailed')
    }
  }
}

async function onDelete(id: string) {
  if (!confirm(t('admin.screens.confirmDelete'))) return
  try {
    await api.delete(`/admin/screens/${id}`)
    if (editingId.value === id) resetForm()
    await loadScreens()
  } catch (error) {
    errorMessage.value =
      error instanceof ApiError
        ? translateApiError(error.code, error.message)
        : t('admin.screens.deleteFailed')
  }
}

onMounted(async () => {
  await loadData()
})
</script>

<template>
  <div class="space-y-8">
    <div>
      <h1 class="text-2xl font-semibold text-copy-primary">{{ t('admin.screens.title') }}</h1>
      <p class="mt-1 text-sm text-copy-secondary">{{ t('admin.screens.subtitle') }}</p>
    </div>

    <Card>
      <CardHeader>
        <CardTitle>{{ editingId ? t('admin.screens.editTitle') : t('admin.screens.addTitle') }}</CardTitle>
      </CardHeader>
      <CardContent>
        <form class="space-y-4" @submit.prevent="onSubmit">
          <div class="space-y-2">
            <label class="text-sm text-copy-secondary" for="cinemaId">{{ t('common.cinema') }}</label>
            <select
              id="cinemaId"
              v-model="form.cinemaId"
              class="w-full rounded-lg border border-surface-border bg-surface px-3 py-2 text-sm text-copy-primary"
              required
            >
              <option v-for="cinema in cinemas" :key="cinema.id" :value="cinema.id">
                {{ cinema.name }}
              </option>
            </select>
          </div>
          <div class="space-y-2">
            <label class="text-sm text-copy-secondary" for="screenName">{{ t('common.name') }}</label>
            <Input id="screenName" v-model="form.name" required />
          </div>
          <div class="space-y-2">
            <label class="text-sm text-copy-secondary" for="layout">{{ t('admin.screens.layoutJson') }}</label>
            <textarea
              id="layout"
              v-model="form.layoutJson"
              rows="10"
              class="w-full rounded-lg border border-surface-border bg-surface px-3 py-2 font-mono text-xs text-copy-primary"
              required
            />
          </div>
          <p v-if="!loading && cinemas.length === 0" class="text-sm text-copy-muted">
            {{ t('admin.screens.addCinemaFirst') }}
          </p>
          <div class="flex gap-2">
            <Button type="submit" :disabled="cinemas.length === 0">
              {{ editingId ? t('common.update') : t('common.create') }}
            </Button>
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
        <div class="flex flex-wrap items-center justify-between gap-3">
          <CardTitle>{{ t('admin.screens.allTitle') }}</CardTitle>
          <select
            v-model="filterCinemaId"
            class="rounded-lg border border-surface-border bg-surface px-3 py-2 text-sm text-copy-primary"
            @change="loadScreens"
          >
            <option value="">{{ t('common.allCinemas') }}</option>
            <option v-for="cinema in cinemas" :key="cinema.id" :value="cinema.id">
              {{ cinema.name }}
            </option>
          </select>
        </div>
      </CardHeader>
      <CardContent>
        <p v-if="loading" class="text-sm text-copy-muted">{{ t('common.loading') }}</p>
        <div v-else class="overflow-x-auto">
          <table class="w-full text-left text-sm">
            <thead class="text-copy-muted">
              <tr>
                <th class="pb-3 pr-4 font-medium">{{ t('common.name') }}</th>
                <th class="pb-3 pr-4 font-medium">{{ t('common.cinema') }}</th>
                <th class="pb-3 pr-4 font-medium">{{ t('common.seats') }}</th>
                <th class="pb-3 font-medium">{{ t('common.actions') }}</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="screen in filteredScreens"
                :key="screen.id"
                class="border-t border-surface-border"
              >
                <td class="py-3 pr-4 text-copy-primary">{{ screen.name }}</td>
                <td class="py-3 pr-4 text-copy-secondary">{{ cinemaName(screen.cinemaId) }}</td>
                <td class="py-3 pr-4 text-copy-secondary">{{ screen.layout.seats.length }}</td>
                <td class="py-3">
                  <div class="flex gap-2">
                    <Button variant="ghost" @click="startEdit(screen)">{{ t('common.edit') }}</Button>
                    <Button variant="destructive" @click="onDelete(screen.id)">{{ t('common.delete') }}</Button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
          <p v-if="!filteredScreens.length" class="text-sm text-copy-muted">{{ t('admin.screens.empty') }}</p>
        </div>
      </CardContent>
    </Card>
  </div>
</template>
