# Feature 04 — MongoDB Data Model & Indexes

Read `context/CONTEXT.md`, `context/architecture-context.md`, and `mongodb-schema-design` skill before starting.

Define Go structs, collection names, and **index creation** for all durable MVP entities. Seed script optional for local dev. No HTTP handlers beyond what spec 03 provides.

**Depends on:** Feature 03 (mongo connection). Domain rules from grill session in `context/CONTEXT.md`.

## Objective

MongoDB collections match `architecture-context.md`. Indexes enforce uniqueness and query paths. Repository interfaces exist for later specs to implement.

**Success looks like:** `go test ./internal/...` covers index definitions; `api/migrations/` or startup hook ensures indexes; seed creates sample cinema, screen, movie, showtime.

## Assumptions

1. Database name: `tbs` (configurable).
2. IDs: MongoDB `ObjectID` as hex string in JSON API responses.
3. Money: `int64` minor units (cents) on `total` and `priceTiers`.
4. Sold seats **not** stored on showtime — derived from `bookings` (see `CONTEXT.md`).
5. Holds are **not** in MongoDB.

## Commands

```bash
cd api
go test ./internal/catalog/... ./internal/booking/...
# Optional seed (implement as cmd/seed or Make target)
go run ./cmd/seed
```

## Project Structure

```
api/internal/
├── catalog/
│   ├── models.go       # Cinema, Screen, Movie, Showtime
│   └── repository.go   # interface + mongo impl
├── booking/
│   ├── models.go       # Booking
│   └── repository.go
├── user/
│   ├── models.go       # User
│   └── repository.go
├── audit/
│   └── models.go       # AuditLog, EmailLog
└── db/
    └── indexes.go      # EnsureIndexes(ctx, db)
api/migrations/
└── README.md           # document index script / EnsureIndexes on boot
api/cmd/seed/main.go    # optional dev seed
```

## Domain models (summary)

Reference `architecture-context.md` tables. Key invariants:

| Collection | Unique indexes | Query indexes |
| ---------- | -------------- | ------------- |
| `users` | `email`; sparse `googleId` | — |
| `movies` | — | `status` |
| `cinemas` | — | — |
| `screens` | — | `cinemaId` |
| `showtimes` | — | `(screenId, startsAt)`; `movieId` |
| `bookings` | `bookingRef` | `userId`; `(userId, showtimeId)` **not unique** |
| `audit_logs` | — | `createdAt` |
| `email_logs` | — | `bookingId` |

### Booking document

```go
type Booking struct {
    ID           primitive.ObjectID `bson:"_id,omitempty"`
    UserID       primitive.ObjectID `bson:"userId"`
    ShowtimeID   primitive.ObjectID `bson:"showtimeId"`
    Seats        []string           `bson:"seats"` // seatIds e.g. "A-12"
    Total        int64              `bson:"total"`
    BookingRef   string             `bson:"bookingRef"`   // TBS-XXXXXX
    TicketToken  string             `bson:"ticketToken"`  // opaque, not equal to bookingRef
    Status       string             `bson:"status"`       // CONFIRMED only in MVP
    ConfirmedAt  time.Time          `bson:"confirmedAt"`
}
```

### Showtime `priceTiers`

```go
type PriceTiers struct {
    Standard   int64 `bson:"standard" json:"standard"`
    VIP        int64 `bson:"vip" json:"vip"`
    Wheelchair int64 `bson:"wheelchair" json:"wheelchair"`
}
```

### Screen layout seat

```go
type LayoutSeat struct {
    SeatID string `bson:"seatId" json:"seatId"`
    Row    int    `bson:"row" json:"row"`
    Col    int    `bson:"col" json:"col"`
    Type   string `bson:"type" json:"type"` // standard, vip, wheelchair, blocked
}
```

## Testing Strategy

- Table-driven tests for `bookingRef` generator format (`TBS-`, no ambiguous chars) when placed in `internal/booking/ref.go`.
- Repository tests with `mongo/integration` tag or `testcontainers` (optional); minimum unit tests on BSON tags and validation helpers.
- Test `EnsureIndexes` does not error on second run (idempotent).

## Boundaries

- **Always:** Run `EnsureIndexes` on server boot or dedicated migrate command.
- **Ask first:** Schema breaking changes, new collections.
- **Never:** Store seat holds in MongoDB; add `soldSeatIds[]` on showtimes.

## Tasks

- [ ] Define all model structs and collection constants
- [ ] Implement `db.EnsureIndexes`
- [ ] Repository interfaces + mongo implementations (CRUD stubs OK)
- [ ] `cmd/seed`: 1 cinema, 2 screens, 2 movies, several showtimes
- [ ] Unit tests for ref generator and price tier sum helper

## Out of scope

- HTTP CRUD routes (spec 06 admin, spec 05 auth for users)
- Redis hold keys
- Booking confirm logic

## Check when done

- [ ] All collections from `architecture-context.md` represented
- [ ] Unique `bookingRef`; non-unique `(userId, showtimeId)` documented
- [ ] Seed data supports manual API testing in spec 06+
- [ ] `go test ./...` passes
- [ ] `progress-tracker.md` updated when implementation lands

## Open Questions

- None — grill decisions captured in `CONTEXT.md`.
