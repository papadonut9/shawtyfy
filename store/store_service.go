package store

import (
	"context"
	"errors"
	"fmt"

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
	// dbstr        = os.Getenv("REDIS_DB")
)

// Cache expiration duration
// const cacheDuration = 24 * time.Hour

// Store service with pointer return
func InitializeStore() *StorageService {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
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

// Returns the Redis client
func (s *StorageService) GetRedisClient() *redis.Client {
	return s.redisClient
}

// Save URL mapping by taking the shortened url, original url and user id
func SaveUrlMapping(shortUrl string, originalUrl string, userid string) {

	// set hash to the database

	err := storeService.redisClient.HSet(ctx, shortUrl, "url", originalUrl, "userid", userid).Err()

	if err != nil {
		panic(fmt.Sprintf("Failed Saving key url | Error: %v - Shorturl: %s\n", err, shortUrl))
	}
	storeService.redisClient.Expire(ctx, shortUrl, cacheDuration)
	// if expErr != nil {
	// 	panic(fmt.Sprintf("Failed setting key expiry on key %s| Error: %v ", shortUrl, expErr))
	// }

	// Publish message on redis Channel
	err = storeService.redisClient.Publish(ctx, "new_url_added", shortUrl).Err()
	if err != nil {
		fmt.Printf("Error publishing message: %v", err)
	}
}

// func RetrieveInitialUrl(shortUrl string, userId string) string {
// 	// res, err := storeService.redisClient.Get(ctx, shortUrl).Result()
// 	res, err := storeService.redisClient.HGet(ctx, userId, shortUrl).Result()

// 	// fetch strings
// 	if err != nil {
// 		panic(fmt.Sprintf("Failed RetrieveInitialUrl url | Error: %v - Shorturl: %s\n", err, shortUrl))
// 	}

// 	return res
// }

func RetrieveInitialUrl(shortUrl string, userId string) (string, error) {
	// res, err := storeService.redisClient.Get(ctx, shortUrl).Result()
	res, err := storeService.redisClient.HGet(ctx, userId, shortUrl).Result()

	if err == redis.Nil {
		// Key not found in Redis
		return "", errors.New("short url not found")
	} else if err != nil {
		// Handle other Redis errors
		return "", err
	}

	return res, nil
}

func RetreiveKeyCount() (int64, error) {
	res, err := storeService.redisClient.DBSize(ctx).Result()

	if err != nil {
		panic(fmt.Sprintf("Failed retreiving keys | error: %v", err))
	}

	return res, err
}


func RetreiveUserId(shortUrl string) string {
	res, err := storeService.redisClient.HGet(ctx, shortUrl, "userid").Result()

	// fetch strings
	if err != nil {
		panic(fmt.Sprintf("Failed RetrieveUserId id | Error: %v - Shorturl: %s\n", err, shortUrl))
	}

	return res
}

func RetreiveKeyCount() (int64, error) {
	res, err := storeService.redisClient.DBSize(ctx).Result()


func FetchUrlsByUserID(userid string) map[string]string {
	// Use the HGetAll function to retrieve all key-value pairs associated with the user ID.
	urls, err := storeService.redisClient.HGetAll(ctx, userid).Result()
	if err != nil {
		panic(fmt.Sprintf("Failed retreiving keys: %v", err))
	}

	return res, err

}
