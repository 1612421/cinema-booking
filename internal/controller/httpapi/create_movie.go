package httpapi

import (
	"github.com/1612421/cinema-booking/internal/entity"
	"github.com/1612421/cinema-booking/pkg/go-kit/errorx"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type CreateMovieRequest struct {
	Title    string `json:"title" binding:"required"`
	Duration int    `json:"duration" binding:"required"`
}

type CreateMovieResponse struct {
	Data *entity.Movie `json:"data"`
}

// CreateMovie godoc
// @Summary      Create a movie
// @Description  Create a movie
// @Tags         movie
// @Accept       json
// @Produce      json
// @Security BearerAuth
// @Param        movie   	body 	CreateMovieRequest  	true  "Create movie request"
// @Success      200  {object}  CreateMovieResponse
// @Failure      400  {object}  errorx.ErrorWrapper
// @Failure      500  {object}  errorx.ErrorWrapper
// @Router       /movies [post]
func (c *Controller) CreateMovie(ctx *gin.Context) {
	request := &CreateMovieRequest{}
	if err := ctx.ShouldBind(request); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errorx.New(http.StatusBadRequest, "Parse create movie request failed"))
		return
	}

	movie := &entity.Movie{
		ID:        uuid.New(),
		Title:     request.Title,
		Duration:  request.Duration,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	movie, err := c.movieService.Create(ctx, movie)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, CreateMovieResponse{Data: movie})
}
