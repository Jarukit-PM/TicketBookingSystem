package catalog

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/audit"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	AuditActionCreate = "create"
	AuditActionUpdate = "update"
	AuditActionDelete = "delete"
)

var (
	ErrNotFound       = errors.New("resource not found")
	ErrInvalidInput   = errors.New("invalid input")
	ErrInvalidStatus  = errors.New("invalid movie status")
	ErrInvalidSeat    = errors.New("invalid seat layout")
	ErrInvalidShowtime = errors.New("invalid showtime")
)

// AdminService handles admin catalog mutations with audit logging.
type AdminService struct {
	repos MongoRepositories
	audit audit.AuditRepository
}

// NewAdminService returns an admin catalog service.
func NewAdminService(repos MongoRepositories, auditRepo audit.AuditRepository) *AdminService {
	return &AdminService{repos: repos, audit: auditRepo}
}

func (s *AdminService) ListMovies(ctx context.Context) ([]Movie, error) {
	return s.repos.Movies.ListMovies(ctx)
}

func (s *AdminService) GetMovie(ctx context.Context, id primitive.ObjectID) (*Movie, error) {
	return s.findMovie(ctx, id)
}

func (s *AdminService) CreateMovie(ctx context.Context, actorID primitive.ObjectID, movie *Movie) (*Movie, error) {
	if err := validateMovie(movie); err != nil {
		return nil, err
	}
	if err := s.repos.Movies.InsertMovie(ctx, movie); err != nil {
		return nil, err
	}
	if err := s.logAudit(ctx, actorID, AuditActionCreate, "movie", movie.ID.Hex(), nil); err != nil {
		return nil, err
	}
	return movie, nil
}

func (s *AdminService) UpdateMovie(ctx context.Context, actorID, id primitive.ObjectID, movie *Movie) (*Movie, error) {
	existing, err := s.findMovie(ctx, id)
	if err != nil {
		return nil, err
	}
	movie.ID = existing.ID
	if err := validateMovie(movie); err != nil {
		return nil, err
	}
	if err := s.repos.Movies.UpdateMovie(ctx, movie); err != nil {
		return nil, err
	}
	if err := s.logAudit(ctx, actorID, AuditActionUpdate, "movie", movie.ID.Hex(), nil); err != nil {
		return nil, err
	}
	return movie, nil
}

func (s *AdminService) DeleteMovie(ctx context.Context, actorID, id primitive.ObjectID) error {
	if _, err := s.findMovie(ctx, id); err != nil {
		return err
	}
	if err := s.repos.Movies.DeleteMovie(ctx, id); err != nil {
		return err
	}
	return s.logAudit(ctx, actorID, AuditActionDelete, "movie", id.Hex(), nil)
}

func (s *AdminService) ListCinemas(ctx context.Context) ([]Cinema, error) {
	return s.repos.Cinemas.ListCinemas(ctx)
}

func (s *AdminService) GetCinema(ctx context.Context, id primitive.ObjectID) (*Cinema, error) {
	return s.findCinema(ctx, id)
}

func (s *AdminService) CreateCinema(ctx context.Context, actorID primitive.ObjectID, cinema *Cinema) (*Cinema, error) {
	if err := validateCinema(cinema); err != nil {
		return nil, err
	}
	if err := s.repos.Cinemas.InsertCinema(ctx, cinema); err != nil {
		return nil, err
	}
	if err := s.logAudit(ctx, actorID, AuditActionCreate, "cinema", cinema.ID.Hex(), nil); err != nil {
		return nil, err
	}
	return cinema, nil
}

func (s *AdminService) UpdateCinema(ctx context.Context, actorID, id primitive.ObjectID, cinema *Cinema) (*Cinema, error) {
	existing, err := s.findCinema(ctx, id)
	if err != nil {
		return nil, err
	}
	cinema.ID = existing.ID
	if err := validateCinema(cinema); err != nil {
		return nil, err
	}
	if err := s.repos.Cinemas.UpdateCinema(ctx, cinema); err != nil {
		return nil, err
	}
	if err := s.logAudit(ctx, actorID, AuditActionUpdate, "cinema", cinema.ID.Hex(), nil); err != nil {
		return nil, err
	}
	return cinema, nil
}

func (s *AdminService) DeleteCinema(ctx context.Context, actorID, id primitive.ObjectID) error {
	if _, err := s.findCinema(ctx, id); err != nil {
		return err
	}
	if err := s.repos.Cinemas.DeleteCinema(ctx, id); err != nil {
		return err
	}
	return s.logAudit(ctx, actorID, AuditActionDelete, "cinema", id.Hex(), nil)
}

