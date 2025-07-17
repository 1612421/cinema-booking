package main

import (
	"github.com/1612421/cinema-booking/cmd/cinema-booking/app"
	_ "github.com/1612421/cinema-booking/docs"
)

// @title Cinema Booking API
// @version 1.0
// @description API for cinema booking
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	app.Run()
}
