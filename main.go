package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context){
		c.JSON(200, gin.H{
			"message": "Hello Shawtyfy!!",
		})
	})

	error := r.Run(":9808")
	if error != nil{
		panic(fmt.Sprintf("Failed to start web server: Error: %v", error))
	}
}
