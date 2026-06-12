<script setup lang="ts">
import { Accessibility, Crown } from 'lucide-vue-next'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useLocaleFormat } from '@/composables/useLocaleFormat'

const props = defineProps<{
  priceTiers?: Record<string, number>
}>()

const { t } = useI18n()
const { formatTHB } = useLocaleFormat()

const statusItems = computed(() => [
  { label: t('seatMap.legend.available'), class: 'bg-subtle border-surface-border' },
  { label: t('seatMap.legend.yourSeats'), class: 'bg-accent-dim border-brand/60' },
  { label: t('seatMap.legend.held'), class: 'bg-state-warning-dim border-state-warning' },
  { label: t('seatMap.legend.sold'), class: 'bg-elevated border-surface-border-subtle' },
  { label: t('seatMap.legend.blocked'), class: 'bg-state-error-dim border-state-error/50' },
])

const typeItems = computed(() => {
  const tiers = props.priceTiers
  return [
    {
      id: 'standard',
      label: t('seatMap.types.standard'),
      price: tiers?.standard,
      swatchClass: 'bg-subtle border-surface-border',
    },
    {
      id: 'vip',
      label: t('seatMap.types.vip'),
      price: tiers?.vip,
      swatchClass: 'bg-accent-dim border-brand/60',
      icon: Crown,
      iconClass: 'text-brand',
    },
    {
      id: 'wheelchair',
      label: t('seatMap.types.wheelchair'),
      price: tiers?.wheelchair,
      swatchClass: 'bg-state-success-dim border-state-success/50',
      icon: Accessibility,
      iconClass: 'text-state-success',
    },
  ]
})
</script>

<template>
  <div class="space-y-4 border-t border-surface-border pt-4">
    <div>
      <p class="mb-2 text-xs font-medium uppercase tracking-wide text-copy-muted">
        {{ t('seatMap.legend.statusHeading') }}
      </p>
      <ul class="flex flex-wrap gap-x-4 gap-y-2 text-sm text-copy-secondary">
        <li v-for="item in statusItems" :key="item.label" class="flex items-center gap-2">
          <span
            class="inline-block h-4 w-4 rounded border"
            :class="item.class"
            aria-hidden="true"
          />
          {{ item.label }}
        </li>
      </ul>
    </div>

    <div>
      <p class="mb-2 text-xs font-medium uppercase tracking-wide text-copy-muted">
        {{ t('seatMap.legend.typesHeading') }}
      </p>
      <ul class="flex flex-wrap gap-x-4 gap-y-2 text-sm text-copy-secondary">
        <li v-for="item in typeItems" :key="item.id" class="flex items-center gap-2">
          <span
            class="relative inline-flex h-4 w-4 items-center justify-center rounded border"
            :class="item.swatchClass"
            aria-hidden="true"
          >
            <component
              :is="item.icon"
              v-if="item.icon"
              class="h-2.5 w-2.5"
              :class="item.iconClass"
            />
          </span>
          <span>
            {{ item.label }}
            <span v-if="item.price != null" class="text-copy-muted">· {{ formatTHB(item.price) }}</span>
          </span>
        </li>
      </ul>
    </div>
  </div>
</template>
