# Feature 01 — Design System

Read `AGENTS.md`, `context/ui-context.md`, and `context/code-standards.md` before starting.

Establish the **dark cinema + gradient orange** design foundation for the Vue 3 SPA: CSS tokens in `app/src/assets/`, Tailwind CSS v4, base UI primitives, and a `cn()` class helper. This is the visual base for catalog, seat map, booking, and admin UI in later specs.

**Depends on:** Vue 3 + Vite starter in `app/`. No Go API, auth, or business logic.

## Replace the Vite starter theme

The default `app/src/assets/base.css` and `main.css` use Inter/green links, light backgrounds, and `prefers-color-scheme: dark`. **Remove all of that** — the app is **dark-only** with black surfaces and orange gradient brand per `ui-context.md`.

### `app/src/assets/tokens.css` (new)

Define CSS custom properties from `ui-context.md`:

- Surfaces: `--bg-base`, `--bg-surface`, `--bg-elevated`, `--bg-subtle`
- Borders: `--border-default`, `--border-subtle`
- Text: `--text-primary`, `--text-secondary`, `--text-muted`, `--text-faint`
- Brand: `--accent-primary`, `--accent-primary-hover`, `--accent-primary-dim`, `--accent-glow`
- Gradients: `--gradient-brand`, `--gradient-brand-hover`, `--gradient-brand-subtle`
- Status: `--state-error`, `--state-success`, `--state-warning`, and `*-dim` variants
- Shadows: `--shadow-1`, `--shadow-2`

### `app/src/assets/main.css`

- `@import './tokens.css'`
- `@import 'tailwindcss'`
- Map tokens in `@theme inline` to Tailwind utilities (see `ui-context.md`):
  - `bg-base`, `bg-surface`, `bg-elevated`, `bg-subtle`
  - `text-copy-primary`, `text-copy-secondary`, `text-copy-muted`
  - `border-surface-border`
  - `text-brand`, `bg-accent-dim`, `shadow-glow-brand`
  - Custom utilities: `bg-gradient-brand`, `bg-gradient-brand-hover`, `bg-gradient-brand-subtle`
- Global `body`: `bg-base text-copy-primary antialiased`
- **Do not** add light-mode overrides or raw palette classes (`zinc-*`, `orange-500`) in app code after this spec.

### `app/index.html`

- Update `<title>` to the cinema product name (align with `project-overview.md`).
- Load **Inter** (or DM Sans) via Google Fonts; set `--font-sans` on `:root`.

## Install Tailwind CSS v4

In `app/`:

```bash
npm install tailwindcss @tailwindcss/vite
```

Wire `@tailwindcss/vite` in `vite.config.ts`. Follow `tailwind-design-system` skill if needed.

Install supporting deps:

```bash
npm install clsx tailwind-merge
npm install lucide-vue-next
```

## Utilities

Create `app/src/lib/cn.ts`:

```ts
import { clsx, type ClassValue } from 'clsx'
import { twMerge } from 'tailwind-merge'

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}
```

## Base UI primitives

Create project-owned components in `app/src/components/ui/` (not shadcn CLI — Vue project):

| Component | File | Minimum variants |
| --------- | ---- | ---------------- |
| Button | `Button.vue` | `primary` (gradient), `secondary` (outlined), `ghost`, `destructive` |
| Card | `Card.vue` | `Card`, `CardHeader`, `CardTitle`, `CardContent` slots or subcomponents |
| Input | `Input.vue` | Dark surface, orange focus ring |
| Badge | `Badge.vue` | Pill; maps booking/seat status colors from `ui-context.md` |

Style per `ui-context.md`:

- Primary button: `rounded-full bg-gradient-brand text-white hover:bg-gradient-brand-hover`
- Card: `rounded-xl bg-surface border border-surface-border shadow-elevation-1`
- Input: `rounded-lg bg-surface border border-surface-border focus:ring-2` using `--accent-glow`

Export a barrel `app/src/components/ui/index.ts` if helpful.

## App shell (token verification)

Replace default `HomeView.vue` content with a minimal **design-system preview** page (dev-only OK to keep as home until catalog spec):

- Page on `bg-base`
- App bar stub: `h-16`, blurred dark bar, gradient logo text sample
- One **gradient primary** button, one **outlined** button
- One **Card** with sample movie title + secondary text
- One **Badge** row: Confirmed, Pending, Expired (MVP statuses — no Cancelled)
- Confirm no green Vite starter link styles remain

Update `App.vue` to full-width dark shell (`min-h-screen bg-base`); remove Vite two-column demo layout from `main.css`.

## Out of scope

- Go API, Pinia auth store, Vue Router guards beyond existing stubs
- Seat map, booking flow, admin tables
- WebSocket composables
- Email or QR ticket views (white QR pad pattern documented in `ui-context.md` for a later spec)

## Check when done

- `tokens.css` + `main.css` use black + gradient orange tokens; no Vite green theme; no light-mode toggle
- Tailwind v4 builds; `npm run build` and `npm run type-check` pass
- `cn()` works; `Button`, `Card`, `Input`, `Badge` render with design tokens
- Home/preview page demonstrates gradient CTA and dark surfaces
- No `zinc-*`, `orange-500`, or hardcoded hex in `app/src/` (tokens only)
- `progress-tracker.md` updated when implementation lands
