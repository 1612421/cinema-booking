package userservice

import (
	"context"
	"github.com/1612421/cinema-booking/internal/entity"
	"github.com/1612421/cinema-booking/pkg/go-kit/log"
	"go.uber.org/zap"
)

func (u *UserService) Create(ctx context.Context, user *entity.User) (*entity.User, error) {
	logger := log.For(ctx)
	logger.Info("Create user", zap.Reflect("user", user))

	user, err := u.userRepo.Create(ctx, user)
	if err != nil {
		logger.Error("Create user failed", log.Error(err))
		return nil, err
	}

	logger.Info("Create user successfully", zap.Reflect("movie", user))

	return user, nil
}
