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

type HoldSeatResponse struct {
	Message string `json:"message"`
}

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
	err := c.bookingService.HoldSeat(ctx, dto)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	c.socketService.BroadcastMessage(gin.H{
		"event": "hold_seat",
		"data": gin.H{
			"showtime_id": dto.ShowtimeID,
			"seat_id":     dto.SeatID,
		},
	})

	ctx.JSON(http.StatusOK, HoldSeatResponse{
		Message: "hold seat successfully",
	})
}
