package email

import "context"

// Message is a transactional email payload.
type Message struct {
	To, Subject, HTMLBody, TextBody string
}

// Sender delivers email via an external provider.
type Sender interface {
	Send(ctx context.Context, msg Message) (string, error)
}
