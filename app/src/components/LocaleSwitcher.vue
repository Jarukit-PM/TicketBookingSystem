<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import {
  persistLocale,
  type AppLocale,
} from '@/i18n'

const { locale, t } = useI18n()

function setLocale(next: AppLocale): void {
  if (locale.value === next) return
  locale.value = next
  persistLocale(next)
}
</script>

<template>
  <div
    class="flex rounded-lg border border-surface-border bg-surface p-0.5 text-xs font-medium"
    role="group"
    :aria-label="t('locale.switcher')"
  >
    <button
      type="button"
      class="rounded-md px-2.5 py-1 transition-colors"
      :class="locale === 'en' ? 'bg-gradient-brand text-white' : 'text-copy-secondary hover:text-copy-primary'"
      :aria-pressed="locale === 'en'"
      @click="setLocale('en')"
    >
      EN
    </button>
    <button
      type="button"
      class="rounded-md px-2.5 py-1 transition-colors"
      :class="locale === 'th' ? 'bg-gradient-brand text-white' : 'text-copy-secondary hover:text-copy-primary'"
      :aria-pressed="locale === 'th'"
      @click="setLocale('th')"
    >
      ไทย
    </button>
  </div>
</template>
