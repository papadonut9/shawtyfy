package generator

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

const userId = "68a099a-610d-4237-89d9-2b6d099a20e2"

// unit test for checking short url
func TestShortLinkGenerator(t *testing.T) {
	originalUrl_1 := "https://youtu.be/dQw4w9WgXcQ"
	shortUrl_1 := generateShortUrl(originalUrl_1, userId)

	originalUrl_2 := "github.com/papadonut9"
	shortUrl_2 := generateShortUrl(originalUrl_2, userId)

	originalUrl_3 := "https://youtu.be/_sPXXL9YXok"
	shortUrl_3 := generateShortUrl(originalUrl_3, userId)

	originalUrl_4 := "https://open.spotify.com/track/4QfVhHGCPVYLiQEnaeknJe"
	shortUrl_4 := generateShortUrl(originalUrl_4, userId)

	originalUrl_5 := "https://open.spotify.com/track/2LMkwUfqC6S6s6qDVlEuzV"
	shortUrl_5 := generateShortUrl(originalUrl_5, userId)

	assert.Equal(t, shortUrl_1, "EUsP4nKs")
	assert.Equal(t, shortUrl_2, "SQkReYrz")
	assert.Equal(t, shortUrl_3, "jQe1gbRB")
	assert.Equal(t, shortUrl_4, "PviA5RT6")
	assert.Equal(t, shortUrl_5, "5rv35y2T")
}
