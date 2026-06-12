<script setup lang="ts">
import {
  Building2,
  Calendar,
  Clapperboard,
  FileText,
  LayoutDashboard,
  LogOut,
  Menu,
  Monitor,
  QrCode,
  Ticket,
  X,
} from 'lucide-vue-next'
import type { Component } from 'vue'
import { computed, onUnmounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { RouterLink, RouterView, useRoute } from 'vue-router'

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
const route = useRoute()

const mobileNavOpen = ref(false)

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

function closeMobileNav(): void {
  mobileNavOpen.value = false
}

function toggleMobileNav(): void {
  mobileNavOpen.value = !mobileNavOpen.value
}

watch(
  () => route.fullPath,
  () => {
    closeMobileNav()
  },
)

watch(mobileNavOpen, (isOpen) => {
  document.body.style.overflow = isOpen ? 'hidden' : ''
})

onUnmounted(() => {
  document.body.style.overflow = ''
})
</script>

<template>
  <div class="flex min-h-screen bg-base">
    <button
      v-if="mobileNavOpen"
      type="button"
      class="fixed inset-0 z-40 bg-base/80 backdrop-blur-sm lg:hidden"
      :aria-label="t('common.closeMenu')"
      @click="closeMobileNav"
    />

    <aside
      class="fixed inset-y-0 left-0 z-50 flex w-[min(16rem,85vw)] flex-col border-r border-surface-border bg-surface transition-transform duration-200 ease-out lg:static lg:z-auto lg:w-56 lg:shrink-0 lg:translate-x-0"
      :class="mobileNavOpen ? 'translate-x-0' : '-translate-x-full lg:translate-x-0'"
    >
      <div class="flex items-start justify-between gap-2 border-b border-surface-border px-4 py-5">
        <div class="min-w-0">
          <RouterLink
            to="/"
            class="inline-flex items-center gap-2 bg-gradient-brand bg-clip-text text-lg font-semibold text-transparent"
            @click="closeMobileNav"
          >
            <Ticket class="h-5 w-5 shrink-0 text-brand" aria-hidden="true" />
            <span class="truncate">{{ t('admin.brand') }}</span>
          </RouterLink>
          <p class="mt-1 text-xs text-copy-muted">{{ t('admin.console') }}</p>
        </div>
        <button
          type="button"
          class="inline-flex h-10 w-10 shrink-0 items-center justify-center rounded-full text-copy-secondary transition-colors hover:bg-subtle hover:text-copy-primary focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-accent-glow lg:hidden"
          :aria-label="t('common.closeMenu')"
          @click="closeMobileNav"
        >
          <X class="h-5 w-5" aria-hidden="true" />
        </button>
      </div>

      <nav class="flex-1 space-y-1 overflow-y-auto p-3">
        <RouterLink
          v-for="item in navItems"
          :key="item.to"
          :to="item.to"
          class="flex items-center gap-2.5 rounded-lg px-3 py-2.5 text-sm font-medium text-copy-secondary transition-colors hover:bg-accent-dim hover:text-brand"
          active-class="!bg-accent-dim !text-brand"
          :exact-active-class="item.exact ? '!bg-accent-dim !text-brand' : undefined"
          @click="closeMobileNav"
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

    <div class="flex min-w-0 flex-1 flex-col">
      <header
        class="sticky top-0 z-30 flex h-14 items-center gap-3 border-b border-surface-border bg-base/90 px-4 backdrop-blur-md lg:hidden"
      >
        <button
          type="button"
          class="inline-flex h-10 w-10 shrink-0 items-center justify-center rounded-full text-copy-secondary transition-colors hover:bg-subtle hover:text-copy-primary focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-accent-glow"
          :aria-label="t('common.openMenu')"
          :aria-expanded="mobileNavOpen"
          @click="toggleMobileNav"
        >
          <Menu class="h-5 w-5" aria-hidden="true" />
        </button>
        <p class="min-w-0 truncate text-sm font-medium text-copy-primary">
          {{ t('admin.console') }}
        </p>
      </header>

      <main class="flex-1 overflow-auto p-4 pb-[max(1rem,env(safe-area-inset-bottom))] md:p-6 lg:p-8">
        <div class="mx-auto max-w-6xl">
          <RouterView />
        </div>
      </main>
    </div>
  </div>
</template>
