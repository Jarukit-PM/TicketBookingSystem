<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { RouterLink, RouterView } from 'vue-router'

import LocaleSwitcher from '@/components/LocaleSwitcher.vue'
import { Button } from '@/components/ui'
import { useAuthStore } from '@/stores/auth'

const { t } = useI18n()
const auth = useAuthStore()

const navItems = computed(() => [
  { to: '/admin', label: t('admin.nav.dashboard'), exact: true },
  { to: '/admin/movies', label: t('admin.nav.movies') },
  { to: '/admin/cinemas', label: t('admin.nav.cinemas') },
  { to: '/admin/screens', label: t('admin.nav.screens') },
  { to: '/admin/showtimes', label: t('admin.nav.showtimes') },
  { to: '/admin/bookings', label: t('admin.nav.bookings') },
  { to: '/admin/scan', label: t('admin.nav.scan') },
  { to: '/admin/logs', label: t('admin.nav.logs') },
])
</script>

<template>
  <div class="flex min-h-screen bg-base">
    <aside class="flex w-56 shrink-0 flex-col border-r border-surface-border bg-surface">
      <div class="border-b border-surface-border px-4 py-5">
        <RouterLink
          to="/"
          class="bg-gradient-brand bg-clip-text text-lg font-semibold text-transparent"
        >
          {{ t('admin.brand') }}
        </RouterLink>
        <p class="mt-1 text-xs text-copy-muted">{{ t('admin.console') }}</p>
      </div>

      <nav class="flex-1 space-y-1 p-3">
        <RouterLink
          v-for="item in navItems"
          :key="item.to"
          :to="item.to"
          class="block rounded-lg px-3 py-2 text-sm font-medium text-copy-secondary transition-colors hover:bg-accent-dim hover:text-brand"
          active-class="!bg-accent-dim !text-brand"
          :exact-active-class="item.exact ? '!bg-accent-dim !text-brand' : undefined"
        >
          {{ item.label }}
        </RouterLink>
      </nav>

      <div class="border-t border-surface-border p-4">
        <p v-if="auth.user" class="truncate text-xs text-copy-muted">{{ auth.user.email }}</p>
        <LocaleSwitcher class="mt-3" />
        <Button variant="secondary" class="mt-3 w-full" @click="auth.logout()">
          {{ t('nav.signOut') }}
        </Button>
      </div>
    </aside>

    <main class="flex-1 overflow-auto p-6 md:p-8">
      <div class="mx-auto max-w-6xl">
        <RouterView />
      </div>
    </main>
  </div>
</template>
