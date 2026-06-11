<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'

import { translateApiError } from '@/api/errors'
import { ApiError, api } from '@/api/client'
import TableSkeleton from '@/components/skeletons/TableSkeleton.vue'
import { Button, Card, CardContent, CardHeader, CardTitle, EmptyState, ErrorAlert, Input } from '@/components/ui'
import { Building2 } from 'lucide-vue-next'
import type { Cinema } from '@/types/catalog'

const { t } = useI18n()

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
    errorMessage.value =
      error instanceof ApiError
        ? translateApiError(error.code, error.message)
        : t('admin.cinemas.loadFailed')
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
    errorMessage.value =
      error instanceof ApiError
        ? translateApiError(error.code, error.message)
        : t('admin.cinemas.saveFailed')
  }
}

async function onDelete(id: string) {
  if (!confirm(t('admin.cinemas.confirmDelete'))) return
  try {
    await api.delete(`/admin/cinemas/${id}`)
    if (editingId.value === id) resetForm()
    await loadCinemas()
  } catch (error) {
    errorMessage.value =
      error instanceof ApiError
        ? translateApiError(error.code, error.message)
        : t('admin.cinemas.deleteFailed')
  }
}

onMounted(loadCinemas)
</script>

<template>
  <div class="space-y-8">
    <div>
      <h1 class="text-2xl font-semibold text-copy-primary">{{ t('admin.cinemas.title') }}</h1>
      <p class="mt-1 text-sm text-copy-secondary">{{ t('admin.cinemas.subtitle') }}</p>
    </div>

    <Card>
      <CardHeader>
        <CardTitle>{{ editingId ? t('admin.cinemas.editTitle') : t('admin.cinemas.addTitle') }}</CardTitle>
      </CardHeader>
      <CardContent>
        <form class="grid gap-4 md:grid-cols-2" @submit.prevent="onSubmit">
          <div class="space-y-2 md:col-span-2">
            <label class="text-sm text-copy-secondary" for="name">{{ t('common.name') }}</label>
            <Input id="name" v-model="form.name" required />
          </div>
          <div class="space-y-2 md:col-span-2">
            <label class="text-sm text-copy-secondary" for="address">{{ t('common.address') }}</label>
            <Input id="address" v-model="form.address" />
          </div>
          <div class="space-y-2 md:col-span-2">
            <label class="text-sm text-copy-secondary" for="timezone">{{ t('admin.cinemas.timezone') }}</label>
            <Input id="timezone" v-model="form.timezone" :placeholder="t('admin.cinemas.timezonePlaceholder')" required />
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

    <ErrorAlert v-if="errorMessage" :message="errorMessage" />

    <Card>
      <CardHeader>
        <CardTitle>{{ t('admin.cinemas.allTitle') }}</CardTitle>
      </CardHeader>
      <CardContent>
        <TableSkeleton v-if="loading" :columns="4" :rows="4" />
        <EmptyState
          v-else-if="!cinemas.length"
          :icon="Building2"
          :title="t('admin.cinemas.empty')"
          class="py-10"
        />
        <div v-else class="overflow-x-auto">
          <table class="w-full text-left text-sm">
            <thead class="text-copy-muted">
              <tr>
                <th class="pb-3 pr-4 font-medium">{{ t('common.name') }}</th>
                <th class="pb-3 pr-4 font-medium">{{ t('admin.cinemas.timezone') }}</th>
                <th class="pb-3 pr-4 font-medium">{{ t('common.address') }}</th>
                <th class="pb-3 font-medium">{{ t('common.actions') }}</th>
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
                    <Button variant="ghost" @click="startEdit(cinema)">{{ t('common.edit') }}</Button>
                    <Button variant="destructive" @click="onDelete(cinema.id)">{{ t('common.delete') }}</Button>
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
