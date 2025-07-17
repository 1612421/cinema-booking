package httpapi

import (
	"github.com/1612421/cinema-booking/internal/entity"
	"github.com/1612421/cinema-booking/pkg/go-kit/errorx"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type CreateShowtimeRequest struct {
	MovieID   uuid.UUID `json:"movie_id" binding:"required,max=255"`
	ScreenID  uuid.UUID `json:"screen_id" binding:"required,max=255"`
	StartTime time.Time `json:"start_time" binding:"required"`
}

type CreateShowtimeResponse struct {
	Data *entity.Showtime `json:"data"`
}

// CreateShowtime godoc
// @Summary      Create a showtime
// @Description  Create a showtime
// @Tags         showtime
// @Accept       json
// @Produce      json
// @Security BearerAuth
// @Param        showtime   	body 	CreateShowtimeRequest  	true  "create showtime request"
// @Success      200  {object}  CreateShowtimeResponse
// @Failure      400  {object}  errorx.ErrorWrapper
// @Failure      500  {object}  errorx.ErrorWrapper
// @Router       /showtimes [post]
func (c *Controller) CreateShowtime(ctx *gin.Context) {
	request := &CreateShowtimeRequest{}
	if err := ctx.ShouldBind(request); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errorx.New(http.StatusBadRequest, "Parse create showtime request failed"))
		return
	}

	showtime := &entity.Showtime{
		ID:        uuid.New(),
		MovieID:   request.MovieID,
		ScreenID:  request.ScreenID,
		StartTime: request.StartTime,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	showtime, err := c.showtimeService.Create(ctx, showtime)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, CreateShowtimeResponse{Data: showtime})
}
