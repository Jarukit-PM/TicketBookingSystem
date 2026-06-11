# UI Context

## Theme

**Dark cinema aesthetic.** Black surfaces, light text, and **gradient orange** for brand accents, primary CTAs, and key highlights. Feels like a premium movie-ticket app — not a generic light admin UI.

No light mode for MVP. Optional `prefers-reduced-motion` should disable animated gradient shifts only; colors stay the same.

All colors are defined as CSS custom properties in `app/src/assets/main.css` (or `app/src/assets/tokens.css`) and mapped to Tailwind via `@theme inline`. Components must use these tokens — no hardcoded hex values or raw Tailwind palette classes like `zinc-*`, `orange-500`, or `black`.

### Core palette

| Role             | CSS Variable             | Value / Notes |
| ---------------- | ------------------------ | ------------- |
| Page background  | `--bg-base`              | `#0a0a0a` — near-black |
| Surface          | `--bg-surface`           | `#141414` — cards, panels |
| Elevated surface | `--bg-elevated`          | `#1c1c1c` — modals, dropdowns |
| Subtle surface   | `--bg-subtle`            | `#242424` — hover rows, chips |
| Default border   | `--border-default`       | `#2e2e2e` |
| Subtle border    | `--border-subtle`        | `#1f1f1f` |
| Primary text     | `--text-primary`         | `#f5f5f5` |
| Secondary text   | `--text-secondary`       | `#a3a3a3` |
| Muted text       | `--text-muted`           | `#737373` |
| Faint text       | `--text-faint`           | `#525252` |
| Brand solid      | `--accent-primary`       | `#ff7a18` — fallback / borders |
| Brand hover      | `--accent-primary-hover` | `#ff9533` |
| Brand dim        | `--accent-primary-dim`   | `rgba(255, 122, 24, 0.12)` — focus / tonal bg |
| Brand glow       | `--accent-glow`          | `rgba(255, 122, 24, 0.35)` — shadows, focus rings |
| Error            | `--state-error`          | `#f87171` |
| Success          | `--state-success`        | `#4ade80` |
| Warning          | `--state-warning`        | `#fbbf24` |
| Success dim      | `--state-success-dim`    | `rgba(74, 222, 128, 0.12)` |
| Error dim        | `--state-error-dim`      | `rgba(248, 113, 113, 0.12)` |
| Warning dim      | `--state-warning-dim`    | `rgba(251, 191, 36, 0.12)` |
| Elevation 1      | `--shadow-1`             | `0 4px 24px rgba(0, 0, 0, 0.45)` |
| Elevation 2      | `--shadow-2`             | `0 8px 40px rgba(0, 0, 0, 0.55)` |

### Brand gradient (orange)

Use for primary buttons, hero accents, active nav, countdown urgency, and promotional strips. Define once:

```css
--gradient-brand: linear-gradient(135deg, #ff6b00 0%, #ff9500 50%, #ffb347 100%);
--gradient-brand-hover: linear-gradient(135deg, #ff7a18 0%, #ffa033 50%, #ffc266 100%);
--gradient-brand-subtle: linear-gradient(180deg, rgba(255, 107, 0, 0.08) 0%, transparent 100%);
```

Tailwind utilities (map in `@theme inline`):

- Surfaces/text: `bg-base`, `bg-surface`, `bg-elevated`, `bg-subtle`, `text-copy-primary`, `text-copy-secondary`, `text-copy-muted`, `border-surface-border`
- Brand: `text-brand`, `bg-brand-solid`, `bg-accent-dim`, `shadow-glow-brand`
- Gradients: `bg-gradient-brand`, `bg-gradient-brand-hover`, `bg-gradient-brand-subtle` (custom utilities wrapping the CSS vars above)

**Rule:** CTAs and hero emphasis use **`bg-gradient-brand`**; secondary actions use outlined orange or `bg-accent-dim`. Do not mix flat blue or Material grey-blue anywhere.

### Seat status colors

Maps to seat inventory in `architecture-context.md`.

