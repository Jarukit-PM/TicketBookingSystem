# Progress Tracker

Update this file whenever the current phase, active feature, or implementation state changes.

## Current Phase

- **Foundation** — Cinema Ticket Booking System (context pivot from prior NovelCraft project)

## Current Goal

- Implement feature spec 01 (design system), then follow specs 02–10 in order for infrastructure through admin.

## Completed

- **GitHub Actions CI (2026-06-11):** Issue #2 — `.github/workflows/ci.yml` (Go `vet`/`test` in `api/`, Vue `lint`/`type-check`/`build` in `app/`); minimal `api/go.mod` + stub test; ESLint override for UI primitive names.
- **MongoDB data model (2026-06-11):** Issue #5 — domain models (users, movies, cinemas, screens, showtimes, bookings, audit_logs, email_logs), `db.EnsureIndexes` on server boot, repository interfaces + mongo implementations, `booking.GenerateBookingRef` (`TBS-` format) with table-driven tests, `catalog.TotalForSeats` pricing helper, `cmd/seed` (1 cinema, 2 screens, 2 movies, 5 showtimes); `go test ./...` passes.
- **Go API scaffold (2026-06-11):** Issue #4 — Gin server + asynq worker (`cmd/server`, `cmd/worker`), Viper config, MongoDB/Redis connect + ping, `GET /api/health` (optional `?deep=1`), middleware (request ID, recovery), `pkg/httputil` JSON errors; `go test ./...` and `go vet` pass.
- **Docker + nginx (2026-06-11):** Issue #3 — `docker-compose.yml` (nginx, app, api, worker, mongo, redis), `nginx/nginx.conf` (SPA + `/api` + `/ws` proxy), `.env.example`, minimal Go api/worker stubs + Dockerfiles; `docker compose config` validates.
- **Mongo GUI access (2026-06-11):** Exposed `mongo` on `localhost:27017` in `docker-compose.yml` for Compass; local Homebrew `mongodb-community` stopped to free the port.
- **Design system (2026-06-11):** Feature spec 01 — `tokens.css` + Tailwind v4 `main.css`, `cn()` helper, UI primitives (`Button`, `Card`, `Input`, `Badge`), dark cinema preview in `HomeView.vue`; Vite green starter theme removed.
- **Context pivot (2026-06-10):** `project-overview.md` and `architecture-context.md` rewritten for cinema ticket booking (Vue SPA, Go/Gin API, MongoDB, Redis holds, WebSocket seat map, SendGrid email, asynq worker).
- **Context pivot (2026-06-10):** `ui-context.md`, `code-standards.md`, `ai-workflow-rules.md`, `current-issue.md`, and this file updated for the new domain.
- **UI theme (2026-06-10):** Black + gradient orange design language in `ui-context.md`.
- **Decisions (2026-06-10):** Multi-cinema MVP; sold seats derived from confirmed `bookings`; admin QR scan → user booking history; feature spec 01 rewritten; NovelCraft specs 02–10 removed.
- **Vue 3 starter:** Vite + TypeScript + Vue Router + Pinia in `app/` (default scaffold — not yet themed or wired to API).
- **Agent skills:** Project skills installed (`.agents/skills/`, `skills-lock.json`); mapping in `AGENTS.md`.
- **Feature specs (2026-06-11):** `context/CONTEXT.md` glossary + grill decisions; feature specs **02–10** authored (`spec-driven-development`).
- **Implementation issues (2026-06-11):** Specs **05–10** broken into 14 vertical-slice GitHub issues **#11–#24** (`to-issues`); label `ready-for-human` added for HITL slices (#13 OAuth, #24 QR scan).
- **Email/password auth (2026-06-11):** Issue #11 on branch `issue-11-auth` (`ce51e91`) — full auth slice (API + SPA guards); Google OAuth deferred to #13.
- **Public catalog browse (2026-06-11):** Issue #12 on branch `issue-12-catalog` (`39588a6`) — public catalog API + customer UI (cinema picker, movie browse, showtime list); pushed (replaced mistaken admin CRUD commit on remote).
- **Read-only seat map (2026-06-11):** Issue #15 on branch `issue-15-seat-map` (`f907bbe`) — `GET /api/showtimes/:id/seats` inventory snapshot (`AVAILABLE`/`HELD`/`SOLD`/`BLOCKED`), `internal/inventory` + Redis hold reader, `SeatMapView` with `SeatMapGrid`/`SeatCell`/`SeatLegend`, public route `/book/:showtimeId`; Go table tests for inventory computation.
- **Redis seat holds (2026-06-11):** Issue #16 on branch `issue-16-redis-holds` (`db7a2df`) — `POST/DELETE /api/showtimes/:id/holds` (auth required), Redis keys `hold:{showtimeId}:{seatId}` + `user_holds:{userId}:{showtimeId}` (5m TTL, refresh on add only), rejects sold/blocked/other-hold/post-cutoff/>10 seats; `expiresAt` in response; `go test ./internal/hold/...` passes.
- **Interactive seat map + WebSocket (2026-06-11):** Issue #17 on branch `issue-17-seat-map-ws` — `GET /ws/showtimes/:id` hub with Redis pub/sub fan-out; `seat_held`/`seat_released` after HTTP mutations; `useShowtimeSocket`, `useHoldCountdown`, `bookingSession` Pinia store; interactive `SeatMapView` (select/deselect, gradient self-held, countdown urgency &lt;60s); guest seat click → login redirect; checkout stub `/book/:showtimeId/checkout`; `go test ./...` + `npm run type-check` + `npm run build` pass.
- **Booking confirm (2026-06-11):** Issue #18 on branch `issue-18-booking-confirm` — `POST /api/bookings/confirm` with `Idempotency-Key`, sorted `lock:confirm:{showtimeId}:{seatId}` acquisition, `bookingRef` + `ticketToken`, Redis idempotency cache (~24h), holds cleared on success, `seat_sold` WS events, asynq `email:send` enqueue + worker stub; `CheckoutView` confirm flow + `BookingConfirmationView`; `go test ./internal/booking/...` + `go test ./...` + `npm run build` pass.
- **My Bookings (2026-06-11):** Issue #19 on branch `issue-19-my-bookings` — `GET /api/bookings/mine?upcoming=true|false` and `GET /api/bookings/:id` (owner or admin), confirmed bookings only with showtime enrichment; `MyBookingsView` (upcoming/history tabs), `BookingDetailView`, `BookingCard`, ticket placeholder route; `go test ./...` + `npm run type-check` + `npm run build` pass.
- **Admin booking search (2026-06-11):** Issue #22 on branch `issue-22-admin-booking-search` — `GET /api/admin/bookings` (email, bookingRef, userId, showtimeId filters) and `GET /api/admin/users/:userId/bookings`; minimal admin shell with `AdminBookingsView`, `AdminUserBookingsView`, `BookingsTable`; customer JWT 403 on admin routes; `go test ./...` + `npm run type-check` pass.

