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

func Instantbuy(ctx context.Context, productCollection, userCollection mongo.Collection, userid, productid string) error {

	var productDetails models.UserCart
	var order models.Order

	order.Order_Id = primitive.NewObjectID()

	ProductIdHex, err := primitive.ObjectIDFromHex(productid)
	if err != nil {
		log.Println(err)
	}

	err = productCollection.FindOne(ctx, bson.M{"product_id": ProductIdHex}).Decode(&productDetails)
	if err != nil {
		return errors.New("unable to find product")
	}

	order.Order_Cart = append(order.Order_Cart, productDetails)
	order.Price = productDetails.Price
	order.Ordered_At = time.Now()
	if order.Payment.COD.IsCod {

		
		order.Payment.COD.Address = models.Address{}
	} else {
		order.Payment.Digital = models.Digital{UPi: &models.UpiDetails{}, Card: &models.CardDetails{}}
	}

	fmt.Println(order)

	return nil

}
