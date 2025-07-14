package mysql

import (
	"context"
	"github.com/1612421/cinema-booking/internal/entity"
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
