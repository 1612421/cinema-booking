package seatservice

import (
	"context"
	"github.com/1612421/cinema-booking/internal/entity"
	"github.com/1612421/cinema-booking/pkg/go-kit/log"
	"go.uber.org/zap"
)

func (s *SeatService) Create(ctx context.Context, seat *entity.Seat) (*entity.Seat, error) {
	logger := log.For(ctx)
	logger.Info("Create seat", zap.Reflect("seat", seat))

	seat, err := s.seatRepo.Create(ctx, seat)
	if err != nil {
		logger.Error("Create seat failed", log.Error(err))
		return nil, err
	}

	logger.Info("Create seat successfully", zap.Reflect("seat", seat))

	return seat, nil
}
