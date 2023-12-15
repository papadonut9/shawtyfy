package main

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/papadonut9/shawtyfy/dynamodb"
	"github.com/papadonut9/shawtyfy/endpoint"
	"github.com/papadonut9/shawtyfy/store"
)

// Main function
func main() {
	route := gin.Default()

	// Middleware to allow CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"} 
	config.AllowHeaders = []string{"Origin", "Content-Type"}
	route.Use(cors.New(config))

	// Calling setupRoutes to define endpoints
	endpoint.SetupRoutes(route)

	// initialize dynamodb service
	dynamodb.InitializeDynamoDB()

	// store initialization
	store.InitializeStore()

	// fetch context from middleware
	route.Use(
		func(ctx *gin.Context) {
			go dynamodb.ListenForNewUrl(ctx, store.InitializeStore().GetRedisClient())
		})

	// start webserver
	error := route.Run(":9808")
	if error != nil {
		panic(fmt.Sprintf("Failed to start web server: Error: %v", error))
	}
}
