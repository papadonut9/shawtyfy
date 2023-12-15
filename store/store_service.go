package store

import (
	"context"
	"fmt"
	"log"
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
	// dbstr        = os.Getenv("REDIS_DB")
)

// Cache expiration duration
const cacheDuration = 24 * time.Hour

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

	// Publish message on redis Channel
	err = storeService.redisClient.Publish(ctx, "new_url_added", shortUrl).Err()
	if err != nil {
		fmt.Printf("Error publishing message: %v", err)
	}
}

func RetrieveInitialUrl(shortUrl string) string {
	// res, err := storeService.redisClient.Get(ctx, shortUrl).Result()
	res, err := storeService.redisClient.HGet(ctx, shortUrl, "url").Result()

	// fetch strings
	if err != nil {
		panic(fmt.Sprintf("Failed RetrieveInitialUrl url | Error: %v - Shorturl: %s\n", err, shortUrl))
	}

	return res
}

func RetreiveUserId(shortUrl string) string {
	res, err := storeService.redisClient.HGet(ctx, shortUrl, "userid").Result()

	// fetch strings
	if err != nil {
		panic(fmt.Sprintf("Failed RetrieveUserId id | Error: %v - Shorturl: %s\n", err, shortUrl))
	}

	return res
}

// Returns total number of keys in the redis.
func RetreiveKeyCount() (int64, error) {
	res, err := storeService.redisClient.DBSize(ctx).Result()

	if err != nil {
		panic(fmt.Sprintf("Failed retreiving keys: %v", err))
	}

	return res, err
}

// Returns all the keys in the cache.
func RetreiveAllKeys() ([]string, error) {
	keys, err := storeService.redisClient.Keys(ctx, "*").Result()

	if err != nil {
		log.Println("Error retrieving keys: ", err)
	}

	return keys, err
}

// Removes the key from cache
func RemoveKey(key string) (int64, error) {
	res, err := storeService.redisClient.Del(ctx, key).Result()

	if err != nil {
		panic(fmt.Sprintf("Error deleting key: %v", err))
	}

	return res, err
}

func FetchMetadata(key string) (userid string, url string) {
	userid = RetreiveUserId(key)
	url = RetrieveInitialUrl(key)

	return userid, url
}
