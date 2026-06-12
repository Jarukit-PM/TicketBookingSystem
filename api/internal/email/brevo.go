package email

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/mail"
)

// BrevoClient sends mail through the Brevo transactional email API.
type BrevoClient struct {
	apiKey, from string
	baseURL      string
	httpClient   *http.Client
}

// NewBrevoClient returns a Brevo-backed sender.
func NewBrevoClient(apiKey, from string) *BrevoClient {
	return &BrevoClient{
		apiKey:     apiKey,
		from:       from,
		baseURL:    "https://api.brevo.com",
		httpClient: http.DefaultClient,
	}
}

func (c *BrevoClient) Send(ctx context.Context, msg Message) (string, error) {
	if c.apiKey == "" || c.from == "" {
		return "", fmt.Errorf("brevo not configured")
	}
	senderName, senderEmail := parseFromAddress(c.from)
	payload := map[string]any{
		"sender":      senderField(senderName, senderEmail),
		"to":          []map[string]string{{"email": msg.To}},
		"subject":     msg.Subject,
		"htmlContent": msg.HTMLBody,
		"textContent": msg.TextBody,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/v3/smtp/email", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("api-key", c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("accept", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		slurp, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return "", fmt.Errorf("brevo status %d: %s", resp.StatusCode, slurp)
	}
	var out struct {
		MessageID string `json:"messageId"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", err
	}
	return out.MessageID, nil
}

func parseFromAddress(from string) (name, email string) {
	addr, err := mail.ParseAddress(from)
	if err != nil {
		return "", from
	}
	return addr.Name, addr.Address
}

func senderField(name, email string) map[string]string {
	sender := map[string]string{"email": email}
	if name != "" {
		sender["name"] = name
	}
	return sender
}
