# Architecture Context

## Stack


| Layer           | Technology                                                 | Role                                                                            |
| --------------- | ---------------------------------------------------------- | ------------------------------------------------------------------------------- |
| Frontend        | Vue 3 + Vite + TypeScript                                  | Customer SPA and admin UI in `app/`                                             |
| UI              | Tailwind CSS                                               | Layout, seat map, admin tables, responsive styling                              |
| Routing / state | Vue Router + Pinia                                         | Pages, guards, booking session and seat-map client state                        |
| HTTP client     | Fetch or Axios                                             | REST calls to Go API; `credentials: 'include'` for httpOnly JWT cookie (MVP)  |
| Backend         | Go + Gin                                                   | REST API, auth, booking logic, WebSocket hub                                    |
| Config          | [Viper](https://github.com/spf13/viper)                    | Env + config files for API and worker (`config.yaml`, `MONGO_URI`, etc.)        |
| Database        | MongoDB                                                    | Persistent domain data: users, movies, showtimes, bookings, logs                |
| Cache / lock    | Redis                                                      | Seat holds (5 min TTL), distributed locks, rate limits, asynq queue, WS pub/sub |
| Background jobs | [hibiken/asynq](https://github.com/hibiken/asynq)          | Async email send with retries (Redis-backed worker)                             |
| Real-time       | WebSocket (Gin + gorilla/websocket or nhooyr.io/websocket) | Per-showtime seat map updates                                                   |
| Auth            | Google OAuth 2.0 + email/password                          | Sign-in; JWT session; roles Customer / Admin                                    |
| Email           | Brevo                                                      | Booking confirmation email only in MVP (HTML + plain-text)                      |
| QR              | `github.com/skip2/go-qrcode`                               | Ticket QR generation server-side                                                |
| Reverse proxy   | nginx                                                      | Single origin: SPA static, `/api` → Gin, `/ws` upgrade                          |
| Deployment      | Docker + `docker-compose.yml`                              | `nginx`, `app`, `api`, `worker`, `mongo`, `redis`                               |
| CI              | GitHub Actions                                             | Lint, `go test`, `npm run build` on push/PR                                     |


## Repository Layout

```
TicketBookingSystem/
├── app/                 # Vue 3 + Vite SPA (customer + admin routes)
├── api/                 # Go Gin service (to be scaffolded)
│   ├── cmd/server/      # API entrypoint
│   ├── cmd/worker/      # asynq worker (email, hold-expiry side effects)
│   ├── internal/
│   │   ├── auth/
│   │   ├── booking/
│   │   ├── catalog/     # movies, cinemas, screens, showtimes
│   │   ├── email/       # Brevo client + templates; enqueue via asynq
│   │   ├── tasks/       # asynq task handlers
│   │   ├── hold/        # Redis seat holds
│   │   ├── ws/          # WebSocket hub
│   │   ├── admin/
│   │   └── middleware/
│   └── pkg/             # Shared types, errors
├── docker-compose.yml
├── nginx/               # nginx.conf for reverse proxy
├── .github/workflows/   # CI: Go + Vue build/test
└── context/
```

## Domain Model

MongoDB document hierarchy (logical aggregates):

```
Cinema
 └── Screen (hall)          ← seat layout template
      └── Showtime          ← one screening; seat inventory snapshot for that event
           └── Booking      ← confirmed purchase; embeds seat refs + ticket token
```

Supporting collections: `users`, `movies`, `audit_logs`, `email_logs`.


| Term             | Meaning                                                                                                                                    |
| ---------------- | ------------------------------------------------------------------------------------------------------------------------------------------ |
| **Cinema**       | Venue (name, address, timezone). **Multi-cinema MVP** — many cinema documents; screens and showtimes are cinema-scoped.                   |
| **Movie**        | **Global catalog** — one document per film, not tied to a cinema. Scheduling is per cinema via showtimes. See `CONTEXT.md`.                  |
| **Screen**       | Physical hall with a **seat layout** (rows, seat labels, types: standard, VIP, wheelchair, blocked). Belongs to one cinema.                  |
| **Showtime**     | Movie + screen + `startsAt` + `priceTiers` (map seat layout `type` → price in minor units). Cinema-scoped via screen. Runtime seat status keyed by `(showtimeId, seatId)`. |
| **Price tier**   | Price for a seat `type` on a given showtime (e.g. `standard: 1200`, `vip: 1800`, `wheelchair: 1200`). Booking `total` = sum of tier prices for selected seats. See `CONTEXT.md`. |
| **Seat**         | Identified by `seatId` within a screen layout (e.g. `A-12`). Runtime status per showtime: `AVAILABLE`, `HELD`, `SOLD`, `BLOCKED`. **Blocked** = layout `type: blocked` on the screen (all showtimes); no per-showtime block list in MVP. |
| **Seat hold**    | Short-lived reservation in **Redis** (5 min TTL). Not a booking — invisible in My Bookings and admin search until confirm. See `CONTEXT.md`. |
| **Booking**      | Durable confirmed purchase in MongoDB: user, showtime, seats, total, `bookingRef` (`TBS-` + short alphanumeric), `ticketToken` (opaque QR secret), status `CONFIRMED`. Created only on confirm. |
| **Hold expiry**  | Hold TTL ended without confirm. Seats released; no `bookings` document. Not an "expired booking."                                            |
| **Ticket token** | Opaque value encoded in QR; admin scan resolves booking → navigates to that user's booking history.                                        |
| **Audit log**    | Append-only admin and system actions.                                                                                                      |
| **Email log**    | Send attempts and provider status per booking event.                                                                                       |


### MongoDB Collections (MVP)


| Collection   | Key fields                                                                                       | Notes                                                |
| ------------ | ------------------------------------------------------------------------------------------------ | ---------------------------------------------------- |
| `users`      | `email`, `passwordHash?`, `googleId?`, `role`, `name`, `createdAt`                               | Unique index on `email`; sparse unique on `googleId` |
| `movies`     | `title`, `posterUrl`, `durationMin`, `rating`, `synopsis`, `status`                              | **Global catalog** (no `cinemaId`). `status`: `NOW_SHOWING` \| `COMING_SOON` \| `ARCHIVED`. **Now Showing** tab: ≥1 future showtime at cinema. **Coming Soon** tab: `COMING_SOON` teasers without showtimes OK. Hide `ARCHIVED` |
| `cinemas`    | `name`, `address`, `timezone`                                                                    |                                                      |
| `screens`    | `cinemaId`, `name`, `layout`                                                                     | `layout.seats[]`: `{ seatId, row, col, type }`       |
| `showtimes`  | `movieId`, `screenId`, `startsAt`, `priceTiers`, `status`                                        | `priceTiers`: `{ standard, vip, wheelchair, ... }` in minor units (cents). Index `(screenId, startsAt)` |
| `bookings`   | `userId`, `showtimeId`, `seats[]`, `total`, `bookingRef`, `ticketToken`, `status`, `confirmedAt` | Unique `bookingRef`; index `userId` + `showtimeId` (not unique — multiple bookings per user per showtime allowed) |
| `audit_logs` | `actorId`, `action`, `entity`, `entityId`, `meta`, `createdAt`                                   | Actions: admin `create`/`update`/`delete`; booking `booking_success`, `booking_timeout`, `seat_released`, `booking_failed`, `system_error`. TTL optional (phase 2) |
| `email_logs` | `bookingId`, `type`, `to`, `providerId`, `status`, `createdAt`                                   | MVP `type`: `CONFIRMATION` only                      |


**Seat inventory for a showtime (MVP — option A):** **Derive sold seats from confirmed `bookings`** — query (or aggregate) all `bookings` where `showtimeId` matches and `status = CONFIRMED`, union their `seats[]`. A seat is **SOLD** if it appears in that set. Holds live only in Redis; do **not** duplicate sold seat IDs on `showtimes.soldSeatIds[]`. **No cancellation in MVP** — sold seats are never released back to available, so the bookings query only ever grows. Optional later: Redis cache of sold set per showtime.

```
AVAILABLE = layout seats − SOLD − BLOCKED − (others' Redis HOLDs)
```

## System Boundaries


| Package / service      | Responsibility                                                                |
| ---------------------- | ----------------------------------------------------------------------------- |
| `app/` (Vue)           | UI, Vue Router guards, Pinia stores, WebSocket client, seat map rendering     |
| `api` HTTP routes      | CRUD catalog, hold/confirm/release booking, My Bookings, admin APIs           |
| `api/internal/ws`      | WebSocket upgrade, subscribe by `showtimeId`, broadcast seat events           |
| `api/internal/hold`    | Redis SET with TTL, hold extension rules, release on abandon/timeout          |
| `api/internal/booking` | Confirm all holds for showtime; distributed lock; idempotency; persist booking + mark seats sold |
| `api/internal/email`   | Build Brevo payloads; enqueue asynq tasks; write `email_logs`                   |
| `api/cmd/worker`       | asynq consumer: send email, optional hold-expiry notifications                |
| `api/internal/auth`    | Register/login, Google OAuth callback, JWT issue/validate, role checks        |
| `mongo`                | Source of truth for all durable entities                                      |
| `redis`                | Holds, locks, rate limit counters                                             |


## Auth and Access Model

### Sign-in methods

1. **Email/password** — register with email + bcrypt hash; login returns session.
2. **Google OAuth** — authorization code flow; create user by `googleId`, or **auto-link** to existing user when Google email matches a registered email (verified). One account per email.

### Session

- **JWT** in `httpOnly` cookie only for MVP (`Secure`, `SameSite=Lax`, path `/api`). SPA uses `credentials: 'include'`. No Bearer token in client JS.
- **Lifetime:** 7 days from issue; configurable via `JWT_EXPIRY` (default `168h`). No refresh tokens in MVP.
- Claims: `sub` (user id), `role` (`customer`  `admin`), `exp`.
- Vue Router `beforeEach`: require auth for `/book/`*, `/my-bookings`, `/admin/*`.
- Gin middleware: validate JWT on protected routes; `RequireAdmin` for `/api/admin/*`.

### Roles


| Role         | Access                                                                                          |
| ------------ | ----------------------------------------------------------------------------------------------- |
| **Customer** | Browse, book, My Bookings, own ticket QR                                                        |
| **Admin**    | **Global (MVP):** all cinemas — dashboard, catalog CRUD, booking search (read-only), audit logs, QR scan → user booking history. No `cinemaId` on user. |


First admin: seed via env `ADMIN_EMAIL` on bootstrap or manual DB role update.

## Seat Hold + Redis

### Hold key pattern

```
hold:{showtimeId}:{seatId}  →  { userId, heldAt }   TTL 5 minutes
user_holds:{userId}:{showtimeId}  →  SET of seatIds   TTL 5 minutes
```

### Rules

- User may hold multiple seats on one showtime on the same showtime session.
- **Concurrent showtimes:** user may hold seats on **multiple showtimes** simultaneously (separate Redis keys per `showtimeId`). Confirming one showtime does not release holds on others.
- **TTL refresh (add only):** each time the user adds a seat (`POST .../holds`), reset **5-minute TTL** on **all** seats that user holds for that showtime (`EXPIRE` on every `hold:{showtimeId}:{seatId}` and `user_holds:{userId}:{showtimeId}`). Response and WebSocket `seat_held` include `expiresAt` for the countdown UI.
- **Deselect:** removing a seat (`DELETE .../holds` with `seatIds`, or per-seat delete) releases that seat immediately (`DEL` hold key, `SREM` from `user_holds`). Remaining held seats **keep** their current TTL — no refresh on remove. Broadcast `seat_released` for the deselected seat.
- Cannot hold `SOLD` or `BLOCKED` seats; reject if another user's hold exists (Redis `SET NX`).
- **Showtime cutoff:** reject hold and confirm if `startsAt <= now` (cinema timezone from `cinemas.timezone`). Browse lists only future showtimes.
- **Seat limit:** max **10 seats** per user per showtime per hold/booking; reject `POST .../holds` and confirm if exceeded.
- On TTL expiry: Redis keys vanish; keyspace notification or asynq scheduled sweep triggers WebSocket `seat_released`.
- On confirm: delete hold keys and write `SOLD` in MongoDB inside a lock.
- On **abandon:** client calls `DELETE /api/showtimes/:id/holds` (all or selected seats). Releases immediately.
- On **navigate away** (close tab, route change without abandon): holds remain until **TTL expiry** — no release on WebSocket disconnect alone. Redis TTL is authoritative; worst case seats are unavailable until the hold expires.

### Distributed lock (confirm booking)

```
lock:confirm:{showtimeId}:{seatId}   TTL ~10s
```

Acquire all seat locks in sorted `seatId` order to avoid deadlock, then transactionally create booking and mark seats sold.

## Real-time (WebSocket)

### Connection

- Endpoint: `GET /ws/showtimes/:showtimeId` (anonymous OK for live map; auth required for `POST .../holds` and confirm).
- Client joins room `showtime:{showtimeId}`.

### Server → client events


| Event           | Payload                                          |
| --------------- | ------------------------------------------------ |
| `seat_held`     | `{ seatId, expiresAt }` (no other user's userId) |
| `seat_released` | `{ seatId }`                                     |
| `seat_sold`     | `{ seatId }`                                     |
| `seat_blocked`  | `{ seatId }` (layout change — rare; blocked seats are static in layout for MVP) |
| `snapshot`      | Full map state on connect                        |


### Multi-instance (docker-compose scale)

Use **Redis pub/sub**: API instance publishes seat events; all instances forward to their local WebSocket clients.

## Booking Lifecycle

```
select seats → Redis seat hold (5 min TTL)
     → POST /api/bookings/confirm (showtimeId; idempotency-key header)
     → Books ALL seats in user_holds for that showtime (no seatIds body)
     → Booking CONFIRMED in MongoDB, holds cleared, seats SOLD
     → enqueue asynq task → Brevo CONFIRMATION + email_log
```

**MVP:** no cancel flow — once **CONFIRMED**, a booking and its seats stay sold. Holds that time out are **hold expiries** — not bookings.


| Status / event | Meaning                                                            | In `bookings` collection? |
| -------------- | ------------------------------------------------------------------ | ------------------------- |
| Seat hold      | Active Redis reservation during checkout                           | No                        |
| `CONFIRMED`    | Persisted booking; ticket valid; seats count as SOLD                 | Yes                       |
| Hold expiry    | Hold timed out without confirm (optional audit log only)           | No                        |
| `CANCELLED`    | Voided; seats released                                             | **Phase 2** — not in MVP  |


**Idempotency:** client sends `Idempotency-Key` (UUID). Server stores successful confirm result in Redis (~24h).
- **Retry after success:** same key → `200` with original booking (no duplicate).
- **Retry after failure / no prior success:** requires active holds; if holds expired → `409` — client re-selects seats and uses a **new** key.

## QR / Digital Ticket

- On confirm: generate `bookingRef` (`TBS-` + 6–8 unambiguous alphanumeric) and `ticketToken` (opaque signed/HMAC secret — not the same as `bookingRef`).
- QR encodes URL `https://{app}/ticket/{bookingRef}?t={ticketToken}` or compact JWT.
- `GET /api/bookings/:id/ticket` — returns QR image or payload for authenticated owner.
- **Admin scan (MVP):** no door check-in or pass/fail validation screen. Admin opens scan UI → camera reads QR → `GET /api/admin/tickets/resolve?ref={bookingRef}&t={ticketToken}` returns `{ userId, bookingId }` → client navigates to `**/admin/users/:userId/bookings`** (that customer's full booking history). Invalid or unknown QR shows an error toast; no alternate flow.

## Email (Brevo + asynq)

- Provider: **Brevo** (`BREVO_API_KEY`, `EMAIL_FROM`).
- Templates: Go `html/template` + plain-text fallback; fields per `project-overview.md` (movie, cinema, screen, showtime, seats, total, booking ref, ticket link/QR).
- **Flow:** confirm handler enqueues asynq task (`email:send`); **worker** calls Brevo API with retries (asynq default backoff). No cancellation email in MVP.
- Worker updates `email_logs` with Brevo message id and delivery status; failures remain retryable without rolling back booking.
- QR image for email: generate via `go-qrcode`, attach inline or link to ticket URL.

## Docker Compose (target)

```yaml
services:
  nginx:      # :80 — static app build, proxy /api and /ws to api
  app:        # build stage: Vite → dist volume consumed by nginx
  api:        # Gin :8080; depends on mongo, redis
  worker:     # asynq worker; same image as api, cmd/worker; depends on redis, mongo
  mongo:      # persistent volume
  redis:      # holds, locks, asynq, pub/sub
```

- API binds `0.0.0.0:8080` (internal); public traffic via **nginx** on `:80`.
- WebSocket: nginx proxies `Upgrade` on `/ws` → `api`.
- Config via **Viper**: `config.yaml` + env overrides (`MONGO_URI`, `REDIS_URL`, `JWT_SECRET`, `GOOGLE_*`, `BREVO_*`, `APP_URL`).

## CI (GitHub Actions)

Workflow `.github/workflows/ci.yml` on `push` / `pull_request`:

1. **api** — `go vet`, `go test ./...`, optional `golangci-lint`.
2. **app** — `npm ci`, `npm run lint`, `npm run type-check`, `npm run build`.
3. Optional: Playwright E2E against `docker compose up` (phase 2).

## API Surface (sketch)


| Method | Path                                | Auth                |
| ------ | ----------------------------------- | ------------------- |
| POST   | `/api/auth/register`                | Public              |
| POST   | `/api/auth/login`                   | Public              |
| GET    | `/api/auth/google`                  | Public              |
| GET    | `/api/auth/google/callback`         | Public              |
| GET    | `/api/movies`                       | Public              |
| GET    | `/api/showtimes`                    | Public              |
| GET    | `/api/showtimes/:id/seats`          | Public              |
| POST   | `/api/showtimes/:id/holds`          | Customer            |
| DELETE | `/api/showtimes/:id/holds`          | Customer            |
| POST   | `/api/bookings/confirm`             | Customer            |
| GET    | `/api/bookings/mine`                | Customer            |
| GET    | `/api/bookings/:id`                 | Customer            |
| GET    | `/api/bookings/:id/ticket`          | Customer            |
| GET    | `/api/admin/tickets/resolve`        | Admin               |
| GET    | `/api/admin/users/:userId/bookings` | Admin               |
| `WS`   | `/ws/showtimes/:id`                 | Optional / Customer |
| `*`    | `/api/admin/`*                      | Admin               |


## Invariants

1. A seat for a given showtime cannot be **CONFIRMED** twice; confirm path uses Redis locks + MongoDB write once.
2. **Holds** exist only in Redis with **5 minute TTL** (refreshed on each new seat in the session); MongoDB never stores `HELD` as durable booking state.
3. WebSocket broadcasts are **advisory**; HTTP hold/confirm responses are **authoritative**.
4. Only **Admin** may scan a ticket QR to open that user's booking history. **No booking cancellation in MVP.**
5. Email send failures do not roll back a confirmed booking; log failure in `email_logs` for retry.
6. Payment is **out of scope**; `total` on booking is informational for MVP.

## Optional Later


| Area           | Option                                                                | Why                             |
| -------------- | --------------------------------------------------------------------- | ------------------------------- |
| Logging        | `log/slog` or [uber-go/zap](https://github.com/uber-go/zap)           | Structured logs                 |
| Validation     | [go-playground/validator](https://github.com/go-playground/validator) | Request DTO validation          |
| API docs       | [swaggo/swag](https://github.com/swaggo/swag)                         | OpenAPI from Gin handlers       |
| Migrations     | Index scripts in `api/migrations/`                                    | Reproducible Mongo indexes      |
| Rate limiting  | Gin middleware + Redis sliding window                                 | Protect hold/confirm from abuse |
| CORS           | `gin-contrib/cors`                                                    | Local dev when not using nginx  |
| Frontend WS    | Vue composable `useShowtimeSocket`                                    | Reconnect + snapshot on resume  |
| Object storage | MinIO or S3                                                           | Movie poster uploads            |
| Observability  | Prometheus + `GET /healthz`                                           | Production metrics              |


**Not recommended for MVP:** Kafka, microservices split, GraphQL, second database for reads.

## Related Context

- `project-overview.md` — features, scope, 5-minute hold, email types.
- `feature-specs/` — detailed specs per feature (to be rewritten for cinema domain).
- `code-standards.md` — naming and PR conventions when added for Go/Vue.

