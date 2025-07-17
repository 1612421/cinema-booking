package httpapi

import (
	authservice "github.com/1612421/cinema-booking/internal/service/auth"
	"github.com/1612421/cinema-booking/pkg/go-kit/errorx"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required,max=24"`
	Password string `json:"password" binding:"required,min=6,max=30"`
}

type LoginResponse struct {
	Data *UserWithAccessToken `json:"data"`
}

// Login godoc
// @Summary      Login
// @Description  Login
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        request   	body 	LoginRequest  	true  "request credentials"
// @Success      200  {object}  LoginResponse
// @Failure      400  {object}  errorx.ErrorWrapper
// @Failure      500  {object}  errorx.ErrorWrapper
// @Router       /users/login [post]
func (c *Controller) Login(ctx *gin.Context) {
	request := &LoginRequest{}
	if err := ctx.ShouldBind(request); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errorx.New(http.StatusBadRequest, "Parse login request failed"))
		return
	}

	user, err := c.userService.GetByUsername(ctx, request.Username)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, "Invalid username or password")
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorx.New(http.StatusUnauthorized, "Invalid username or password"))
		return
	}

	accessToken, _ := c.authService.CreateAccessToken(authservice.AuthPayload{
		Username: user.Username,
		UserId:   user.ID,
	})

	// ⏩ Start bot once if not already running
	go func() {
		c.viewerSeatBotRunner.StartIfNotRunning(
			c.showtimeService.GetShowtimeRepo(),
			c.seatService.GetSeatRepo(),
			c.userService.GetUserRepo(),
			c.bookingService,
			c.socketService,
			user.ID,
		)
	}()

	ctx.JSON(http.StatusOK, UserRegisterResponse{Data: &UserWithAccessToken{
		User:        user,
		AccessToken: accessToken,
	}})
}
