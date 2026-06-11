# Progress Tracker

Update this file whenever the current phase, active feature, or implementation state changes.

## Current Phase

- **Foundation** — Cinema Ticket Booking System (context pivot from prior NovelCraft project)

## Current Goal

- Implement feature spec 01 (design system), then follow specs 02–10 in order for infrastructure through admin.

## Completed

- **Design system (2026-06-11):** Feature spec 01 — `tokens.css` + Tailwind v4 `main.css`, `cn()` helper, UI primitives (`Button`, `Card`, `Input`, `Badge`), dark cinema preview in `HomeView.vue`; Vite green starter theme removed.
- **Context pivot (2026-06-10):** `project-overview.md` and `architecture-context.md` rewritten for cinema ticket booking (Vue SPA, Go/Gin API, MongoDB, Redis holds, WebSocket seat map, SendGrid email, asynq worker).
- **Context pivot (2026-06-10):** `ui-context.md`, `code-standards.md`, `ai-workflow-rules.md`, `current-issue.md`, and this file updated for the new domain.
- **UI theme (2026-06-10):** Black + gradient orange design language in `ui-context.md`.
- **Decisions (2026-06-10):** Multi-cinema MVP; sold seats derived from confirmed `bookings`; admin QR scan → user booking history; feature spec 01 rewritten; NovelCraft specs 02–10 removed.
- **Vue 3 starter:** Vite + TypeScript + Vue Router + Pinia in `app/` (default scaffold — not yet themed or wired to API).
- **Agent skills:** Project skills installed (`.agents/skills/`, `skills-lock.json`); mapping in `AGENTS.md`.
- **Feature specs (2026-06-11):** `context/CONTEXT.md` glossary + grill decisions; feature specs **02–10** authored (`spec-driven-development`).

## In Progress

- None — start with [`feature-specs/02-infrastructure.md`](feature-specs/02-infrastructure.md).

## Next Up

| Order | Feature | Spec |
| ----- | ------- | ---- |
| 1 | Docker + nginx + CI | [`02-infrastructure.md`](feature-specs/02-infrastructure.md) |
| 2 | Go API scaffold | [`03-api-scaffold.md`](feature-specs/03-api-scaffold.md) |
| 3 | MongoDB models + indexes | [`04-data-model.md`](feature-specs/04-data-model.md) |
| 4 | Authentication | [`05-authentication.md`](feature-specs/05-authentication.md) |
| 5 | Catalog (browse + admin CRUD) | [`06-catalog.md`](feature-specs/06-catalog.md) |
| 6 | Seat map + holds + WebSocket | [`07-seat-map-and-holds.md`](feature-specs/07-seat-map-and-holds.md) |
| 7 | Booking confirm + My Bookings | [`08-booking.md`](feature-specs/08-booking.md) |
| 8 | QR ticket + email | [`09-ticket-and-email.md`](feature-specs/09-ticket-and-email.md) |
| 9 | Admin dashboard + QR scan | [`10-admin.md`](feature-specs/10-admin.md) |

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
- **Auth:** JWT httpOnly cookie only (MVP); roles Customer / Admin.
- **Real-time:** WebSocket advisory; HTTP hold/confirm authoritative.
- **Email:** SendGrid via asynq worker; failures do not roll back confirmed bookings.
- **Payment:** out of scope for MVP — confirm-only bookings.
- **Invariant:** No double booking — Redis locks + single MongoDB write on confirm.

## Session Notes

- **2026-06-11:** Feature spec 01 design system implemented on branch `issue-1-design-system` (tokens, Tailwind v4, UI primitives, preview page).
- **2026-06-11:** Grill-with-docs session — 19 domain decisions captured in `CONTEXT.md`; specs 02–10 written per `spec-driven-development`.
- **2026-06-10:** Seat inventory option A confirmed; **no booking cancellation in MVP** (sold seats never released).
- **2026-06-10:** User migrated product context from NovelCraft to Cinema Ticket Booking System. Codebase is still Vue Vite starter + no `api/` yet.
