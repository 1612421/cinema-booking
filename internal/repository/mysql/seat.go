package mysql

import (
	"context"
	"github.com/1612421/cinema-booking/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SeatRepository struct {
	db *gorm.DB
}

func NewSeatRepository(db *gorm.DB) *SeatRepository {
	return &SeatRepository{
		db: db,
	}
}

func (s *SeatRepository) Create(ctx context.Context, seat *entity.Seat) (*entity.Seat, error) {
	tx := s.db.WithContext(ctx).Create(seat)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return seat, nil
}

func (s *SeatRepository) GetSeatsByShowtimeID(ctx context.Context, showtimeID uuid.UUID) ([]*entity.SeatWithStatus, error) {
	var seats []*entity.SeatWithStatus

	tx := s.db.WithContext(ctx).
		Table("seats s").
		Select("s.*, (bs.id IS NULL) is_available").
		Joins("JOIN showtime st ON st.screen_id = s.screen_id").
		Joins("LEFT JOIN booking_seats bs ON bs.seat_id = s.id AND bs.showtime_id = st.id").
		Where("st.id = ?", showtimeID).
		Order("s.seat_row ASC, s.number ASC").
		Find(&seats)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return seats, nil
}

func (s *SeatRepository) GetRandomSeatByShowtimeID(ctx context.Context, showtimeID uuid.UUID) (*entity.Seat, error) {
	var seat *entity.Seat
	tx := s.db.WithContext(ctx).
		Table("seats s").
		Select("s.*").
		Joins("JOIN showtime st ON st.screen_id = s.screen_id").
		Joins("LEFT JOIN booking_seats bs ON bs.seat_id = s.id AND bs.showtime_id = st.id").
		Where("st.id = ? AND bs.id IS NULL", showtimeID).
		Order("RAND()").
		First(&seat)

	return seat, tx.Error
}
