# MongoDB indexes

MVP indexes are created idempotently on API server boot via `internal/db.EnsureIndexes`.

## Collections and indexes

| Collection   | Index                         | Unique | Notes                                      |
| ------------ | ----------------------------- | ------ | ------------------------------------------ |
| `users`      | `email`                       | yes    |                                            |
| `users`      | `googleId`                    | yes    | sparse                                     |
| `movies`     | `status`                      | no     | browse filters                             |
| `screens`    | `cinemaId`                    | no     | list halls per cinema                      |
| `showtimes`  | `(screenId, startsAt)`        | no     | schedule queries                           |
| `showtimes`  | `movieId`                     | no     | film showtimes                             |
| `bookings`   | `bookingRef`                  | yes    | customer-facing ref (`TBS-…`)              |
| `bookings`   | `userId`                      | no     | My Bookings                                |
| `bookings`   | `(userId, showtimeId)`        | no     | multiple bookings per user per showtime OK |
| `audit_logs` | `createdAt`                   | no     | admin log listing                          |
| `email_logs` | `bookingId`                   | no     | delivery history per booking               |

## Manual verification

With Mongo running locally:

```bash
cd api
MONGO_URI=mongodb://localhost:27017/tbs go test ./internal/db/... -run TestEnsureIndexes_Idempotent -v
```

## Seed data

```bash
cd api
go run ./cmd/seed
```

Creates one cinema, two screens, two movies, and several future showtimes for local API testing.
