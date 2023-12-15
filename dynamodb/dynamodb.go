package dynamodb

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/papadonut9/shawtyfy/store"
)

// DynamoDB wrapper
type DynamoDBService struct {
	dynamoDBClient *dynamodb.DynamoDB
}

// DynamoDBEntry represents the structure of data in DynamoDB
type DynamoDBEntry struct {
	ShortURL    string
	OriginalURL string
	UserID      string
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

	/////////////////////////////////////////////////
	// Call FetchAndPopulateRedis to populate Redis with DynamoDB data
	err := FetchAndPopulateRedis()
	if err != nil {
		log.Printf("Error fetching and populating Redis: %v\n", err)
	}

	/////////////////////////////////////////////////

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

// Fetches data from dynamodb on cold restart and populates data to redis
func FetchAndPopulateRedis() error {
	// Create a new DynamoDB session
	sess := session.Must(session.NewSessionWithOptions(
		session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))

	// Create a DynamoDB client
	dynamoDBClient := dynamodb.New(sess)

	// Define the projection expression to get only specific attributes
	projection := expression.NamesList(expression.Name("shortUrl"), expression.Name("url"), expression.Name("userid"))

	// Build the DynamoDB expression
	expr, err := expression.NewBuilder().WithProjection(projection).Build()
	if err != nil {
		log.Printf("Error building expression: %v\n", err)
		return err
	}

	// Scan the DynamoDB table
	input := &dynamodb.ScanInput{
		TableName:                 aws.String(tableName),
		ProjectionExpression:      expr.Projection(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	}

	result, err := dynamoDBClient.ScanWithContext(context.Background(), input)
	if err != nil {
		log.Printf("Error scanning DynamoDB table: %v\n", err)
		return err
	}

	// Iterate through the scan result and populate Redis
	for _, item := range result.Items {
		shortURL := *item["shortUrl"].S
		originalURL := *item["url"].S
		userID := *item["userid"].S

		// Save to Redis
		store.SaveUrlMapping(shortURL, originalURL, userID)
	}

	return nil
}

func RemoveUrl(shortUrl string) (string, error) {

	if shortUrl == "" {
		return "EMPTY STRING", fmt.Errorf("Short URL cannot be empty")
	}

	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"shortUrl": {
				S: aws.String(shortUrl),
			},
		},
		TableName: aws.String(tableName),
	}

	_, err := dynamoDBService.dynamoDBClient.DeleteItemWithContext(context.Background(), input)
	if err != nil {
		log.Printf("Failed Deleting from DynamoDB: %v\n", err)
		return "", err
	}

	return "OK", nil
}
