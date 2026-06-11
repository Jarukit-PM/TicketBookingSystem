package tasks

import (
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
)

const (
	// TypeEmailSend is the asynq task type for booking confirmation emails.
	TypeEmailSend = "email:send"
)

// EmailSendPayload is enqueued after a successful booking confirm.
type EmailSendPayload struct {
	BookingID string `json:"bookingId"`
}

// NewEmailSendTask builds an email:send task for the worker.
func NewEmailSendTask(bookingID string) (*asynq.Task, error) {
	payload, err := json.Marshal(EmailSendPayload{BookingID: bookingID})
	if err != nil {
		return nil, fmt.Errorf("marshal email task: %w", err)
	}
	return asynq.NewTask(TypeEmailSend, payload), nil
}
