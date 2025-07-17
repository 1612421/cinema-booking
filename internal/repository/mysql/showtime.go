package mysql

import (
	"context"
	"github.com/1612421/cinema-booking/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ShowtimeRepository struct {
	db *gorm.DB
}

func NewShowtimeRepository(db *gorm.DB) *ShowtimeRepository {
	return &ShowtimeRepository{
		db: db,
	}
}

func (s *ShowtimeRepository) Create(ctx context.Context, showtime *entity.Showtime) (*entity.Showtime, error) {
	tx := s.db.WithContext(ctx).Create(showtime)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return showtime, nil
}

func (s *ShowtimeRepository) GetShowtimesByMovieID(ctx context.Context, movieID uuid.UUID) ([]*entity.Showtime, error) {
	var showtimes []*entity.Showtime

	tx := s.db.WithContext(ctx).
		Model(&entity.Showtime{}).
		Where("movie_id = ?", movieID).
		Find(&showtimes)

	return showtimes, tx.Error
}

func (s *ShowtimeRepository) GetRandomShowtime(ctx context.Context) (*entity.Showtime, error) {
	var showtime *entity.Showtime

	tx := s.db.WithContext(ctx).
		Model(&entity.Showtime{}).
		Order("RAND()").
		First(&showtime)

	return showtime, tx.Error
}
