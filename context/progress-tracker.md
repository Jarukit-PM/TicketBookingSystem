# Progress Tracker

Update this file whenever the current phase, active feature, or implementation state changes.

## Current Phase

- **Foundation** — Cinema Ticket Booking System (context pivot from prior NovelCraft project)

## Current Goal

- Implement feature spec 01 (design system), then scaffold Go API and vertical-slice features.

## Completed

- **Context pivot (2026-06-10):** `project-overview.md` and `architecture-context.md` rewritten for cinema ticket booking (Vue SPA, Go/Gin API, MongoDB, Redis holds, WebSocket seat map, SendGrid email, asynq worker).
- **Context pivot (2026-06-10):** `ui-context.md`, `code-standards.md`, `ai-workflow-rules.md`, `current-issue.md`, and this file updated for the new domain.
- **UI theme (2026-06-10):** Black + gradient orange design language in `ui-context.md`.
- **Decisions (2026-06-10):** Multi-cinema MVP; sold seats derived from confirmed `bookings`; admin QR scan → user booking history; feature spec 01 rewritten; NovelCraft specs 02–10 removed.
- **Vue 3 starter:** Vite + TypeScript + Vue Router + Pinia in `app/` (default scaffold — not yet themed or wired to API).
- **Agent skills:** Project skills installed (`.agents/skills/`, `skills-lock.json`); mapping in `AGENTS.md`.

## In Progress

- None — start with [`feature-specs/01-design-system.md`](feature-specs/01-design-system.md).

## Next Up

| Order | Feature | Notes |
| ----- | ------- | ----- |
| 1 | Design system + Tailwind | [`feature-specs/01-design-system.md`](feature-specs/01-design-system.md) — black + gradient orange tokens |
| 2 | Docker + nginx + CI skeleton | `docker-compose.yml`, `.github/workflows/ci.yml` |
| 3 | Go API scaffold | `api/cmd/server`, Viper config, health route |
| 4 | MongoDB models + indexes | Users, movies, **cinemas**, screens, showtimes, bookings |
| 5 | Authentication | Email/password + Google OAuth, JWT, Vue Router guards |
| 6 | Catalog (movies + showtimes) | Multi-cinema browse + admin CRUD |
| 7 | Seat map + Redis holds + WebSocket | 5-minute TTL, countdown, real-time sync |
| 8 | Booking confirm + My Bookings | Idempotent confirm; sold seats from bookings query |
| 9 | QR ticket + email | go-qrcode, SendGrid via asynq worker |
| 10 | Admin dashboard + QR scan | Scan QR → navigate to user's booking history |

Author new specs (02+) with `spec-driven-development` before coding each row.

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
- **Auth:** JWT (httpOnly cookie or Bearer); roles Customer / Admin.
- **Real-time:** WebSocket advisory; HTTP hold/confirm authoritative.
- **Email:** SendGrid via asynq worker; failures do not roll back confirmed bookings.
- **Payment:** out of scope for MVP — confirm-only bookings.
- **Invariant:** No double booking — Redis locks + single MongoDB write on confirm.

## Session Notes

- **2026-06-10:** Seat inventory option A confirmed; **no booking cancellation in MVP** (sold seats never released).
- **2026-06-10:** User migrated product context from NovelCraft to Cinema Ticket Booking System. Codebase is still Vue Vite starter + no `api/` yet.
