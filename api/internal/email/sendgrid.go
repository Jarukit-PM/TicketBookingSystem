package email

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Message is a transactional email payload.
type Message struct {
	To, Subject, HTMLBody, TextBody string
}

// Sender delivers email via an external provider.
type Sender interface {
	Send(ctx context.Context, msg Message) (string, error)
}

// SendGridClient sends mail through SendGrid v3 API.
type SendGridClient struct {
	apiKey, from string
}

// NewSendGridClient returns a SendGrid-backed sender.
func NewSendGridClient(apiKey, from string) *SendGridClient {
	return &SendGridClient{apiKey: apiKey, from: from}
}

func (c *SendGridClient) Send(ctx context.Context, msg Message) (string, error) {
	if c.apiKey == "" || c.from == "" {
		return "", fmt.Errorf("sendgrid not configured")
	}
	body, err := json.Marshal(map[string]any{
		"personalizations": []map[string]any{{"to": []map[string]string{{"email": msg.To}}}},
		"from":             map[string]string{"email": c.from},
		"subject":          msg.Subject,
		"content": []map[string]string{
			{"type": "text/plain", "value": msg.TextBody},
			{"type": "text/html", "value": msg.HTMLBody},
		},
	})
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.sendgrid.com/v3/mail/send", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		slurp, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return "", fmt.Errorf("sendgrid status %d: %s", resp.StatusCode, slurp)
	}
	return resp.Header.Get("X-Message-Id"), nil
}
