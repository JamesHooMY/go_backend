package user

import (
	"errors"

	"go_backend/api/service/user"
	"go_backend/model"

	"gorm.io/gorm"
)

var ErrUserNotFound = errors.New("user not found")

type userQueryRepo struct {
	db *gorm.DB
}

func NewUserQueryRepo(db *gorm.DB) user.IUserQueryRepo {
	return &userQueryRepo{
		db: db,
	}
}

func (q *userQueryRepo) GetUserByEmail(email string) (user *model.User, err error) {
	err = q.db.Model(&model.User{}).Where("email = ?", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}
