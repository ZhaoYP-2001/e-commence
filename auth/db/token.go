package db

import (
	"context"
	"e_commerce/auth/config"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

func SetToken(tokenHash string) error {
	_, err := redisClient.Set(context.Background(), tokenHash, 1, config.RefreshTokenExpTime+time.Minute*30).Result()
	if err != nil {
		return fmt.Errorf("set token in redis error: %w", err)
	}
	return nil
}

func GetToken(tokenHash string) (int, error) {
	result, err := redisClient.Get(context.Background(), tokenHash).Int()
	if err != nil {
		return 0, fmt.Errorf("get token from redis error: %w", err)
	}
	return result, nil
}

func DelToken(tokenHash string) error {
	_, err := redisClient.Del(context.Background(), tokenHash).Result()
	if err != nil && err != redis.Nil {
		return fmt.Errorf("delete token from redis error: %w", err)
	}
	return nil
}
