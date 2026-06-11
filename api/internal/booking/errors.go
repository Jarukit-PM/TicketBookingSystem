package booking

import "errors"

var (
	ErrNoActiveHolds       = errors.New("no active holds")
	ErrIdempotencyRequired = errors.New("idempotency key required")
	ErrSeatConflict        = errors.New("seat conflict")
	ErrShowtimeNotFound    = errors.New("showtime not found")
	ErrShowtimeStarted     = errors.New("showtime already started")
	ErrSeatLimitExceeded   = errors.New("seat limit exceeded")
)
