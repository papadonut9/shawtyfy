package dynamodb

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/papadonut9/shawtyfy/store"
)

// ListenForNewURL listens for new URL added events from Redis and saves to DynamoDB
func ListenForNewURL(redisClient *redis.Client) {
	pubsub := redisClient.Subscribe(context.Background(), "new_url_added")
	defer pubsub.Close()

	_, err := pubsub.Receive(context.Background())
	if err != nil {
		log.Fatalf("Error subscribing to Redis channel: %v\n", err)
	}

	ch := pubsub.Channel()

	for msg := range ch {
		shortURL := msg.Payload
		originalURL := store.RetrieveInitialUrl(shortURL) // Assuming this function fetches the original URL

		// Fetch user ID based on short URL from Redis or any other source
		userID := store.RetreiveUserId(shortURL)

		// Save to DynamoDB
		err := SaveUrlMapping( shortURL, originalURL, userID)
		if err != nil {
			log.Printf("Error saving to DynamoDB: %v\n", err)
		}
	}
}
