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

type UrlFetchRequest struct {
	UserId string `json:"user_id" binding:"required"`
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

	cookie := http.Cookie{
		Name:     "user_id",
		Value:    creationRequest.UserId,
		HttpOnly: true,
	}

	http.SetCookie(ctx.Writer, &cookie)

	host := "http://localhost:9808/"
	ctx.JSON(200, gin.H{
		"message":   "Short url creation successful!",
		"short_url": host + shortUrl,
	})
}

// HTTP redirection function
func HandleShortUrlRedirect(ctx *gin.Context) {
	var fetchRequest UrlFetchRequest

	// if err := ctx.ShouldBindJSON(&fetchRequest); err != nil {
	// 	ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }

	// Read the user_id from the cookie
	cookie, err := ctx.Request.Cookie("user_id")
	if err != nil {
		panic(fmt.Sprint(err))
	}

	fetchRequest.UserId = cookie.Value

	shortUrl := ctx.Param("shortUrl")
	initialUrl, err := store.RetrieveInitialUrl(shortUrl, fetchRequest.UserId)

	if err != nil {
		panic(fmt.Sprint(err))
	}
	// Check if the initial URL is absolute or relative
	if !strings.HasPrefix(initialUrl, "http://") && !strings.HasPrefix(initialUrl, "https://") {
		initialUrl = "http://" + initialUrl
	}

	ctx.Redirect(302, initialUrl)
}

// Fetch total key count in the database
func HandleKeyCount(ctx *gin.Context) {
	count, err := store.RetreiveKeyCount()

	if err != nil {
		panic(fmt.Sprintf("Error retreiving key count | %v", err))
	}

	ctx.JSON(200, gin.H{
		"key_count": count,
	})
}

// Handle Fetch Url by id
func HandleFetchUrlbyId(ctx *gin.Context) {
	var fetchRequest UrlFetchRequest
	if err := ctx.ShouldBindJSON(&fetchRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{
		"user_id": fetchRequest.UserId,
		"urls":    store.FetchUrlsByUserID(fetchRequest.UserId),
	})
}
