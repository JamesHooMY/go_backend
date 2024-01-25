package user

import (
	"go_backend/api/rest/service/user"
	"go_backend/model"

	"gorm.io/gorm"
)

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

func (q *userQueryRepo) GetUserByID(id uint) (user *model.User, err error) {
	err = q.db.Model(&model.User{}).Where("id = ?", id).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}