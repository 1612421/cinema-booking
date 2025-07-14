package mysql

import (
	"context"
	"github.com/1612421/cinema-booking/internal/entity"
	"gorm.io/gorm"
)

type BookingSeatRepository struct {
	db *gorm.DB
}

func NewBookingSeatRepository(db *gorm.DB) *BookingRepository {
	return &BookingRepository{
		db: db,
	}
}

func (bs *BookingSeatRepository) CreateBulk(ctx context.Context, bookingSeats []*entity.BookingSeat) ([]*entity.BookingSeat, error) {
	tx := bs.db.WithContext(ctx).Create(bookingSeats)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return bookingSeats, nil
}
