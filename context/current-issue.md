# Current Issue

Active gaps and blockers for the Cinema Ticket Booking System.

---

### Go API not scaffolded

**Status:** Open  
**Recorded:** 2026-06-10

#### Summary

`architecture-context.md` defines `api/` layout (`cmd/server`, `cmd/worker`, `internal/*`) but the directory does not exist yet. No Gin server, Viper config, MongoDB/Redis clients, or health endpoint.

First backend slice: scaffold module, `GET /healthz`, Docker Compose services for `mongo` and `redis`.

---

### Frontend not aligned with UI context

**Status:** Open  
**Recorded:** 2026-06-10

#### Summary

`app/` is the default Vue 3 + Vite starter (Inter font, green links, no Tailwind). Black + gradient orange tokens and Tailwind v4 from `ui-context.md` are not applied. No customer or admin routes beyond Home/About stubs.

Design system feature spec should land before booking UI.

---

### Real-time seat map + hold TTL (core risk)

**Status:** Open — design documented, not implemented  
**Recorded:** 2026-06-10

#### Summary

Core product value depends on Redis seat holds (5-minute TTL, refresh on add), WebSocket broadcast, and confirm path with distributed locks. Highest concurrency and correctness risk; implement with TDD on Go booking/hold packages and integration tests before polishing Vue seat map.

See invariants in `architecture-context.md`.

---

### Documentation ahead of implementation

**Status:** Open (expected at greenfield)  
**Recorded:** 2026-06-10

#### Summary

Context files describe the full MVP product. Implementation has not started beyond the Vue starter and context docs. Do not assume completed work from old NovelCraft progress notes.

Implement against **new** feature specs in vertical slices; update `progress-tracker.md` as each slice lands.
