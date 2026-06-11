package audit

// Booking and seat lifecycle actions (system + customer).
const (
	ActionBookingSuccess = "booking_success"
	ActionBookingTimeout = "booking_timeout"
	ActionSeatReleased   = "seat_released"
	ActionBookingFailed  = "booking_failed"
	ActionSystemError    = "system_error"
)

// Admin catalog mutation actions.
const (
	ActionCreate = "create"
	ActionUpdate = "update"
	ActionDelete = "delete"
)
