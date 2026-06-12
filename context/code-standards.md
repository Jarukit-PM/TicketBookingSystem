# Code Standards

## General

- Keep modules small and single-purpose.
- Fix root causes — do not layer workarounds.
- Do not mix unrelated concerns in one handler, service, or component.
- Respect the system boundaries defined in `architecture-context.md`.

## TypeScript (Vue frontend — `app/`)

- Strict mode is required (`vue-tsc`).
- Avoid `any`; use explicit interfaces or narrowly scoped types.
- Validate unknown external input at API boundaries (fetch responses) before trusting it.
- Use `interface` for object contracts shared across components and stores.
- Prefer **Composition API** with `<script setup lang="ts">` unless the project explicitly requires otherwise.

## Vue

- Route-level views in `app/src/views/`; reusable UI in `app/src/components/`.
- **Pinia** for client session state (auth user, booking draft, seat selection); server state is source of truth for inventory.
- **Composables** in `app/src/composables/` for WebSocket (`useShowtimeSocket`), hold countdown, and shared async logic — use `MaybeRef` / `toValue()` for adaptable inputs.
- Vue Router `beforeEach` guards for `/book/*`, `/my-bookings`, `/admin/*`.
- Do not call hold/confirm APIs from watchers without debouncing or explicit user action — avoid hold storms.

## Go (API — `api/`)

- Follow standard layout: `cmd/server`, `cmd/worker`, `internal/<domain>`, `pkg/`.
- Handlers stay thin: parse/validate → call service → map errors to HTTP status.
- Business logic in `internal/` packages (`booking`, `hold`, `catalog`, `auth`, etc.), not in route handlers.
- Return typed errors from services; map to consistent JSON error bodies at the handler layer.
- Use `context.Context` on all I/O; respect cancellation and timeouts.
- Seat confirm and hold paths must use Redis locks and documented ordering (sorted `seatId`) to avoid deadlock.
- Long-running work (email send) goes to **asynq** tasks in `cmd/worker`, not in HTTP handlers.

## Styling (Vue)

- Dark cinema UI per `ui-context.md` — black surfaces, light text, gradient orange brand (`bg-gradient-brand`). No light-mode alternate theme for MVP.
- Use CSS custom property tokens — no raw Tailwind palette classes like `zinc-*`, `orange-500`, or hardcoded hex in components.
- Tailwind utilities via design tokens: `bg-base`, `bg-surface`, `text-copy-primary`, `border-surface-border`, `bg-gradient-brand`, `text-brand`, etc.
- Primary CTAs use gradient orange pills; cards use `bg-surface` + subtle border; QR ticket view may use a white pad for scan contrast only.

## HTTP API (Go)

- Validate request DTOs at the boundary (manual checks or `go-playground/validator` when added).
- JWT middleware on protected routes; `RequireAdmin` on `/api/admin/*`.
- Consistent JSON response shapes; use appropriate status codes (`409` for seat conflict, `402` not used — no payment in MVP).
- **Idempotency-Key** header on `POST /api/bookings/confirm`.
- WebSocket events are advisory; REST hold/confirm responses are authoritative.

## Data and Storage

- **MongoDB:** durable entities (users, movies, cinemas, screens, showtimes, bookings, audit_logs, email_logs). Holds are **not** stored in MongoDB.
- **Redis:** seat holds (TTL 5 minutes), confirm locks, idempotency cache, rate limits, asynq queue, WebSocket pub/sub across API instances.
- Do not store passwords in plain text — bcrypt for email/password auth.
- Do not store email API payloads with secrets in logs — log booking id and provider message id only.

## Real-time

- Client composable reconnects and applies `snapshot` on connect before merging events.
- Never update seat map to SOLD/HELD from WebSocket alone without reconciling with last REST hold/confirm response on conflict.

## File Organization

| Path | Responsibility |
| ---- | -------------- |
| `app/src/views/` | Route pages (catalog, seat map, my bookings, admin) |
| `app/src/components/` | Feature and UI components |
| `app/src/composables/` | Shared reactive logic (WebSocket, countdown) |
| `app/src/stores/` | Pinia stores |
| `app/src/lib/` or `app/src/api/` | Typed fetch client for REST API |
| `api/internal/<domain>/` | Domain services and repositories |
| `api/internal/middleware/` | Auth, logging, CORS (dev) |
| `api/internal/tasks/` | asynq task definitions and handlers |
| `api/cmd/server/` | HTTP + WebSocket entrypoint |
| `api/cmd/worker/` | asynq worker entrypoint |

Name files after the responsibility they contain, not the technology.

## Testing

- Go: table-driven tests for hold TTL rules, confirm idempotency, and seat conflict paths (`go test ./...`).
- Vue: unit test composables and pure helpers; Playwright E2E for critical booking happy path (phase 2).
- Prefer TDD for booking invariants — see `tdd` skill.

## Related Context

- `architecture-context.md` — invariants, Redis key patterns, API sketch.
- `ui-context.md` — tokens, seat map UX, component layers.
