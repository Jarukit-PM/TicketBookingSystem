# Feature 09 — Digital Ticket (QR) & Confirmation Email

Read `context/CONTEXT.md`, `context/architecture-context.md`, and `context/ui-context.md`.

Generate **QR digital tickets** per booking, in-app ticket view, and **SendGrid** confirmation email via **asynq** worker. Email failure does not roll back booking.

**Depends on:** Features 03 (worker), 04 (`email_logs`), 08 (confirmed bookings).

## Objective

Customer receives scannable ticket in app and email. Admin QR resolves to user booking history (spec 10). Worker retries failed sends.

**Success looks like:** After confirm, email enqueued → worker sends HTML+text with ticket link/QR; `GET /api/bookings/:id/ticket` returns QR for owner; `email_logs` records status.

## Ticket & QR

| Field | Rule |
| ----- | ---- |
| `bookingRef` | Human-facing `TBS-XXXXXX` |
| `ticketToken` | HMAC/signed opaque secret |
| QR payload | URL `https://{APP_URL}/ticket/{bookingRef}?t={ticketToken}` |
| One QR | Per booking — lists all seats in that booking |

### API

| Method | Path | Auth | Response |
| ------ | ---- | ---- | -------- |
| GET | `/api/bookings/:id/ticket` | Owner | `{ bookingRef, qrPngBase64?, ticketUrl, seats, showtime... }` |
| GET | `/api/admin/tickets/resolve` | Admin | `?ref=&t=` → `{ userId, bookingId }` |

Use `github.com/skip2/go-qrcode` server-side.

## Email (SendGrid + asynq)

| Item | Detail |
| ---- | ------ |
| Trigger | Confirm handler enqueues `email:send` with `bookingId` |
| Worker | `cmd/worker` registers handler; retries with asynq backoff |
| Templates | Go `html/template` + plain text |
| Content | Movie, cinema, screen, showtime, seats, total, `bookingRef`, ticket URL, inline QR optional |
| Log | `email_logs`: `bookingId`, `type: CONFIRMATION`, `status`, SendGrid message id |
| Failure | Log + retry; **do not** delete booking |

Env: `SENDGRID_API_KEY`, `EMAIL_FROM`.

## Vue Structure

```
app/src/views/TicketView.vue    # white QR pad per ui-context.md
app/src/components/TicketCard.vue
```

### Ticket UI

- Dark page; QR on **white pad** for scan contrast (`ui-context.md`)
- Show movie, cinema, screen, time, seats, `bookingRef`
- Print-friendly CSS optional

## Code Style

Task handler:

```go
func HandleEmailSend(ctx context.Context, t *asynq.Task) error {
    var p EmailPayload
    if err := json.Unmarshal(t.Payload(), &p); err != nil { return err }
    // load booking, render templates, SendGrid API, upsert email_log
}
```

Token validation for admin resolve:

```go
func ValidateTicketToken(ref, token string, booking *Booking) bool
```

## Testing Strategy

- Go: QR generation smoke test; token validate/reject tampered token.
- Go: email handler with mocked SendGrid client.
- Go: confirm still succeeds if enqueue fails (log error) — email is async.
- Manual: Mailtrap or SendGrid sandbox for HTML preview.

## Boundaries

- **Always:** Enqueue email in confirm path, not inline HTTP to SendGrid.
- **Ask first:** PDF attachment, SMS.
- **Never:** Roll back booking on email failure; log full API keys.

## Tasks

- [ ] `ticketToken` sign/verify helpers
- [ ] QR generation endpoint + TicketView
- [ ] Email templates (HTML + text)
- [ ] asynq task `email:send` + worker handler
- [ ] `email_logs` repository writes
- [ ] Wire confirm → enqueue task
- [ ] Admin resolve endpoint (navigation UI in spec 10)

## Out of scope

- Cancellation email
- Door check-in scan validation UI
- PDF ticket download (optional nice-to-have)

## Check when done

- [ ] Owner can open ticket with QR from My Bookings
- [ ] Confirmation email received in dev/sandbox
- [ ] `email_logs` row per send attempt
- [ ] Invalid `ticketToken` rejected on admin resolve
- [ ] `go test` for token + email handler
- [ ] `progress-tracker.md` updated when implementation lands

## Open Questions

- None.
