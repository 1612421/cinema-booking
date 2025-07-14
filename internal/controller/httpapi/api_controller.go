package httpapi

import (
	"context"
	"github.com/1612421/cinema-booking/internal/controller/httpapi/middleware"
	"github.com/1612421/cinema-booking/internal/entity"
	authservice "github.com/1612421/cinema-booking/internal/service/auth"
	bookingservice "github.com/1612421/cinema-booking/internal/service/booking"
	"github.com/gin-gonic/gin"
	"net/http"
)

type IAuthService interface {
	Authenticate(req *http.Request) (*authservice.AuthPayload, error)
	CreateAccessToken(payload authservice.AuthPayload) (string, error)
	VerifyAccessToken(tokenString string) (*authservice.AuthPayload, error)
}

type IMovieService interface {
	Create(ctx context.Context, movie *entity.Movie) (*entity.Movie, error)
}

type IShowtimeService interface {
	Create(ctx context.Context, showtime *entity.Showtime) (*entity.Showtime, error)
}

type ISeatService interface {
	Create(ctx context.Context, seat *entity.Seat) (*entity.Seat, error)
}

type IBookingService interface {
	HoldSeat(ctx context.Context, dto *bookingservice.HoldSeatDTO) error
	ReleaseSeat(ctx context.Context, dto *bookingservice.ReleaseSeatDTO) error
	CreateBooking(ctx context.Context, dto *bookingservice.CreateBookingDTO) (*entity.Booking, *[]entity.BookingSeat, error)
}

type IUserService interface {
	Create(ctx context.Context, user *entity.User) (*entity.User, error)
	GetByUsername(ctx context.Context, username string) (*entity.User, error)
}

type ISocketService interface {
	HandleWebsocket(ctx *gin.Context)
	BroadcastMessage(message any)
}

type Controller struct {
	authService     IAuthService
	movieService    IMovieService
	userService     IUserService
	showtimeService IShowtimeService
	seatService     ISeatService
	bookingService  IBookingService
	socketService   ISocketService
}

func NewController(
	authService IAuthService,
	movieService IMovieService,
	userService IUserService,
	showtimeService IShowtimeService,
	seatService ISeatService,
	bookingService IBookingService,
	socketService ISocketService,
) *Controller {
	return &Controller{
		authService:     authService,
		movieService:    movieService,
		userService:     userService,
		showtimeService: showtimeService,
		seatService:     seatService,
		bookingService:  bookingService,
		socketService:   socketService,
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
		movieV1Group.POST("/", c.CreateMovie)
	}

	userV1Group := apiV1.Group("/users")
	{
		userV1Group.POST("/register", c.RegisterUser)
		userV1Group.POST("/login", c.Login)
	}

	showtimeV1Group := apiV1.Group("/showtimes")
	{
		showtimeV1Group.POST("/", c.CreateShowtime)
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

	r.GET("/ws", middleware.Authenticate(c.authService), c.socketService.HandleWebsocket)
}
