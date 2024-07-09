package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rodoben/ecommerce/controllers"
)

func UserRoutes(incomingUserRoutes *gin.Engine) {
	incomingUserRoutes.POST("/users/signin", controllers.SignIn())
	incomingUserRoutes.POST("/users/signup", controllers.SignUp())
	incomingUserRoutes.POST("/admin/addproduct", controllers.Addproduct())
	incomingUserRoutes.GET("/users/productview", controllers.Productview())
	incomingUserRoutes.GET("/users/search", controllers.Search())
}

func CartRoutes(incomingCartRoutes *gin.Engine) {
	incomingCartRoutes.POST("/users/addtocart", controllers.AddToCart())
	incomingCartRoutes.DELETE("/users/deletefromcart/{id}", controllers.DeleteFromCartByID())
	incomingCartRoutes.DELETE("/users/deletefromcart/{id}", controllers.DeleteFromCart())
}
