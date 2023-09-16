package handler

import (
	"net/http"
	"github.com/papadonut9/shawtyfy/generator"
	"github.com/papadonut9/shawtyfy/store"
	"github.com/gin-gonic/gin"
)

// Request model definition
type urlCreateRequest struct {
	longUrl string `json: "long_url" binding:"required"`
	userId  string `json: "user_id" binding:"required"`
}

// Handler stubs
func CreateShortUrl(ctx *gin.Context) {
	var creationRequest urlCreateRequest
	if err := ctx.ShouldBindJSON(&creationRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	shortUrl := generator.GenerateShortLink(creationRequest.longUrl, creationRequest.userId)
	store.SaveUrlMapping(shortUrl, creationRequest.longUrl, creationRequest.userId)

	host := "http://localhost:9808"
	ctx.JSON(200, gin.H{
		"message": "Short url creation successful!",
		"short_url": host + shortUrl,
	})
}

func HandleShortUrlRedirect(ctx *gin.Context) {
	// TODO: add implementation
}
