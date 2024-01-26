package user

import (
	"context"

	"go_backend/api/rest/service/user"
	"go_backend/model"

	"gorm.io/gorm"
)

type userCommandRepo struct {
	db *gorm.DB
}

func NewUserCommandRepo(db *gorm.DB) user.IUserCommandRepo {
	return &userCommandRepo{
		db: db,
	}
}

func (r *userCommandRepo) CreateUser(ctx context.Context, user *model.User) (err error) {
	var existedUser *model.User
	err = r.db.WithContext(ctx).Model(&model.User{}).Where("email = ?", user.Email).First(&existedUser).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	if existedUser != nil && existedUser.Email == user.Email {
		return ErrUserExisted
	}

	if err := r.db.Model(&model.User{}).Create(user).Error; err != nil {
		return err
	}

	return nil
}
