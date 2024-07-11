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

var productCollection *mongo.Collection = database.ProductData(database.DBConnect(), "productdata")

func SignIn() gin.HandlerFunc {

	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(c, 20*time.Second)
		defer cancel()
		var user, founduser models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "unable to parse the body in the model"})
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&founduser)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "error fetching information from database"})
			return
		}
		fmt.Println("founduser", founduser.FirstName)

		passwordverified, _ := VerifyPassword_v1(*user.Password, *founduser.Password)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "error verifying password"})
			return
		}
		if !passwordverified {
			c.JSON(http.StatusBadRequest, gin.H{"error": "password is incorrect"})
			return
		}

		c.JSON(http.StatusAccepted, gin.H{"Success": "signed in!",
			"UserName":    *founduser.FirstName,
			"status code": http.StatusAccepted})
	}
}

func (app *Application) SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		//context deadline

		var ctx, cancel = context.WithTimeout(context.Background(), time.Second*100)
		defer cancel()

		// bind json in model
		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//validate struct
		validate := validator.New()
		validateError := validate.Struct(user)
		if validateError != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validateError.Error()})
			return
		}

		//check if the email exists in database

		count, err := app.UserCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			log.Fatal("error which checking email")
			c.JSON(http.StatusBadRequest, gin.H{"email": "unable to check if the email address is present"})
			return
		}

		if count > 0 {
			//log.Fatal()
			c.JSON(http.StatusBadRequest, gin.H{"email": "email already exist"})
			return
		}
		//checks if the phone number exists

		phone, err := app.UserCollection.CountDocuments(ctx, bson.M{"phonenumber": user.PhoneNumber})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"phonenumber": "unable to verify phone number"})
			return
		}

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

		userid := "1234"
		// prepare the model
		user.Password = &hashedPassword
		user.ID = primitive.NewObjectID()
		user.User_Id = &userid
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
	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(c, 10*time.Second)
		defer cancel()

		var product, foundProduct models.Product

		if err := c.BindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		count, err := productCollection.CountDocuments(ctx, bson.M{"productname": product.ProductName})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		}
		fmt.Println(count)

		if count > 0 {

			err := productCollection.FindOne(ctx, bson.M{"productname": *product.ProductName}).Decode(&foundProduct)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

			}
			foundProduct.Quantity++
			update := bson.D{
				{Key: "$set", Value: bson.D{
					{Key: "quantity", Value: foundProduct.Quantity},
				}},
			}

			productCollection.UpdateOne(ctx, bson.M{"productname": product.ProductName}, update)
			c.JSON(http.StatusCreated, gin.H{"Success": "Succesfully Updated!", "Quantity": foundProduct.Quantity})
			return
		}
		product.Product_id = primitive.NewObjectID()

		_, err = productCollection.InsertOne(ctx, product)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		product.Quantity = 1
		c.JSON(http.StatusCreated, gin.H{"Success": "Succesfully created!", "Quantity": product.Quantity})
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
