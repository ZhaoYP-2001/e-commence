package db

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB
var redisClient *redis.Client

func InitDB(dsn string) {
	var err error
	db, err = gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic(fmt.Errorf("init db error: %w", err))
	}
	db.AutoMigrate(&PreKey{})
}

func InitRedis(addr string, password string) {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})
	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
}
