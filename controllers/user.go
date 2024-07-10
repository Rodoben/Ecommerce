package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rodoben/ecommerce/database"
	"github.com/rodoben/ecommerce/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

//-------------------user routes controllers-----------------

var userCollection *mongo.Collection = database.UserData(database.DBConnect(), "userdata")

func SignIn() gin.HandlerFunc {

	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(c, 20*time.Second)
		defer cancel()
		var user models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "unable to parse the body in the model"})
			return
		}

		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email, "phonenumber": user.PhoneNumber})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "error fetching information from database"})
			return
		}

		fmt.Println("count", count)
		if count != 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user does not exists"})
			return
		}

		c.JSON(http.StatusAccepted, gin.H{"Success": "signed in!",
			"status code": http.StatusAccepted})
	}
}

func (app *Application) SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		//context deadline
		fmt.Println("1")
		var ctx, cancel = context.WithTimeout(context.Background(), time.Second*100)
		defer cancel()

		// bind json in model
		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		fmt.Println("2")
		//validate struct
		validate := validator.New()
		validateError := validate.Struct(user)
		if validateError != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validateError.Error()})
			return
		}
		fmt.Println("3")

		//check if the email exists in database

		fmt.Println(user.Email)

		count, err := app.UserCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			log.Fatal("error which checking email")
			c.JSON(http.StatusBadRequest, gin.H{"email": "unable to check if the email address is present"})
			return
		}

		fmt.Println("4")
		if count > 0 {
			//log.Fatal()
			c.JSON(http.StatusBadRequest, gin.H{"email": "email already exist"})
			return
		}
		//checks if the phone number exists
		fmt.Println(user.PhoneNumber)

		fmt.Println("5")
		phone, err := app.UserCollection.CountDocuments(ctx, bson.M{"phonenumber": user.PhoneNumber})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"phonenumber": "unable to verify phone number"})
			return
		}
		fmt.Println("6")

		if phone > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"phonenumber": "phone number already exist"})
			return
		}
		// hash the password

		if len(*user.Password) < 6 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "password must be of more then 6 characters",
				"status": http.StatusBadRequest,
			})
		}

		hashedPassword := HashPassword(*user.Password)
		fmt.Println("7")
		// prepare the model
		user.Password = &hashedPassword
		user.ID = primitive.NewObjectID()
		user.Created_At = time.Now()
		user.Updated_At = time.Now()
		user.WishList = make([]models.WishList, 0)
		user.Address = make([]models.Address, 0)
		user.UserCart = make([]models.UserCart, 0)

		// save it to mongo database

		_, err = app.UserCollection.InsertOne(ctx, user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "not created"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"Success": "Succesfully created!", "userID": *user.User_Id})
	}
}

func Addproduct() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

func Productview() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

func SearchProduct() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}
