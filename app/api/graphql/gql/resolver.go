package gql

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

import (
	userSrv "go_backend/app/service/user"
)

type Resolver struct {
	UserService userSrv.IUserService
}
