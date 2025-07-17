package bookingservice

import (
	"context"
	"github.com/1612421/cinema-booking/internal/repository/redis"
	"github.com/1612421/cinema-booking/pkg/go-kit/log"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type SeatInterestDTO struct {
	ShowtimeID uuid.UUID `json:"showtime_id"`
	SeatID     uuid.UUID `json:"seat_id"`
}

func (b *BookingService) GetSeatInterest(ctx context.Context, dto *SeatInterestDTO) (int64, error) {
	logger := log.For(ctx)

	qty, err := b.seatCache.GetSeatInterest(ctx, redis.SeatInterestDTO{
		SeatID:     dto.SeatID,
		ShowtimeID: dto.ShowtimeID,
	})
	if err != nil {
		logger.Error("error get seat interest", zap.Error(err))
		return 0, err
	}

	logger.Info("get seat interest", zap.Int64("qty", qty))

	return qty, nil
}
