package mysql

import (
	"context"
	"github.com/1612421/cinema-booking/internal/entity"
	"gorm.io/gorm"
)

type MovieRepository struct {
	db *gorm.DB
}

func NewMovieRepository(db *gorm.DB) *MovieRepository {
	return &MovieRepository{
		db: db,
	}
}

func (m *MovieRepository) Create(ctx context.Context, movie *entity.Movie) (*entity.Movie, error) {
	tx := m.db.WithContext(ctx).Create(movie)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return movie, nil
}

func (m *MovieRepository) GetMovies(ctx context.Context) ([]*entity.Movie, error) {
	var movies []*entity.Movie

	tx := m.db.WithContext(ctx).
		Model(&entity.Movie{}).
		Find(&movies)

	return movies, tx.Error
}

func (m *MovieRepository) Get(ctx context.Context, movie *entity.Movie) (*entity.Movie, error) {
	tx := m.db.WithContext(ctx).First(movie)

	if tx.Error != nil {
		return nil, tx.Error
	}

	return movie, nil
}
