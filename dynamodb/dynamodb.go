package dynamodb

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/papadonut9/shawtyfy/store"
)

// DynamoDB wrapper
type DynamoDBService struct {
	dynamoDBClient *dynamodb.DynamoDB
}

// High level declaration
var (
	dynamoDBService = &DynamoDBService{}
	// ctx = context.Background()
)

const tableName = "shawtyfy-dev"

func InitializeDynamoDB() *DynamoDBService {
	sesh := session.Must(session.NewSessionWithOptions(
		session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))

	dynamoDBClient := dynamodb.New(sesh)

	dynamoDBService.dynamoDBClient = dynamoDBClient

	// listen for new urls in background
	go ListenForNewURL(store.InitializeStore().GetRedisClient())

	return dynamoDBService
}

// save url mapping to dynamodb
func SaveUrlMapping(shortUrl, originalUrl, userid string) error {
	// stdContext := ctx.Request.Context()

	input := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"shortUrl": {
				S: aws.String(shortUrl),
			},
			"url": {
				S: aws.String(originalUrl),
			},
			"userid": {
				S: aws.String(userid),
			},
		},
		TableName: aws.String(tableName),
	}

	_, err := dynamoDBService.dynamoDBClient.PutItemWithContext(context.Background(), input)
	if err != nil {
		log.Printf("Failed saving to DynamoDB: %v\n", err)
		return err
	}

	return nil
}

// Listens to redis subscription for new push events
func ListenForNewUrl(ctx *gin.Context, redisClient *redis.Client) {

	pubsub := redisClient.Subscribe(ctx, "new_url_added")
	defer pubsub.Close()

	_, err := pubsub.Receive(ctx)
	if err != nil {
		log.Panic("Error subscribing to Redis Channel: ", err)
		// Make sure it isn't fatal
	}

	ch := pubsub.Channel()

	for msg := range ch {
		shortUrl := msg.Payload
		originalUrl := store.RetrieveInitialUrl(shortUrl)

		userid := store.RetreiveUserId(shortUrl)

		err := SaveUrlMapping(shortUrl, originalUrl, userid)

		if err != nil {
			log.Printf("Error saving to DynamoDb: %v\n", err)
		}

	}

}
