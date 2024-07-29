package database

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/rodoben/ecommerce/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func Instantbuy(ctx context.Context, productCollection, userCollection mongo.Collection, userid, productid, addressId string, parameters models.Payment) (models.Order, error) {

	var productDetails models.UserCart
	var order models.Order
	var userDetails models.User

	ProductIdHex, err := primitive.ObjectIDFromHex(productid)
	if err != nil {
		log.Println(err)
	}

	addressIdHex, err := primitive.ObjectIDFromHex(addressId)
	if err != nil {
		log.Println(err)
	}

	fmt.Println("Instantbuy1")

	err = userCollection.FindOne(ctx, bson.M{"user_id": userid}).Decode(&userDetails)
	if err != nil {
		return models.Order{}, errors.New("unable to find user")
	}
	fmt.Println("Instantbuy2")

	err = productCollection.FindOne(ctx, bson.M{"product_id": ProductIdHex}).Decode(&productDetails)
	if err != nil {
		return models.Order{}, errors.New("unable to find product")
	}
	fmt.Println("Instantbuy3")
	order.Order_Id = primitive.NewObjectID()
	order.Order_Cart = append(order.Order_Cart, productDetails)
	order.Price = productDetails.Price
	order.Ordered_At = time.Now()

	foundAdress := false
	fmt.Println("Instantbuy4", addressId)
	for _, v := range userDetails.Address {
		if addressIdHex == v.Address_ID {
			fmt.Println(v.Address_ID.String())
			foundAdress = true
			order.Address.Address_ID = v.Address_ID
			order.Address.Column1 = v.Column1
			order.Address.Column2 = v.Column2
			order.Address.Landmark = v.Landmark
			order.Address.Pincode = v.Pincode

		}
	}
	fmt.Println("Instantbuy4a")
	if !foundAdress {
		return models.Order{}, errors.New("adress not found in user details")
	}
	fmt.Println("Instantbuy4b")

	if !parameters.COD.IsCod && parameters.Digital.UPi.IsUpi {
		order.Payment.Digital.UPi = parameters.Digital.UPi
	} else if !parameters.COD.IsCod && parameters.Digital.Card.IsCard {
		order.Payment.Digital.Card = parameters.Digital.Card
	} else {
		order.Payment.COD.IsCod = parameters.COD.IsCod
	}
	fmt.Println("Instantbuy5")
	fmt.Println("__________", order.Order_Id)
	fmt.Println("__________", order.Address)

	fmt.Println("__________", order.Payment)

	if _, ok := userDetails.OrderStatus[order.Order_Id]; !ok {
		userDetails.OrderStatus[order.Order_Id] = order
	} else {
		return models.Order{}, errors.New("order already exists")
	}

	updateorder := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "orderstatus", Value: userDetails.OrderStatus},
		}},
	}

	_, err = userCollection.UpdateOne(ctx, bson.M{"user_id": userid}, updateorder)
	if err != nil {
		return models.Order{}, fmt.Errorf("unabale to insert into user orders")
	}

	return order, nil

}
