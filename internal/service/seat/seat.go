package seatservice

import (
	"context"
	"github.com/1612421/cinema-booking/config"
	"github.com/1612421/cinema-booking/internal/entity"
)

type ISeatRepository interface {
	Create(ctx context.Context, seat *entity.Seat) (*entity.Seat, error)
}

type SeatService struct {
	cfg      *config.Config
	seatRepo ISeatRepository
}

func NewSeatService(cfg *config.Config, seatRepo ISeatRepository) *SeatService {
	return &SeatService{
		cfg:      cfg,
		seatRepo: seatRepo,
	}
}
