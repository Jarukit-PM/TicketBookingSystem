<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { fetchSeatMap } from '@/api/seats'
import SeatLegend from '@/components/seat-map/SeatLegend.vue'
import SeatMapGrid from '@/components/seat-map/SeatMapGrid.vue'
import { Button, Card, CardContent, CardHeader, CardTitle } from '@/components/ui'
import { formatShowtime } from '@/lib/format'
import type { SeatMapSnapshot } from '@/types/seats'

const route = useRoute()
const router = useRouter()

const snapshot = ref<SeatMapSnapshot | null>(null)
const loading = ref(true)
const error = ref<string | null>(null)

const showtimeId = computed(() => route.params.showtimeId as string)

async function loadSeatMap(): Promise<void> {
  loading.value = true
  error.value = null
  try {
    snapshot.value = await fetchSeatMap(showtimeId.value)
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Failed to load seat map'
    snapshot.value = null
  } finally {
    loading.value = false
  }
}

function goBack(): void {
  if (window.history.length > 1) {
    router.back()
    return
  }
  router.push({ name: 'home' })
}

onMounted(loadSeatMap)
watch(showtimeId, loadSeatMap)
</script>

<template>
  <div class="min-h-screen bg-base">
    <header
      class="sticky top-0 z-10 flex h-16 items-center border-b border-surface-border bg-base/80 px-4 backdrop-blur-md md:px-6"
    >
      <Button variant="ghost" type="button" @click="goBack">Back</Button>
      <span class="ml-4 text-sm text-copy-secondary">Select seats</span>
    </header>

    <main class="mx-auto max-w-4xl px-4 py-8 md:px-6">
      <p v-if="loading" class="text-copy-secondary">Loading seat map…</p>
      <p v-else-if="error" class="text-state-error">{{ error }}</p>

      <template v-else-if="snapshot">
        <Card>
          <CardHeader>
            <CardTitle>{{ snapshot.screenName }}</CardTitle>
            <p class="text-sm text-copy-secondary">
              {{ formatShowtime(snapshot.startsAt) }}
            </p>
          </CardHeader>
          <CardContent class="space-y-6">
            <div
              class="rounded-lg border border-surface-border bg-elevated px-4 py-2 text-center text-xs uppercase tracking-widest text-copy-muted"
            >
              Screen
            </div>

            <SeatMapGrid :seats="snapshot.seats" />
            <SeatLegend />
          </CardContent>
        </Card>
      </template>
    </main>
  </div>
</template>
