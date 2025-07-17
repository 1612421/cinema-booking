package httpapi

import (
	"github.com/1612421/cinema-booking/internal/entity"
	"github.com/1612421/cinema-booking/pkg/go-kit/errorx"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type CreateSeatRequest struct {
	ScreenID uuid.UUID `json:"screen_id" binding:"required,max=255"`
	Row      string    `json:"row" binding:"required,max=255"`
	Number   int       `json:"number" binding:"required"`
	Class    string    `json:"class" binding:"required,max=255"`
}

type CreateSeatResponse struct {
	Data *entity.Seat `json:"data"`
}

// CreateSeat godoc
// @Summary      Create a seat
// @Description  Create a seat
// @Tags         seat
// @Accept       json
// @Produce      json
// @Security BearerAuth
// @Param        seat   	body 	CreateSeatRequest  	true  "create seat request"
// @Success      200  {object}  CreateSeatResponse
// @Failure      400  {object}  errorx.ErrorWrapper
// @Failure      500  {object}  errorx.ErrorWrapper
// @Router       /seats [post]
func (c *Controller) CreateSeat(ctx *gin.Context) {
	request := &CreateSeatRequest{}
	if err := ctx.ShouldBind(request); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errorx.New(http.StatusBadRequest, "Parse create seat request failed"))
		return
	}

	seat := &entity.Seat{
		ID:        uuid.New(),
		ScreenID:  request.ScreenID,
		Row:       request.Row,
		Number:    request.Number,
		Class:     request.Class,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	seat, err := c.seatService.Create(ctx, seat)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, CreateSeatResponse{Data: seat})
}
