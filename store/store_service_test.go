package store

import (
	"testing"
	"time"

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
	const cacheDuration = 24 * time.Hour

	// persistent data mapping
	// SaveUrlMapping(originalUrl, userid, shortUrl)
	SaveUrlMapping(shortUrl, originalUrl, userid)

	// fetch original url
	receivedUrl := RetrieveInitialUrl(shortUrl)

	assert.Equal(t, originalUrl, receivedUrl)
}

func TestFetchMetadata(t *testing.T) {
	originalUrl := "https://youtu.be/dQw4w9WgXcQ"
	userid := "950c182b-1745-4aa9-b872-d8c558fadc8d"
	shortUrl := "lkoyw"

	SaveUrlMapping(shortUrl, originalUrl, userid)

	receivedUserId, receivedUrl := FetchMetadata(shortUrl)

	assert.Equal(t, userid, receivedUserId)
	assert.Equal(t, originalUrl, receivedUrl)
}
