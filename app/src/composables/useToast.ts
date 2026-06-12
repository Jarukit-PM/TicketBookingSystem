import { onUnmounted, ref } from 'vue'

const DEFAULT_DURATION_MS = 4000

export function useToast(durationMs = DEFAULT_DURATION_MS) {
  const message = ref('')
  let timer: ReturnType<typeof setTimeout> | null = null

  function show(text: string): void {
    message.value = text
    if (timer) clearTimeout(timer)
    timer = setTimeout(() => {
      message.value = ''
    }, durationMs)
  }

  onUnmounted(() => {
    if (timer) clearTimeout(timer)
  })

  return { message, show }
}
