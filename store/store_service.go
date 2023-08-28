package store

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis"
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
)


// Cache expiration duration
const cacheDuration = 6 * time.Hour