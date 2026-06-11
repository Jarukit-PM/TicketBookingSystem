## Application Building Context

Read the following files in order before implementing or making any architectural decision:

1. `context/project-overview.md` — product definition, goals, features, and scope
2. `context/architecture-context.md` — system structure, boundaries, storage model, and invariants
3. `context/ui-context.md` — theme, colors, typography, and component conventions
4. `context/code-standards.md` — implementation rules and conventions
5. `context/ai-workflow-rules.md` — development workflow, scoping rules, and delivery approach
6. `context/progress-tracker.md` — current phase, completed work, open questions, and next steps

Update `context/progress-tracker.md` after each meaningful implementation change.

If implementation changes the architecture, scope, or standards documented in the context files, update the relevant file before continuing.

## Agent Skills

Project skills live in `.agents/skills/` (see `skills-lock.json`; update with `npx skills update -p -y`).

Run **`setup-matt-pocock-skills`** once per repo before using the other Matt Pocock engineering skills.

### Domain documentation mapping

Skills that reference `CONTEXT.md` or `docs/adr/` use this repo's layout instead:

| Skill expectation | This repo |
| --- | --- |
| `CONTEXT.md` (glossary / ubiquitous language) | `context/project-overview.md` + relevant `context/feature-specs/*.md` |
| Architecture / boundaries | `context/architecture-context.md` |
| UI / design language | `context/ui-context.md` |
| Code conventions | `context/code-standards.md` |
| Workflow / scoping | `context/ai-workflow-rules.md` |
| Feature specs | `context/feature-specs/` |
| ADRs | Add under `context/adr/` when a decision is hard to reverse |

### Installed starter kit (12 skills)

**Matt Pocock** ([mattpocock/skills](https://github.com/mattpocock/skills)):

- `setup-matt-pocock-skills` — one-time repo config (issue tracker, triage labels, domain doc layout)
- `grill-with-docs` — align on domain language before building (seat hold, showtime, booking ref)
- `to-issues` — break specs into vertical-slice GitHub issues
- `tdd` — red-green-refactor for Go booking logic and Vue components
- `diagnose` — hard bugs (double booking, WebSocket desync, hold TTL)
- `improve-codebase-architecture` — periodic architecture cleanup

**Addy Osmani** ([addyosmani/agent-skills](https://github.com/addyosmani/agent-skills)):

- `using-agent-skills` — route work to the right skill workflow
- `spec-driven-development` — write specs under `context/feature-specs/` before code
- `api-and-interface-design` — Go REST + WebSocket contracts (`/holds`, `/confirm`, `/ws`)
- `frontend-ui-engineering` — seat map, countdown timer, admin UI, accessibility
- `security-and-hardening` — JWT auth, Google OAuth, admin routes, idempotent confirm
- `browser-testing-with-devtools` — seat map UX, real-time UI, Playwright + DevTools

### Stack skills (Vue, Go, MongoDB, Tailwind)

**Vue** ([antfu/skills](https://github.com/antfu/skills), [hyf0/vue-skills](https://github.com/hyf0/vue-skills), [vuejs-ai/skills](https://github.com/vuejs-ai/skills)):

- `vue` — core Vue 3 + Vite patterns
- `vue-best-practices` — components, Pinia, Vue Router
- `vue-debug-guides` — debugging reactive/WebSocket UI issues
- `create-adaptable-composable` — composables (`useShowtimeSocket`, hold countdown)

**Go API** ([samber/cc-skills-golang](https://github.com/samber/cc-skills-golang)):

- `golang-project-layout` — scaffold `api/cmd`, `internal/`
- `golang-code-style` — idiomatic Go for Gin handlers
- `golang-design-patterns` — booking locks, worker patterns

**MongoDB** ([mongodb/agent-skills](https://github.com/mongodb/agent-skills)):

- `mongodb-schema-design` — collections, embed vs reference, indexes
- `mongodb-connection` — driver pools, connection config
- `mongodb-query-optimizer` — slow showtime/booking queries

**Tailwind** ([wshobson/agents](https://github.com/wshobson/agents)):

- `tailwind-design-system` — tokens, layout, seat map styling

### Invoke when relevant

| Task | Skill |
| --- | --- |
| New feature or major change | `grill-with-docs` → `spec-driven-development` |
| Break a spec into issues | `to-issues` |
| Go API / WebSocket design | `api-and-interface-design` + `golang-project-layout` |
| Scaffold or structure `api/` | `golang-project-layout`, `golang-code-style` |
| Go handler/service logic | `golang-code-style`, `golang-design-patterns` |
| MongoDB collections / indexes | `mongodb-schema-design` |
| MongoDB driver / connection issues | `mongodb-connection` |
| Slow booking or showtime queries | `mongodb-query-optimizer` |
| Vue components, Pinia, routing | `vue`, `vue-best-practices` |
| Vue composables (WebSocket, holds) | `create-adaptable-composable` |
| Vue UI bugs | `vue-debug-guides` |
| Seat map, booking flow, admin UI | `frontend-ui-engineering` + `tailwind-design-system` |
| Tailwind tokens and layout | `tailwind-design-system` |
| Auth, JWT, booking confirm, admin | `security-and-hardening` |
| Logic or behavior changes | `tdd` |
| Browser UX verification | `browser-testing-with-devtools` |
| Bugs, failed builds, unexpected behavior | `diagnose` |
| Periodic architecture cleanup | `improve-codebase-architecture` |
| Unsure which skill applies | `using-agent-skills` |
