package seatservice

import (
	"context"
	"github.com/1612421/cinema-booking/internal/entity"
	"github.com/1612421/cinema-booking/pkg/go-kit/log"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (s *SeatService) GetSeatsByShowtimeID(ctx context.Context, showtimeID uuid.UUID) ([]*entity.SeatWithStatus, error) {
	logger := log.For(ctx)
	logger.Info("Get seats by showtimeID", zap.Reflect("showtimeID", showtimeID.String()))

	seats, err := s.seatRepo.GetSeatsByShowtimeID(ctx, showtimeID)
	if err != nil {
		logger.Error("Get seats by showtimeID failed", log.Error(err))
		return nil, err
	}

	logger.Info("Get seats by showtimeID successfully", zap.Reflect("seat", seats))

	return seats, nil
}
