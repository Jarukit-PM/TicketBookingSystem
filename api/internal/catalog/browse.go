package catalog

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BrowseTab identifies the customer catalog browse mode.
type BrowseTab string

const (
	BrowseTabNowShowing BrowseTab = "now_showing"
	BrowseTabComingSoon BrowseTab = "coming_soon"
)

// ParseBrowseTab validates a tab query parameter.
func ParseBrowseTab(raw string) (BrowseTab, bool) {
	switch BrowseTab(raw) {
	case BrowseTabNowShowing, BrowseTabComingSoon:
		return BrowseTab(raw), true
	default:
		return "", false
	}
}

// FilterNowShowing returns movies that are not archived and have at least one
// future showtime at the selected cinema.
func FilterNowShowing(movies []Movie, futureShowtimeMovieIDs map[primitive.ObjectID]struct{}) []Movie {
	out := make([]Movie, 0, len(movies))
	for _, m := range movies {
		if m.Status == MovieStatusArchived {
			continue
		}
		if _, ok := futureShowtimeMovieIDs[m.ID]; !ok {
			continue
		}
		out = append(out, m)
	}
	return out
}

// FilterComingSoon returns COMING_SOON teasers (showtimes optional).
func FilterComingSoon(movies []Movie) []Movie {
	out := make([]Movie, 0)
	for _, m := range movies {
		if m.Status == MovieStatusComingSoon {
			out = append(out, m)
		}
	}
	return out
}

// FutureShowtimeMovieIDs builds a set of movie IDs from showtimes after cutoff.
func FutureShowtimeMovieIDs(showtimes []Showtime, after time.Time) map[primitive.ObjectID]struct{} {
	ids := make(map[primitive.ObjectID]struct{})
	for _, st := range showtimes {
		if st.StartsAt.After(after) && st.Status != ShowtimeStatusCancelled {
			ids[st.MovieID] = struct{}{}
		}
	}
	return ids
}

// FilterFutureShowtimes returns showtimes strictly after cutoff at an open cinema.
func FilterFutureShowtimes(showtimes []Showtime, after time.Time) []Showtime {
	out := make([]Showtime, 0, len(showtimes))
	for _, st := range showtimes {
		if st.StartsAt.After(after) && st.Status != ShowtimeStatusCancelled {
			out = append(out, st)
		}
	}
	return out
}
