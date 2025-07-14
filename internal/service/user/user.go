package userservice

import (
	"context"
	"github.com/1612421/cinema-booking/config"
	"github.com/1612421/cinema-booking/internal/entity"
)

type IUserRepository interface {
	Create(ctx context.Context, user *entity.User) (*entity.User, error)
	GetByUsername(ctx context.Context, username string) (*entity.User, error)
}

type UserService struct {
	cfg      *config.Config
	userRepo IUserRepository
}

func NewUserService(cfg *config.Config, userRepo IUserRepository) *UserService {
	return &UserService{
		cfg:      cfg,
		userRepo: userRepo,
	}
}
