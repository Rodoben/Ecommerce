package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rodoben/ecommerce/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (app *Application) AddToWishList() gin.HandlerFunc {

	return func(c *gin.Context) {
		var foundUser models.User
		var productDetails models.WishList
		var ctx, cancel = context.WithTimeout(c, 10*time.Second)

		defer cancel()

		userId := c.Query("userid")
		productId := c.Query("productid")
		if userId == " " {
			c.JSON(http.StatusBadRequest, gin.H{"error": "userid is empty"})
			return
		}
		if productId == " " {
			c.JSON(http.StatusBadRequest, gin.H{"error": "product id is empty"})
			return
		}
		err := app.UserCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to parse user details"})
			return
		}

		productIdHex, err := primitive.ObjectIDFromHex(productId)
		if err != nil {
			log.Println("error unable to convert productid to product Hex")
		}

		err = app.ProductCollection.FindOne(ctx, bson.M{"product_id": productIdHex}).Decode(&productDetails)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to parse product details"})
		}

		if _, ok := foundUser.WishList[productIdHex]; ok {
			c.JSON(http.StatusOK, gin.H{"message": "product already added to wish List"})
			return
		} else {
			foundUser.WishList[productIdHex] = productDetails
		}

		updateWishList := bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "wishlist", Value: foundUser.WishList},
			}},
		}

		_, err = app.UserCollection.UpdateOne(ctx, bson.M{"user_id": userId}, updateWishList)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to update"})
			return

		}
		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": foundUser.WishList})
	}

}

func (app *Application) DeleteFromWishlist() gin.HandlerFunc {

	return func(c *gin.Context) {

		var foundUser models.User
		var ctx, cancel = context.WithTimeout(c, 10*time.Second)
		defer cancel()

		userId := c.Param("userid")
		productId := c.Query("productid")

		if userId == " " {
			c.JSON(http.StatusBadRequest, gin.H{"error": " userid is missing"})
			return
		}

		if productId == " " {
			c.JSON(http.StatusBadRequest, gin.H{"error": "product id is missing"})
			return
		}

		err := app.UserCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to parse user details"})
			return
		}

		productIdHex, err := primitive.ObjectIDFromHex(productId)
		if err != nil {
			log.Println("Unable to convert productid to productidHex")
		}

		if _, ok := foundUser.WishList[productIdHex]; !ok {
			c.JSON(http.StatusBadRequest, gin.H{"message": "wishlist already exist"})
			return
		}

		fmt.Println("productid", productId)
		fmt.Println("userid", userId)

		deletewishlist := bson.M{"$unset": bson.M{"wishlist." + productId: ""}}

		_, err = app.UserCollection.UpdateOne(ctx, bson.M{"user_id": userId}, deletewishlist)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to delete the record"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"deleteduser": productId})

	}
}
