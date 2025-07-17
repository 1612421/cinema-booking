package bot

import (
	"context"
	"fmt"
	"github.com/1612421/cinema-booking/config"
	"github.com/1612421/cinema-booking/internal/controller/httpapi"
	"github.com/1612421/cinema-booking/internal/entity"
	bookingservice "github.com/1612421/cinema-booking/internal/service/booking"
	seatservice "github.com/1612421/cinema-booking/internal/service/seat"
	showtimeservice "github.com/1612421/cinema-booking/internal/service/showtime"
	userservice "github.com/1612421/cinema-booking/internal/service/user"
	"github.com/1612421/cinema-booking/pkg/go-kit/log"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"math/rand"
	"sync"
	"time"
)

type ViewerSeatBot struct {
	Enabled        bool
	Frequency      time.Duration
	BookingChange  float64
	Workers        int
	showtimeRepo   showtimeservice.IShowtimeRepository
	seatRepo       seatservice.ISeatRepository
	userRepo       userservice.IUserRepository
	bookingService httpapi.IBookingService
	socketService  httpapi.ISocketService
	exceptedUserID uuid.UUID
}

func NewViewSeatBot(
	showtimeRepo showtimeservice.IShowtimeRepository,
	seatRepo seatservice.ISeatRepository,
	userRepo userservice.IUserRepository,
	bookingService httpapi.IBookingService,
	socketService httpapi.ISocketService,
	exceptedUserID uuid.UUID,
) *ViewerSeatBot {
	configConfig := config.GetConfig()

	return &ViewerSeatBot{
		Enabled:        configConfig.Bot.IsEnabled,
		Frequency:      time.Duration(configConfig.Bot.Frequency) * time.Second,
		Workers:        configConfig.Bot.Workers,
		BookingChange:  configConfig.Bot.BookingChange, // Odds to booking
		showtimeRepo:   showtimeRepo,
		seatRepo:       seatRepo,
		bookingService: bookingService,
		socketService:  socketService,
		exceptedUserID: exceptedUserID,
		userRepo:       userRepo,
	}
}

func (v *ViewerSeatBot) Start() {
	if !v.Enabled {
		log.Bg().Info("ViewerSeatBot is disabled")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	v.StartWithContext(ctx)
	log.Bg().Info(fmt.Sprintf("ViewerSeatBot started with %d workers (freq %s)", v.Workers, v.Frequency))

	// Optional: Wait for expiration (useful if you want to log or clean)
	go func() {
		<-ctx.Done()
		log.Bg().Info("ViewerSeatBot expired after 10 minutes")
		cancel()
	}()
}

func (v *ViewerSeatBot) StartWithContext(ctx context.Context) {
	var wg sync.WaitGroup
	users, err := v.getUsers()
	if err != nil {
		log.Bg().Error("ViewerSeatBot: failed to get users:", log.Error(err))
		return
	}

	for _, user := range users {
		wg.Add(1)
		go func(userID uuid.UUID) {
			defer wg.Done()
			v.workerLoop(ctx, userID)
		}(user.ID)
	}

	wg.Wait()
}

func (v *ViewerSeatBot) runOnce(userID uuid.UUID) {
	// Get random showtime
	showtimeID, err := v.getRandomShowtimeID()
	if err != nil {
		log.Bg().Error("ViewerSeatBot: failed to get showtime:", log.Error(err))
		return
	}

	// Get random seat
	seatID, err := v.getRandomSeatIDByShowtime(showtimeID)
	if err != nil {
		log.Bg().Error("ViewerSeatBot: failed to get seat:", log.Error(err))
		return
	}

	ctx := context.Background()

	// Hold seat
	holdSeatDTO := &bookingservice.HoldSeatDTO{
		SeatID:     seatID,
		ShowtimeID: showtimeID,
		UserID:     userID,
	}
	_, err = v.bookingService.HoldSeat(ctx, holdSeatDTO)
	if err != nil {
		log.Bg().Error("ViewerSeatBot: failed to hold seat:", log.Error(err))
		return
	}

	log.Bg().Info("ViewerSeatBot: hold seat successfully", zap.Reflect("holdSeatDTO", holdSeatDTO))

	//region Book seat
	chance := rand.Float64()
	if chance < v.BookingChange {
		// Delay 2–5 seconds
		log.Bg().Info("ViewerSeatBot: waiting for booking seat...")
		time.Sleep(time.Duration(rand.Intn(2000)+5000) * time.Millisecond)
		createBookingDTO := &bookingservice.CreateBookingDTO{
			UserID:     userID,
			ShowtimeID: showtimeID,
			SeatIDs:    []uuid.UUID{seatID},
		}
		booking, seats, bookingErr := v.bookSeats(ctx, createBookingDTO)
		if bookingErr != nil {
			log.Bg().Error("ViewerSeatBot: failed to book", zap.Error(err))
			return
		}

		log.Bg().Info("ViewerSeatBot: booked seat", zap.Reflect("booking", booking), zap.Reflect("seats", seats))

		// Delete booking after 1 minute
		go func(bookingID uuid.UUID) {
			time.Sleep(2 * time.Minute)
			_, delErr := v.bookingService.Delete(ctx, bookingID)
			if delErr != nil {
				log.Bg().Warn("ViewerSeatBot: failed to cancel booking", zap.Error(err))
				return
			}
			log.Bg().Info("ViewerSeatBot: canceled booking", zap.String("booking_id", bookingID.String()))
		}(booking.ID)
	}
	//endregion
}

func (v *ViewerSeatBot) workerLoop(ctx context.Context, userId uuid.UUID) {
	ticker := time.NewTicker(v.Frequency)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Bg().Info(fmt.Sprintf("ViewerSeatBot worker of user %s stopped", userId))
			return
		case <-ticker.C:
			v.runOnce(userId)
		}
	}
}

func (v *ViewerSeatBot) getRandomShowtimeID() (uuid.UUID, error) {
	showtime, err := v.showtimeRepo.GetRandomShowtime(context.Background())
	if err != nil {
		return uuid.Nil, err
	}

	return showtime.ID, nil
}

func (v *ViewerSeatBot) getRandomSeatIDByShowtime(showtimeID uuid.UUID) (uuid.UUID, error) {
	seat, err := v.seatRepo.GetRandomSeatByShowtimeID(context.Background(), showtimeID)
	if err != nil {
		return uuid.Nil, err
	}

	return seat.ID, err
}

func (v *ViewerSeatBot) getUsers() ([]*entity.User, error) {
	users, err := v.userRepo.GetUsersExceptID(context.Background(), v.Workers, v.exceptedUserID)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (v *ViewerSeatBot) bookSeats(ctx context.Context, dto *bookingservice.CreateBookingDTO) (*entity.Booking, []*entity.BookingSeat, error) {
	booking, seats, err := v.bookingService.CreateBooking(ctx, dto)
	if err != nil {
		return nil, nil, err
	}

	v.socketService.BroadcastMessageExceptID(gin.H{
		"event": "booking_created",
		"data": gin.H{
			"showtime_id": dto.ShowtimeID,
			"seat_ids":    dto.SeatIDs,
		},
	}, dto.UserID)

	return booking, seats, nil
}
