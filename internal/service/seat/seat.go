package seatservice

import (
	"context"
	"github.com/1612421/cinema-booking/config"
	"github.com/1612421/cinema-booking/internal/entity"
	"github.com/google/uuid"
)

type ISeatRepository interface {
	Create(ctx context.Context, seat *entity.Seat) (*entity.Seat, error)
	GetSeatsByShowtimeID(ctx context.Context, showtimeID uuid.UUID) ([]*entity.SeatWithStatus, error)
	GetRandomSeatByShowtimeID(ctx context.Context, showtimeID uuid.UUID) (*entity.Seat, error)
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

func (s *SeatService) GetSeatRepo() ISeatRepository {
	return s.seatRepo
}
