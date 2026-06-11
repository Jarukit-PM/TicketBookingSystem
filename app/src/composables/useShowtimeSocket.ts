import { onScopeDispose, ref, shallowRef, toValue, watch } from 'vue'
import type { MaybeRefOrGetter } from 'vue'
import type { Seat, SeatMapSnapshot } from '@/types/seats'
import type {
  WsMessage,
  WsSeatHeldPayload,
  WsSeatReleasedPayload,
  WsSeatSoldPayload,
  WsSnapshotPayload,
} from '@/types/ws'

function wsUrl(showtimeId: string): string {
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  return `${protocol}//${window.location.host}/ws/showtimes/${showtimeId}`
}

function applySeatEvent(
  seats: Seat[],
  seatId: string,
  status: Seat['status'],
): Seat[] {
  return seats.map((seat) => (seat.seatId === seatId ? { ...seat, status } : seat))
}

export function useShowtimeSocket(showtimeId: MaybeRefOrGetter<string | null | undefined>) {
  const snapshot = shallowRef<SeatMapSnapshot | null>(null)
  const connected = ref(false)
  const error = ref<string | null>(null)

  let socket: WebSocket | null = null
  let reconnectTimer: ReturnType<typeof setTimeout> | null = null
  let disposed = false

  function clearReconnect() {
    if (reconnectTimer !== null) {
      clearTimeout(reconnectTimer)
      reconnectTimer = null
    }
  }

  function disconnect() {
    clearReconnect()
    if (socket) {
      socket.onopen = null
      socket.onclose = null
      socket.onerror = null
      socket.onmessage = null
      socket.close()
      socket = null
    }
    connected.value = false
  }

  function scheduleReconnect(id: string) {
    if (disposed) {
      return
    }
    clearReconnect()
    reconnectTimer = setTimeout(() => connect(id), 2000)
  }

  function handleMessage(event: MessageEvent<string>) {
    let msg: WsMessage
    try {
      msg = JSON.parse(event.data) as WsMessage
    } catch {
      return
    }

    if (!snapshot.value && msg.type !== 'snapshot') {
      return
    }

    switch (msg.type) {
      case 'snapshot': {
        const payload = msg.payload as WsSnapshotPayload
        snapshot.value = payload.snapshot
        break
      }
      case 'seat_held': {
        const payload = msg.payload as WsSeatHeldPayload
        if (snapshot.value) {
          snapshot.value = {
            ...snapshot.value,
            seats: applySeatEvent(snapshot.value.seats, payload.seatId, 'HELD'),
          }
        }
        break
      }
      case 'seat_released': {
        const payload = msg.payload as WsSeatReleasedPayload
        if (snapshot.value) {
          snapshot.value = {
            ...snapshot.value,
            seats: applySeatEvent(snapshot.value.seats, payload.seatId, 'AVAILABLE'),
          }
        }
        break
      }
      case 'seat_sold': {
        const payload = msg.payload as WsSeatSoldPayload
        if (snapshot.value) {
          snapshot.value = {
            ...snapshot.value,
            seats: applySeatEvent(snapshot.value.seats, payload.seatId, 'SOLD'),
          }
        }
        break
      }
      default:
        break
    }
  }

  function connect(id: string) {
    disconnect()
    error.value = null

    try {
      socket = new WebSocket(wsUrl(id))
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'WebSocket failed'
      scheduleReconnect(id)
      return
    }

    socket.onopen = () => {
      connected.value = true
      error.value = null
    }

    socket.onmessage = handleMessage

    socket.onerror = () => {
      error.value = 'Connection error'
    }

    socket.onclose = () => {
      connected.value = false
      if (!disposed) {
        scheduleReconnect(id)
      }
    }
  }

  watch(
    () => toValue(showtimeId),
    (id) => {
      disconnect()
      snapshot.value = null
      if (!id) {
        return
      }
      connect(id)
    },
    { immediate: true },
  )

  onScopeDispose(() => {
    disposed = true
    disconnect()
  })

  return {
    snapshot,
    connected,
    error,
  }
}
