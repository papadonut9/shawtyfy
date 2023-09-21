package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/papadonut9/shawtyfy/endpoint"
	"github.com/papadonut9/shawtyfy/store"
)

// Main function
func main() {
	route := gin.Default()
	
	// Calling setupRoutes to define endpoints
	endpoint.SetupRoutes(route)

	// store initialization
	store.InitializeStore()

	error := route.Run(":9808")
	if error != nil{
		panic(fmt.Sprintf("Failed to start web server: Error: %v", error))
	}
}
