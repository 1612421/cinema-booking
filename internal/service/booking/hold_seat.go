package bookingservice

import (
	"context"
	"github.com/1612421/cinema-booking/internal/repository/redis"
	"github.com/1612421/cinema-booking/pkg/go-kit/log"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type HoldSeatDTO struct {
	ShowtimeID uuid.UUID `json:"showtime_id"`
	UserID     uuid.UUID `json:"user_id"`
	SeatID     uuid.UUID `json:"seat_id,max=255"`
}

func (b *BookingService) HoldSeat(ctx context.Context, dto *HoldSeatDTO) error {
	logger := log.For(ctx)

	// Temporary lock seats
	if err := b.seatCache.HoldSeat(ctx, redis.HoldSeatCacheDTO{
		SeatID:     dto.SeatID,
		UserID:     dto.UserID,
		ShowtimeID: dto.ShowtimeID,
	}); err != nil {
		return err
	}

	logger.Info("hold seat", zap.Reflect("seat", dto))

	return nil
}
