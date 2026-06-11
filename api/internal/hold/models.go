package hold

import "time"

// Record is the JSON payload stored at hold:{showtimeId}:{seatId}.
type Record struct {
	UserID string    `json:"userId"`
	HeldAt time.Time `json:"heldAt"`
}

// Result is returned after add or remove operations.
type Result struct {
	Holds     []string   `json:"holds"`
	ExpiresAt *time.Time `json:"expiresAt,omitempty"`
}
