# Feature 06 — Catalog (Browse & Admin CRUD)

Read `context/CONTEXT.md`, `context/project-overview.md`, `context/ui-context.md`, and `api-and-interface-design` skill.

**Customer:** pick cinema, browse **Now Showing** / **Coming Soon**, movie detail, showtime list. **Admin:** global movie CRUD, cinema/screen/showtime management.

**Depends on:** Features 01, 03, 04, 05 (admin routes).

## Objective

Customers discover films and showtimes per cinema with correct browse rules (`CONTEXT.md`). Admins manage global movies and per-cinema scheduling without code changes.

**Success looks like:** Customer selects cinema → Now Showing lists only movies with future showtimes; Coming Soon shows `COMING_SOON` teasers without showtimes; admin creates movie + showtime → appears in browse after refresh.

## Browse rules (from `CONTEXT.md`)

| Tab | Inclusion rule |
| --- | -------------- |
| **Now Showing** | `status != ARCHIVED` AND ≥1 **future** showtime at selected cinema |
| **Coming Soon** | `status == COMING_SOON` (teaser OK with **zero** showtimes) |
| **Hidden** | `ARCHIVED` |

Movies are **global**; showtimes are cinema-scoped via `screen.cinemaId`.

## API Surface

### Public

| Method | Path | Query | Response |
| ------ | ---- | ----- | -------- |
| GET | `/api/cinemas` | — | `[{ id, name, address, timezone }]` |
| GET | `/api/movies` | `cinemaId`, `tab=now_showing\|coming_soon` | Movie cards + metadata |
| GET | `/api/movies/:id` | `cinemaId` | Detail + upcoming showtimes at cinema |
| GET | `/api/showtimes` | `cinemaId`, `movieId`, `date?` | Future showtimes only; `startsAt > now` in cinema TZ |

### Admin (`RequireAdmin`)

| Method | Path | Notes |
| ------ | ---- | ----- |
| CRUD | `/api/admin/movies` | Global catalog |
| CRUD | `/api/admin/cinemas` | Venues |
| CRUD | `/api/admin/screens` | Layout with seats `{ seatId, row, col, type }` |
| CRUD | `/api/admin/showtimes` | `priceTiers`, `startsAt`, `movieId`, `screenId` |
| GET | `/api/admin/showtimes` | Filters for dashboard (spec 10) |

Audit log entry on admin mutations (`audit_logs`).

## Vue Structure

```
app/src/views/
├── HomeView.vue              # cinema picker + tabs
├── MovieDetailView.vue
├── ShowtimesView.vue         # or combined in detail
app/src/views/admin/
├── AdminMoviesView.vue
├── AdminCinemasView.vue
├── AdminScreensView.vue
└── AdminShowtimesView.vue
app/src/stores/catalog.ts     # selectedCinemaId persisted (localStorage OK)
```

### UI notes (`ui-context.md`)

- Movie cards: poster, title, rating, gradient CTA
- Cinema selector: sticky header or home step
- Coming Soon detail: show "Showtimes not yet announced" when no showtimes

## Code Style

Showtime query uses cinema timezone for `startsAt` comparison:

```go
// Reject listing past showtimes for customer browse
filter := bson.M{
    "screenId": bson.M{"$in": screenIDs},
    "startsAt": bson.M{"$gt": nowInCinemaTZ(cinema.Timezone)},
}
```

## Testing Strategy

- Go table tests: `now_showing` filter excludes movie with no future showtimes at cinema.
- Go: `coming_soon` includes `COMING_SOON` without showtimes.
- Go: archived movies never returned.
- Vue: catalog store persists cinema selection.

## Boundaries

- **Always:** Use cinema `timezone` for showtime cutoff in listings.
- **Ask first:** Poster upload to S3 (use URL string in MVP).
- **Never:** Per-cinema movie documents; list bookable films without showtimes in Now Showing.

## Tasks

- [ ] Public catalog handlers + repository queries
- [ ] Admin CRUD handlers with audit log
- [ ] Screen layout editor (JSON or simple grid UI — minimum viable form)
- [ ] Customer views: cinema picker, tabs, movie detail, showtime list
- [ ] Admin views: movies, cinemas, screens, showtimes tables
- [ ] Link from showtime → `/book/:showtimeId` (seat map spec 07)

## Out of scope

- Seat map, holds, booking
- Dynamic pricing beyond `priceTiers`
- Genre filter (optional nice-to-have — skip unless quick)

## Check when done

- [ ] Now Showing / Coming Soon rules match `CONTEXT.md`
- [ ] Admin can create full chain: cinema → screen → movie → showtime
- [ ] Customer flow: cinema → movie → pick showtime → route to booking
- [ ] `go test` for browse filters; `npm run build` passes
- [ ] `progress-tracker.md` updated when implementation lands

## Open Questions

- None.
