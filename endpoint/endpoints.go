package endpoint

import (
    "github.com/gin-gonic/gin"
    "github.com/papadonut9/shawtyfy/handler"
)

// Initialize and configure all the routes for the application.
func SetupRoutes(r *gin.Engine) {
    // Root endpoint
    r.GET("/", func(ctx *gin.Context) {
        ctx.JSON(200, gin.H{
            "message": "Hello Shawtyfy!!",
        })
    })

    // Endpoint to create a short URL
    r.POST("/create-short-url", func(ctx *gin.Context) {
        handler.CreateShortUrl(ctx)
    })

    // Endpoint to handle short URL redirection
    r.GET("/shortUrl", func(ctx *gin.Context) {
        handler.HandleShortUrlRedirect(ctx)
    })

    // Add more endpoints as needed
}