func (s *AdminService) ListScreens(ctx context.Context, cinemaID *primitive.ObjectID) ([]Screen, error) {
	return s.repos.Screens.ListScreens(ctx, cinemaID)
}

func (s *AdminService) GetScreen(ctx context.Context, id primitive.ObjectID) (*Screen, error) {
	return s.findScreen(ctx, id)
}

func (s *AdminService) CreateScreen(ctx context.Context, actorID primitive.ObjectID, screen *Screen) (*Screen, error) {
	if err := validateScreen(ctx, s.repos, screen); err != nil {
		return nil, err
	}
	if err := s.repos.Screens.InsertScreen(ctx, screen); err != nil {
		return nil, err
	}
	if err := s.logAudit(ctx, actorID, AuditActionCreate, "screen", screen.ID.Hex(), nil); err != nil {
		return nil, err
	}
	return screen, nil
}

func (s *AdminService) UpdateScreen(ctx context.Context, actorID, id primitive.ObjectID, screen *Screen) (*Screen, error) {
	existing, err := s.findScreen(ctx, id)
	if err != nil {
		return nil, err
	}
	screen.ID = existing.ID
	if err := validateScreen(ctx, s.repos, screen); err != nil {
		return nil, err
	}
	if err := s.repos.Screens.UpdateScreen(ctx, screen); err != nil {
		return nil, err
	}
	if err := s.logAudit(ctx, actorID, AuditActionUpdate, "screen", screen.ID.Hex(), nil); err != nil {
		return nil, err
	}
	return screen, nil
}

func (s *AdminService) DeleteScreen(ctx context.Context, actorID, id primitive.ObjectID) error {
	if _, err := s.findScreen(ctx, id); err != nil {
		return err
	}
	if err := s.repos.Screens.DeleteScreen(ctx, id); err != nil {
		return err
	}
	return s.logAudit(ctx, actorID, AuditActionDelete, "screen", id.Hex(), nil)
}

func (s *AdminService) ListShowtimes(ctx context.Context, filter AdminShowtimeFilter) ([]Showtime, error) {
	return s.repos.Showtimes.ListAdminShowtimes(ctx, filter)
}

func (s *AdminService) GetShowtime(ctx context.Context, id primitive.ObjectID) (*Showtime, error) {
	return s.findShowtime(ctx, id)
}

func (s *AdminService) CreateShowtime(ctx context.Context, actorID primitive.ObjectID, showtime *Showtime) (*Showtime, error) {
	if err := validateShowtime(ctx, s.repos, showtime); err != nil {
		return nil, err
	}
	if showtime.Status == "" {
		showtime.Status = ShowtimeStatusOpen
	}
	if err := s.repos.Showtimes.InsertShowtime(ctx, showtime); err != nil {
		return nil, err
	}
	if err := s.logAudit(ctx, actorID, AuditActionCreate, "showtime", showtime.ID.Hex(), nil); err != nil {
		return nil, err
	}
	return showtime, nil
}

func (s *AdminService) UpdateShowtime(ctx context.Context, actorID, id primitive.ObjectID, showtime *Showtime) (*Showtime, error) {
	existing, err := s.findShowtime(ctx, id)
	if err != nil {
		return nil, err
	}
	showtime.ID = existing.ID
	if err := validateShowtime(ctx, s.repos, showtime); err != nil {
		return nil, err
	}
	if err := s.repos.Showtimes.UpdateShowtime(ctx, showtime); err != nil {
		return nil, err
	}
	if err := s.logAudit(ctx, actorID, AuditActionUpdate, "showtime", showtime.ID.Hex(), nil); err != nil {
		return nil, err
	}
	return showtime, nil
}

func (s *AdminService) DeleteShowtime(ctx context.Context, actorID, id primitive.ObjectID) error {
	if _, err := s.findShowtime(ctx, id); err != nil {
		return err
	}
	if err := s.repos.Showtimes.DeleteShowtime(ctx, id); err != nil {
		return err
	}
	return s.logAudit(ctx, actorID, AuditActionDelete, "showtime", id.Hex(), nil)
}

func (s *AdminService) findMovie(ctx context.Context, id primitive.ObjectID) (*Movie, error) {
	movie, err := s.repos.Movies.FindMovieByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if movie == nil {
		return nil, ErrNotFound
	}
	return movie, nil
}

func (s *AdminService) findCinema(ctx context.Context, id primitive.ObjectID) (*Cinema, error) {
	cinema, err := s.repos.Cinemas.FindCinemaByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if cinema == nil {
		return nil, ErrNotFound
	}
	return cinema, nil
}

func (s *AdminService) findScreen(ctx context.Context, id primitive.ObjectID) (*Screen, error) {
	screen, err := s.repos.Screens.FindScreenByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if screen == nil {
		return nil, ErrNotFound
	}
	return screen, nil
}

