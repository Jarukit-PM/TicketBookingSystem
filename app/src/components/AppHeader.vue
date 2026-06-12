<script setup lang="ts">
import { Film, LayoutDashboard, LogIn, LogOut, Ticket, UserPlus } from 'lucide-vue-next'
import { useI18n } from 'vue-i18n'
import { RouterLink } from 'vue-router'
import LocaleSwitcher from '@/components/LocaleSwitcher.vue'
import { Button } from '@/components/ui'
import { useAuthStore } from '@/stores/auth'

withDefaults(
  defineProps<{
    showBack?: boolean
    backLabel?: string
    subtitle?: string
  }>(),
  {
    showBack: false,
    backLabel: undefined,
    subtitle: undefined,
  },
)

const emit = defineEmits<{
  back: []
}>()

const { t } = useI18n()
const auth = useAuthStore()

function handleBack(): void {
  emit('back')
}
</script>

<template>
  <header
    class="sticky top-0 z-10 border-b border-surface-border bg-base/80 backdrop-blur-md"
  >
    <div
      class="bg-gradient-brand-subtle h-px w-full"
      aria-hidden="true"
    />
    <div class="flex h-16 min-w-0 items-center gap-2 px-4 sm:gap-3 md:px-6">
      <button
        v-if="showBack"
        type="button"
        class="inline-flex items-center gap-1.5 rounded-full px-2 py-1.5 text-sm text-copy-secondary transition-colors hover:bg-subtle hover:text-copy-primary"
        @click="handleBack"
      >
        <span aria-hidden="true">←</span>
        {{ backLabel ?? t('common.back') }}
      </button>

      <RouterLink
        to="/"
        class="inline-flex min-w-0 shrink items-center gap-2 transition-opacity hover:opacity-90"
      >
        <Film class="h-5 w-5 shrink-0 text-brand" aria-hidden="true" />
        <span
          class="truncate bg-gradient-brand bg-clip-text text-lg font-semibold tracking-tight text-transparent sm:text-xl"
        >
          {{ t('common.appName') }}
        </span>
      </RouterLink>

      <p v-if="subtitle" class="hidden text-sm text-copy-secondary sm:inline">
        {{ subtitle }}
      </p>

      <nav class="ml-auto flex shrink-0 items-center gap-0.5 sm:gap-2">
        <LocaleSwitcher />

        <template v-if="auth.isAuthenticated">
          <RouterLink
            to="/my-bookings"
            class="inline-flex items-center gap-1.5 rounded-full px-3 py-2 text-sm text-copy-secondary transition-colors hover:bg-subtle hover:text-copy-primary"
          >
            <Ticket class="h-4 w-4" aria-hidden="true" />
            <span class="hidden sm:inline">{{ t('nav.myBookings') }}</span>
          </RouterLink>
          <RouterLink
            v-if="auth.isAdmin"
            to="/admin"
            class="inline-flex items-center gap-1.5 rounded-full px-3 py-2 text-sm text-copy-secondary transition-colors hover:bg-subtle hover:text-copy-primary"
          >
            <LayoutDashboard class="h-4 w-4" aria-hidden="true" />
            <span class="hidden sm:inline">{{ t('nav.admin') }}</span>
          </RouterLink>
          <span class="hidden max-w-[10rem] truncate text-sm text-copy-muted lg:inline">
            {{ auth.user?.email }}
          </span>
          <Button variant="secondary" class="gap-1.5" @click="auth.logout()">
            <LogOut class="h-4 w-4" aria-hidden="true" />
            <span class="hidden sm:inline">{{ t('nav.signOut') }}</span>
          </Button>
        </template>
        <template v-else>
          <RouterLink
            to="/login"
            class="inline-flex items-center gap-1.5 rounded-full px-3 py-2 text-sm text-copy-secondary transition-colors hover:bg-subtle hover:text-copy-primary"
          >
            <LogIn class="h-4 w-4" aria-hidden="true" />
            <span class="hidden sm:inline">{{ t('nav.signIn') }}</span>
          </RouterLink>
          <RouterLink to="/register">
            <Button variant="secondary" class="gap-1.5">
              <UserPlus class="h-4 w-4" aria-hidden="true" />
              <span class="hidden sm:inline">{{ t('nav.register') }}</span>
            </Button>
          </RouterLink>
        </template>
      </nav>
    </div>
  </header>
</template>
