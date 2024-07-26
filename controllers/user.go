package controllers

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"encoding/json"

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
			c.JSON(http.StatusBadRequest, gin.H{"error": "User not Exists! please signup"})
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
		user.WishList = make(map[primitive.ObjectID]models.WishList, 0)
		user.Address = make([]models.Address, 0)
		user.UserCart = make(map[primitive.ObjectID]models.UserCart, 0)
		user.OrderStatus = make(map[primitive.ObjectID]models.Order, 0)

		// save it to mongo database

		_, err = app.UserCollection.InsertOne(ctx, user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "not created"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"Success": "Succesfully created!", "userID": *user.User_Id})
	}
}

func validateProduct(product models.Product) error {

	if product.ProductName == nil {
		return errors.New("product name is required")
	} else if product.Price == nil {
		return errors.New("price is required")
	} else if product.Image == nil {
		return errors.New("image is required")
	} else if product.Rating == nil {
		return errors.New("rating is required")
	} else if product.Quantity == 0 {
		product.Quantity = 1
		return errors.New("quantity is required")
	}
	return nil

}

func Addproduct() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		fmt.Println("a")
		var product models.Product
		var foundProduct models.Product
		jsonData, _ := io.ReadAll(c.Request.Body)

		fmt.Println(string(jsonData))
		if err := json.Unmarshal(jsonData, &product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "err": "unable to parse json into struct"})
			return
		}

		err := validateProduct(product)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "err": "unable to parse json into struct"})
			return
		}

		if product.ProductName == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "product_name cannot be nil"})
			return
		}
		fmt.Println("b")

		fmt.Println("c")
		count, err := productCollection.CountDocuments(ctx, bson.M{"productname": product.ProductName})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return

		}
		fmt.Println(count, "aaa")

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
		fmt.Println("d")
		product.Product_id = primitive.NewObjectID()
		fmt.Println("e")
		_, err = productCollection.InsertOne(ctx, product)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		fmt.Println("f")
		product.Quantity = 1
		c.JSON(http.StatusCreated, gin.H{"Success": "Succesfully created!", "Quantity": product.Quantity})
	}
}

func ProductViewerAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var products models.Product

		defer cancel()
		if err := c.ShouldBindJSON(&products); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "err": "unable to parse json into struct"})
			return
		}
		fmt.Println(products)
		products.Product_id = primitive.NewObjectID()
		_, anyerr := productCollection.InsertOne(ctx, products)
		if anyerr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Not Created"})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, "Successfully added our Product Admin!!")
	}
}
func (app *Application) Productview() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx := context.Background()

		var products []models.Product
		fmt.Println("1")

		cur, err := app.ProductCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to fetch documents"})
		}
		fmt.Println("2", cur)
		for cur.Next(context.TODO()) {
			var product models.Product
			fmt.Println("a")
			err = cur.Decode(&product)
			fmt.Println("b")
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to fetch documents"})
			}
			fmt.Println("c")
			products = append(products, product)

		}
		fmt.Println("3")

		c.JSON(http.StatusOK, gin.H{"success": "sucesfully fetched!", "data": products})

	}
}

func SearchProduct() gin.HandlerFunc {
	return func(c *gin.Context) {

		query := c.Query("name")
		if query == "" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Query is empty"})
			return
		}
		var products []models.Product
		cur, err := productCollection.Find(c, bson.M{"productname": bson.M{"$regex": query}})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to fetch record from database"})
			return
		}
		err = cur.All(c, &products)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to fetch record from database"})
			return
		}
		defer cur.Close(c)

		c.JSON(http.StatusOK, gin.H{"Success": http.StatusOK, "data": products})
	}
}
