package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rodoben/ecommerce/controllers"
)

func UserRoutes(incomingUserRoutes *gin.Engine, app *controllers.Application) {
	incomingUserRoutes.POST("/users/signin", controllers.SignIn())
	incomingUserRoutes.POST("/users/signup", app.SignUp())
	incomingUserRoutes.POST("/admin/addproduct", controllers.Addproduct())
	incomingUserRoutes.GET("/users/productview", app.Productview())
	incomingUserRoutes.GET("/users/search", controllers.SearchProduct())
}

func CartRoutes(incomingCartRoutes *gin.Engine) {
	incomingCartRoutes.POST("/users/addtocart", controllers.AddToCart())
	incomingCartRoutes.DELETE("/users/deletefromcart/{id}", controllers.DeleteFromCartByID())
	incomingCartRoutes.DELETE("/users/deletefromcart", controllers.DeleteFromCart())
}
