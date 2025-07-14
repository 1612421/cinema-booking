//go:build wireinject
// +build wireinject

//go:generate wire gen

package internal

import (
	"github.com/1612421/cinema-booking/config"
	infrastructurecache "github.com/1612421/cinema-booking/internal/infrastructure/cache"
	"github.com/1612421/cinema-booking/internal/infrastructure/database"
	"github.com/1612421/cinema-booking/internal/repository/mysql"
	"github.com/1612421/cinema-booking/internal/repository/redis"
	bookingservice "github.com/1612421/cinema-booking/internal/service/booking"
	movieservice "github.com/1612421/cinema-booking/internal/service/movie"
	seatservice "github.com/1612421/cinema-booking/internal/service/seat"
	showtimeservice "github.com/1612421/cinema-booking/internal/service/showtime"
	socketservice "github.com/1612421/cinema-booking/internal/service/socket"
	userservice "github.com/1612421/cinema-booking/internal/service/user"
	"github.com/google/wire"

	"github.com/1612421/cinema-booking/internal/controller/httpapi"
	authservice "github.com/1612421/cinema-booking/internal/service/auth"
)

var (
	configSet          = wire.NewSet(config.GetConfig)
	databaseSet        = wire.NewSet(database.NewMySQLDB)
	redisSet           = wire.NewSet(infrastructurecache.NewRedisClient, infrastructurecache.NewRedLockClient)
	movieRepositorySet = wire.NewSet(
		mysql.NewMovieRepository,
		wire.Bind(new(movieservice.IMovieRepository), new(*mysql.MovieRepository)),
	)
	userRepositorySet = wire.NewSet(
		mysql.NewUserRepository,
		wire.Bind(new(userservice.IUserRepository), new(*mysql.UserRepository)),
	)
	showtimeRepositorySet = wire.NewSet(
		mysql.NewShowtimeRepository,
		wire.Bind(new(showtimeservice.IShowtimeRepository), new(*mysql.ShowtimeRepository)),
	)
	seatRepositorySet = wire.NewSet(
		mysql.NewSeatRepository,
		wire.Bind(new(seatservice.ISeatRepository), new(*mysql.SeatRepository)),
	)
	bookingRepositorySet = wire.NewSet(
		mysql.NewBookingRepository,
		wire.Bind(new(bookingservice.IBookingRepository), new(*mysql.BookingRepository)),
	)
	authServiceSet = wire.NewSet(
		authservice.ProvideHeaders,
		authservice.NewAuthService,
		wire.Bind(new(httpapi.IAuthService), new(*authservice.AuthService)),
	)
	movieServiceSet = wire.NewSet(
		movieservice.NewMovieService,
		wire.Bind(new(httpapi.IMovieService), new(*movieservice.MovieService)),
	)
	userServiceSet = wire.NewSet(
		userservice.NewUserService,
		wire.Bind(new(httpapi.IUserService), new(*userservice.UserService)),
	)
	showtimeServiceSet = wire.NewSet(
		showtimeservice.NewShowtimeService,
		wire.Bind(new(httpapi.IShowtimeService), new(*showtimeservice.ShowtimeService)),
	)
	seatServiceSet = wire.NewSet(
		seatservice.NewSeatService,
		wire.Bind(new(httpapi.ISeatService), new(*seatservice.SeatService)),
	)
	bookingServiceSet = wire.NewSet(
		bookingservice.NewBookingService,
		redis.NewSeatCache,

		wire.Bind(new(httpapi.IBookingService), new(*bookingservice.BookingService)),
	)
	socketServiceSet = wire.NewSet(
		socketservice.NewSocketService,
		wire.Bind(new(httpapi.ISocketService), new(*socketservice.SocketService)),
	)
	httpapiControllerSet = wire.NewSet(httpapi.NewController)
)

func InitializeHTTPAPIController() (*httpapi.Controller, func(), error) {
	panic(
		wire.Build(
			configSet,
			redisSet,
			databaseSet,
			movieRepositorySet,
			userRepositorySet,
			showtimeRepositorySet,
			seatRepositorySet,
			bookingRepositorySet,
			authServiceSet,
			movieServiceSet,
			userServiceSet,
			showtimeServiceSet,
			seatServiceSet,
			bookingServiceSet,
			socketServiceSet,
			httpapiControllerSet,
		),
	)
}
