package httpapi

import (
	"github.com/1612421/cinema-booking/internal/entity"
	"github.com/1612421/cinema-booking/pkg/go-kit/errorx"
	"github.com/1612421/cinema-booking/pkg/go-kit/log"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type GetSeatsRequest struct {
	ShowtimeID string `uri:"showtime_id" binding:"required,uuid"`
}

type GetSeatsResponse struct {
	Data []*entity.SeatWithStatus `json:"data,omitempty"`
}

// getSeats godoc
// @Summary      Get seats of a showtime
// @Description  Get seats of a showtime
// @Tags         seat
// @Accept       json
// @Produce      json
// @Param        request	path 	GetSeatsRequest  	true  "request credential"
// @Success      200  {object}  GetSeatsResponse
// @Failure      400  {object}  errorx.ErrorWrapper
// @Failure      500  {object}  errorx.ErrorWrapper
// @Router       /showtimes/{showtime_id}/seats [get]
func (c *Controller) getSeats(ctx *gin.Context) {
	request := &GetSeatsRequest{}
	if err := ctx.ShouldBindUri(&request); err != nil {
		log.For(ctx).Error("Parse params failed", log.Error(err))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errorx.New(http.StatusBadRequest, "Parse params failed"))
		return
	}

	ID, _ := uuid.Parse(request.ShowtimeID)
	seats, err := c.seatService.GetSeatsByShowtimeID(ctx, ID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	if seats == nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, errorx.New(http.StatusNotFound, "Movie not found"))
		return
	}

	ctx.JSON(http.StatusOK, GetSeatsResponse{Data: seats})
}
