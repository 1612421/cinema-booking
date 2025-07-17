package httpapi

import (
	"github.com/1612421/cinema-booking/internal/entity"
	"github.com/1612421/cinema-booking/pkg/go-kit/errorx"
	"github.com/1612421/cinema-booking/pkg/go-kit/log"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type GetMovieRequest struct {
	ID string `uri:"movie_id" binding:"required,uuid"`
}

type GetMovieResponse struct {
	Data *entity.Movie `json:"data,omitempty"`
}

// getMovie godoc
// @Summary      Get movie detail
// @Description  Get movie detail
// @Tags         movie
// @Accept       json
// @Produce      json
// @Param        movie	path 	GetMovieRequest  	true  "movie ID"
// @Success      200  {object}  GetMovieResponse
// @Failure      400  {object}  errorx.ErrorWrapper
// @Failure      500  {object}  errorx.ErrorWrapper
// @Router       /movies/{movie_id} [get]
func (c *Controller) getMovie(ctx *gin.Context) {
	request := &GetMovieRequest{}
	if err := ctx.ShouldBindUri(&request); err != nil {
		log.For(ctx).Error("Parse params failed", log.Error(err))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errorx.New(http.StatusBadRequest, "Parse params failed"))
		return
	}

	ID, _ := uuid.Parse(request.ID)
	movie, err := c.movieService.GetMovie(ctx, ID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	if movie == nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, errorx.New(http.StatusNotFound, "Movie not found"))
		return
	}

	ctx.JSON(http.StatusOK, GetMovieResponse{Data: movie})
}
