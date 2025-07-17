package showtimeservice

import (
	"context"
	"github.com/1612421/cinema-booking/config"
	"github.com/1612421/cinema-booking/internal/entity"
	"github.com/google/uuid"
)

type IShowtimeRepository interface {
	Create(ctx context.Context, showtime *entity.Showtime) (*entity.Showtime, error)
	GetShowtimesByMovieID(ctx context.Context, movieID uuid.UUID) ([]*entity.Showtime, error)
	GetRandomShowtime(ctx context.Context) (*entity.Showtime, error)
}

type ShowtimeService struct {
	cfg          *config.Config
	showtimeRepo IShowtimeRepository
}

func NewShowtimeService(cfg *config.Config, showtimeRepo IShowtimeRepository) *ShowtimeService {
	return &ShowtimeService{
		cfg:          cfg,
		showtimeRepo: showtimeRepo,
	}
}

func (s *ShowtimeService) GetShowtimeRepo() IShowtimeRepository {
	return s.showtimeRepo
}
