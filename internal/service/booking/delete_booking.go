package bookingservice

import (
	"context"
	"github.com/1612421/cinema-booking/internal/entity"
	"github.com/google/uuid"
)

func (b *BookingService) Delete(ctx context.Context, bookingID uuid.UUID) (success bool, err error) {
	booking := &entity.Booking{
		ID: bookingID,
	}

	return b.bookingRepo.Delete(ctx, booking)
}
