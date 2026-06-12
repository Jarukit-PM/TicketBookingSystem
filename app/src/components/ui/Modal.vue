<script setup lang="ts">
import { X } from 'lucide-vue-next'
import { computed, nextTick, onMounted, onUnmounted, ref, useId, watch } from 'vue'
import { useI18n } from 'vue-i18n'

import { cn } from '@/lib/cn'

type ModalSize = 'md' | 'lg' | 'xl'

const props = withDefaults(
  defineProps<{
    open: boolean
    title?: string
    titleId?: string
    size?: ModalSize
  }>(),
  {
    title: undefined,
    titleId: undefined,
    size: 'lg',
  },
)

const emit = defineEmits<{
  'update:open': [value: boolean]
}>()

const { t } = useI18n()
const closeButtonRef = ref<HTMLButtonElement | null>(null)
const generatedTitleId = useId()
const resolvedTitleId = computed(() => props.titleId ?? generatedTitleId)

const sizeClasses: Record<ModalSize, string> = {
  md: 'max-w-md',
  lg: 'max-w-lg',
  xl: 'max-w-2xl',
}

function close() {
  emit('update:open', false)
}

function onBackdropClick(event: MouseEvent) {
  if (event.target === event.currentTarget) close()
}

function onKeydown(event: KeyboardEvent) {
  if (event.key === 'Escape' && props.open) {
    event.preventDefault()
    close()
  }
}

onMounted(() => document.addEventListener('keydown', onKeydown))

watch(
  () => props.open,
  async (isOpen) => {
    if (isOpen) {
      document.body.style.overflow = 'hidden'
      await nextTick()
      closeButtonRef.value?.focus()
      return
    }
    document.body.style.overflow = ''
  },
)

onUnmounted(() => {
  document.removeEventListener('keydown', onKeydown)
  document.body.style.overflow = ''
})
</script>

<template>
  <Teleport to="body">
    <Transition
      enter-active-class="transition-opacity duration-200 ease-out"
      enter-from-class="opacity-0"
      enter-to-class="opacity-100"
      leave-active-class="transition-opacity duration-150 ease-in"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
      <div
        v-if="open"
        class="fixed inset-0 z-50 flex items-end justify-center bg-base/80 p-4 backdrop-blur-sm sm:items-center"
        role="presentation"
        @click="onBackdropClick"
      >
        <Transition
          enter-active-class="transition-all duration-200 ease-out"
          enter-from-class="opacity-0 translate-y-4 sm:translate-y-2 sm:scale-95"
          enter-to-class="opacity-100 translate-y-0 sm:scale-100"
          leave-active-class="transition-all duration-150 ease-in"
          leave-from-class="opacity-100 translate-y-0 sm:scale-100"
          leave-to-class="opacity-0 translate-y-4 sm:translate-y-2 sm:scale-95"
        >
          <div
            v-if="open"
            role="dialog"
            aria-modal="true"
            :aria-labelledby="title ? resolvedTitleId : undefined"
            :class="
              cn(
                'flex max-h-[min(90vh,48rem)] w-full flex-col overflow-hidden rounded-2xl border border-surface-border bg-elevated shadow-elevation-2',
                sizeClasses[size],
              )
            "
            @click.stop
          >
            <div
              v-if="title || $slots.header"
              class="flex items-start justify-between gap-3 border-b border-surface-border px-5 py-4"
            >
              <slot name="header">
                <h2 :id="resolvedTitleId" class="text-lg font-semibold text-copy-primary">
                  {{ title }}
                </h2>
              </slot>
              <button
                ref="closeButtonRef"
                type="button"
                class="inline-flex h-10 w-10 shrink-0 items-center justify-center rounded-full text-copy-secondary transition-colors hover:bg-subtle hover:text-copy-primary focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-accent-glow"
                :aria-label="t('common.closeDialog')"
                @click="close"
              >
                <X class="h-5 w-5" aria-hidden="true" />
              </button>
            </div>

            <div class="overflow-y-auto px-5 py-4">
              <slot />
            </div>

            <div
              v-if="$slots.footer"
              class="border-t border-surface-border px-5 py-4"
            >
              <slot name="footer" />
            </div>
          </div>
        </Transition>
      </div>
    </Transition>
  </Teleport>
</template>
