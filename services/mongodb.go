package services

import (
	"context"
	"os"

	"github.com/mongodb/mongo-go-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// read the config
// var v = utils.ReadConfig("config")
var mongoURL = os.Getenv("MONGO_URL")

func mongoClient(url string) *mongo.Client {
	client, _ := mongo.Connect(context.TODO(), options.Client().ApplyURI(url))
	return client
}

var client = mongoClient(mongoURL)
var Db_user = client.Database("basic_information").Collection("users")
