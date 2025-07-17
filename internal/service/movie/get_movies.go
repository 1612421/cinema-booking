package movieservice

import (
	"context"
	"github.com/1612421/cinema-booking/internal/entity"
	"github.com/1612421/cinema-booking/pkg/go-kit/log"
	"go.uber.org/zap"
)

func (m *MovieService) GetMovies(ctx context.Context) ([]*entity.Movie, error) {
	logger := log.For(ctx)

	logger.Info("Get movies")

	movies, err := m.movieRepo.GetMovies(ctx)
	if err != nil {
		logger.Error("Failed to get movies", log.Error(err))
		return nil, err
	}

	logger.Info("Get movies success", zap.Reflect("movies", movies))

	return movies, nil
}
