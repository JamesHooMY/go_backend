package user

import (
	"go_backend/api/service/user"

	"gorm.io/gorm"
)

type UserQuery struct {
	db *gorm.DB
}

func NewUserQuery(db *gorm.DB) user.UserQuery {
	return &UserQuery{
		db: db,
	}
}

func (q *UserQuery) Login(username, password string) (token string, err error) {
	return "", nil
}

func (q *UserQuery) Info(token string) (name string, err error) {
	return "UserInfo", nil
}