| Status | Seat fill / border | Notes |
| ------ | ------------------ | ----- |
| Available | `--bg-subtle` / `--border-default` | Neutral dark tile |
| Selected (self) | `bg-gradient-brand` / white text | Strong orange gradient |
| Held (other user) | `--state-warning-dim` / `--state-warning` border | Amber, not orange — distinct from self |
| Sold | `--bg-elevated` + `--text-faint`, disabled | Muted, no gradient |
| Blocked | `--state-error-dim` + diagonal hatch | Admin maintenance |

Seat buttons: minimum **40×40px** touch target; show `seatId` (e.g. `A-12`) on hover/focus. VIP / wheelchair: small icon badge — do not rely on color alone.

### Booking status badges

Pill shape (`rounded-full px-3 py-1 text-xs font-medium`).

| Status | Background | Text |
| ------ | ---------- | ---- |
| Confirmed | `--state-success-dim` | `--state-success` |
| Pending (hold active) | `--accent-primary-dim` | `--accent-primary` + optional gradient dot |
| Expired | `--bg-subtle` | `--text-muted` |
| Cancelled | `--state-error-dim` | `--state-error` | *(phase 2 — not in MVP)* |

## Typography

| Role | Font | CSS Variable |
| ---- | ---- | ------------ |
| UI text | Inter + **Noto Sans Thai** | `--font-sans` |
| Display / hero | Same family, tighter tracking | `--font-display` |
| Code / mono | JetBrains Mono or Roboto Mono | `--font-mono` |

**Bilingual stack:** `'Inter', 'Noto Sans Thai', system-ui, sans-serif` — Inter for Latin UI; Noto Sans Thai for Thai script (spec 11). Load both from Google Fonts in `index.html`. White/light grey text on black — never dark text on dark surfaces.

### Type scale

| Use | Classes |
| --- | ------- |
| Hero / page title | `text-3xl md:text-4xl font-semibold tracking-tight text-copy-primary` |
| Section title | `text-lg font-medium text-copy-primary` |
| Body | `text-sm font-normal leading-6 text-copy-primary` |
| Secondary body | `text-sm text-copy-secondary` |
| Label / caption | `text-xs font-medium uppercase tracking-wide text-copy-muted` |
| Movie title (card) | `text-lg font-semibold text-copy-primary` |
| Showtime meta | `text-sm text-copy-secondary` |
| Price / total | `text-base font-semibold text-copy-primary`; accent total may use `text-brand` |

Optional hero headline: one word or phrase with `bg-gradient-brand bg-clip-text text-transparent` for orange gradient text (use sparingly — logo or "Now Showing" only).

## Spacing

**8px grid.** Page padding: `px-4 md:px-6`. Section gaps: `gap-4`, `gap-6`, `gap-8`. App bar height: `h-16` (64px).

## Border radius and elevation

Dark UI uses **subtle borders + glow** on focus/active, not heavy Material shadows alone.

| Context | Class / token |
| ------- | ------------- |
| Primary CTA | `rounded-full bg-gradient-brand` |
| Secondary / outlined | `rounded-full border border-[--accent-primary]/40 text-brand hover:bg-accent-dim` |
| Inputs | `rounded-lg bg-surface border border-surface-border focus:ring-2 focus:ring-[--accent-glow]` |
| Cards / panels | `rounded-xl bg-surface border border-surface-border shadow-elevation-1` |
| Modal | `rounded-2xl bg-elevated border border-surface-border shadow-elevation-2` |
| Seat cells | `rounded-md` |
| Movie poster | `rounded-lg ring-1 ring-white/10` |

Hover on list rows: `hover:bg-subtle`. Active nav link: gradient underline or `text-brand` + bottom border gradient.

## Component Library

**Vue 3 + Tailwind CSS v4** in `app/`. Composition API with `<script setup lang="ts">`.

| Layer | Location | Notes |
| ----- | -------- | ----- |
| Primitives | `app/src/components/ui/` | Button, Input, Dialog, Card — dark + orange tokens |
| Feature | `app/src/components/` | Seat map, movie cards, booking summary, admin tables |
| Composables | `app/src/composables/` | `useShowtimeSocket`, hold countdown, auth |
| Views | `app/src/views/` | Customer + admin routes |

