package grpc

import (
	userRepo "go_backend/app/repo/mysql/user"
	userRdsRepo "go_backend/app/repo/redis/user"
	userSrv "go_backend/app/service/user"
	"go_backend/internal/grpc/user"

	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

func InitGrpcServer(grpcSvr *grpc.Server, db *gorm.DB, rds *redis.Client) *grpc.Server {
	userServer := user.NewUserServer(userSrv.NewUserService(
		userRepo.NewUserQueryRepo(db),
		userRepo.NewUserCommandRepo(db),
		userRdsRepo.NewUserRedisRepo(rds),
	))
	user.RegisterUserServiceServer(grpcSvr, userServer)

	return grpcSvr
}
