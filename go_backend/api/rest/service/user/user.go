package user

import (
	"context"
	"errors"
	"strings"
	"time"

	"go_backend/model"
	"go_backend/util"

	"golang.org/x/crypto/bcrypt"
)

var ErrPasswordIncorrect = errors.New("password incorrect")

type IUserService interface {
	Login(ctx context.Context, name, password string) (loginResp *LoginResp, err error)
	Register(ctx context.Context, email, password string) (err error)
	GetUserByID(ctx context.Context, id uint) (userResp *UserResp, err error)
	GetUserList(ctx context.Context, page, limit int) (userListResp *UserListResp, err error)
	UpdateUser(ctx context.Context, user *model.User) (err error)
	DeleteUserByID(ctx context.Context, id uint) (err error)
}

type IUserQueryRepo interface {
	GetUserByEmail(ctx context.Context, email string) (user *model.User, err error)
	GetUserByID(ctx context.Context, id uint) (user *model.User, err error)
	GetUserList(ctx context.Context, page, limit int) (userList []*model.User, total int64, err error)
}

type IUserCommandRepo interface {
	CreateUser(ctx context.Context, user *model.User) (err error)
	UpdateUser(ctx context.Context, user *model.User) (err error)
	DeleteUser(ctx context.Context, id uint) (err error)
}

type IUserRedisRepo interface {
	SetLoginToken(ctx context.Context, userID uint, token string) (err error)
	GetLoginToken(ctx context.Context, userID uint) (tokenUserID string, err error)
}

type userService struct {
	userQryRepo IUserQueryRepo
	userCmdRepo IUserCommandRepo
	userRdsRepo IUserRedisRepo
}

// add database repo here
func NewUserService(userQryRepo IUserQueryRepo, userCmdRepo IUserCommandRepo, userRdsRepo IUserRedisRepo) IUserService {
	return &userService{
		userQryRepo: userQryRepo,
		userCmdRepo: userCmdRepo,
		userRdsRepo: userRdsRepo,
	}
}

func (s *userService) Login(ctx context.Context, email, password string) (loginResp *LoginResp, err error) {
	user, err := s.userQryRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, ErrPasswordIncorrect
	}

	token, err := util.GenerateJwtToken(user.ID, user.Name)
	if err != nil {
		return nil, err
	}

	if err := s.userRdsRepo.SetLoginToken(ctx, user.ID, token); err != nil {
		return nil, err
	}

	return &LoginResp{
		ID:    user.ID,
		Token: token,
	}, nil
}

type LoginResp struct {
	ID    uint   `json:"id"`
	Token string `json:"token"`
}

func (s *userService) Register(ctx context.Context, email, password string) (err error) {
	user := &model.User{
		Name:     strings.Split(email, "@")[0],
		Email:    email,
		Password: password,
	}

	err = s.userCmdRepo.CreateUser(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (s *userService) GetUserByID(ctx context.Context, id uint) (userResp *UserResp, err error) {
	user, err := s.userQryRepo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	userResp = &UserResp{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		DeletedAt: user.DeletedAt.Time,
		Mobile:    user.Mobile,
		Name:      user.Name,
		Age:       user.Age,
	}

	return userResp, nil
}

type UserResp struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	DeletedAt time.Time `json:"deletedAt"`
	Email     string    `json:"email"`
	Mobile    string    `json:"mobile"`
	Name      string    `json:"name"`
	Age       int       `json:"age"`
}

func (s *userService) GetUserList(ctx context.Context, page, limit int) (userListResp *UserListResp, err error) {
	userList, total, err := s.userQryRepo.GetUserList(ctx, page, limit)
	if err != nil {
		return nil, err
	}

	userListResp = &UserListResp{
		UserList: make([]*UserResp, 0, len(userList)),
		Total:    total,
		Page:     page,
		Limit:    limit,
	}

	for _, user := range userList {
		userListResp.UserList = append(userListResp.UserList, &UserResp{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			DeletedAt: user.DeletedAt.Time,
			Mobile:    user.Mobile,
			Name:      user.Name,
			Age:       user.Age,
		})
	}

	return userListResp, nil
}

type UserListResp struct {
	UserList []*UserResp `json:"userList"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	Limit    int         `json:"limit"`
}

func (s *userService) UpdateUser(ctx context.Context, user *model.User) (err error) {
	err = s.userCmdRepo.UpdateUser(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (s *userService) DeleteUserByID(ctx context.Context, id uint) (err error) {
	err = s.userCmdRepo.DeleteUser(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
