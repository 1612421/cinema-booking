package httpapi

import (
	authservice "github.com/1612421/cinema-booking/internal/service/auth"
	bookingservice "github.com/1612421/cinema-booking/internal/service/booking"
	"github.com/1612421/cinema-booking/pkg/go-kit/errorx"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type ReleaseSeatRequest struct {
	ShowtimeID uuid.UUID `json:"showtime_id" binding:"required,max=255"`
	SeatID     uuid.UUID `json:"seat_id" binding:"required,max=255"`
}

type ReleaseSeatResponse struct {
	Message string `json:"message"`
}

// ReleaseSeat godoc
// @Summary      Release a seat
// @Description  Unselect a seat
// @Tags         seat
// @Accept       json
// @Produce      json
// @Security BearerAuth
// @Param        request   	body 	ReleaseSeatRequest  	true  "request credentials"
// @Success      200  {object}  ReleaseSeatResponse
// @Failure      400  {object}  errorx.ErrorWrapper
// @Failure      500  {object}  errorx.ErrorWrapper
// @Router       /bookings/release-seat [post]
func (c *Controller) ReleaseSeat(ctx *gin.Context) {
	request := &ReleaseSeatRequest{}

	if err := ctx.ShouldBind(request); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errorx.New(http.StatusBadRequest, "Parse release seat request failed"))
		return
	}

	userSessionValue := authservice.GetUserSessionValue(ctx).(*authservice.AuthPayload)
	dto := &bookingservice.ReleaseSeatDTO{
		ShowtimeID: request.ShowtimeID,
		UserID:     userSessionValue.UserId,
		SeatID:     request.SeatID,
	}
	err := c.bookingService.ReleaseSeat(ctx, dto)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, ReleaseSeatResponse{
		Message: "release seat successfully",
	})
}
