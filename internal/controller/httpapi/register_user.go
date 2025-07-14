package httpapi

import (
	"github.com/1612421/cinema-booking/internal/entity"
	authservice "github.com/1612421/cinema-booking/internal/service/auth"
	"github.com/1612421/cinema-booking/pkg/go-kit/errorx"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
	"time"
)

type UserRegisterRequest struct {
	Username    string `json:"username" binding:"required,min=6,max=24"`
	Password    string `json:"password" binding:"required,min=6,max=30"`
	Address     string `json:"address" binding:"required,max=255"`
	PhoneNumber string `json:"phone_number" binding:"required,max=11"`
}

type UserWithAccessToken struct {
	*entity.User
	AccessToken string `json:"access_token"`
}

type UserRegisterResponse struct {
	Data *UserWithAccessToken `json:"data"`
}

func (c *Controller) RegisterUser(ctx *gin.Context) {
	request := &UserRegisterRequest{}
	if err := ctx.ShouldBind(request); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errorx.New(http.StatusBadRequest, "Parse user register request failed"))
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(request.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errorx.New(http.StatusBadRequest, "Password hashing failed"))
	}

	user := &entity.User{
		ID:          uuid.New(),
		Username:    request.Username,
		Password:    string(hashedPassword),
		Address:     request.Address,
		PhoneNumber: request.PhoneNumber,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	user, err = c.userService.Create(ctx, user)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, errorx.New(http.StatusBadRequest, "User already exists"))
			return
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	accessToken, _ := c.authService.CreateAccessToken(authservice.AuthPayload{
		Username: request.Username,
		UserId:   user.ID,
	})

	ctx.JSON(http.StatusOK, UserRegisterResponse{Data: &UserWithAccessToken{
		User:        user,
		AccessToken: accessToken,
	}})
}
