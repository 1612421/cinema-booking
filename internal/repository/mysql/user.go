package mysql

import (
	"context"
	"github.com/1612421/cinema-booking/internal/entity"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (m *UserRepository) Create(ctx context.Context, user *entity.User) (*entity.User, error) {
	tx := m.db.WithContext(ctx).Create(user)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return user, nil
}

func (m *UserRepository) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	var user *entity.User
	tx := m.db.WithContext(ctx).Model(&entity.User{}).First(&user, "username = ?", username)

	return user, tx.Error
}
