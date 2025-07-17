package httpapi

import (
	"github.com/1612421/cinema-booking/internal/entity"
	authservice "github.com/1612421/cinema-booking/internal/service/auth"
	bookingservice "github.com/1612421/cinema-booking/internal/service/booking"
	"github.com/1612421/cinema-booking/pkg/go-kit/errorx"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strings"
)

type CreateBookingRequest struct {
	ShowtimeID uuid.UUID   `json:"showtime_id" binding:"required,max=255"`
	SeatIDs    []uuid.UUID `json:"seat_ids" binding:"required,max=5"`
}

type BookingWithSeats struct {
	*entity.Booking
	BookingSeats []*entity.BookingSeat `json:"booking_seats"`
}

type CreateBookingResponse struct {
	Data BookingWithSeats `json:"data"`
}

// CreateBooking godoc
// @Summary      Create a booking
// @Description  Create a booking from selected seats
// @Tags         booking
// @Accept       json
// @Produce      json
// @Security BearerAuth
// @Param        booking   	body 	CreateBookingRequest  	true  "create booking request"
// @Success      200  {object}  CreateBookingResponse
// @Failure      400  {object}  errorx.ErrorWrapper
// @Failure      500  {object}  errorx.ErrorWrapper
// @Router       /bookings [post]
func (c *Controller) CreateBooking(ctx *gin.Context) {
	request := &CreateBookingRequest{}

	if err := ctx.ShouldBind(request); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errorx.New(http.StatusBadRequest, "Parse create booking request failed"))
		return
	}

	userSessionValue := authservice.GetUserSessionValue(ctx).(*authservice.AuthPayload)
	dto := &bookingservice.CreateBookingDTO{
		ShowtimeID: request.ShowtimeID,
		UserID:     userSessionValue.UserId,
		SeatIDs:    request.SeatIDs,
	}
	booking, bookingSeats, err := c.bookingService.CreateBooking(ctx, dto)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, errorx.New(http.StatusBadRequest, "Your seat have already been booking"))
			return
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	c.socketService.BroadcastMessageExceptID(gin.H{
		"event": "booking_created",
		"data": gin.H{
			"showtime_id": dto.ShowtimeID,
			"seat_ids":    dto.SeatIDs,
		},
	}, dto.UserID)

	ctx.JSON(http.StatusOK, CreateBookingResponse{
		Data: BookingWithSeats{
			Booking:      booking,
			BookingSeats: bookingSeats,
		},
	})
}
