package movieservice

import (
	"context"
	"github.com/1612421/cinema-booking/internal/entity"
	"github.com/1612421/cinema-booking/pkg/go-kit/log"
	"go.uber.org/zap"
)

func (m *MovieService) Create(ctx context.Context, movie *entity.Movie) (*entity.Movie, error) {
	logger := log.For(ctx)
	logger.Info("Create movie", zap.Reflect("movie", movie))

	movie, err := m.movieRepo.Create(ctx, movie)
	if err != nil {
		logger.Error("Create movie failed", log.Error(err))
		return nil, err
	}

	logger.Info("Create movie successfully", zap.Reflect("movie", movie))

	return movie, nil
}
