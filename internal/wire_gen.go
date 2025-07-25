// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package internal

import (
	"github.com/1612421/cinema-booking/config"
	"github.com/1612421/cinema-booking/internal/bot"
	"github.com/1612421/cinema-booking/internal/controller/httpapi"
	"github.com/1612421/cinema-booking/internal/infrastructure/cache"
	"github.com/1612421/cinema-booking/internal/infrastructure/database"
	"github.com/1612421/cinema-booking/internal/repository/mysql"
	"github.com/1612421/cinema-booking/internal/repository/redis"
	"github.com/1612421/cinema-booking/internal/service/auth"
	"github.com/1612421/cinema-booking/internal/service/booking"
	"github.com/1612421/cinema-booking/internal/service/movie"
	"github.com/1612421/cinema-booking/internal/service/seat"
	"github.com/1612421/cinema-booking/internal/service/showtime"
	"github.com/1612421/cinema-booking/internal/service/socket"
	"github.com/1612421/cinema-booking/internal/service/user"
	"github.com/google/wire"
)

// Injectors from wire.go:

func InitializeHTTPAPIController() (*httpapi.Controller, func(), error) {
	configConfig := config.GetConfig()
	parser := authservice.ProvideHeaders()
	authService := authservice.NewAuthService(configConfig, parser)
	db, cleanup, err := database.NewMySQLDB(configConfig)
	if err != nil {
		return nil, nil, err
	}
	movieRepository := mysql.NewMovieRepository(db)
	movieService := movieservice.NewMovieService(configConfig, movieRepository)
	userRepository := mysql.NewUserRepository(db)
	userService := userservice.NewUserService(configConfig, userRepository)
	showtimeRepository := mysql.NewShowtimeRepository(db)
	showtimeService := showtimeservice.NewShowtimeService(configConfig, showtimeRepository)
	seatRepository := mysql.NewSeatRepository(db)
	seatService := seatservice.NewSeatService(configConfig, seatRepository)
	bookingRepository := mysql.NewBookingRepository(db)
	universalClient, cleanup2, err := infrastructurecache.NewRedisClient(configConfig)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	iSeatCache := redis.NewSeatCache(universalClient)
	bookingService := bookingservice.NewBookingService(configConfig, bookingRepository, iSeatCache)
	socketService := socketservice.NewSocketService()
	globalViewerSeatBotRunner := bot.NewGlobalViewerSeatBotRunner()
	controller := httpapi.NewController(authService, movieService, userService, showtimeService, seatService, bookingService, socketService, globalViewerSeatBotRunner)
	return controller, func() {
		cleanup2()
		cleanup()
	}, nil
}

// wire.go:

var (
	configSet              = wire.NewSet(config.GetConfig)
	databaseSet            = wire.NewSet(database.NewMySQLDB)
	redisSet               = wire.NewSet(infrastructurecache.NewRedisClient, infrastructurecache.NewRedLockClient)
	movieRepositorySet     = wire.NewSet(mysql.NewMovieRepository, wire.Bind(new(movieservice.IMovieRepository), new(*mysql.MovieRepository)))
	userRepositorySet      = wire.NewSet(mysql.NewUserRepository, wire.Bind(new(userservice.IUserRepository), new(*mysql.UserRepository)))
	showtimeRepositorySet  = wire.NewSet(mysql.NewShowtimeRepository, wire.Bind(new(showtimeservice.IShowtimeRepository), new(*mysql.ShowtimeRepository)))
	seatRepositorySet      = wire.NewSet(mysql.NewSeatRepository, wire.Bind(new(seatservice.ISeatRepository), new(*mysql.SeatRepository)))
	bookingRepositorySet   = wire.NewSet(mysql.NewBookingRepository, wire.Bind(new(bookingservice.IBookingRepository), new(*mysql.BookingRepository)))
	authServiceSet         = wire.NewSet(authservice.ProvideHeaders, authservice.NewAuthService, wire.Bind(new(httpapi.IAuthService), new(*authservice.AuthService)))
	movieServiceSet        = wire.NewSet(movieservice.NewMovieService, wire.Bind(new(httpapi.IMovieService), new(*movieservice.MovieService)))
	userServiceSet         = wire.NewSet(userservice.NewUserService, wire.Bind(new(httpapi.IUserService), new(*userservice.UserService)))
	showtimeServiceSet     = wire.NewSet(showtimeservice.NewShowtimeService, wire.Bind(new(httpapi.IShowtimeService), new(*showtimeservice.ShowtimeService)))
	seatServiceSet         = wire.NewSet(seatservice.NewSeatService, wire.Bind(new(httpapi.ISeatService), new(*seatservice.SeatService)))
	bookingServiceSet      = wire.NewSet(bookingservice.NewBookingService, redis.NewSeatCache, wire.Bind(new(httpapi.IBookingService), new(*bookingservice.BookingService)))
	socketServiceSet       = wire.NewSet(socketservice.NewSocketService, wire.Bind(new(httpapi.ISocketService), new(*socketservice.SocketService)))
	viewerSeatBotRunnerSet = wire.NewSet(bot.NewGlobalViewerSeatBotRunner, wire.Bind(new(httpapi.IGlobalViewerSeatBotRunner), new(*bot.GlobalViewerSeatBotRunner)))
	httpapiControllerSet   = wire.NewSet(httpapi.NewController)
)
