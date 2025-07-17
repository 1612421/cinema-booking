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

func (b *BookingService) HoldSeat(ctx context.Context, dto *HoldSeatDTO) (int64, error) {
	logger := log.For(ctx)

	qty, err := b.seatCache.HoldSeat(ctx, redis.HoldSeatCacheDTO{
		SeatID:     dto.SeatID,
		UserID:     dto.UserID,
		ShowtimeID: dto.ShowtimeID,
	})
	if err != nil {
		logger.Error("error hold seat", zap.Error(err))
		return 0, err
	}

	logger.Info("hold seat", zap.Reflect("seat", dto))

	return qty, nil
}
