package user

import (
	"context"
	"errors"

	"go_backend/api/rest/service/user"
	"go_backend/model"

	"github.com/go-sql-driver/mysql"
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
	// method 1: one more I/O to check if user existed
	// var existedUser *model.User
	// err = r.db.WithContext(ctx).Model(&model.User{}).Where("email = ?", user.Email).First(&existedUser).Error
	// if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
	// 	return err
	// }

	// if existedUser != nil && existedUser.Email == user.Email {
	// 	return ErrUserExisted
	// }

	// if err := r.db.WithContext(ctx).Model(&model.User{}).Create(user).Error; err != nil {
	// 	return err
	// }

	// method 2: use mysql error code to check if user existed
	if err := r.db.WithContext(ctx).Model(&model.User{}).Create(user).Error; err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			switch mysqlErr.Number {
			case 1062:
				return ErrUserExisted
			}
		}

		return err
	}

	return nil
}

func (r *userCommandRepo) UpdateUser(ctx context.Context, user *model.User) (err error) {
	result := r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", user.ID).Updates(user)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (r *userCommandRepo) DeleteUser(ctx context.Context, id uint) (err error) {
	var user *model.User
	if err := r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	if err := r.db.Model(&model.User{}).Where("id = ?", id).Delete(&model.User{}).Error; err != nil {
		return err
	}

	return nil
}
