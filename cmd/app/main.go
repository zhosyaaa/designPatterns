package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"pattern/internal/app"
)

func main() {
	routes := app.SetupApp()
	app := gin.Default()
	routes.SetupRouter(app)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	err := app.Run(":" + port)
	if err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
	fmt.Printf("Server is running on port %s\n", port)
}
