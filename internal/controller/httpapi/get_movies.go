package httpapi

import (
	"github.com/1612421/cinema-booking/internal/entity"
	"github.com/gin-gonic/gin"
	"net/http"
)

type GetMoviesResponse struct {
	Data []*entity.Movie `json:"data,omitempty"`
}

// getMovies godoc
// @Summary      Get list movies
// @Description  Get list movies
// @Tags         movie
// @Accept       json
// @Produce      json
// @Success      200  {object}  GetMoviesResponse
// @Failure      400  {object}  errorx.ErrorWrapper
// @Failure      500  {object}  errorx.ErrorWrapper
// @Router       /movies [get]
func (c *Controller) getMovies(ctx *gin.Context) {
	movies, err := c.movieService.GetMovies(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, GetMoviesResponse{Data: movies})
}
