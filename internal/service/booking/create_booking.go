package bookingservice

import (
	"context"
	"github.com/1612421/cinema-booking/internal/entity"
	"github.com/1612421/cinema-booking/internal/repository/redis"
	"github.com/1612421/cinema-booking/pkg/go-kit/log"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"time"
)

type CreateBookingDTO struct {
	ShowtimeID uuid.UUID   `json:"showtime_id"`
	UserID     uuid.UUID   `json:"user_id"`
	SeatIDs    []uuid.UUID `json:"seat_ids"`
}

func (b *BookingService) CreateBooking(ctx context.Context, dto *CreateBookingDTO) (*entity.Booking, []*entity.BookingSeat, error) {
	logger := log.For(ctx)

	booking := &entity.Booking{
		ID:         uuid.New(),
		ShowtimeID: dto.ShowtimeID,
		UserId:     dto.UserID,
		Status:     entity.BookingStatusConfirmed,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	booking, bookingSeats, err := b.bookingRepo.Create(ctx, booking, dto.SeatIDs, func() error {
		return b.seatCache.HDelHoldSeatsAll(ctx, redis.ReleaseSeatsBulkDTO{
			SeatIDs:    dto.SeatIDs,
			ShowtimeID: dto.ShowtimeID,
		})
	})

	if err != nil {
		return nil, nil, err
	}

	logger.Info(
		"create booking successfully",
		zap.Reflect("entity", booking),
		zap.Reflect("bookingSeats", bookingSeats),
	)

	return booking, bookingSeats, nil
}
