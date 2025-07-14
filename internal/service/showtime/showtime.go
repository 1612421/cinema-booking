package showtimeservice

import (
	"context"
	"github.com/1612421/cinema-booking/config"
	"github.com/1612421/cinema-booking/internal/entity"
)

type IShowtimeRepository interface {
	Create(ctx context.Context, showtime *entity.Showtime) (*entity.Showtime, error)
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
