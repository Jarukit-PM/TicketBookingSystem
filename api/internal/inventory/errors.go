package inventory

import "errors"

var (
	ErrShowtimeNotFound = errors.New("showtime not found")
	ErrScreenNotFound   = errors.New("screen not found")
)
