# Feature 12 — Public Ticket Access & Email Link Reliability

Read `context/CONTEXT.md`, `context/architecture-context.md`, `context/ui-context.md`, and **Feature 09** (digital ticket baseline).

**Goal:** Let customers open their ticket from the **email link** or **QR code** without signing in, while keeping ticket data scoped to holders of the signed `ticketToken`. Fix token lifecycle so links in confirmation emails remain valid after deploy/secret rotation.

**Depends on:** Features 08 (confirmed bookings), 09 (ticket tokens, QR, email worker), 11 (locale on booking — display only).

**Status:** Implemented (2026-06-12). HITL: incognito email-link smoke test remains.

## Objective

Confirmation emails include a link `https://{APP_URL}/ticket/{bookingRef}?t={ticketToken}`. Today that link must work when opened on a phone without an active session. Admins scan the same QR for support lookup (spec 10). Tokens must survive HMAC secret rotation because the **persisted** `bookings.ticketToken` is the source of truth for email links.

**Success looks like:** Valid `ref` + `t` returns ticket JSON with QR; invalid or missing token returns 404; owner logged in can still view ticket when token query is absent or stale; worker backfills missing tokens before send; `go test` covers token validation, service, and HTTP handler paths.

## Assumptions

1. **Public read-only** — `GET /api/tickets/:ref` exposes the same ticket fields as the owner endpoint, but only when `?t=` validates. No PII beyond movie/showtime/seats/ref.
2. **No auth required** for public ticket route; owner fallback is **client-side** (Vue) when token fails and session exists.
3. **Token persistence** — `ticketToken` is written to Mongo **after** insert (booking ID known); email worker uses stored value.
4. **Validation order** — HMAC recompute **or** constant-time compare against stored `ticketToken` (legacy links keep working if secret rotates).
5. **QR payload** — same public URL as email link (`TicketURL` helper).
6. **Cancelled bookings** — not in MVP; only `status: confirmed` tickets resolve.

## API Surface

| Method | Path | Auth | Query | Response |
| ------ | ---- | ---- | ----- | -------- |
| GET | `/api/tickets/:ref` | None | `t` (required server-side) | `TicketDetail` — same shape as owner ticket |
| GET | `/api/bookings/:id/ticket` | Owner JWT | — | `TicketDetail` (existing, spec 09) |
| GET | `/api/admin/tickets/resolve` | Admin JWT | `ref`, `t` | `{ userId, bookingId }` (spec 09/10) |

### Error mapping

| Condition | HTTP | Code |
| --------- | ---- | ---- |
| Missing/empty `t`, unknown ref, wrong token, non-confirmed | 404 | `INVALID_TICKET` |
| Catalog lookup failure | 500 | `TICKET_ERROR` |

## Token Lifecycle

```
Confirm booking
  → Insert booking (no token yet)
  → SignTicketToken(secret, ref, bookingId)
  → UpdateTicketToken(bookingId, token)
  → Enqueue email:send

Worker email:send
  → If ticketToken empty → backfill SignTicketToken + UpdateTicketToken
  → Build TicketURL(appURL, ref, storedToken)
  → Brevo send with link + optional inline QR

Customer opens /ticket/:ref?t=...
  → GET /api/tickets/:ref?t=...
  → ValidateTicketToken(ref, t, booking, secret)
  → buildTicketDetail → JSON + qrPngBase64
```

### ValidateTicketToken rules

```go
// 1. ref, token, booking non-nil; bookingRef match; status confirmed
// 2. Accept HMAC(secret, ref, bookingId)
// 3. Else accept booking.TicketToken (persisted email value)
```

## Vue Structure

```
app/src/views/PublicTicketView.vue   # route /ticket/:bookingRef (public)
app/src/views/TicketView.vue         # route /bookings/:bookingId/ticket (auth)
app/src/api/tickets.ts               # fetchPublicTicket, fetchBookingTicket
app/src/router/index.ts              # public + authenticated ticket routes
```

### Public ticket UX

1. Read `bookingRef` from route, `t` from query.
2. If `t` present → `fetchPublicTicket(ref, t)`.
3. On failure → if authenticated, find owned booking by ref → `fetchBookingTicket(id)`.
4. If no `t` and not owner → show `booking.ticket.missingToken` error.
5. Render `TicketCard` with white QR pad (spec 09 / `ui-context.md`).
6. Back link → My Bookings when logged in, Home when guest.

