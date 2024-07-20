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

func (app *Application) AddAdress() gin.HandlerFunc {
	return func(c *gin.Context) {

		var foundUser models.User
		var address models.Address
		var ctx, cancel = context.WithTimeout(c, 10*time.Second)
		defer cancel()
		userid := c.Query("userid")

		if userid == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Userid is empty"})
			return
		}
		fmt.Println("userid", userid)
		fmt.Println("founduser", foundUser)

		err := userCollection.FindOne(ctx, bson.M{"user_id": userid}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to fetch user record!"})
			return
		}
		fmt.Println("founduser", foundUser)

		if err := c.BindJSON(&address); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "unable to bind json to the struct!"})
			return
		}

		address.Address_ID = primitive.NewObjectID()

		foundUser.Address = append(foundUser.Address, address)

		updateAddress := bson.M{"$push": bson.M{"address": address}}

		_, err = app.UserCollection.UpdateOne(ctx, bson.M{"user_id": userid}, updateAddress)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to update adress"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success!": foundUser})
	}
}

func DeleteAdress() gin.HandlerFunc {

	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(c, 10*time.Second)
		defer cancel()
		var foundUser models.User

		addressid := c.Param("id")
		userid := c.Query("userid")

		err := userCollection.FindOne(ctx, bson.M{"user_id": userid}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to marshal in struct"})
			return
		}
		addressFound := false
		addressidHex, err := primitive.ObjectIDFromHex(addressid)
		if err != nil {
			log.Println("caanot convert it to objectidHex")
		}

		for i, v := range foundUser.Address {
			fmt.Println(v.Address_ID)
			if v.Address_ID == addressidHex {
				addressFound = true
				foundUser.Address = append(foundUser.Address[:i], foundUser.Address[i+1:]...)
				break
			}
		}
		if !addressFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": "unable to find address"})
			return
		}
		fmt.Println(addressFound)

		//deleteAddress := bson.M{"$pull": bson.M{"address_id": addressidHex}}
		deleteAddress := bson.M{"$pull": bson.M{"address": bson.M{"address_id": addressidHex}}}
		_, err = userCollection.UpdateOne(ctx, bson.M{"user_id": userid}, deleteAddress)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to delete address"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"deleted": addressid})

	}
}