func (s *AdminService) findShowtime(ctx context.Context, id primitive.ObjectID) (*Showtime, error) {
	showtime, err := s.repos.Showtimes.FindShowtimeByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if showtime == nil {
		return nil, ErrNotFound
	}
	return showtime, nil
}

func (s *AdminService) logAudit(
	ctx context.Context,
	actorID primitive.ObjectID,
	action, entity, entityID string,
	meta map[string]any,
) error {
	return s.audit.InsertAuditLog(ctx, &audit.AuditLog{
		ActorID:   audit.ActorIDPtr(actorID),
		Action:    action,
		Entity:    entity,
		EntityID:  entityID,
		Meta:      meta,
		CreatedAt: time.Now().UTC(),
	})
}

func validateMovie(movie *Movie) error {
	if strings.TrimSpace(movie.Title) == "" {
		return fmt.Errorf("%w: title required", ErrInvalidInput)
	}
	switch movie.Status {
	case MovieStatusNowShowing, MovieStatusComingSoon, MovieStatusArchived, "":
		if movie.Status == "" {
			movie.Status = MovieStatusNowShowing
		}
	default:
		return ErrInvalidStatus
	}
	return nil
}

func validateCinema(cinema *Cinema) error {
	if strings.TrimSpace(cinema.Name) == "" {
		return fmt.Errorf("%w: name required", ErrInvalidInput)
	}
	if strings.TrimSpace(cinema.Timezone) == "" {
		return fmt.Errorf("%w: timezone required", ErrInvalidInput)
	}
	if _, err := time.LoadLocation(cinema.Timezone); err != nil {
		return fmt.Errorf("%w: invalid timezone", ErrInvalidInput)
	}
	return nil
}

func validateScreen(ctx context.Context, repos MongoRepositories, screen *Screen) error {
	if screen.CinemaID.IsZero() {
		return fmt.Errorf("%w: cinemaId required", ErrInvalidInput)
	}
	if strings.TrimSpace(screen.Name) == "" {
		return fmt.Errorf("%w: name required", ErrInvalidInput)
	}
	cinema, err := repos.Cinemas.FindCinemaByID(ctx, screen.CinemaID)
	if err != nil {
		return err
	}
	if cinema == nil {
		return fmt.Errorf("%w: cinema not found", ErrInvalidInput)
	}
	return validateLayout(screen.Layout)
}

func validateLayout(layout ScreenLayout) error {
	if len(layout.Seats) == 0 {
		return fmt.Errorf("%w: at least one seat required", ErrInvalidSeat)
	}
	seen := make(map[string]struct{}, len(layout.Seats))
	for _, seat := range layout.Seats {
		if strings.TrimSpace(seat.SeatID) == "" {
			return fmt.Errorf("%w: seatId required", ErrInvalidSeat)
		}
		if _, ok := seen[seat.SeatID]; ok {
			return fmt.Errorf("%w: duplicate seatId %s", ErrInvalidSeat, seat.SeatID)
		}
		seen[seat.SeatID] = struct{}{}
		switch seat.Type {
		case SeatTypeStandard, SeatTypeVIP, SeatTypeWheelchair, SeatTypeBlocked:
		default:
			return fmt.Errorf("%w: invalid seat type %s", ErrInvalidSeat, seat.Type)
		}
	}
	return nil
}

func validateShowtime(ctx context.Context, repos MongoRepositories, showtime *Showtime) error {
	if showtime.MovieID.IsZero() {
		return fmt.Errorf("%w: movieId required", ErrInvalidShowtime)
	}
	if showtime.ScreenID.IsZero() {
		return fmt.Errorf("%w: screenId required", ErrInvalidShowtime)
	}
	if showtime.StartsAt.IsZero() {
		return fmt.Errorf("%w: startsAt required", ErrInvalidShowtime)
	}
	movie, err := repos.Movies.FindMovieByID(ctx, showtime.MovieID)
	if err != nil {
		return err
	}
	if movie == nil {
		return fmt.Errorf("%w: movie not found", ErrInvalidShowtime)
	}
	screen, err := repos.Screens.FindScreenByID(ctx, showtime.ScreenID)
	if err != nil {
		return err
	}
	if screen == nil {
		return fmt.Errorf("%w: screen not found", ErrInvalidShowtime)
	}
	if showtime.Status == "" {
		showtime.Status = ShowtimeStatusOpen
	}
	switch showtime.Status {
	case ShowtimeStatusOpen, ShowtimeStatusCancelled:
	default:
		return fmt.Errorf("%w: invalid status", ErrInvalidShowtime)
	}
	return nil
}
