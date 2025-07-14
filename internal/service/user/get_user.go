package userservice

import (
	"context"
	"github.com/1612421/cinema-booking/internal/entity"
	"github.com/1612421/cinema-booking/pkg/go-kit/log"
	"go.uber.org/zap"
)

func (u *UserService) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	logger := log.For(ctx)

	logger.Info("Get user by username", zap.String("username", username))

	user, err := u.userRepo.GetByUsername(ctx, username)
	if err != nil {
		logger.Error("Failed to get user", log.Error(err))
		return nil, err
	}

	logger.Info("Get user success", zap.Reflect("user", user))

	return user, nil
}
