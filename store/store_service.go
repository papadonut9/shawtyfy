package store

import (
	"context"
	"fmt"
	"os"
	"time"
	"github.com/go-redis/redis/v8"
)

// Raw wrapper
type StorageService struct {
	redisClient *redis.Client
}

// high level declaration
var (
	storeService = &StorageService{}
	ctx          = context.Background()
	dbstr = os.Getenv("REDIS_DB")
)

// environment variable string to integer conversion
// redis_db := strconv.Atoi(dbstr);

// Cache expiration duration
const cacheDuration = 6 * time.Hour

// Store service with pointer return
func InitializeStore() *StorageService {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: "",
		// DB:       strconv.Atoi(os.Getenv(REDIS_DB)),
		DB:       0,
	})

	pong, err := redisClient.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("Error init redis: %v", err))
	}

	fmt.Printf("\nRedis Started successfully: pong message = {%s}", pong)
	storeService.redisClient = redisClient

	return storeService
}

// Storage APIs
func saveUrlMapping(shortUrl string, originalUrl string, userid string) {
	err := storeService.redisClient.Set(ctx, shortUrl, originalUrl, cacheDuration).Err()
	if err != nil {
		panic(fmt.Sprintf("Failed Saving key url | Error: %v - Shorturl: %s\n", err, shortUrl))

	}
}

func retrieveInitialUrl(shortUrl string) string {
	res, err := storeService.redisClient.Get(ctx, shortUrl).Result()

	if err != nil {
		panic(fmt.Sprintf("Failed RetrieveInitialUrl url | Error: %v - Shorturl: %s\n", err, shortUrl))
	}

	return res
}
