package booking

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ConfirmedFilter narrows confirmed booking queries for admin search.
type ConfirmedFilter struct {
	BookingRef    string
	UserID        primitive.ObjectID
	ShowtimeID    primitive.ObjectID
	ShowtimeIDs   []primitive.ObjectID
	Locale        string
	ConfirmedFrom *time.Time
	ConfirmedTo   *time.Time // exclusive upper bound
}

// Matches reports whether a booking satisfies the filter.
func (f ConfirmedFilter) Matches(b Booking) bool {
	if b.Status != StatusConfirmed {
		return false
	}
	if ref := f.BookingRef; ref != "" && b.BookingRef != ref {
		return false
	}
	if !f.UserID.IsZero() && b.UserID != f.UserID {
		return false
	}
	if !f.ShowtimeID.IsZero() {
		if b.ShowtimeID != f.ShowtimeID {
			return false
		}
	} else if len(f.ShowtimeIDs) > 0 {
		found := false
		for _, id := range f.ShowtimeIDs {
			if b.ShowtimeID == id {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	if locale := f.Locale; locale != "" && ParseLocale(b.Locale) != locale {
		return false
	}
	if f.ConfirmedFrom != nil && b.ConfirmedAt.Before(*f.ConfirmedFrom) {
		return false
	}
	if f.ConfirmedTo != nil && !b.ConfirmedAt.Before(*f.ConfirmedTo) {
		return false
	}
	return true
}

func confirmedFilterBSON(f ConfirmedFilter) bson.M {
	filter := bson.M{"status": StatusConfirmed}
	if ref := f.BookingRef; ref != "" {
		filter["bookingRef"] = ref
	}
	if !f.UserID.IsZero() {
		filter["userId"] = f.UserID
	}
	if !f.ShowtimeID.IsZero() {
		filter["showtimeId"] = f.ShowtimeID
	} else if len(f.ShowtimeIDs) > 0 {
		filter["showtimeId"] = bson.M{"$in": f.ShowtimeIDs}
	}
	if locale := f.Locale; locale != "" {
		filter["locale"] = locale
	}
	if f.ConfirmedFrom != nil || f.ConfirmedTo != nil {
		rangeFilter := bson.M{}
		if f.ConfirmedFrom != nil {
			rangeFilter["$gte"] = *f.ConfirmedFrom
		}
		if f.ConfirmedTo != nil {
			rangeFilter["$lt"] = *f.ConfirmedTo
		}
		filter["confirmedAt"] = rangeFilter
	}
	return filter
}