## Code Style

Handler stays thin:

```go
func GetPublicTicket(deps BookingsDeps) gin.HandlerFunc {
    return func(c *gin.Context) {
        ticket, err := deps.Bookings.GetTicketByRef(ctx, c.Param("ref"), c.Query("t"))
        if err != nil { writePublicTicketError(c, err); return }
        httputil.OK(c, ticket)
    }
}
```

Service validates before catalog hydration:

```go
func (s *Service) GetTicketByRef(ctx context.Context, ref, token string) (*TicketDetail, error) {
    if ref == "" || token == "" { return nil, ErrInvalidTicket }
    b, _ := s.bookings.FindByBookingRef(ctx, ref)
    if b == nil || !ValidateTicketToken(ref, token, b, s.ticketSecret) {
        return nil, ErrInvalidTicket
    }
    return s.buildTicketDetail(ctx, b)
}
```

## Commands

```bash
# Backend unit tests (token, service, handler, email backfill)
cd api && go test ./internal/booking/... ./internal/handler/... ./internal/email/... -count=1

# Frontend type-check (PublicTicketView)
cd app && npm run type-check && npm run build

# Manual smoke
cd app && npm run dev
# Confirm booking → open email link in incognito → QR loads
# Tamper ?t= → 404; log in as owner → fallback ticket loads
```

## Project Structure

```
api/internal/booking/
├── token.go           # SignTicketToken, ValidateTicketToken, TicketURL
├── ticket.go          # GetTicket, GetTicketByRef, buildTicketDetail
├── token_test.go
└── ticket_test.go

api/internal/handler/
├── tickets.go         # GetPublicTicket
└── tickets_test.go    # HTTP 200/404 cases

api/internal/email/
├── service.go         # backfill ticketToken before send
└── service_test.go    # TestHandleEmailSendBackfillsMissingToken

app/src/views/PublicTicketView.vue
app/src/api/tickets.ts
```

## Testing Strategy

| Layer | Framework | Cases |
| ----- | --------- | ----- |
| Token helpers | Go `testing` | HMAC valid/tampered; stored token when secret differs; `TicketURL` encoding |
| `GetTicketByRef` | Go + stub repos | Valid token returns QR; tampered → `ErrInvalidTicket` |
| `GetPublicTicket` handler | Go + `httptest` | 200 with body fields; 404 invalid token; 404 missing `t` |
| Email worker | Go + fake sender | Backfill empty `ticketToken`; ticket URL in body |
| `ConfirmedFilter` | Go table tests | Spec 10 admin filters — locale, date range, showtime IDs |
| Vue | Manual / optional Vitest | Parser for route query — low priority |

Run `go test ./...` in `api/` before merge.

## Boundaries

- **Always:** Persist `ticketToken` before enqueueing confirmation email; use constant-time token compare; map invalid ticket to 404 (no enumeration).
- **Ask first:** Rate limiting public ticket endpoint; ticket PDF download.
- **Never:** Return ticket without token validation; expose user email on public route; roll back booking on email failure.

## Success Criteria

- [x] `GET /api/tickets/:ref?t=` returns `TicketDetail` with non-empty `qrPngBase64` for valid link
- [x] Invalid or missing `t` returns 404 `INVALID_TICKET`
- [x] `ValidateTicketToken` accepts persisted `ticketToken` after secret rotation
- [x] Worker backfills empty `ticketToken` before send
- [x] `PublicTicketView` loads ticket from email link without login
- [x] Owner session fallback works when link token is wrong but booking is theirs
- [x] `go test` passes for `token`, `ticket`, `tickets` handler, `ConfirmedFilter`, and email backfill
- [x] `progress-tracker.md` updated when implementation lands

## Tasks

- [x] Persist `ticketToken` after booking insert
- [x] `ValidateTicketToken` accepts stored token
- [x] `GET /api/tickets/:ref` handler + route
- [x] `PublicTicketView` + `fetchPublicTicket`
- [x] Worker token backfill before email
- [x] Handler tests `tickets_test.go`
- [x] Expanded `ConfirmedFilter` unit tests
- [x] Email backfill test `TestHandleEmailSendBackfillsMissingToken`
- [ ] Manual incognito email-link smoke (HITL)

## Out of Scope

- Ticket transfer to another user
- Time-window gate (showtime started/expired)
- Rate limiting / CAPTCHA on public ticket endpoint

## Open Questions

- None.
