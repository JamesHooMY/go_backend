package user

import (
	"context"

	"gorm.io/gorm"
)

type UserService interface {
	Login(ctx context.Context, username, password string) (token string, err error)
	Info(ctx context.Context, token string) (name string, err error)
}

// add database repo here
func NewUserService(db *gorm.DB) UserService {
	return &userService{
		db: db,
	}
}

type userService struct {
	db *gorm.DB
}

func (s *userService) Login(ctx context.Context, username, password string) (token string, err error) {
	return "", nil
}

func (s *userService) Info(ctx context.Context, token string) (name string, err error) {
	return "UserInfo", nil
}
