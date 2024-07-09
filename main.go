package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rodoben/ecommerce/routes"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	port := os.Getenv("Port")
	if port == "" {
		log.Fatal("Cannot fetch the port form the env", port)
	}

	router := gin.New()

	routes.UserRoutes(router)
	routes.CartRoutes(router)

	panic(router.Run(":" + port))

}
