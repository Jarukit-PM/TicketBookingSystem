# Feature 10 — Admin Dashboard & QR Scan

Read `context/CONTEXT.md`, `context/project-overview.md`, and `context/ui-context.md`.

**Admin dashboard:** today's metrics, schedule, recent bookings. **QR scan:** camera reads customer ticket → navigate to that user's **booking history** (support lookup — not door check-in).

**Depends on:** Features 05 (admin auth), 06 (catalog admin), 08 (bookings), 09 (ticket resolve API).

## Objective

Global admins operate all cinemas from one console. Scanning a ticket QR quickly opens the customer's booking history for support.

**Success looks like:** Dashboard shows live-ish stats; booking search works read-only; scan valid QR → `/admin/users/:userId/bookings`; invalid QR → error toast.

## Admin scope (`CONTEXT.md`)

- **Global admin** — all cinemas, no `cinemaId` on user
- **Read-only bookings** in MVP — no cancel/refund
- **QR scan** → user booking history only — no pass/fail check-in screen

## API Surface

| Method | Path | Notes |
| ------ | ---- | ----- |
| GET | `/api/admin/dashboard` | Counts: bookings today, upcoming showtimes, occupancy approx |
| GET | `/api/admin/bookings` | Search: `email`, `bookingRef`, `userId`, `showtimeId` |
| GET | `/api/admin/users/:userId/bookings` | Full history for support |
| GET | `/api/admin/tickets/resolve` | From spec 09 — `ref` + `t` |
| GET | `/api/admin/audit-logs` | Paginated, newest first |
| GET | `/api/admin/email-logs` | Filter by `bookingId` |

Occupancy %: for today's showtimes, `sold seats / (layout capacity - blocked)` aggregated — approximate OK for MVP.

## Vue Structure

```
app/src/views/admin/
├── AdminDashboardView.vue
├── AdminBookingsView.vue
├── AdminUserBookingsView.vue   # target after QR scan
├── AdminScanView.vue           # camera QR
├── AdminAuditLogsView.vue
└── AdminLayout.vue             # sidebar nav, admin shell
app/src/components/admin/
├── StatsCard.vue
├── BookingsTable.vue
└── QrScanner.vue               # browser getUserMedia or file input fallback
```

### Admin shell (`ui-context.md`)

- Dark theme consistent with customer app
- Sidebar: Dashboard, Movies, Cinemas, Screens, Showtimes, Bookings, Scan, Logs
- Tables: `bg-surface`, sticky header, muted secondary columns

### QR scan flow

1. Admin opens **Scan** → camera permission
2. Decode URL `.../ticket/{bookingRef}?t={token}`
3. `GET /api/admin/tickets/resolve?ref=...&t=...`
4. On success → `router.push(`/admin/users/${userId}/bookings`)`
5. On failure → toast error; stay on scan

**No** showtime window validation, **no** "admit/deny" UI.

## Code Style

Dashboard handler returns DTO:

```go
type Dashboard struct {
    BookingsToday   int     `json:"bookingsToday"`
    ShowtimesToday  int     `json:"showtimesToday"`
    AvgOccupancyPct float64 `json:"avgOccupancyPct"`
    RecentBookings  []BookingSummary `json:"recentBookings"`
}
```

## Testing Strategy

- Go: admin routes return 403 for customer JWT.
- Go: booking search filters.
- Vue: QrScanner parses sample URL (unit test pure parser).
- Manual: scan QR from spec 09 ticket on second device.

## Boundaries

- **Always:** Audit log for destructive admin actions (catalog mutations from spec 06).
- **Ask first:** Export CSV, pagination defaults.
- **Never:** Booking cancel/refund UI; door check-in mode.

## Tasks

- [ ] Dashboard aggregation queries
- [ ] Booking search + user bookings list handlers
- [ ] Audit log + email log list endpoints
- [ ] AdminLayout + dashboard + bookings table views
- [x] QrScanner + scan flow to user bookings
- [ ] Link catalog admin views into sidebar (spec 06)
- [ ] Seed data demo for dashboard metrics

## Out of scope

- Cancel/refund
- Per-cinema admin roles
- POS hardware
- Real-time dashboard WebSocket

## Check when done

- [ ] Customer cannot access `/admin/*` UI or API
- [ ] Dashboard loads with seed data
- [ ] Booking search by `bookingRef` and email works
- [ ] QR scan navigates to correct user history
- [ ] Invalid QR shows error
- [ ] `progress-tracker.md` updated when implementation lands

## Open Questions

- Camera library choice (`@zxing/browser` or similar) — pick at implementation.
