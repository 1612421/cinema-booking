package httpapi

import (
	"github.com/1612421/cinema-booking/internal/entity"
	"github.com/1612421/cinema-booking/pkg/go-kit/errorx"
	"github.com/1612421/cinema-booking/pkg/go-kit/log"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type GetShowtimesRequest struct {
	MovieID string `uri:"movie_id" binding:"required,uuid"`
}

type GetShowtimesResponse struct {
	Data []*entity.Showtime `json:"data,omitempty"`
}

// getShowtimes godoc
// @Summary      Get list showtimes of a movie
// @Description  Get list showtimes of a movie
// @Tags         showtime
// @Accept       json
// @Produce      json
// @Param        request	path 	GetShowtimesRequest  	true  "request credential"
// @Success      200  {object}  GetShowtimesResponse
// @Failure      400  {object}  errorx.ErrorWrapper
// @Failure      500  {object}  errorx.ErrorWrapper
// @Router       /movies/{movie_id}/showtimes [get]
func (c *Controller) getShowtimes(ctx *gin.Context) {
	request := &GetShowtimesRequest{}
	if err := ctx.ShouldBindUri(&request); err != nil {
		log.For(ctx).Error("Parse params failed", log.Error(err))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errorx.New(http.StatusBadRequest, "Parse params failed"))
		return
	}

	ID, _ := uuid.Parse(request.MovieID)
	showtimes, err := c.showtimeService.GetShowtimesByMovieID(ctx, ID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	if showtimes == nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, errorx.New(http.StatusNotFound, "Movie not found"))
		return
	}

	ctx.JSON(http.StatusOK, GetShowtimesResponse{Data: showtimes})
}
