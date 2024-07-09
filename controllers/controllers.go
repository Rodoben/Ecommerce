package controllers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rodoben/ecommerce/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

//-----------------generic functions----------------------

func HashPassword(password string) string {
	passbytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err.Error()
	}
	return string(passbytes)
}

func VerifyPassword(UserPassword string, Givenpassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(UserPassword), []byte(Givenpassword))
	valid := true
	if err != nil {
		valid = false
		err = errors.New("unable to match the passwords")
	}
	return valid, err
}

//-------------------user routes controllers-----------------

func SignIn() gin.HandlerFunc {

	return func(ctx *gin.Context) {

	}
}

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		//context deadline

		var _, cancel = context.WithTimeout(context.Background(), time.Second*100)
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

		fmt.Println(user.Email)
		//checks if the phone number exists
		fmt.Println(user.PhoneNumber)

		// hash the password

		hashedPassword := HashPassword(*user.Password)

		// prepare the model
		user.Password = &hashedPassword
		user.ID = primitive.NewObjectID()
		user.Created_At = time.Now()
		user.Updated_At = time.Now()
		user.WishList = make([]models.WishList, 0)
		user.Address = make([]models.Address, 0)
		user.UserCart = make([]models.UserCart, 0)

		// save it to mongo database

		fmt.Println(user)
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

//-------------------cart routes controllers-----------------

func AddToCart() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

func DeleteFromCartByID() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

func DeleteFromCart() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}
