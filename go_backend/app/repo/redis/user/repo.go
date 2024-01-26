package user

import (
	"context"
	"fmt"
	"time"

	"go_backend/app/service/user"

	"github.com/redis/go-redis/v9"
)

var duration = 1 * time.Hour

type userRedisRepo struct {
	rds *redis.Client
}

func NewUserRedisRepo(rds *redis.Client) user.IUserRedisRepo {
	return &userRedisRepo{
		rds: rds,
	}
}

func (r *userRedisRepo) SetLoginToken(ctx context.Context, userID uint, token string) (err error) {
	return r.rds.Set(ctx, fmt.Sprintf("go_backend_login_token_%d", userID), token, duration).Err()
}

func (r *userRedisRepo) GetLoginToken(ctx context.Context, userID uint) (token string, err error) {
	return r.rds.Get(ctx, fmt.Sprintf("go_backend_login_token_%d", userID)).Result()
}
