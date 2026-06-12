package email

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBrevoClient_SendNotConfigured(t *testing.T) {
	t.Parallel()

	c := NewBrevoClient("", "from@example.com")
	_, err := c.Send(context.Background(), Message{To: "to@example.com", Subject: "Hi"})
	if err == nil {
		t.Fatal("expected error when brevo is not configured")
	}
}

func TestBrevoClient_SendOK(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v3/smtp/email" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		if got := r.Header.Get("api-key"); got != "xkeysib-test" {
			t.Fatalf("api-key = %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"messageId":"<abc@brevo.com>"}`))
	}))
	defer srv.Close()

	c := NewBrevoClient("xkeysib-test", "TBS <tickets@example.com>")
	c.baseURL = srv.URL
	id, err := c.Send(context.Background(), Message{
		To:       "customer@example.com",
		Subject:  "Booking confirmed",
		HTMLBody: "<p>Thanks</p>",
		TextBody: "Thanks",
	})
	if err != nil {
		t.Fatal(err)
	}
	if id != "<abc@brevo.com>" {
		t.Fatalf("id = %q, want <abc@brevo.com>", id)
	}
}

func TestParseFromAddress(t *testing.T) {
	t.Parallel()

	name, email := parseFromAddress("Ticket Booking <tickets@example.com>")
	if name != "Ticket Booking" || email != "tickets@example.com" {
		t.Fatalf("parsed = %q <%q>", name, email)
	}
	name, email = parseFromAddress("plain@example.com")
	if name != "" || email != "plain@example.com" {
		t.Fatalf("parsed plain = %q <%q>", name, email)
	}
}
