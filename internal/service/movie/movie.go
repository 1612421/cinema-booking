package movieservice

import (
	"context"
	"github.com/1612421/cinema-booking/config"
	"github.com/1612421/cinema-booking/internal/entity"
)

type IMovieRepository interface {
	Create(ctx context.Context, movie *entity.Movie) (*entity.Movie, error)
	Get(ctx context.Context, movie *entity.Movie) (*entity.Movie, error)
	GetMovies(ctx context.Context) ([]*entity.Movie, error)
}

type MovieService struct {
	cfg       *config.Config
	movieRepo IMovieRepository
}

func NewMovieService(cfg *config.Config, movieRepo IMovieRepository) *MovieService {
	return &MovieService{
		cfg:       cfg,
		movieRepo: movieRepo,
	}
}
