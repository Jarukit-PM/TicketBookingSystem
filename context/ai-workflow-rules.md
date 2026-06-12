# Development Workflow

## Approach

Build this project incrementally using a spec-driven workflow. Context files define what to build, how to build it, and what the current state of progress is. Always implement against these specs — do not infer or invent behavior from scratch.

Read context in order before implementing:

1. `project-overview.md`
2. `architecture-context.md`
3. `ui-context.md`
4. `code-standards.md`
5. `progress-tracker.md`

## Scoping Rules

- Work on one feature unit or vertical slice at a time (e.g. "hold API + Redis TTL" before "full seat map polish").
- Prefer small, verifiable increments over large speculative changes.
- Do not combine unrelated system boundaries in a single implementation step (e.g. admin movie CRUD + WebSocket hub in one PR).

## When To Split Work

Split an implementation step if it combines:

- Vue UI and Go API changes that cannot be verified without both (prefer API-first with curl/integration tests, then wire UI).
- Pinia client state and MongoDB persistence in one ambiguous step.
- WebSocket broadcast logic and email worker in the same change.
- Multiple unrelated API routes without a shared spec.
- Behavior that is not clearly defined in the context files.

If a change cannot be verified end to end quickly, the scope is too broad — split it.

## Handling Missing Requirements

- Do not invent product behavior that is not defined in the context files.
- If a requirement is ambiguous, resolve it in the relevant context file (or new feature spec) before implementing.
- If a requirement is missing, add it as an open question in `progress-tracker.md` before continuing.
- Legacy NovelCraft specs in `context/feature-specs/` are **not** authoritative — rewrite for cinema booking before use.

## Protected Foundation Components

Do not modify generated or vendor foundation artifacts unless explicitly instructed.

This includes:

- Third-party library internals
- Default shadcn or CLI-generated files if added later — adapt at the feature/block layer, not by forking primitives unnecessarily

Project-specific styling, layout, and feature logic belong in `app/src/components/`, `app/src/views/`, and `app/src/composables/` — not scattered across unrelated modules.

## Keeping Docs In Sync

Update the relevant context file whenever implementation changes:

- System architecture or boundaries → `architecture-context.md`
- Product scope or features → `project-overview.md`
- UI tokens or patterns → `ui-context.md`
- Conventions → `code-standards.md`
- Blockers → `current-issue.md`

Progress state must reflect the **actual** state of the implementation, not the intended state. After each meaningful slice, update `progress-tracker.md`.

## Booking-Specific Workflow Notes

- **Holds before confirm:** implement and test Redis hold TTL + extension before building the full Vue seat map.
- **Authority:** HTTP responses win over WebSocket events when they disagree.
- **Idempotency:** confirm booking must be safe to retry — test with duplicate `Idempotency-Key`.
- **Email:** enqueue via asynq; never block confirm response on Brevo.
- **Admin vs customer:** keep `/api/admin/*` behind `RequireAdmin`; do not leak admin fields in public catalog DTOs.

## Before Moving To The Next Unit

1. The current unit works end to end within its defined scope (or API-only slice is proven with tests).
2. No invariant defined in `architecture-context.md` was violated.
3. `progress-tracker.md` reflects the completed work.
4. Open blockers in `current-issue.md` are updated or resolved.

## Skills (when relevant)

| Task | Skill |
| ---- | ----- |
| New feature | `grill-with-docs` → `spec-driven-development` |
| Go API / WebSocket | `api-and-interface-design`, `golang-project-layout` |
| Seat map / booking UI | `frontend-ui-engineering`, `vue-best-practices`, `tailwind-design-system` |
| Holds / confirm bugs | `diagnose`, `tdd` |
| MongoDB schema | `mongodb-schema-design` |
| Auth / JWT | `security-and-hardening` |
