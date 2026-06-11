<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router'

import { Button } from '@/components/ui'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()

const navItems = [
  { to: '/admin/bookings', label: 'Bookings' },
]
</script>

<template>
  <div class="flex min-h-screen bg-base">
    <aside class="flex w-56 shrink-0 flex-col border-r border-surface-border bg-surface">
      <div class="border-b border-surface-border px-4 py-5">
        <RouterLink
          to="/"
          class="bg-gradient-brand bg-clip-text text-lg font-semibold text-transparent"
        >
          Cinema Tickets
        </RouterLink>
        <p class="mt-1 text-xs text-copy-muted">Admin console</p>
      </div>

      <nav class="flex-1 space-y-1 p-3">
        <RouterLink
          v-for="item in navItems"
          :key="item.to"
          :to="item.to"
          class="block rounded-lg px-3 py-2 text-sm font-medium text-copy-secondary transition-colors hover:bg-accent-dim hover:text-brand"
          active-class="!bg-accent-dim !text-brand"
        >
          {{ item.label }}
        </RouterLink>
      </nav>

      <div class="border-t border-surface-border p-4">
        <p v-if="auth.user" class="truncate text-xs text-copy-muted">{{ auth.user.email }}</p>
        <Button variant="secondary" class="mt-3 w-full" @click="auth.logout()">Sign out</Button>
      </div>
    </aside>

    <main class="flex-1 overflow-auto p-6 md:p-8">
      <div class="mx-auto max-w-6xl">
        <RouterView />
      </div>
    </main>
  </div>
</template>
