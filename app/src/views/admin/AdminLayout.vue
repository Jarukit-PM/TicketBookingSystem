<script setup lang="ts">
import {
  Building2,
  Calendar,
  Clapperboard,
  FileText,
  LayoutDashboard,
  LogOut,
  Monitor,
  QrCode,
  Ticket,
} from 'lucide-vue-next'
import type { Component } from 'vue'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { RouterLink, RouterView } from 'vue-router'

import LocaleSwitcher from '@/components/LocaleSwitcher.vue'
import { Button } from '@/components/ui'
import { useAuthStore } from '@/stores/auth'

type NavItem = {
  to: string
  label: string
  icon: Component
  exact?: boolean
}

const { t } = useI18n()
const auth = useAuthStore()

const navItems = computed<NavItem[]>(() => [
  { to: '/admin', label: t('admin.nav.dashboard'), icon: LayoutDashboard, exact: true },
  { to: '/admin/movies', label: t('admin.nav.movies'), icon: Clapperboard },
  { to: '/admin/cinemas', label: t('admin.nav.cinemas'), icon: Building2 },
  { to: '/admin/screens', label: t('admin.nav.screens'), icon: Monitor },
  { to: '/admin/showtimes', label: t('admin.nav.showtimes'), icon: Calendar },
  { to: '/admin/bookings', label: t('admin.nav.bookings'), icon: Ticket },
  { to: '/admin/scan', label: t('admin.nav.scan'), icon: QrCode },
  { to: '/admin/logs', label: t('admin.nav.logs'), icon: FileText },
])
</script>

<template>
  <div class="flex min-h-screen bg-base">
    <aside class="flex w-56 shrink-0 flex-col border-r border-surface-border bg-surface">
      <div class="border-b border-surface-border px-4 py-5">
        <RouterLink
          to="/"
          class="inline-flex items-center gap-2 bg-gradient-brand bg-clip-text text-lg font-semibold text-transparent"
        >
          <Ticket class="h-5 w-5 text-brand" aria-hidden="true" />
          {{ t('admin.brand') }}
        </RouterLink>
        <p class="mt-1 text-xs text-copy-muted">{{ t('admin.console') }}</p>
      </div>

      <nav class="flex-1 space-y-1 p-3">
        <RouterLink
          v-for="item in navItems"
          :key="item.to"
          :to="item.to"
          class="flex items-center gap-2.5 rounded-lg px-3 py-2 text-sm font-medium text-copy-secondary transition-colors hover:bg-accent-dim hover:text-brand"
          active-class="!bg-accent-dim !text-brand"
          :exact-active-class="item.exact ? '!bg-accent-dim !text-brand' : undefined"
        >
          <component :is="item.icon" class="h-4 w-4 shrink-0" aria-hidden="true" />
          {{ item.label }}
        </RouterLink>
      </nav>

      <div class="border-t border-surface-border p-4">
        <p v-if="auth.user" class="truncate text-xs text-copy-muted">{{ auth.user.email }}</p>
        <LocaleSwitcher class="mt-3" />
        <Button variant="secondary" class="mt-3 w-full gap-1.5" @click="auth.logout()">
          <LogOut class="h-4 w-4" aria-hidden="true" />
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
