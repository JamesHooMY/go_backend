package grpc

import (
	"fmt"

	"go_backend/app/api/grpc/user"
	"go_backend/global"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcClient struct {
	UserSrvClient user.UserServiceClient
}

func NewGrpcClient() *GrpcClient {
	conn, err := grpc.Dial(
		fmt.Sprintf("localhost:%s", viper.GetString("server.grpcPort")),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		global.Logger.Error(fmt.Sprintf("grpc.Dial error: %s\n", err))
		panic(fmt.Sprintf("grpc.Dial error: %s\n", err))
	}

	return &GrpcClient{
		UserSrvClient: user.NewUserServiceClient(conn),
	}
}
