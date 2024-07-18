package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rodoben/ecommerce/models"
	"go.mongodb.org/mongo-driver/bson"
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

		err := app.UserCollection.FindOne(ctx, bson.M{"userid": userid}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found!"})
			return
		}

		if err := c.BindJSON(&address); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "unable to bind json to the struct!"})
			return
		}

	}
}

func DeleteAdress() gin.HandlerFunc {

	return func(ctx *gin.Context) {

	}
}
