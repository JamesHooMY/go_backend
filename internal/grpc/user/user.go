package user

import (
	context "context"

	userSrv "go_backend/app/service/user"
)

type UserServer struct {
	UserService userSrv.IUserService
}

func NewUserServer(userSrv userSrv.IUserService) UserServiceServer {
	return &UserServer{
		UserService: userSrv,
	}
}

func (u *UserServer) CreateUser(ctx context.Context, req *CreateUserRequest) (*CreateUserResponse, error) {
	return &CreateUserResponse{
		User: &User{
			Id:    "1",
			Email: req.Email,
		},
	}, nil
}

func (u *UserServer) GetUsers(ctx context.Context, req *GetUsersRequest) (*GetUsersResponse, error) {
	return &GetUsersResponse{
		Users: &UserList{
			Users: []*User{
				{
					Id:    "1",
					Email: "james@gmail.com",
				},
			},
		},
	}, nil
}

func (u *UserServer) mustEmbedUnimplementedUserServiceServer() {}
