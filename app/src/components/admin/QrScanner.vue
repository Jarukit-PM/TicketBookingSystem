<script setup lang="ts">
import { BrowserQRCodeReader, type IScannerControls } from '@zxing/browser'
import { onBeforeUnmount, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'

const emit = defineEmits<{
  scan: [value: string]
  error: [message: string]
}>()

const { t } = useI18n()

const videoRef = ref<HTMLVideoElement | null>(null)
const fileInputRef = ref<HTMLInputElement | null>(null)
const cameraActive = ref(false)
const cameraError = ref('')
const scanning = ref(false)

let controls: IScannerControls | null = null
const reader = new BrowserQRCodeReader()

function stopCamera(): void {
  controls?.stop()
  controls = null
  cameraActive.value = false
}

async function startCamera(): Promise<void> {
  if (!videoRef.value || scanning.value) {
    return
  }

  cameraError.value = ''
  scanning.value = true

  try {
    controls = await reader.decodeFromVideoDevice(undefined, videoRef.value, (result, error) => {
      if (result) {
        emit('scan', result.getText())
        stopCamera()
      } else if (error && error.name !== 'NotFoundException') {
        cameraError.value = t('admin.qrScanner.readFailed')
      }
    })
    cameraActive.value = true
  } catch {
    cameraError.value = t('admin.qrScanner.cameraDenied')
    emit('error', cameraError.value)
  } finally {
    scanning.value = false
  }
}

async function onFileSelected(event: Event): Promise<void> {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  input.value = ''
  if (!file) {
    return
  }

  try {
    const url = URL.createObjectURL(file)
    try {
      const result = await reader.decodeFromImageUrl(url)
      emit('scan', result.getText())
    } finally {
      URL.revokeObjectURL(url)
    }
  } catch {
    emit('error', t('admin.qrScanner.imageReadFailed'))
  }
}

function openFilePicker(): void {
  fileInputRef.value?.click()
}

onMounted(() => {
  void startCamera()
})

onBeforeUnmount(() => {
  stopCamera()
})

defineExpose({ startCamera, stopCamera })
</script>

<template>
  <div class="space-y-4">
    <div
      class="relative overflow-hidden rounded-xl border border-surface-border bg-elevated"
      :aria-label="t('admin.qrScanner.previewLabel')"
    >
      <video
        ref="videoRef"
        class="aspect-[4/3] w-full bg-base object-cover"
        muted
        playsinline
      />
      <div
        v-if="!cameraActive && !scanning"
        class="absolute inset-0 flex items-center justify-center bg-base/80 px-6 text-center text-sm text-copy-secondary"
      >
        {{ cameraError || t('admin.qrScanner.startingCamera') }}
      </div>
    </div>

    <p v-if="cameraError" class="text-sm text-state-warning" role="status">{{ cameraError }}</p>

    <div class="flex flex-wrap gap-3">
      <button
        type="button"
        class="rounded-lg bg-gradient-brand px-4 py-2 text-sm font-medium text-copy-primary transition-opacity hover:opacity-90"
        @click="startCamera"
      >
        {{ cameraActive ? t('admin.qrScanner.restartCamera') : t('admin.qrScanner.tryCameraAgain') }}
      </button>
      <button
        type="button"
        class="rounded-lg border border-surface-border bg-surface px-4 py-2 text-sm font-medium text-copy-secondary transition-colors hover:bg-subtle hover:text-copy-primary"
        @click="openFilePicker"
      >
        {{ t('admin.qrScanner.uploadImage') }}
      </button>
    </div>

    <input
      ref="fileInputRef"
      type="file"
      accept="image/*"
      capture="environment"
      class="sr-only"
      @change="onFileSelected"
    />
  </div>
</template>
