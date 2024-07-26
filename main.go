package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rodoben/ecommerce/controllers"
	"github.com/rodoben/ecommerce/database"
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

	dbConnect := database.DBConnect()

	app := controllers.NewApplication(database.UserData(dbConnect, "userdata"), database.ProductData(dbConnect, "productdata"))
	router := gin.New()
	routes.UserRoutes(router, app)
	routes.CartRoutes(router)
	routes.AddressRoutes(router, app)
	routes.WishListRoutes(router, app)
	routes.OrderRoutes(router, app)

	panic(router.Run(":" + port))

}
