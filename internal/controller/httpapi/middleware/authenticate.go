package middleware

import (
	authservice "github.com/1612421/cinema-booking/internal/service/auth"
	"net/http"

	"github.com/1612421/cinema-booking/pkg/go-kit/errorx"
	"github.com/gin-gonic/gin"
)

type IAuthService interface {
	Authenticate(req *http.Request) (*authservice.AuthPayload, error)
	CreateAccessToken(payload authservice.AuthPayload) (string, error)
	VerifyAccessToken(tokenString string) (*authservice.AuthPayload, error)
}

func Authenticate(authService IAuthService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authPayload, err := authService.Authenticate(ctx.Request)
		if err != nil {
			ctx.AbortWithStatusJSON(errorx.GetHTTPCode(err), err)
			return
		}

		ctx.Request = ctx.Request.WithContext(authservice.NewContextWithUserSession(ctx.Request.Context(), authPayload))
		ctx.Next()
	}
}