Replace default Vite starter styles when implementing the design system.

### Button variants

| Variant | Style |
| ------- | ----- |
| Filled primary | `rounded-full bg-gradient-brand text-white font-medium hover:bg-gradient-brand-hover shadow-glow-brand` |
| Text / link | `text-brand hover:bg-accent-dim rounded-full px-4` |
| Outlined | `rounded-full border border-white/20 bg-transparent text-copy-primary hover:border-brand/50 hover:bg-accent-dim` |
| Tonal | `rounded-full bg-accent-dim text-brand` |
| Destructive | `text-state-error`; solid red fill only on confirm dialogs |

Icon buttons: `rounded-full p-2 hover:bg-subtle text-copy-secondary hover:text-copy-primary`, `h-10 w-10` min touch target.

## Layout Patterns

### App shell (customer)

Top app bar: `h-16`, `bg-base/80 backdrop-blur-md`, `border-b border-surface-border`. Logo may use gradient text. Nav: **Now Showing**, **Coming Soon**, **My Bookings**, sign-in. Optional thin **`bg-gradient-brand-subtle`** strip under app bar on marketing pages.

Content `max-w-6xl mx-auto` for catalog; full-width for seat map.

### Movie catalog

Card grid on `bg-base`. Poster: `aspect-[2/3]`, `rounded-lg`, ring. Card: `bg-surface`, hover lift + subtle orange glow on hover (`shadow-glow-brand` at low opacity). CTA: **View showtimes** — filled gradient pill.

Filter chips: `bg-subtle` default; active chip `bg-accent-dim text-brand border border-brand/30`.

### Showtime picker

Horizontal **date strip** above showtime cards: **7 days per page** with prev/next paging (not a long scroll of every date). Day chips show abbreviated weekday + day number; month/year label above. Selected chip uses `bg-gradient-brand`; unselected uses `border-brand/30` on `bg-surface`. Showtime cards for the selected day show **time only** (date is implied by the filter).

Rows on `bg-surface` cards. Selected row: left border or ring in gradient orange. Time and price prominent; screen name secondary.

### Seat map (core UX)

- Page background `bg-base`; map panel `bg-surface rounded-xl`.
- **Screen label** toward top; optional gradient divider.
- **Legend** with dark chips matching seat states.
- **Order panel**: sticky / bottom sheet, `bg-elevated`. **Hold countdown** in `text-brand` or gradient when under 1 minute.
- **Confirm booking**: full-width gradient primary on mobile.

### Booking confirmation + ticket

Success icon with green check on dark surface; booking ref in `text-brand`. QR on white pad (`bg-white p-4 rounded-lg`) for scan reliability — only place pure white is required.

### My Bookings

List rows `bg-surface` with poster thumb. Upcoming: subtle orange left accent bar on next showtime.

### Admin shell

Same dark tokens. Sidebar `bg-surface border-r border-surface-border`; active item `bg-accent-dim text-brand`. Tables: zebra optional via `bg-subtle/50` on alternate rows.

### Empty states

Muted icon in `bg-subtle rounded-full p-4`, title + `text-copy-secondary`, single **Browse movies** gradient CTA.

## Icons

**Lucide Vue** or similar outlined icons. Default `text-copy-secondary`; active/brand contexts `text-brand`. Sizes: `h-4 w-4` inline, `h-5 w-5` buttons, `h-6 w-6` app bar.

## Accessibility

- Contrast: body text `#f5f5f5` on `#0a0a0a` meets WCAG AA; muted text only for non-essential copy.
- Focus rings: `ring-2 ring-[--accent-glow]` on buttons, seats, inputs — visible on black.
- Seat map: `aria-label` with row, number, status.
- Countdown: `aria-live="polite"` under 1 minute.
- Do not use gradient orange as the only indicator for seat state — pair with icons or labels.

## Related Context

- `project-overview.md` — booking flow, hold TTL, ticket QR.
- `architecture-context.md` — WebSocket events, seat statuses, API sketch.
