<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import LocaleSwitcher from '@/components/LocaleSwitcher.vue'

const { t, locale } = useI18n()

const currentLabel = computed(() =>
  locale.value === 'th' ? t('locale.th') : t('locale.en'),
)
</script>

<template>
  <div class="rounded-lg border border-surface-border bg-subtle/40 px-3 py-3">
    <p class="text-sm text-copy-secondary">{{ t('booking.checkout.emailLanguageHint') }}</p>
    <div class="mt-3 flex flex-wrap items-center justify-between gap-3">
      <div>
        <p class="text-sm font-medium text-copy-primary">{{ t('booking.checkout.emailLanguage') }}</p>
        <p class="text-xs text-copy-muted">{{ currentLabel }}</p>
      </div>
      <LocaleSwitcher />
    </div>
  </div>
</template>
