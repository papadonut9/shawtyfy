package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/papadonut9/shawtyfy/dynamodb"
	"github.com/papadonut9/shawtyfy/generator"
	"github.com/papadonut9/shawtyfy/store"
)

type UrlCreateRequest struct {
	LongUrl string `json:"long_url" binding:"required"`
	UserId  string `json:"user_id" binding:"required"`
}

type UrlDeleteRequest struct {
	ShortUrl string `json:"short_url" binding:"required"`
}

// Handler stubs
func CreateShortUrl(ctx *gin.Context) {
	var creationRequest UrlCreateRequest
	if err := ctx.ShouldBindJSON(&creationRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	shortUrl := generator.GenerateShortLink(creationRequest.LongUrl, creationRequest.UserId)
	store.SaveUrlMapping(shortUrl, creationRequest.LongUrl, creationRequest.UserId)

	host := "http://localhost:9808/"
	ctx.JSON(200, gin.H{
		"message":   "Short url creation successful!",
		"short_url": host + shortUrl,
	})
}

// HTTP redirection function
func HandleShortUrlRedirect(ctx *gin.Context) {
	shortUrl := ctx.Param("shortUrl")
	initialUrl := store.RetrieveInitialUrl(shortUrl)

	// Check if the initial URL is absolute or relative
	if !strings.HasPrefix(initialUrl, "http://") && !strings.HasPrefix(initialUrl, "https://") {
		initialUrl = "http://" + initialUrl
	}
	ctx.Redirect(302, initialUrl)
}

// fetch total number of keys
func HandleKeyCount(ctx *gin.Context) {
	count, err := store.RetreiveKeyCount()

	if err != nil {
		panic(fmt.Sprintf("Error retreiving key count: %v\n", err))
	}

	// host := "http://localhost:9808/"
	ctx.JSON(200, gin.H{
		"key_count": count,
	})
}

// fetch all keys
func RetreiveAllKeys(ctx *gin.Context) {
	keys, err := store.RetreiveAllKeys()

	if err != nil {
		panic(fmt.Sprintf("Error retreiving keys: %v\n", err))
	}

	ctx.JSON(200, gin.H{
		"keys": keys,
	})
}

func DeleteKey(ctx *gin.Context) {
	var deletionRequest UrlDeleteRequest

	if err := ctx.ShouldBindJSON(&deletionRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var message string
	var redisStatus, dynamodbStatus string

	key := deletionRequest.ShortUrl
	dynamodbResult, dynamodbErr := dynamodb.RemoveUrl(key)
	redisResult, redisErr := store.RemoveKey(key)

	// successful deletes
	if redisErr == nil && dynamodbErr == nil {
		message = "deletion successful"
		redisStatus = "OK"
		dynamodbStatus = "OK"

		// if redis fails
	} else if redisErr != nil && dynamodbErr == nil {
		message = "deletion unsuccessful"
		redisStatus = "FAIL"
		dynamodbStatus = "OK"

		// if dynamodb fails
	} else if redisErr == nil && dynamodbErr != nil {
		message = "deletion unsuccessful"
		redisStatus = "OK"
		dynamodbStatus = "FAIL"

		// everything goes haywire
	} else {
		message = "deletion unsuccessful"
		redisStatus = "FAIL"
		dynamodbStatus = "FAIL"
	}

	ctx.JSON(200, gin.H{
		"message":         message,
		"redis":           redisStatus,
		"dynamodb":        dynamodbStatus,
		"redis_result":    redisResult,
		"dynamodb_result": dynamodbResult,
	})
}

/*	MESSAGE STRUCTURE
on success
{
	"message": "deletion successful",
	"redis": "OK",
	"dynamodb": "OK"
}

partial delete(if deletion does not occur from either redis or dynamodb)
{
	"message": "deletion unsuccessful",
	"redis": "FAIL",
	"dynamodb": "OK"
}

unsuccessful delete
{
	"message": "deletion unsuccessful",
	"redis": "FAIL",
	"dynamodb": "FAIL"
}
*/
