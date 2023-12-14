package store

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testStoreService = &StorageService{}

// Set up test shell
func init() {
	testStoreService = InitializeStore()
}

// Test store service init
func testStoreInit(t *testing.T) {
	assert.True(t, testStoreService.redisClient != nil)
}

// Storage API Tests
func TestSetandGet(t *testing.T) {
	originalUrl := "https://youtu.be/dQw4w9WgXcQ"
	userid := "950c182b-1745-4aa9-b872-d8c558fadc8d"
	shortUrl := "lkoyw"

	// persistent data mapping
	// SaveUrlMapping(originalUrl, userid, shortUrl)
	SaveUrlMapping(shortUrl, originalUrl, userid)

	// fetch original url
	receivedUrl, err := RetrieveInitialUrl(shortUrl, userid)

	if err != nil{
		panic(fmt.Sprint(err))
	}

	assert.Equal(t, originalUrl, receivedUrl)
}
