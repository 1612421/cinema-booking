package bookingservice

import (
	"context"
	"github.com/1612421/cinema-booking/internal/repository/redis"
	"github.com/1612421/cinema-booking/pkg/go-kit/log"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ReleaseSeatDTO struct {
	ShowtimeID uuid.UUID `json:"showtime_id"`
	UserID     uuid.UUID `json:"user_id"`
	SeatID     uuid.UUID `json:"seat_id"`
}

func (b *BookingService) ReleaseSeat(ctx context.Context, dto *ReleaseSeatDTO) error {
	logger := log.For(ctx)

	// Temporary lock seats
	if err := b.seatCache.ReleaseSeat(ctx, redis.ReleaseSeatDTO{
		SeatID:     dto.SeatID,
		UserID:     dto.UserID,
		ShowtimeID: dto.ShowtimeID,
	}); err != nil {
		logger.Error("error release seat", zap.Error(err))
		return err
	}

	logger.Info("release seat", zap.Reflect("seat", dto))

	return nil
}
