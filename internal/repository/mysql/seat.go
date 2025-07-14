package mysql

import (
	"context"
	"github.com/1612421/cinema-booking/internal/entity"
	"gorm.io/gorm"
)

type SeatRepository struct {
	db *gorm.DB
}

func NewSeatRepository(db *gorm.DB) *SeatRepository {
	return &SeatRepository{
		db: db,
	}
}

func (s *SeatRepository) Create(ctx context.Context, seat *entity.Seat) (*entity.Seat, error) {
	tx := s.db.WithContext(ctx).Create(seat)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return seat, nil
}
