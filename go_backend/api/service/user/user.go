package user

import (
	"context"

	"go_backend/model"
	"go_backend/util"
)

type IUserService interface {
	Login(ctx context.Context, username, password string) (loginResp *LoginResp, err error)
}

type IUserQueryRepo interface {
	GetUserByEmail(email string) (user *model.User, err error)
}

// add database repo here
func NewUserService(userQueryRepo IUserQueryRepo) IUserService {
	return &userService{
		userQuery: userQueryRepo,
	}
}

type userService struct {
	userQuery IUserQueryRepo
}

func (s *userService) Login(ctx context.Context, email, password string) (loginResp *LoginResp, err error) {
	user, err := s.userQuery.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	token, err := util.GenerateJwtToken(user.ID, user.Name)
	if err != nil {
		return nil, err
	}

	return &LoginResp{
		Username: user.Name,
		Token:    token,
	}, nil
}

type LoginResp struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}
