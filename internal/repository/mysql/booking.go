package mysql

import (
	"context"
	"github.com/1612421/cinema-booking/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type BookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) *BookingRepository {
	return &BookingRepository{
		db: db,
	}
}

func (b *BookingRepository) Create(ctx context.Context, booking *entity.Booking, seatIDs []uuid.UUID, callback func() error) (*entity.Booking, []*entity.BookingSeat, error) {
	bookingSeats := make([]*entity.BookingSeat, len(seatIDs))
	for i, seatID := range seatIDs {
		bookingSeats[i] = &entity.BookingSeat{
			ID:         uuid.New(),
			BookingID:  booking.ID,
			ShowtimeID: booking.ShowtimeID,
			SeatID:     seatID,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
	}

	err := b.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Insert Booking
		if err := tx.Create(&booking).Error; err != nil {
			return err
		}

		// 2. Insert booking_seats
		if err := tx.Create(&bookingSeats).Error; err != nil {
			return err
		}

		callbackErr := callback()
		if callbackErr != nil {
			return callbackErr
		}

		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	return booking, bookingSeats, nil
}

func (b *BookingRepository) Delete(ctx context.Context, booking *entity.Booking) (success bool, err error) {
	result := b.db.WithContext(ctx).Delete(booking)
	if result.Error != nil {
		return false, result.Error
	}

	return result.RowsAffected > 0, nil
}
