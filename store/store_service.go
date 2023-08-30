package store

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis"
	// "github.com/go-redis/redis/v8"
)

// Raw wrapper
type StorageService struct {
	redisClient *redis.Client
}

// high level declaration
var (
	storeService = &StorageService{}
	ctx          = context.Background()
)

// Cache expiration duration
const cacheDuration = 6 * time.Hour

// Store service with pointer return
func InitializeStore() *StorageService{
	redisClient := redis.NewClient(&redis.Options{
		Addr: "172.28.83.36:6379",
		Password: "",
		DB: 0,
	})

	pong, error := redisClient.Ping(ctx).Result()
	if error != nil{
		panic(fmt.Sprintf("Error init redis: %v", error))
	}

	fmt.Printf("\nRedis Started successfully: pong message = {%s}", &pong)
	storeService.redisClient = redisClient

	return storeService
}

