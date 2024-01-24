package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func InitRedisClient(c context.Context) (*redis.ClusterClient, error) {
	rdClient := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    []string{viper.GetString("redis.node1"), viper.GetString("redis.node2"), viper.GetString("redis.node3")},
		Password: viper.GetString("redis.password"),
	})

	_, err := rdClient.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return rdClient, nil
}
