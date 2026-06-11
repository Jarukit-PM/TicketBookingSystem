# Feature 07 — Seat Map, Redis Holds & WebSocket

Read `context/CONTEXT.md`, `context/architecture-context.md`, `context/ui-context.md`, `create-adaptable-composable`, and `tdd` skills.

Interactive **seat map** per showtime, **Redis seat holds** (5 min TTL, refresh on add), **WebSocket** real-time sync, hold **countdown** UI. Guest can view map; auth required to hold (`CONTEXT.md`).

**Depends on:** Features 01, 03, 04, 05, 06 (showtime + screen layout).

## Objective

Multiple users see live seat availability without double-selecting. Holds enforce all rules in `CONTEXT.md`. HTTP hold/confirm responses are authoritative; WebSocket is advisory.

**Success looks like:** Two browsers on same showtime see holds/solds update sub-second; TTL countdown accurate; deselect releases seat; max 10 seats; no hold after `startsAt`.

## Hold rules (authoritative — `CONTEXT.md`)

| Rule | Behavior |
| ---- | -------- |
| TTL | 5 minutes; **refresh all holds on add** |
| Deselect | Immediate release; **no TTL refresh on remove** |
| Navigate away | Holds remain until TTL (no WS disconnect release) |
| Concurrent showtimes | Independent hold sets per showtime |
| Seat limit | Max **10** per user per showtime |
| Cutoff | Reject if `startsAt <= now` (cinema TZ) |
| Inventory | `AVAILABLE = layout − SOLD − BLOCKED − others' holds` |
| SOLD | From confirmed `bookings` query |
| BLOCKED | Layout `type: blocked` only |

### Redis keys

```
hold:{showtimeId}:{seatId}  →  { userId, heldAt }   TTL 5m
user_holds:{userId}:{showtimeId}  →  SET seatIds     TTL 5m
```

Use `SET NX` on hold; `EXPIRE` all user holds on add.

## API Surface

| Method | Path | Auth | Body |
| ------ | ---- | ---- | ---- |
| GET | `/api/showtimes/:id/seats` | Public | Snapshot: seats with status, `priceTiers`, layout |
| POST | `/api/showtimes/:id/holds` | Customer | `{ seatIds: ["A-1"] }` → `{ holds, expiresAt }` |
| DELETE | `/api/showtimes/:id/holds` | Customer | `{ seatIds?: [...] }` — omit = release all |
| WS | `/ws/showtimes/:id` | Optional | Events below |

### WebSocket events (server → client)

| Event | Payload |
| ----- | ------- |
| `snapshot` | Full map on connect |
| `seat_held` | `{ seatId, expiresAt }` |
| `seat_released` | `{ seatId }` |
| `seat_sold` | `{ seatId }` |

Multi-instance: Redis pub/sub fan-out (`architecture-context.md`).

## Vue Structure

```
app/src/views/SeatMapView.vue
app/src/components/seat-map/
├── SeatMapGrid.vue
├── SeatCell.vue
├── HoldCountdown.vue
└── SeatLegend.vue
app/src/composables/
├── useShowtimeSocket.ts    # MaybeRef showtimeId, reconnect, snapshot
└── useHoldCountdown.ts     # expiresAt from server
app/src/stores/bookingSession.ts  # selected seats, expiresAt for showtime
```

### UX (`ui-context.md`)

- Seat colors: available, held (self/other), sold, blocked, selected
- Self-held seats: gradient border
- Countdown: prominent; urgency styling when &lt; 60s
- Guest click seat → redirect login with `redirect` back
- Optimistic UI OK; reconcile on REST error

## Code Style

Composable pattern:

```ts
export function useShowtimeSocket(showtimeId: MaybeRefOrGetter<string>) {
  watchEffect(() => {
    const id = toValue(showtimeId)
    if (!id) return
    connect(`/ws/showtimes/${id}`)
  })
}
```

Go hold service — sorted seat IDs when multi-seat:

```go
func (s *HoldService) AddSeats(ctx context.Context, userID, showtimeID string, seatIDs []string) (expiresAt time.Time, err error)
```

## Testing Strategy (**tdd** for hold rules)

- Go table tests: TTL refresh on add; no refresh on remove; max 10; conflict on other's hold; sold/blocked rejection; cutoff after `startsAt`.
- Go: `SET NX` race simulation with miniredis.
- Vue: unit test `useHoldCountdown` ticks from `expiresAt`.
- Manual two-browser test required before closing spec.

## Boundaries

- **Always:** REST authoritative over WS; publish WS after successful hold/release.
- **Ask first:** Changing TTL duration; keyspace notifications vs polling.
- **Never:** Store holds in MongoDB; release on WS disconnect only.

## Tasks

- [ ] `internal/hold` service + Redis implementation
- [ ] `internal/inventory` — compute seat statuses from bookings + holds + layout
- [ ] HTTP handlers GET seats, POST/DELETE holds
- [ ] `internal/ws` hub + Redis pub/sub
- [ ] nginx WS proxy verified (spec 02)
- [ ] `useShowtimeSocket`, `useHoldCountdown`, seat map UI
- [ ] Login gate on first seat select for guests
- [ ] Order summary stub route → spec 08 confirm

## Out of scope

- Booking confirm (spec 08)
- Payment
- Per-showtime seat blocks

## Check when done

- [ ] All hold rules in `CONTEXT.md` covered by tests
- [ ] Two concurrent users cannot hold same seat
- [ ] Countdown matches server `expiresAt`
- [ ] Sub-second WS updates on local compose
- [ ] `progress-tracker.md` updated when implementation lands

## Open Questions

- None.
