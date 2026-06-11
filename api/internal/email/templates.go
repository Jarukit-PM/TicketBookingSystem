package email

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"
	texttemplate "text/template"

	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/booking"
)

type confirmationData struct {
	BookingRef, MovieTitle, CinemaName, ScreenName, StartsAt, Seats, Total, TicketURL string
}

var (
	htmlTmplEN = template.Must(template.New("html-en").Parse(`<h1>Booking confirmed</h1><p>Ref: {{.BookingRef}}</p><p>{{.MovieTitle}} — {{.CinemaName}} / {{.ScreenName}}</p><p>{{.StartsAt}}</p><p>Seats: {{.Seats}} · {{.Total}} THB</p><p><a href="{{.TicketURL}}">View ticket</a></p>`))
	textTmplEN = texttemplate.Must(texttemplate.New("text-en").Parse("Booking {{.BookingRef}}\n{{.MovieTitle}}\n{{.TicketURL}}\n"))
	htmlTmplTH = template.Must(template.New("html-th").Parse(`<h1>ยืนยันการจองแล้ว</h1><p>เลขที่: {{.BookingRef}}</p><p>{{.MovieTitle}} — {{.CinemaName}} / {{.ScreenName}}</p><p>{{.StartsAt}}</p><p>ที่นั่ง: {{.Seats}} · {{.Total}} บาท</p><p><a href="{{.TicketURL}}">ดูตั๋ว</a></p>`))
	textTmplTH = texttemplate.Must(texttemplate.New("text-th").Parse("การจอง {{.BookingRef}}\n{{.MovieTitle}}\n{{.TicketURL}}\n"))
)

func renderConfirmation(locale string, d confirmationData) (string, string, error) {
	htmlTmpl, textTmpl := htmlTmplEN, textTmplEN
	if booking.ParseLocale(locale) == booking.LocaleTH {
		htmlTmpl, textTmpl = htmlTmplTH, textTmplTH
	}

	var h, t bytes.Buffer
	if err := htmlTmpl.Execute(&h, d); err != nil {
		return "", "", fmt.Errorf("html: %w", err)
	}
	if err := textTmpl.Execute(&t, d); err != nil {
		return "", "", fmt.Errorf("text: %w", err)
	}
	return h.String(), t.String(), nil
}

func confirmationSubject(locale, movieTitle string) string {
	if booking.ParseLocale(locale) == booking.LocaleTH {
		return "ตั๋วของคุณ — " + movieTitle
	}
	return "Your tickets — " + movieTitle
}

func formatSeats(seats []string) string { return strings.Join(seats, ", ") }

// formatTHBAmount converts satang (minor units) to a THB display string.
func formatTHBAmount(satang int64) string {
	baht := satang / 100
	remainder := satang % 100
	if remainder == 0 {
		return fmt.Sprintf("%d", baht)
	}
	return fmt.Sprintf("%d.%02d", baht, remainder)
}
