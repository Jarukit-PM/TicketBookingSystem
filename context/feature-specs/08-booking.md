# Feature 08 — Booking Confirm & My Bookings

Read `context/CONTEXT.md`, `context/architecture-context.md`, and `tdd` skill.

**Confirm** books the user's **entire hold set** for one showtime as a single **CONFIRMED** MongoDB booking. **My Bookings** lists upcoming and past tickets. Idempotent confirm via `Idempotency-Key`.

**Depends on:** Features 04, 05, 07 (holds must exist).

## Objective

Customer completes checkout without duplicate bookings under retry or concurrency. Booking history reflects **confirmed bookings only** — not Redis holds (`CONTEXT.md`).

**Success looks like:** Confirm → booking row + seats SOLD + holds cleared; retry same idempotency key returns same booking; My Bookings shows entry; second confirm same showtime creates second booking (multiple allowed).

## Confirm rules (`CONTEXT.md`)

| Rule | Behavior |
| ---- | -------- |
| Scope | All seats in `user_holds:{userId}:{showtimeId}` |
| Partial | Not supported — deselect before confirm |
| Status | `CONFIRMED` only in `bookings` |
| `bookingRef` | `TBS-` + 6–8 unambiguous alphanumeric |
| `ticketToken` | Opaque secret ≠ `bookingRef` |
| Idempotency | Same key + prior success → 200 same body; no holds + no prior → **409** |
| Locks | `lock:confirm:{showtimeId}:{seatId}` sorted order |
| Total | Sum `priceTiers[seat.type]` for each seat |
| Limit | Max 10 seats |
| Cutoff | `startsAt <= now` → reject |

## API Surface

| Method | Path | Headers | Body |
| ------ | ---- | ------- | ---- |
| POST | `/api/bookings/confirm` | `Idempotency-Key` (UUID) | `{ showtimeId }` |
| GET | `/api/bookings/mine` | Cookie | `?upcoming=true\|false` optional |
| GET | `/api/bookings/:id` | Cookie | Owner or admin only |

### Confirm response

```json
{
  "id": "...",
  "bookingRef": "TBS-7K2M9P",
  "showtimeId": "...",
  "seats": ["A-1", "A-2"],
  "total": 2400,
  "status": "CONFIRMED",
  "confirmedAt": "2026-06-11T19:00:00Z"
}
```

Enqueue email task (spec 09) — do not block HTTP on SendGrid.

## Vue Structure

```
app/src/views/
├── CheckoutView.vue         # order summary + confirm CTA
├── BookingConfirmationView.vue
├── MyBookingsView.vue
└── BookingDetailView.vue
app/src/stores/bookingSession.ts  # idempotency key per checkout attempt
```

### My Bookings UX

- Tabs or filter: **Upcoming** vs **History** (by showtime `startsAt` vs now)
- Show: movie title, cinema, screen, showtime, seats, `bookingRef`, status badge **Confirmed**
- Link to ticket view (spec 09)
- Do **not** show active Redis holds

## Code Style

Confirm service outline:

```go
func (s *BookingService) Confirm(ctx context.Context, userID, showtimeID, idempotencyKey string) (*Booking, error) {
    if cached, ok := s.idempotency.Get(idempotencyKey); ok { return cached, nil }
    seatIDs, err := s.holds.AllForUserShowtime(ctx, userID, showtimeID)
    if len(seatIDs) == 0 { return nil, ErrNoActiveHolds } // 409
    // acquire locks in sorted order, insert booking, clear holds, cache idempotency
}
```

Client generates fresh `Idempotency-Key` per checkout attempt; reuse only on retry of same attempt.

## Testing Strategy (**tdd**)

- Go: idempotency hit/miss; 409 when holds expired; double confirm same seats blocked.
- Go: concurrent confirm two users one seat — one wins, one 409.
- Go: `bookingRef` format; total price calculation.
- Go: multiple bookings same user same showtime allowed.
- Vue: checkout sends `Idempotency-Key` header.

## Boundaries

- **Always:** Delete Redis holds after successful confirm; broadcast `seat_sold` via WS.
- **Ask first:** Changing idempotency TTL (default ~24h).
- **Never:** Create booking without clearing holds; partial confirm body.

## Tasks

- [ ] `internal/booking` Confirm with locks + idempotency cache (Redis)
- [ ] `bookingRef` + `ticketToken` generators
- [ ] POST confirm, GET mine, GET by id handlers
- [ ] Checkout + confirmation + My Bookings views
- [ ] Wire confirm → enqueue `email:send` task (handler stub OK until spec 09)
- [ ] WS `seat_sold` events on confirm

## Out of scope

- QR image generation (spec 09)
- Cancellation
- Payment

## Check when done

- [ ] Confirm books full hold set; holds cleared
- [ ] Idempotency and 409 paths tested
- [ ] My Bookings shows only CONFIRMED bookings
- [ ] Multiple bookings per showtime work
- [ ] `go test ./internal/booking/...` passes
- [ ] `progress-tracker.md` updated when implementation lands

## Open Questions

- None.
