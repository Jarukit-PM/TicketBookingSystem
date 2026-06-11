package email

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"
	texttemplate "text/template"
)

type confirmationData struct {
	BookingRef, MovieTitle, CinemaName, ScreenName, StartsAt, Seats, Total, TicketURL string
}

var (
	htmlTmpl = template.Must(template.New("html").Parse(`<h1>Booking confirmed</h1><p>Ref: {{.BookingRef}}</p><p>{{.MovieTitle}} — {{.CinemaName}} / {{.ScreenName}}</p><p>{{.StartsAt}}</p><p>Seats: {{.Seats}} · {{.Total}} THB</p><p><a href="{{.TicketURL}}">View ticket</a></p>`))
	textTmpl = texttemplate.Must(texttemplate.New("text").Parse("Booking {{.BookingRef}}\n{{.MovieTitle}}\n{{.TicketURL}}\n"))
)

func renderConfirmation(d confirmationData) (string, string, error) {
	var h, t bytes.Buffer
	if err := htmlTmpl.Execute(&h, d); err != nil {
		return "", "", fmt.Errorf("html: %w", err)
	}
	if err := textTmpl.Execute(&t, d); err != nil {
		return "", "", fmt.Errorf("text: %w", err)
	}
	return h.String(), t.String(), nil
}

func formatSeats(seats []string) string { return strings.Join(seats, ", ") }
