<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router'

import { Button } from '@/components/ui'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()

const navItems = [
  { to: '/admin/movies', label: 'Movies' },
  { to: '/admin/cinemas', label: 'Cinemas' },
  { to: '/admin/screens', label: 'Screens' },
  { to: '/admin/showtimes', label: 'Showtimes' },
]
</script>

<template>
  <div class="min-h-screen bg-base">
    <header class="border-b border-surface-border bg-surface">
      <div class="mx-auto flex max-w-6xl flex-wrap items-center justify-between gap-4 px-4 py-4">
        <div>
          <RouterLink
            to="/"
            class="bg-gradient-brand bg-clip-text text-lg font-semibold text-transparent"
          >
            Cinema Tickets
          </RouterLink>
          <p class="text-sm text-copy-secondary">Admin catalog</p>
        </div>
        <div class="flex items-center gap-3">
          <span v-if="auth.user" class="text-sm text-copy-muted">{{ auth.user.email }}</span>
          <Button variant="secondary" @click="auth.logout()">Sign out</Button>
        </div>
      </div>
      <nav class="mx-auto flex max-w-6xl gap-1 overflow-x-auto px-4 pb-3">
        <RouterLink
          v-for="item in navItems"
          :key="item.to"
          :to="item.to"
          class="rounded-full px-4 py-2 text-sm font-medium text-copy-secondary transition-colors hover:bg-accent-dim hover:text-brand"
          active-class="!bg-accent-dim !text-brand"
        >
          {{ item.label }}
        </RouterLink>
      </nav>
    </header>

    <main class="mx-auto max-w-6xl px-4 py-8">
      <RouterView />
    </main>
  </div>
</template>
