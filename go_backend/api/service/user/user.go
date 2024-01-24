package user

import (
	"context"
)

type UserService interface {
	Login(ctx context.Context, username, password string) (token string, err error)
	Info(ctx context.Context, token string) (name string, err error)
}

type UserQuery interface {
	Login(username, password string) (token string, err error)
	Info(token string) (name string, err error)
}

// add database repo here
func NewUserService(userQuery UserQuery) UserService {
	return &userService{
		userQuery: userQuery,
	}
}

type userService struct {
	userQuery UserQuery
}

func (s *userService) Login(ctx context.Context, username, password string) (token string, err error) {
	return "", nil
}

func (s *userService) Info(ctx context.Context, token string) (name string, err error) {
	return "UserInfo", nil
}
