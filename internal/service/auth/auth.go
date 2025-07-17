package authservice

import (
	"github.com/1612421/cinema-booking/config"
	"github.com/1612421/cinema-booking/pkg/go-kit/errorx"
	"github.com/1612421/cinema-booking/pkg/go-kit/log"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"net/http"
	"time"

	"github.com/1612421/cinema-booking/pkg/go-kit/header"
)

type AuthService struct {
	cfg           *config.Config
	requestParser header.Parser
}

type AuthPayload struct {
	Username string    `json:"username"`
	UserId   uuid.UUID `json:"user_id"`
}

type AuthSocketRequest struct {
	AccessToken string `form:"access_token" binding:"required"`
}

func ProvideHeaders() header.Parser {
	return header.NewParser(
		header.WithParser(
			header.ParseAuthorization,
			header.ParseCookies,
		),
	)
}

func NewAuthService(cfg *config.Config, requestParser header.Parser) *AuthService {
	return &AuthService{
		cfg:           cfg,
		requestParser: requestParser,
	}
}

func (a *AuthService) CreateAccessToken(payload AuthPayload) (string, error) {
	// Create a new JWT token with claims
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":     payload.Username,
		"iss":     a.cfg.Service.Name,
		"exp":     time.Now().Add(time.Duration(a.cfg.Auth.ExpireIn) * time.Second).Unix(),
		"iat":     time.Now().Unix(),
		"user_id": payload.UserId,
	})

	tokenString, err := claims.SignedString([]byte(a.cfg.Auth.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (a *AuthService) VerifyAccessToken(tokenString string) (*AuthPayload, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate signing algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errorx.New(http.StatusUnauthorized, "Invalid token")
		}

		return []byte(a.cfg.Auth.Secret), nil
	})

	if err != nil || !token.Valid {
		return nil, errorx.New(http.StatusUnauthorized, "Invalid token")
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return nil, errorx.New(http.StatusUnauthorized, "Invalid token")
	}

	userIDString := claims["user_id"].(string)
	userID, err := uuid.Parse(userIDString)
	if err != nil {
		return nil, err
	}

	payload := &AuthPayload{
		Username: claims["sub"].(string),
		UserId:   userID,
	}

	return payload, nil
}

func (a *AuthService) Authenticate(req *http.Request) (*AuthPayload, error) {
	metadata, err := a.requestParser.ParseHeader(req.Header)
	if err != nil {
		log.For(req.Context()).Error("Parse header failed", log.Error(err))
		return nil, errorx.New(http.StatusUnauthorized, "Authentication parse metadata failed")
	}

	const missingAuthMsg = "Authorization header missing or malformed"
	if metadata == nil {
		return nil, errorx.New(http.StatusUnauthorized, missingAuthMsg)
	}

	// Check header Authorization
	tokenString := metadata.Authorization
	if tokenString == "" {
		return nil, errorx.New(http.StatusUnauthorized, missingAuthMsg)
	}

	payload, err := a.VerifyAccessToken(tokenString)
	if err != nil {
		return nil, errorx.New(http.StatusUnauthorized, err.Error())
	}

	return payload, nil
}

func (a *AuthService) AuthenticateSocket(ctx *gin.Context) (*AuthPayload, error) {
	request := &AuthSocketRequest{}
	if err := ctx.ShouldBindQuery(request); err != nil {
		return nil, errorx.New(http.StatusBadRequest, "Parse auth socket request failed")
	}

	// Check header Authorization
	tokenString := request.AccessToken
	if tokenString == "" {
		return nil, errorx.New(http.StatusUnauthorized, "Access token missing or malformed")
	}

	payload, err := a.VerifyAccessToken(tokenString)
	if err != nil {
		return nil, errorx.New(http.StatusUnauthorized, err.Error())
	}

	return payload, nil
}
