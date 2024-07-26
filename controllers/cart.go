package controllers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rodoben/ecommerce/database"
	"github.com/rodoben/ecommerce/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//-------------------cart routes controllers-----------------

func AddToCart() gin.HandlerFunc {
	return func(c *gin.Context) {

		userId := c.Query("userid")
		productId := c.Query("productid")

		ctx, cancel := context.WithTimeout(c, 20*time.Second)
		defer cancel()

		if userId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "UserId not found"})
			return
		}
		if productId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "productId not found"})
			return
		}

		var foundUser models.User
		var productDetails models.UserCart

		fmt.Println("userid", userId)
		fmt.Println("productid", productId)

		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&foundUser)
		if err != nil {
			fmt.Println("dscvsdvfvf")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
			return
		}
		fmt.Println(foundUser)

		productidHex, err := primitive.ObjectIDFromHex(productId)
		if err != nil {
			log.Println("caanot convert it to objectidHex")
		}
		fmt.Println("dscvsdvrfgrr1111fvf")
		err = productCollection.FindOne(ctx, bson.M{"product_id": productidHex}).Decode(&productDetails)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "product not found"})
			return
		}
		productExist := false

		for id, _ := range foundUser.UserCart {
			if item, ok := foundUser.UserCart[productidHex]; ok {
				fmt.Println("item exist and updating count", id)
				item.Quantity++
				productExist = true
				break
			}
			// else {
			// 	foundUser.UserCart[productidHex] = productDetails
			// }

		}

		if !productExist {
			fmt.Println("I am insertinng a record with id", productidHex)
			foundUser.UserCart[productidHex] = productDetails

		}

		// for _, cartItem := range foundUser.UserCart {
		// 	if item, ok := cartItem[productidHex]; ok {
		// 		// Product already exists in the cart, update the quantity
		// 		item.Quantity++
		// 		productExists = true
		// 		break
		// 	}
		// }

		// if !productExists {
		// 	// If the product does not exist in the cart, add it with quantity 1
		// 	newCartItem := productDetails
		// 	newCartMap := map[primitive.ObjectID]models.UserCart{productidHex: newCartItem}
		// 	foundUser.UserCart = append(foundUser.UserCart, newCartMap)
		// } else {
		// 	// Update the quantity of the existing product in the cart
		// 	for i, cartItem := range foundUser.UserCart {
		// 		if item, ok := cartItem[productidHex]; ok {
		// 			// Update the quantity
		// 			foundUser.UserCart[i][productidHex] = item
		// 			break
		// 		}
		// 	}
		// }

		// Update the cart in the database
		updateCart := bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "usercart", Value: foundUser.UserCart},
			}},
		}

		fmt.Println(foundUser.UserCart)

		_, err = userCollection.UpdateOne(ctx, bson.M{"user_id": userId}, updateCart)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user cart not updated"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": foundUser.UserCart})

	}
}

func UpdateCartQuantity() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(c, 10*time.Second)
		defer cancel()
		userid := c.Param("userid")
		//	productId := c.Param("productid")
		var foundUser models.User
		err := userCollection.FindOne(ctx, bson.M{"user_id": userid}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
			return
		}
		//productidHex, err := primitive.ObjectIDFromHex(productId)
		if err != nil {
			log.Println("caanot convert it to objectidHex")
		}

		// for _, v := range foundUser.UserCart {
		// 	if item, ok := v[productidHex]; ok {
		// 		item.Quantity++
		// 		c.JSON(http.StatusBadRequest, gin.H{"error": "product already added to cart"})
		// 		return
		// 	}

		// }

		updateCart := bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "usercart", Value: foundUser.UserCart},
			}},
		}

		_, err = userCollection.UpdateOne(ctx, bson.M{"user_id": userid}, updateCart)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "usercart not Updated"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": foundUser.UserCart})
	}
}

func DeleteFromCart() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

func DeleteFromCartByID() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(c, 10*time.Second)
		defer cancel()
		productId := c.Query("productid")
		userid := c.Param("userid")
		var foundUser models.User
		err := userCollection.FindOne(ctx, bson.M{"user_id": userid}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
			return
		}
		// Update the user's cart in the database
		update := bson.M{"$unset": bson.M{"usercart." + productId: ""}}
		_, err = userCollection.UpdateOne(ctx, bson.M{"user_id": userid}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot update user cart"})
			return
		}
		c.JSON(http.StatusAccepted, gin.H{"succesful": "susseful", "datadeleted": productId})
	}
}

func (app *Application) InstantBuy() gin.HandlerFunc {
	return func(c *gin.Context) {

		var parameters models.Payment

		if err := c.BindJSON(&parameters); err != nil {
			c.AbortWithError(http.StatusBadRequest, errors.New("unbale to unmarshal the payload"))
		}

		userId := c.Query("userid")
		if userId == " " {
			c.AbortWithError(http.StatusBadRequest, errors.New("UserId is empty"))
		}

		productId := c.Query("productid")
		if productId == " " {
			c.AbortWithError(http.StatusBadRequest, errors.New("Product id is empty"))
		}

		var ctx, cancel = context.WithTimeout(c, 10*time.Second)
		defer cancel()

		err := database.Instantbuy(ctx, *app.ProductCollection, *app.UserCollection, userId, productId)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, errors.New("unable to process for Instant Buy"))
		}

		c.IndentedJSON(http.StatusOK, "succesful")

	}

}
