<script setup lang="ts">
import { onMounted, ref } from 'vue'

import { ApiError, api } from '@/api/client'
import { Button, Card, CardContent, CardHeader, CardTitle, Input } from '@/components/ui'
import type { Cinema } from '@/types/catalog'

const cinemas = ref<Cinema[]>([])
const loading = ref(true)
const errorMessage = ref('')
const editingId = ref<string | null>(null)

const form = ref({
  name: '',
  address: '',
  timezone: 'Asia/Bangkok',
})

function resetForm() {
  editingId.value = null
  form.value = { name: '', address: '', timezone: 'Asia/Bangkok' }
}

async function loadCinemas() {
  loading.value = true
  errorMessage.value = ''
  try {
    const response = await api.get<{ cinemas: Cinema[] }>('/admin/cinemas')
    cinemas.value = response.cinemas ?? []
  } catch (error) {
    errorMessage.value = error instanceof ApiError ? error.message : 'Failed to load cinemas'
  } finally {
    loading.value = false
  }
}

function startEdit(cinema: Cinema) {
  editingId.value = cinema.id
  form.value = {
    name: cinema.name,
    address: cinema.address,
    timezone: cinema.timezone,
  }
}

async function onSubmit() {
  errorMessage.value = ''
  try {
    if (editingId.value) {
      await api.put(`/admin/cinemas/${editingId.value}`, form.value)
    } else {
      await api.post('/admin/cinemas', form.value)
    }
    resetForm()
    await loadCinemas()
  } catch (error) {
    errorMessage.value = error instanceof ApiError ? error.message : 'Save failed'
  }
}

async function onDelete(id: string) {
  if (!confirm('Delete this cinema?')) return
  try {
    await api.delete(`/admin/cinemas/${id}`)
    if (editingId.value === id) resetForm()
    await loadCinemas()
  } catch (error) {
    errorMessage.value = error instanceof ApiError ? error.message : 'Delete failed'
  }
}

onMounted(loadCinemas)
</script>

<template>
  <div class="space-y-8">
    <div>
      <h1 class="text-2xl font-semibold text-copy-primary">Cinemas</h1>
      <p class="mt-1 text-sm text-copy-secondary">Manage venues and their timezones.</p>
    </div>

    <Card>
      <CardHeader>
        <CardTitle>{{ editingId ? 'Edit cinema' : 'Add cinema' }}</CardTitle>
      </CardHeader>
      <CardContent>
        <form class="grid gap-4 md:grid-cols-2" @submit.prevent="onSubmit">
          <div class="space-y-2 md:col-span-2">
            <label class="text-sm text-copy-secondary" for="name">Name</label>
            <Input id="name" v-model="form.name" required />
          </div>
          <div class="space-y-2 md:col-span-2">
            <label class="text-sm text-copy-secondary" for="address">Address</label>
            <Input id="address" v-model="form.address" />
          </div>
          <div class="space-y-2 md:col-span-2">
            <label class="text-sm text-copy-secondary" for="timezone">Timezone (IANA)</label>
            <Input id="timezone" v-model="form.timezone" placeholder="Asia/Bangkok" required />
          </div>
          <div class="flex gap-2 md:col-span-2">
            <Button type="submit">{{ editingId ? 'Update' : 'Create' }}</Button>
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
        <CardTitle>All cinemas</CardTitle>
      </CardHeader>
      <CardContent>
        <p v-if="loading" class="text-sm text-copy-muted">Loading…</p>
        <div v-else class="overflow-x-auto">
          <table class="w-full text-left text-sm">
            <thead class="text-copy-muted">
              <tr>
                <th class="pb-3 pr-4 font-medium">Name</th>
                <th class="pb-3 pr-4 font-medium">Timezone</th>
                <th class="pb-3 pr-4 font-medium">Address</th>
                <th class="pb-3 font-medium">Actions</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="cinema in cinemas"
                :key="cinema.id"
                class="border-t border-surface-border"
              >
                <td class="py-3 pr-4 text-copy-primary">{{ cinema.name }}</td>
                <td class="py-3 pr-4 text-copy-secondary">{{ cinema.timezone }}</td>
                <td class="py-3 pr-4 text-copy-secondary">{{ cinema.address }}</td>
                <td class="py-3">
                  <div class="flex gap-2">
                    <Button variant="ghost" @click="startEdit(cinema)">Edit</Button>
                    <Button variant="destructive" @click="onDelete(cinema.id)">Delete</Button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
          <p v-if="!cinemas.length" class="text-sm text-copy-muted">No cinemas yet.</p>
        </div>
      </CardContent>
    </Card>
  </div>
</template>
