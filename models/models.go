package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID            primitive.ObjectID              `json:"_id" bson:"_id"`
	FirstName     *string                         `json:"first_name"  validate:"required,min=2,max=30"`
	Lastname      *string                         `json:"last_name" validate:"required,min=2,max=30"`
	Email         *string                         `json:"email" validate:"email,required"`
	PhoneNumber   *string                         `json:"phonenumber" validate:"required"`
	Password      *string                         `json:"password" validate:"required,min=6"`
	Dob           time.Time                       `json:"dob"`
	Token         *string                         `json:"token"`
	Refresh_Token *string                         `json:"refresh_token"`
	Created_At    time.Time                       `json:"created_at"`
	Updated_At    time.Time                       `json:"updated_at"`
	User_Id       *string                         `json:"user_id"`
	WishList      map[primitive.ObjectID]WishList `json:"wishlist"`
	Address       []Address                       `json:"address"`
	UserCart      map[primitive.ObjectID]UserCart `json:"usercart"`
	OrderStatus   map[primitive.ObjectID]Order    `json:"order"`
}

type WishList struct {
	Product_id  primitive.ObjectID `json:"product_id" bson:"product_id"`
	ProductName *string            `json:"product_name"`
	Price       *uint64            `json:"price"`
	Rating      *uint8             `json:"rating"`
	Image       *string            `json:"image"`
}

type Address struct {
	Address_ID primitive.ObjectID `json:"address_id" bson:"address_id"`
	Column1    *string            `json:"column1" validate:"required, min 20" bson:"column1"`
	Column2    *string            `json:"column2" bson:"column2"`
	Landmark   *string            `json:"landmark" bson:"landmark"`
	Pincode    *string            `json:"pincode" validate:"required, len=6, numeric" bson:"pincode"`
}

type UserCart struct {
	Product_id  primitive.ObjectID `json:"product_id" bson:"product_id"`
	ProductName *string            `json:"product_name"`
	Price       *uint64            `json:"price"`
	Quantity    int                `json:"quantity"`
	Rating      *uint8             `json:"rating"`
	Image       *string            `json:"image"`
}

type Product struct {
	Product_id  primitive.ObjectID `json:"product_id" bson:"product_id"`
	ProductName *string            `json:"product_name"`
	Price       *uint64            `json:"price"`
	Quantity    int                `json:"quantity"`
	Rating      *uint8             `json:"rating"`
	Image       *string            `json:"image"`
}

type Order struct {
	Order_Id   primitive.ObjectID `json:"_id" bson:"_id"`
	Order_Cart []UserCart         `json:"order_cart"`
	Price      *uint64            `json:"price"`
	Discount   *uint8             `json:"discount"`
	Address    Address            `json:"address" bson:"address"`
	Ordered_At time.Time          `json:"ordered_at"`
	Payment    Payment            `json:"payment"`
}

type Payment struct {
	COD     COD     `json:"cod"`
	Digital Digital `json:"digital"`
}

type COD struct {
	IsCod bool `json:"is_cod"`
}
type Digital struct {
	UPi  *UpiDetails  `json:"upi"`
	Card *CardDetails `json:"cardDetails"`
}

type UpiDetails struct {
	IsUpi bool    `json:"isupi"`
	UpiID *string `json:"upi_id"`
}

type CardDetails struct {
	IsCard         bool    `json:"iscard"`
	CardType       *string `json:"cardtype"`
	CardHolderName *string `json:"cardholdername"`
	CardNumber     *string `json:"cardnumber"`
	CardExpiry     *string `json:"cardexpiry"`
	CardCVV        *string `json:"cardcvv"`
}
