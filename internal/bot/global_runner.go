package bot

import (
	"context"
	"github.com/1612421/cinema-booking/internal/controller/httpapi"
	seatservice "github.com/1612421/cinema-booking/internal/service/seat"
	showtimeservice "github.com/1612421/cinema-booking/internal/service/showtime"
	userservice "github.com/1612421/cinema-booking/internal/service/user"
	"github.com/1612421/cinema-booking/pkg/go-kit/log"
	"github.com/google/uuid"
	"sync"
	"time"
)

type GlobalViewerSeatBotRunner struct {
	isRunning bool
	mu        sync.Mutex
}

func NewGlobalViewerSeatBotRunner() *GlobalViewerSeatBotRunner {
	return &GlobalViewerSeatBotRunner{}
}

func (g *GlobalViewerSeatBotRunner) StartIfNotRunning(
	showtimeRepo showtimeservice.IShowtimeRepository,
	seatRepo seatservice.ISeatRepository,
	userRepo userservice.IUserRepository,
	bookingService httpapi.IBookingService,
	socketService httpapi.ISocketService,
	exceptedUserID uuid.UUID,
) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.isRunning {
		log.Bg().Info("Global ViewerSeatBot already running. Skipping start.")
		return
	}

	log.Bg().Info("Starting global ViewerSeatBot for 10 minutes.")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	g.isRunning = true
	bot := NewViewSeatBot(showtimeRepo, seatRepo, userRepo, bookingService, socketService, exceptedUserID)

	go func() {
		defer func() {
			cancel()
			g.mu.Lock()
			g.isRunning = false
			g.mu.Unlock()
			log.Bg().Info("Global ViewerBot stopped after 10 minutes.")
		}()

		bot.StartWithContext(ctx)
	}()
}
