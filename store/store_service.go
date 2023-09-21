package store

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// Raw wrapper
type StorageService struct {
	redisClient *redis.Client
}

// type UrlData struct {
// 	url    string
// 	userid string
// }

// high level declaration
var (
	storeService = &StorageService{}
	ctx          = context.Background()
	// dbstr        = os.Getenv("REDIS_DB")
)

// environment variable string to integer conversion
// redis_db := strconv.Atoi(dbstr);

// Cache expiration duration
const cacheDuration = 24 * time.Hour

// Store service with pointer return
func InitializeStore() *StorageService {
	redisClient := redis.NewClient(&redis.Options{
		// Addr:     os.Getenv("REDIS_ADDR"),
		Addr:     "127.0.0.1:6379",
		Password: "",
		// DB:       strconv.Atoi(os.Getenv(REDIS_DB)),
		DB: 0,
	})

	pong, err := redisClient.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("Error init redis: %v", err))
	}

	fmt.Printf("\nRedis Started successfully: pong message = {%s}", pong)
	storeService.redisClient = redisClient

	return storeService
}

// Save URL mapping by taking the shortened url, original url and user id
func SaveUrlMapping(shortUrl string, originalUrl string, userid string) {

	// data := map[string]interface{}{
	// 	"url":    originalUrl,
	// 	"userid": userid,
	// }
	// jsonData, errj := json.Marshal(data)
	//   if errj != nil {
	//       panic(fmt.Sprintf("Failed to marshal data to JSON: %v", errj))
	//   }

	// set hash to the database
	// func
	err := storeService.redisClient.HSet(ctx, shortUrl, "url", originalUrl, "userid", userid).Err()
	// err := storeService.redisClient.HMSet(ctx, shortUrl, data, cacheDuration).Err()
	// converting map to json string
	if err != nil {

		panic(fmt.Sprintf("Failed Saving key url | Error: %v - Shorturl: %s\n", err, shortUrl))

	}

	// Set key
	// err := storeService.redisClient.Set(ctx, shortUrl, originalUrl, cacheDuration).Err()
	// if err != nil {
	// 	panic(fmt.Sprintf("Failed Saving key url | Error: %v - Shorturl: %s\n", err, shortUrl))
	// }
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

func RetreiveKeyCount() (int64, error){
	res, err := storeService.redisClient.DBSize(ctx).Result()

	if err != nil {
		panic(fmt.Sprintf("Failed retreiving keys"))
	}

	return res, err;
}