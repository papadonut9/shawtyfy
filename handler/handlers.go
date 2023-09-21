package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/papadonut9/shawtyfy/generator"
	"github.com/papadonut9/shawtyfy/store"
)

type UrlCreateRequest struct {
	LongUrl string `json:"long_url" binding:"required"`
	UserId  string `json:"user_id" binding:"required"`
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

func HandleKeyCount(ctx *gin.Context){
	count, err := store.RetreiveKeyCount()
	
	if err != nil{
		panic(fmt.Sprintf("Error retreiving key count"))
	}
	
	// host := "http://localhost:9808/"
	ctx.JSON(200, gin.H{
		"key_count": count,
	})
}
