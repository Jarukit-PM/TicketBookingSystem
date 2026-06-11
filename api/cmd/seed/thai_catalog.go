package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/catalog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func seedThaiCatalog(ctx context.Context, repos catalog.MongoRepositories) error {
	bangkok, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		return fmt.Errorf("load Asia/Bangkok: %w", err)
	}

	movies := make([]*catalog.Movie, 0, len(thaiMovies))
	for _, m := range thaiMovies {
		movie := &catalog.Movie{
			Title:       m.Title,
			PosterURL:   m.PosterURL,
			DurationMin: m.DurationMin,
			Rating:      m.Rating,
			Synopsis:    m.Synopsis,
			Status:      m.Status,
		}
		if err := repos.Movies.InsertMovie(ctx, movie); err != nil {
			return fmt.Errorf("insert movie %q: %w", m.Title, err)
		}
		movies = append(movies, movie)
	}

	nowShowing := filterMoviesByStatus(movies, catalog.MovieStatusNowShowing)
	if len(nowShowing) == 0 {
		return fmt.Errorf("no now-showing movies for showtime generation")
	}

	screenCount := 0
	showtimeCount := 0
	dayStart := time.Now().In(bangkok).Truncate(24 * time.Hour)

	for _, c := range thaiCinemas {
		cinema := &catalog.Cinema{
			Name:     c.Name,
			Address:  c.Address,
			Timezone: "Asia/Bangkok",
		}
		if err := repos.Cinemas.InsertCinema(ctx, cinema); err != nil {
			return fmt.Errorf("insert cinema %q: %w", c.Name, err)
		}

		for hallIdx, hallName := range c.Screens {
			screen := &catalog.Screen{
				CinemaID: cinema.ID,
				Name:     hallName,
				Layout:   thaiHallLayout(c.Name, hallName, hallIdx),
			}
			if err := repos.Screens.InsertScreen(ctx, screen); err != nil {
				return fmt.Errorf("insert screen %q: %w", hallName, err)
			}
			screenCount++

			for day := 0; day < thaiShowtimeDays; day++ {
				for slotIdx, hour := range dailyShowtimeHours {
					movie := nowShowing[(day+hallIdx+slotIdx+screenCount)%len(nowShowing)]
					startsLocal := dayStart.AddDate(0, 0, day).Add(time.Duration(hour) * time.Hour)
					showtime := &catalog.Showtime{
						MovieID:    movie.ID,
						ScreenID:   screen.ID,
						StartsAt:   startsLocal.UTC(),
						PriceTiers: thaiPriceTiers,
						Status:     catalog.ShowtimeStatusOpen,
					}
					if err := repos.Showtimes.InsertShowtime(ctx, showtime); err != nil {
						return fmt.Errorf("insert showtime: %w", err)
					}
					showtimeCount++
				}
			}
		}
	}

	log.Printf(
		"seed: created %d Major/SF cinemas, %d screens, %d movies, %d showtimes (%d days)",
		len(thaiCinemas),
		screenCount,
		len(movies),
		showtimeCount,
		thaiShowtimeDays,
	)
	return nil
}

func filterMoviesByStatus(movies []*catalog.Movie, status string) []*catalog.Movie {
	out := make([]*catalog.Movie, 0)
	for _, m := range movies {
		if m.Status == status {
			out = append(out, m)
		}
	}
	return out
}

func clearCatalog(ctx context.Context, database *mongo.Database) error {
	collections := []string{
		catalog.CollectionShowtimes,
		catalog.CollectionScreens,
		catalog.CollectionMovies,
		catalog.CollectionCinemas,
	}
	for _, name := range collections {
		if _, err := database.Collection(name).DeleteMany(ctx, bson.M{}); err != nil {
			return fmt.Errorf("clear %s: %w", name, err)
		}
	}
	log.Println("seed: cleared existing catalog collections")
	return nil
}

// thaiHallLayout builds a realistic auditorium map (8 rows × 10 seats).
func thaiHallLayout(cinemaName, hallName string, hallIdx int) catalog.ScreenLayout {
	const rows = 8
	const cols = 10
	prefix := fmt.Sprintf("%c%d", 'A'+hallIdx, hallIdx+1)

	seats := make([]catalog.LayoutSeat, 0, rows*cols)
	for row := 1; row <= rows; row++ {
		for col := 1; col <= cols; col++ {
			seatType := catalog.SeatTypeStandard
			switch {
			case row <= 2 && col >= 4 && col <= 7:
				seatType = catalog.SeatTypeVIP
			case row == rows && (col == 1 || col == cols):
				seatType = catalog.SeatTypeWheelchair
			case row == 4 && col == 5:
				seatType = catalog.SeatTypeBlocked
			case row == 5 && col == 5:
				seatType = catalog.SeatTypeBlocked
			}

			rowLabel := string(rune('A' + row - 1))
			seats = append(seats, catalog.LayoutSeat{
				SeatID: fmt.Sprintf("%s-%s%d", prefix, rowLabel, col),
				Row:    row,
				Col:    col,
				Type:   seatType,
			})
		}
	}
	_ = cinemaName
	_ = hallName
	return catalog.ScreenLayout{Seats: seats}
}
