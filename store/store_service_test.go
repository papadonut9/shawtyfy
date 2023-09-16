package store

import (
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
func testSetandGet(t *testing.T) {
	originalUrl := "https://youtu.be/dQw4w9WgXcQ"
	userid := "950c182b-1745-4aa9-b872-d8c558fadc8d"
	shortUrl := "https://rb.gy/lkoyw"

	// persistent data mapping
	SaveUrlMapping(originalUrl, userid, shortUrl)

	// fetch original url
	receivedUrl := RetrieveInitialUrl(shortUrl)

	assert.Equal(t, originalUrl, receivedUrl)
}
