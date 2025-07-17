package bookingservice

import (
	"context"
	"github.com/1612421/cinema-booking/config"
	"github.com/1612421/cinema-booking/internal/entity"
	"github.com/1612421/cinema-booking/internal/repository/redis"
	"github.com/google/uuid"
)

type IBookingRepository interface {
	Create(ctx context.Context, booking *entity.Booking, seatIDs []uuid.UUID, callback func() error) (*entity.Booking, []*entity.BookingSeat, error)
	Delete(ctx context.Context, booking *entity.Booking) (success bool, err error)
}

type BookingService struct {
	cfg         *config.Config
	bookingRepo IBookingRepository
	seatCache   redis.ISeatCache
}

func NewBookingService(cfg *config.Config, bookingRepo IBookingRepository, seatCache redis.ISeatCache) *BookingService {
	return &BookingService{
		cfg:         cfg,
		bookingRepo: bookingRepo,
		seatCache:   seatCache,
	}
}