## In Progress

- None — ready for #24 (Admin QR scan).

## Next Up

Specs 05–10 broken into **14 vertical-slice issues** (GitHub #11–#24). HITL: #13 (Google OAuth), #24 (QR scan).

| Order | Issue | Slice | Spec |
| ----- | ----- | ----- | ---- |
| ~~1~~ | [#11](https://github.com/Jarukit-PM/TicketBookingSystem/issues/11) | Email/password auth + middleware + route guards ✅ `issue-11-auth` | 05 |
| 2 | ~~[#12](https://github.com/Jarukit-PM/TicketBookingSystem/issues/12)~~ | ~~Public catalog browse~~ ✅ `issue-12-catalog` (`39588a6`) | 06 |
| 3 | [#13](https://github.com/Jarukit-PM/TicketBookingSystem/issues/13) | Google OAuth (HITL) | 05 |
| 4 | ~~[#14](https://github.com/Jarukit-PM/TicketBookingSystem/issues/14)~~ | ~~Admin catalog CRUD~~ ✅ `issue-14-admin-catalog` (`c51aec3`) | 06 |
| 5 | ~~[#15](https://github.com/Jarukit-PM/TicketBookingSystem/issues/15)~~ | ~~Read-only seat map~~ ✅ `issue-15-seat-map` (`f907bbe`) | 07 |
| 6 | ~~[#16](https://github.com/Jarukit-PM/TicketBookingSystem/issues/16)~~ | ~~Redis seat holds API~~ ✅ `issue-16-redis-holds` (`db7a2df`) | 07 |
| 7 | ~~[#17](https://github.com/Jarukit-PM/TicketBookingSystem/issues/17)~~ | ~~Interactive seat map + WebSocket~~ ✅ `issue-17-seat-map-ws` | 07 |
| 8 | ~~[#18](https://github.com/Jarukit-PM/TicketBookingSystem/issues/18)~~ | ~~Booking confirm~~ ✅ `issue-18-booking-confirm` | 08 |
| 9 | ~~[#19](https://github.com/Jarukit-PM/TicketBookingSystem/issues/19)~~ | ~~My Bookings~~ ✅ `issue-19-my-bookings` | 08 |
| 10 | [#20](https://github.com/Jarukit-PM/TicketBookingSystem/issues/20) | Digital ticket + confirmation email | 09 |
| 11 | [#21](https://github.com/Jarukit-PM/TicketBookingSystem/issues/21) | Admin shell + dashboard | 10 |
| 12 | ~~[#22](https://github.com/Jarukit-PM/TicketBookingSystem/issues/22)~~ | ~~Admin booking search~~ ✅ `issue-22-admin-booking-search` | 10 |
| 13 | [#23](https://github.com/Jarukit-PM/TicketBookingSystem/issues/23) | Admin audit + email logs | 10 |
| 14 | [#24](https://github.com/Jarukit-PM/TicketBookingSystem/issues/24) | Admin QR scan (HITL) | 10 |

**Start immediately (no blockers):** #20 (Digital ticket + confirmation email).

**Local stack:** `cp .env.example .env && docker compose up --build` → SPA at `http://localhost`, `/api/health` via nginx proxy.

## Resolved Decisions

- **Booking vs seat hold:** A **booking** exists only after confirm (MongoDB, `CONFIRMED`). Redis holds are **seat holds**, not bookings — My Bookings and admin search reflect confirmed bookings only. Glossary in `context/CONTEXT.md`.
- **Multiple bookings per showtime:** A user may confirm multiple separate bookings for the same showtime (each with its own `bookingRef`, ticket, and email).
- **Seat deselect:** Deselecting a seat during checkout releases it immediately; TTL refreshes on add only, not on remove.
- **Navigate away:** Holds survive tab close / route change until Redis TTL expires; explicit `DELETE` abandon or confirm clears them. No release on WebSocket disconnect.
- **Movie catalog:** Global `movies` collection — one film document shared across cinemas; scheduling and browse-by-cinema via showtimes on screens.
- **Now showing browse:** Per cinema, Now Showing = ≥1 future showtime (`status != ARCHIVED`). Coming Soon tab = `COMING_SOON` teasers even without showtimes.
- **Blocked seats:** Screen layout only (`type: blocked` on seat). No per-showtime block overrides in MVP.
- **Pricing:** Showtime `priceTiers` maps layout seat `type` → price; booking total = sum per selected seat.
- **Concurrent holds:** User may hold seats on multiple showtimes at once; independent TTL per showtime.
- **Confirm:** Books entire hold set for one showtime; deselect to reduce before confirm.
- **Idempotency:** Same key returns cached booking on success; expired holds on failed retry → `409`, new key for fresh checkout.
- **Admin scope:** Global — any admin manages all cinemas in MVP; no per-cinema RBAC.
- **Auth session:** httpOnly JWT cookie only for SPA (same origin via nginx); no Bearer in client for MVP.
- **JWT expiry:** **7 days** (`JWT_EXPIRY`, default `168h`); no refresh tokens in MVP.
- **OAuth link:** Google sign-in auto-links to existing user when verified email matches.
- **Showtime cutoff:** No hold/confirm after `startsAt`; cinema timezone for comparison.
- **Seat limit:** Max 10 seats per hold/booking per user per showtime.
- **Booking ref:** `TBS-` + short alphanumeric (no ambiguous chars); separate from `ticketToken`.
- **Guest browse:** Anonymous seat map OK; auth required at hold/confirm only.
- **Showtime seat inventory (option A):** Derive sold seats from confirmed `bookings` only. No `showtimes.soldSeatIds[]`. **No cancellation in MVP** — sold seats never return to available, so inventory is append-only.
- **Multi-cinema:** **In scope for MVP.** Multiple `cinemas` documents; customer browses/filters by cinema; admin manages venues, screens, and showtimes per cinema.
- **Admin QR scan:** **Scan only → open that customer's booking history** in admin UI (`/admin/users/:userId/bookings`). No door check-in validation, showtime window gate, or pass/fail scan result screen.

## Architecture Decisions

See `context/architecture-context.md`. Summary:

- **Domain:** Cinema → Screen → Showtime → Booking; seat holds in **Redis** (5 min TTL, refresh on each new seat).
- **Sold seats:** Computed from confirmed bookings (append-only in MVP — no cancellation).
- **Stack:** Vue 3 SPA (`app/`), Go Gin API (`api/`), MongoDB durable data, Redis holds/locks/asynq, WebSocket per showtime.
- **Auth:** JWT httpOnly cookie only (MVP); roles Customer / Admin.
- **Real-time:** WebSocket advisory; HTTP hold/confirm authoritative.
- **Email:** SendGrid via asynq worker; failures do not roll back confirmed bookings.
- **Payment:** out of scope for MVP — confirm-only bookings.
- **Invariant:** No double booking — Redis locks + single MongoDB write on confirm.

## Session Notes

- **2026-06-11:** Issue #2 CI merged to `main` — workflow + minimal `api/` module before full API scaffold.
- **2026-06-11:** Issue #5 — MongoDB data model on branch `issue-5-data-model` (models, indexes, repos, seed, booking ref generator).
- **2026-06-11:** Issue #4 — Go API scaffold on branch `issue-4-api-scaffold` (Gin, Viper, Mongo/Redis, asynq worker stub, health route).
- **2026-06-11:** Issue #3 — Docker Compose + nginx on branch `issue-3-docker-nginx` (six services, SPA volume, api/worker stubs).
- **2026-06-11:** Feature spec 01 design system implemented on branch `issue-1-design-system` (tokens, Tailwind v4, UI primitives, preview page).
- **2026-06-11:** Grill-with-docs session — 19 domain decisions captured in `CONTEXT.md`; specs 02–10 written per `spec-driven-development`.
- **2026-06-11:** Auth JWT session lifetime confirmed — **7 days** (`JWT_EXPIRY`, default `168h`); spec 05 open question closed.
- **2026-06-10:** Seat inventory option A confirmed; **no booking cancellation in MVP** (sold seats never released).
- **2026-06-10:** User migrated product context from NovelCraft to Cinema Ticket Booking System. `api/` scaffolded per spec 03.
- **2026-06-11:** Branch hygiene — split catalog (#12), seat map (#15), and Redis holds (#16) onto correct branches; `issue-12-catalog` force-pushed to remove mistaken admin CRUD commit from remote.
- **2026-06-11:** Issue #16 — Redis seat holds on branch `issue-16-redis-holds` (`db7a2df`).
- **2026-06-11:** Issue #15 — read-only seat map on branch `issue-15-seat-map` (`f907bbe`).
- **2026-06-11:** Issue #17 — interactive seat map + WebSocket on branch `issue-17-seat-map-ws`.
- **2026-06-11:** Issue #18 — booking confirm on branch `issue-18-booking-confirm`.
