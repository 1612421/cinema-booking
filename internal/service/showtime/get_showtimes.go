package showtimeservice

import (
	"context"
	"github.com/1612421/cinema-booking/internal/entity"
	"github.com/1612421/cinema-booking/pkg/go-kit/log"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (s *ShowtimeService) GetShowtimesByMovieID(ctx context.Context, movieID uuid.UUID) ([]*entity.Showtime, error) {
	logger := log.For(ctx)
	logger.Info("Get showtimes by movieID", zap.String("movieID", movieID.String()))

	showtimes, err := s.showtimeRepo.GetShowtimesByMovieID(ctx, movieID)
	if err != nil {
		logger.Error("Get showtimes by movieID failed", log.Error(err))
		return nil, err
	}

	logger.Info("Get showtimes by movieID successfully", zap.Reflect("showtimes", showtimes))

	return showtimes, nil
}
