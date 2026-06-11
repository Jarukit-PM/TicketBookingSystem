package catalog_test

import (
	"testing"
	"time"

	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/catalog"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestFilterNowShowing(t *testing.T) {
	t.Parallel()

	nowShowingID := primitive.NewObjectID()
	comingSoonID := primitive.NewObjectID()
	archivedID := primitive.NewObjectID()
	noShowtimeID := primitive.NewObjectID()

	movies := []catalog.Movie{
		{ID: nowShowingID, Title: "Now Playing", Status: catalog.MovieStatusNowShowing},
		{ID: comingSoonID, Title: "Coming Teaser", Status: catalog.MovieStatusComingSoon},
		{ID: archivedID, Title: "Old Film", Status: catalog.MovieStatusArchived},
		{ID: noShowtimeID, Title: "No Local Showtimes", Status: catalog.MovieStatusNowShowing},
	}

	withShowtimes := map[primitive.ObjectID]struct{}{
		nowShowingID: {},
		comingSoonID: {},
	}

	tests := []struct {
		name     string
		movies   []catalog.Movie
		ids      map[primitive.ObjectID]struct{}
		wantLen  int
		wantIDs  []primitive.ObjectID
	}{
		{
			name:    "includes non-archived movies with future showtimes",
			movies:  movies,
			ids:     withShowtimes,
			wantLen: 2,
			wantIDs: []primitive.ObjectID{nowShowingID, comingSoonID},
		},
		{
			name:    "excludes archived even with showtimes",
			movies:  movies,
			ids:     map[primitive.ObjectID]struct{}{archivedID: {}},
			wantLen: 0,
		},
		{
			name:    "excludes movie without future showtimes at cinema",
			movies:  movies,
			ids:     map[primitive.ObjectID]struct{}{nowShowingID: {}},
			wantLen: 1,
			wantIDs: []primitive.ObjectID{nowShowingID},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := catalog.FilterNowShowing(tt.movies, tt.ids)
			if len(got) != tt.wantLen {
				t.Fatalf("len = %d, want %d", len(got), tt.wantLen)
			}
			for i, id := range tt.wantIDs {
				if got[i].ID != id {
					t.Fatalf("got[%d].ID = %v, want %v", i, got[i].ID, id)
				}
			}
		})
	}
}

func TestFilterComingSoon(t *testing.T) {
	t.Parallel()

	comingID := primitive.NewObjectID()
	movies := []catalog.Movie{
		{ID: comingID, Title: "Teaser", Status: catalog.MovieStatusComingSoon},
		{ID: primitive.NewObjectID(), Title: "Playing", Status: catalog.MovieStatusNowShowing},
		{ID: primitive.NewObjectID(), Title: "Archived", Status: catalog.MovieStatusArchived},
	}

	got := catalog.FilterComingSoon(movies)
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
	if got[0].ID != comingID {
		t.Fatalf("got id %v, want %v", got[0].ID, comingID)
	}
}

func TestFilterFutureShowtimes(t *testing.T) {
	t.Parallel()

	cutoff := time.Date(2026, 6, 11, 18, 0, 0, 0, time.UTC)
	past := cutoff.Add(-time.Hour)
	future := cutoff.Add(time.Hour)

	showtimes := []catalog.Showtime{
		{StartsAt: past, Status: catalog.ShowtimeStatusOpen},
		{StartsAt: future, Status: catalog.ShowtimeStatusOpen},
		{StartsAt: future, Status: catalog.ShowtimeStatusCancelled},
	}

	got := catalog.FilterFutureShowtimes(showtimes, cutoff)
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
	if !got[0].StartsAt.Equal(future) {
		t.Fatalf("got startsAt %v, want %v", got[0].StartsAt, future)
	}
}

func TestFutureShowtimeMovieIDs(t *testing.T) {
	t.Parallel()

	movieID := primitive.NewObjectID()
	cutoff := time.Now().UTC()
	showtimes := []catalog.Showtime{
		{MovieID: movieID, StartsAt: cutoff.Add(time.Hour), Status: catalog.ShowtimeStatusOpen},
		{MovieID: primitive.NewObjectID(), StartsAt: cutoff.Add(-time.Hour), Status: catalog.ShowtimeStatusOpen},
	}

	ids := catalog.FutureShowtimeMovieIDs(showtimes, cutoff)
	if len(ids) != 1 {
		t.Fatalf("len = %d, want 1", len(ids))
	}
	if _, ok := ids[movieID]; !ok {
		t.Fatalf("expected movie %v in set", movieID)
	}
}

func TestParseBrowseTab(t *testing.T) {
	t.Parallel()

	tests := []struct {
		raw   string
		valid bool
		tab   catalog.BrowseTab
	}{
		{"now_showing", true, catalog.BrowseTabNowShowing},
		{"coming_soon", true, catalog.BrowseTabComingSoon},
		{"invalid", false, ""},
		{"", false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.raw, func(t *testing.T) {
			t.Parallel()
			tab, ok := catalog.ParseBrowseTab(tt.raw)
			if ok != tt.valid {
				t.Fatalf("valid = %v, want %v", ok, tt.valid)
			}
			if tab != tt.tab {
				t.Fatalf("tab = %q, want %q", tab, tt.tab)
			}
		})
	}
}
