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

// populate redis cache with data from dynamodb on cold restart
func PopulateRedisFromDynamoDb(redisClient *redis.Client) error {
	dynamoDBService := InitializeDynamoDB()

	// Querying whole table
	res, err := dynamoDBService.dynamoDBClient.Scan(&dynamodb.ScanInput{
		TableName: aws.String(tableName),
	})
	if err != nil {
		log.Printf("Failed scanning DynamoDb table: %v\n", err)
		return err
	}

	// Iterate and store each item in Redis
	for _, item := range res.Items {
		shortUrl := aws.StringValue(item["shortUrl"].S)
		originalUrl := aws.StringValue(item["originalUrl"].S)
		userid := aws.StringValue(item["userid"].S)

		// Map to redis
		store.SaveUrlMapping(shortUrl, originalUrl, userid)
	}

	log.Println("Successfully populated Redis from Remote")
	return nil
}

func InitializeDynamoDB() *DynamoDBService {
	sesh := session.Must(session.NewSessionWithOptions(
		session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))

	dynamoDBClient := dynamodb.New(sesh)

	// initialize dynamoDB service
	dynamoDBService.dynamoDBClient = dynamoDBClient
	// dynamoDBService = &DynamoDBService{
	// 	dynamoDBClient: dynamoDBClient,
	// }
	// Populate Redis from DynamoDB on cold boot
	storeService := store.InitializeStore()
	if storeService == nil {
		log.Println("Error initializing store")
		return dynamoDBService
	}

	// populate redis from dynamodb on cold boot
	err := PopulateRedisFromDynamoDb(storeService.GetRedisClient())
	if err != nil {
		log.Printf("Error populating Redis from Remote: %v\n", err)
	}

	// listen for new urls in background
	// go ListenForNewURL(store.InitializeStore().GetRedisClient())

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
