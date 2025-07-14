package showtimeservice

import (
	"context"
	"github.com/1612421/cinema-booking/internal/entity"
	"github.com/1612421/cinema-booking/pkg/go-kit/log"
	"go.uber.org/zap"
)

func (s *ShowtimeService) Create(ctx context.Context, showtime *entity.Showtime) (*entity.Showtime, error) {
	logger := log.For(ctx)
	logger.Info("Create showtime", zap.Reflect("showtime", showtime))

	movie, err := s.showtimeRepo.Create(ctx, showtime)
	if err != nil {
		logger.Error("Create showtime failed", log.Error(err))
		return nil, err
	}

	logger.Info("Create showtime successfully", zap.Reflect("showtime", movie))

	return movie, nil
}
