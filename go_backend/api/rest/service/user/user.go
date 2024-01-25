package user

import (
	"context"
	"strings"

	"go_backend/model"
	"go_backend/util"
)

type IUserService interface {
	Login(ctx context.Context, username, password string) (loginResp *LoginResp, err error)
	Register(ctx context.Context, email, password string) (err error)
	GetUserByID(ctx context.Context, id uint) (user *model.User, err error)
}

type IUserQueryRepo interface {
	GetUserByEmail(email string) (user *model.User, err error)
	GetUserByID(id uint) (user *model.User, err error)
}

type IUserCommandRepo interface {
	CreateUser(user *model.User) (err error)
}

type userService struct {
	userQryRepo IUserQueryRepo
	userCmdRepo IUserCommandRepo
}

// add database repo here
func NewUserService(userQryRepo IUserQueryRepo, userCmdRepo IUserCommandRepo) IUserService {
	return &userService{
		userQryRepo: userQryRepo,
		userCmdRepo: userCmdRepo,
	}
}

func (s *userService) Login(ctx context.Context, email, password string) (loginResp *LoginResp, err error) {
	user, err := s.userQryRepo.GetUserByEmail(email)
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

func (s *userService) Register(ctx context.Context, email, password string) (err error) {
	user := &model.User{
		Name:     strings.Split(email, "@")[0],
		Email:    email,
		Password: password,
	}

	err = s.userCmdRepo.CreateUser(user)
	if err != nil {
		return err
	}

	return nil
}

func (s *userService) GetUserByID(ctx context.Context, id uint) (user *model.User, err error) {
	return s.userQryRepo.GetUserByID(id)
}

type LoginResp struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}
