package httpapi

import (
	authservice "github.com/1612421/cinema-booking/internal/service/auth"
	bookingservice "github.com/1612421/cinema-booking/internal/service/booking"
	"github.com/1612421/cinema-booking/pkg/go-kit/errorx"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type HoldSeatRequest struct {
	ShowtimeID uuid.UUID `json:"showtime_id" binding:"required,max=255"`
	SeatID     uuid.UUID `json:"seat_id" binding:"required,max=255"`
}

type HoldSeatResponseData struct {
	// Number of users currently holding seats
	Qty int64 `json:"qty"`
}

type HoldSeatResponse struct {
	Data *HoldSeatResponseData `json:"data"`
}

// HoldSeat godoc
// @Summary      Hold a seat
// @Description  Hold a seat of a showtime (not lock seat other user still can book this seat)
// @Tags         seat
// @Accept       json
// @Produce      json
// @Security BearerAuth
// @Param        request   	body 	HoldSeatRequest  	true  "request credentials"
// @Success      200  {object}  HoldSeatResponse
// @Failure      400  {object}  errorx.ErrorWrapper
// @Failure      500  {object}  errorx.ErrorWrapper
// @Router       /bookings/hold-seat [post]
func (c *Controller) HoldSeat(ctx *gin.Context) {
	request := &HoldSeatRequest{}

	if err := ctx.ShouldBind(request); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errorx.New(http.StatusBadRequest, "Parse hold seat request failed"))
		return
	}

	userSessionValue := authservice.GetUserSessionValue(ctx).(*authservice.AuthPayload)
	dto := &bookingservice.HoldSeatDTO{
		ShowtimeID: request.ShowtimeID,
		UserID:     userSessionValue.UserId,
		SeatID:     request.SeatID,
	}
	qty, err := c.bookingService.HoldSeat(ctx, dto)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, HoldSeatResponse{
		Data: &HoldSeatResponseData{
			Qty: qty,
		},
	})
}
