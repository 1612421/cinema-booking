package httpapi

import (
	bookingservice "github.com/1612421/cinema-booking/internal/service/booking"
	"github.com/1612421/cinema-booking/pkg/go-kit/errorx"
	"github.com/1612421/cinema-booking/pkg/go-kit/log"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type GetSeatInterestRequest struct {
	ShowtimeID string `uri:"showtime_id" binding:"required,uuid"`
	SeatID     string `uri:"seat_id" binding:"required,uuid"`
}

type GetSeatInterestResponseData struct {
	Qty int64 `json:"qty"`
}

type GetSeatInterestResponse struct {
	Data *GetSeatInterestResponseData `json:"data"`
}

// getSeatInterest godoc
// @Summary      Get quantity of interested user of this seat
// @Description  Get quantity of interested user of this seat
// @Tags         seat
// @Accept       json
// @Produce      json
// @Param        request	path 	GetSeatInterestRequest  	true  "request credential"
// @Success      200  {object}  GetSeatInterestResponse
// @Failure      400  {object}  errorx.ErrorWrapper
// @Failure      500  {object}  errorx.ErrorWrapper
// @Router       /showtimes/{showtime_id}/seats/{seat_id}/interest [get]
func (c *Controller) getSeatInterest(ctx *gin.Context) {
	request := &GetSeatInterestRequest{}
	if err := ctx.ShouldBindUri(request); err != nil {
		log.For(ctx).Error("Parse params failed", log.Error(err))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errorx.New(http.StatusBadRequest, "Parse request failed"))
		return
	}

	showtimeID, _ := uuid.Parse(request.ShowtimeID)
	seatID, _ := uuid.Parse(request.SeatID)
	dto := &bookingservice.SeatInterestDTO{
		ShowtimeID: showtimeID,
		SeatID:     seatID,
	}
	qty, err := c.bookingService.GetSeatInterest(ctx, dto)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, &GetSeatInterestResponse{
		Data: &GetSeatInterestResponseData{Qty: qty},
	})
}
