package controllers

import (
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type Application struct {
	UserCollection    *mongo.Collection
	ProductCollection *mongo.Collection
}

func NewApplication(userCollection, productCollection *mongo.Collection) *Application {
	return &Application{
		UserCollection:    userCollection,
		ProductCollection: productCollection,
	}
}

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

func VerifyPassword_v1(userpassword string, givenpassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(givenpassword), []byte(userpassword))
	valid := true
	msg := ""
	if err != nil {
		msg = "Login Or Passowrd is Incorerct"
		valid = false
	}
	return valid, msg
}
