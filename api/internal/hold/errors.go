package hold

import "errors"

var (
	ErrShowtimeNotFound  = errors.New("showtime not found")
	ErrScreenNotFound    = errors.New("screen not found")
	ErrCinemaNotFound    = errors.New("cinema not found")
	ErrShowtimeStarted   = errors.New("showtime already started")
	ErrSeatNotFound      = errors.New("seat not found in layout")
	ErrSeatBlocked       = errors.New("seat is blocked")
	ErrSeatSold          = errors.New("seat is sold")
	ErrSeatHeldByOther   = errors.New("seat held by another user")
	ErrSeatLimitExceeded = errors.New("seat limit exceeded")
	ErrSeatNotHeld       = errors.New("seat not in user holds")
)
