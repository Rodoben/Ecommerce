package database

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type DbProperties struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	DbName   string `json:"dbname"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func loadMongoUri() (string, error) {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("unable to load the properties for mongodb")
	}
	mongoProperties := os.Getenv("APP_MONGO_PROPERTIES")
	var mongoDBProperies DbProperties
	json.Unmarshal([]byte(mongoProperties), &mongoDBProperies)
	uri := fmt.Sprintf("%s://%s:%s@%s:%v", mongoDBProperies.DbName, mongoDBProperies.Username, mongoDBProperies.Password, mongoDBProperies.Host, mongoDBProperies.Port)
	return uri, nil
}

func DBConnect() *mongo.Client {
	uri, err := loadMongoUri()
	if err != nil {
		log.Fatalf("error loading the mongo uri")
	}
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("connected")
	return client
}

func UserData(client *mongo.Client, collectionName string) *mongo.Collection {
	var collection *mongo.Collection = client.Database("Ecommerce").Collection(collectionName)
	return collection
}

func ProductData(client *mongo.Client, collectionName string) *mongo.Collection {
	var collection *mongo.Collection = client.Database("Ecommerce").Collection(collectionName)
	return collection
}
