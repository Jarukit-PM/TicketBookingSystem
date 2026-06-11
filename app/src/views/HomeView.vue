<script setup lang="ts">
import { ref } from 'vue'
import {
  Badge,
  Button,
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  Input,
} from '@/components/ui'

const searchQuery = ref('')
</script>

<template>
  <div class="min-h-screen bg-base">
    <header
      class="sticky top-0 z-10 flex h-16 items-center border-b border-surface-border bg-base/80 px-4 backdrop-blur-md md:px-6"
    >
      <span
        class="bg-gradient-brand bg-clip-text text-xl font-semibold tracking-tight text-transparent"
      >
        Cinema Tickets
      </span>
    </header>

    <main class="mx-auto max-w-6xl px-4 py-8 md:px-6 md:py-12">
      <div class="mb-8">
        <h1 class="text-3xl font-semibold tracking-tight text-copy-primary md:text-4xl">
          Design System Preview
        </h1>
        <p class="mt-2 text-sm text-copy-secondary">
          Dark cinema theme with gradient orange brand tokens.
        </p>
      </div>

      <section class="mb-8 flex flex-wrap gap-4">
        <Button variant="primary">View showtimes</Button>
        <Button variant="secondary">Browse movies</Button>
      </section>

      <section class="mb-8 grid gap-6 md:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>Dune: Part Three</CardTitle>
          </CardHeader>
          <CardContent>
            <p class="text-sm text-copy-secondary">
              Sci-fi epic · 2h 46m · IMAX · Starts 7:30 PM
            </p>
          </CardContent>
        </Card>

        <div class="flex flex-col gap-4">
          <label class="text-xs font-medium uppercase tracking-wide text-copy-muted" for="search">
            Search movies
          </label>
          <Input
            id="search"
            v-model="searchQuery"
            placeholder="Find a title..."
          />
        </div>
      </section>

      <section>
        <h2 class="mb-4 text-lg font-medium text-copy-primary">Booking status badges</h2>
        <div class="flex flex-wrap gap-3">
          <Badge variant="confirmed">Confirmed</Badge>
          <Badge variant="hold-active">Hold active</Badge>
          <Badge variant="hold-expired">Hold expired</Badge>
        </div>
      </section>
    </main>
  </div>
</template>
