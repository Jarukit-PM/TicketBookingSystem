package catalog

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionCinemas   = "cinemas"
	CollectionScreens   = "screens"
	CollectionMovies    = "movies"
	CollectionShowtimes = "showtimes"
)

const (
	MovieStatusNowShowing = "NOW_SHOWING"
	MovieStatusComingSoon = "COMING_SOON"
	MovieStatusArchived   = "ARCHIVED"
)

const (
	ShowtimeStatusOpen      = "OPEN"
	ShowtimeStatusCancelled = "CANCELLED"
)

const (
	SeatTypeStandard   = "standard"
	SeatTypeVIP        = "vip"
	SeatTypeWheelchair = "wheelchair"
	SeatTypeBlocked    = "blocked"
)

// Cinema is a venue with timezone for showtime cutoff rules.
type Cinema struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name     string             `bson:"name" json:"name"`
	Address  string             `bson:"address" json:"address"`
	Timezone string             `bson:"timezone" json:"timezone"`
}

// LayoutSeat is one seat in a screen layout template.
type LayoutSeat struct {
	SeatID string `bson:"seatId" json:"seatId"`
	Row    int    `bson:"row" json:"row"`
	Col    int    `bson:"col" json:"col"`
	Type   string `bson:"type" json:"type"`
}

// ScreenLayout holds the seat map for a hall.
type ScreenLayout struct {
	Seats []LayoutSeat `bson:"seats" json:"seats"`
}

// Screen is a physical hall belonging to one cinema.
type Screen struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CinemaID primitive.ObjectID `bson:"cinemaId" json:"cinemaId"`
	Name     string             `bson:"name" json:"name"`
	Layout   ScreenLayout       `bson:"layout" json:"layout"`
}

// Movie is a global catalog entry shared across cinemas.
type Movie struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title       string             `bson:"title" json:"title"`
	PosterURL   string             `bson:"posterUrl" json:"posterUrl"`
	DurationMin int                `bson:"durationMin" json:"durationMin"`
	Rating      string             `bson:"rating" json:"rating"`
	Synopsis    string             `bson:"synopsis" json:"synopsis"`
	Status      string             `bson:"status" json:"status"`
}

// PriceTiers maps seat layout types to prices in minor units (cents).
type PriceTiers struct {
	Standard   int64 `bson:"standard" json:"standard"`
	VIP        int64 `bson:"vip" json:"vip"`
	Wheelchair int64 `bson:"wheelchair" json:"wheelchair"`
}

// Showtime is a scheduled screening on a screen.
type Showtime struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	MovieID    primitive.ObjectID `bson:"movieId" json:"movieId"`
	ScreenID   primitive.ObjectID `bson:"screenId" json:"screenId"`
	StartsAt   time.Time          `bson:"startsAt" json:"startsAt"`
	PriceTiers PriceTiers         `bson:"priceTiers" json:"priceTiers"`
	Status     string             `bson:"status" json:"status"`
}
