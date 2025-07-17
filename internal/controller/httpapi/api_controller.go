package httpapi

import (
	"context"
	"github.com/1612421/cinema-booking/internal/controller/httpapi/middleware"
	"github.com/1612421/cinema-booking/internal/entity"
	authservice "github.com/1612421/cinema-booking/internal/service/auth"
	bookingservice "github.com/1612421/cinema-booking/internal/service/booking"
	seatservice "github.com/1612421/cinema-booking/internal/service/seat"
	showtimeservice "github.com/1612421/cinema-booking/internal/service/showtime"
	userservice "github.com/1612421/cinema-booking/internal/service/user"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/swaggo/files"       // swagger embed files
	"github.com/swaggo/gin-swagger" // gin-swagger middleware
	"gorm.io/gorm"
	"net/http"
)

type IAuthService interface {
	Authenticate(req *http.Request) (*authservice.AuthPayload, error)
	CreateAccessToken(payload authservice.AuthPayload) (string, error)
	VerifyAccessToken(tokenString string) (*authservice.AuthPayload, error)
	AuthenticateSocket(ctx *gin.Context) (*authservice.AuthPayload, error)
}

type IMovieService interface {
	Create(ctx context.Context, movie *entity.Movie) (*entity.Movie, error)
	GetMovie(ctx context.Context, id uuid.UUID) (*entity.Movie, error)
	GetMovies(ctx context.Context) ([]*entity.Movie, error)
}

type IShowtimeService interface {
	Create(ctx context.Context, showtime *entity.Showtime) (*entity.Showtime, error)
	GetShowtimesByMovieID(ctx context.Context, movieID uuid.UUID) ([]*entity.Showtime, error)
	GetShowtimeRepo() showtimeservice.IShowtimeRepository
}

type ISeatService interface {
	Create(ctx context.Context, seat *entity.Seat) (*entity.Seat, error)
	GetSeatsByShowtimeID(ctx context.Context, showtimeID uuid.UUID) ([]*entity.SeatWithStatus, error)
	GetSeatRepo() seatservice.ISeatRepository
}

type IBookingService interface {
	HoldSeat(ctx context.Context, dto *bookingservice.HoldSeatDTO) (int64, error)
	ReleaseSeat(ctx context.Context, dto *bookingservice.ReleaseSeatDTO) error
	CreateBooking(ctx context.Context, dto *bookingservice.CreateBookingDTO) (*entity.Booking, []*entity.BookingSeat, error)
	GetSeatInterest(ctx context.Context, dto *bookingservice.SeatInterestDTO) (int64, error)
	Delete(ctx context.Context, bookingID uuid.UUID) (success bool, err error)
}

type IUserService interface {
	Create(ctx context.Context, user *entity.User) (*entity.User, error)
	GetByUsername(ctx context.Context, username string) (*entity.User, error)
	GetUserRepo() userservice.IUserRepository
}

type ISocketService interface {
	HandleWebsocket(authService IAuthService) func(ctx *gin.Context)
	BroadcastMessage(message any)
	BroadcastMessageExceptID(message any, exceptedID uuid.UUID)
}

type IGlobalViewerSeatBotRunner interface {
	StartIfNotRunning(
		showtimeRepo showtimeservice.IShowtimeRepository,
		seatRepo seatservice.ISeatRepository,
		userRepo userservice.IUserRepository,
		bookingService IBookingService,
		socketService ISocketService,
		exceptedUserID uuid.UUID,
	)
}

type Controller struct {
	db                  *gorm.DB
	authService         IAuthService
	movieService        IMovieService
	userService         IUserService
	showtimeService     IShowtimeService
	seatService         ISeatService
	bookingService      IBookingService
	socketService       ISocketService
	viewerSeatBotRunner IGlobalViewerSeatBotRunner
}

func NewController(
	authService IAuthService,
	movieService IMovieService,
	userService IUserService,
	showtimeService IShowtimeService,
	seatService ISeatService,
	bookingService IBookingService,
	socketService ISocketService,
	viewerSeatBotRunner IGlobalViewerSeatBotRunner,
) *Controller {
	return &Controller{
		authService:         authService,
		movieService:        movieService,
		userService:         userService,
		showtimeService:     showtimeService,
		seatService:         seatService,
		bookingService:      bookingService,
		socketService:       socketService,
		viewerSeatBotRunner: viewerSeatBotRunner,
	}
}

func (c *Controller) GetAuthService() IAuthService {
	return c.authService
}

func (c *Controller) SetupRouter(r *gin.Engine) {
	r.Group("/api/v1")
	apiV1 := r.Group("/api/v1")
	//apiV1.Use(middleware.Authenticate(c.authService))
	movieV1Group := apiV1.Group("/movies")
	{
		movieV1Group.GET("/", c.getMovies)
		movieV1Group.POST("/", c.CreateMovie)
		movieV1Group.GET("/:movie_id", c.getMovie)
		movieV1Group.GET("/:movie_id/showtimes", c.getShowtimes)
	}

	userV1Group := apiV1.Group("/users")
	{
		userV1Group.POST("/register", c.RegisterUser)
		userV1Group.POST("/login", c.Login)
	}

	showtimeV1Group := apiV1.Group("/showtimes")
	{
		showtimeV1Group.POST("/", c.CreateShowtime)
		showtimeV1Group.GET("/:showtime_id/seats", c.getSeats)
		showtimeV1Group.GET("/:showtime_id/seats/:seat_id/interest", c.getSeatInterest)
	}

	seatV1Group := apiV1.Group("/seats")
	{
		seatV1Group.POST("/", c.CreateSeat)
	}

	bookingV1Group := apiV1.Group("/bookings")
	{
		bookingV1Group.POST("/", middleware.Authenticate(c.authService), c.CreateBooking)
		bookingV1Group.POST("/hold-seat", middleware.Authenticate(c.authService), c.HoldSeat)
		bookingV1Group.POST("/release-seat", middleware.Authenticate(c.authService), c.ReleaseSeat)
	}

	r.GET("/ws", c.socketService.HandleWebsocket(c.authService))
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
