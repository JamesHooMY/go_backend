package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func InitRedisClient(c context.Context) (*redis.Client, error) {
	rdClient := redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.host"),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
	})

	if _, err := rdClient.Ping(c).Result(); err != nil {
		return nil, err
	}

	return rdClient, nil
}
