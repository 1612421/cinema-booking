package movieservice

import (
	"context"
	"github.com/1612421/cinema-booking/internal/entity"
	"github.com/1612421/cinema-booking/pkg/go-kit/log"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (m *MovieService) GetMovie(ctx context.Context, id uuid.UUID) (*entity.Movie, error) {
	logger := log.For(ctx)

	logger.Info("Get movie", zap.String("id", id.String()))

	result, err := m.movieRepo.Get(ctx, &entity.Movie{ID: id})
	if err != nil {
		logger.Error("Failed to get movie", zap.Error(err))
		return nil, err
	}

	logger.Info("Get movie success", zap.Reflect("result", result))

	return result, nil
}
