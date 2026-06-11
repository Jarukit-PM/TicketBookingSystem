<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { Button, Card, CardContent } from '@/components/ui'
import { formatDuration } from '@/lib/format'
import type { Movie } from '@/types/catalog'

const props = defineProps<{
  movie: Movie
  tab: 'now_showing' | 'coming_soon'
}>()

const router = useRouter()

const posterStyle = computed(() => ({
  backgroundImage: props.movie.posterUrl ? `url(${props.movie.posterUrl})` : undefined,
}))

function openMovie(): void {
  router.push({ name: 'movie-detail', params: { id: props.movie.id } })
}
</script>

<template>
  <Card class="group overflow-hidden transition-shadow hover:shadow-glow-brand/20">
    <div
      class="aspect-[2/3] w-full rounded-t-xl bg-subtle bg-cover bg-center ring-1 ring-white/10"
      :style="posterStyle"
      role="img"
      :aria-label="`${movie.title} poster`"
    />
    <CardContent class="flex flex-col gap-3 p-4">
      <div>
        <h3 class="text-lg font-semibold text-copy-primary">{{ movie.title }}</h3>
        <p class="mt-1 text-sm text-copy-secondary">
          {{ movie.rating }} · {{ formatDuration(movie.durationMin) }}
        </p>
      </div>
      <Button variant="primary" class="w-full" @click="openMovie">
        {{ tab === 'coming_soon' ? 'View details' : 'View showtimes' }}
      </Button>
    </CardContent>
  </Card>
</template>
