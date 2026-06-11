<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'

import { ApiError, api } from '@/api/client'
import { Button, Card, CardContent, CardHeader, CardTitle, Input } from '@/components/ui'
import type { Cinema, Screen, ScreenLayout } from '@/types/catalog'

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
    throw new Error('Layout must include a seats array')
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
    errorMessage.value = error instanceof ApiError ? error.message : 'Failed to load screens'
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
    errorMessage.value =
      error instanceof ApiError ? error.message : error instanceof Error ? error.message : 'Save failed'
  }
}

async function onDelete(id: string) {
  if (!confirm('Delete this screen?')) return
  try {
    await api.delete(`/admin/screens/${id}`)
    if (editingId.value === id) resetForm()
    await loadScreens()
  } catch (error) {
    errorMessage.value = error instanceof ApiError ? error.message : 'Delete failed'
  }
}

onMounted(async () => {
  await loadData()
})
</script>

<template>
  <div class="space-y-8">
    <div>
      <h1 class="text-2xl font-semibold text-copy-primary">Screens</h1>
      <p class="mt-1 text-sm text-copy-secondary">
        Halls with seat layouts. Edit layout as JSON ({ seatId, row, col, type }).
      </p>
    </div>

    <Card>
      <CardHeader>
        <CardTitle>{{ editingId ? 'Edit screen' : 'Add screen' }}</CardTitle>
      </CardHeader>
      <CardContent>
        <form class="space-y-4" @submit.prevent="onSubmit">
          <div class="space-y-2">
            <label class="text-sm text-copy-secondary" for="cinemaId">Cinema</label>
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
            <label class="text-sm text-copy-secondary" for="screenName">Name</label>
            <Input id="screenName" v-model="form.name" required />
          </div>
          <div class="space-y-2">
            <label class="text-sm text-copy-secondary" for="layout">Layout JSON</label>
            <textarea
              id="layout"
              v-model="form.layoutJson"
              rows="10"
              class="w-full rounded-lg border border-surface-border bg-surface px-3 py-2 font-mono text-xs text-copy-primary"
              required
            />
          </div>
          <p v-if="!loading && cinemas.length === 0" class="text-sm text-copy-muted">
            Add a cinema first under Cinemas before creating a screen.
          </p>
          <div class="flex gap-2">
            <Button type="submit" :disabled="cinemas.length === 0">
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
        <div class="flex flex-wrap items-center justify-between gap-3">
          <CardTitle>All screens</CardTitle>
          <select
            v-model="filterCinemaId"
            class="rounded-lg border border-surface-border bg-surface px-3 py-2 text-sm text-copy-primary"
            @change="loadScreens"
          >
            <option value="">All cinemas</option>
            <option v-for="cinema in cinemas" :key="cinema.id" :value="cinema.id">
              {{ cinema.name }}
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
                <th class="pb-3 pr-4 font-medium">Name</th>
                <th class="pb-3 pr-4 font-medium">Cinema</th>
                <th class="pb-3 pr-4 font-medium">Seats</th>
                <th class="pb-3 font-medium">Actions</th>
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
                    <Button variant="ghost" @click="startEdit(screen)">Edit</Button>
                    <Button variant="destructive" @click="onDelete(screen.id)">Delete</Button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
          <p v-if="!filteredScreens.length" class="text-sm text-copy-muted">No screens yet.</p>
        </div>
      </CardContent>
    </Card>
  </div>
</template>
