package infra

import (
	"context"
	"fmt"
	"log"
	"message_app/internal/logger"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoClient struct {
	C *mongo.Client
}

func NewMongoClient(logger *logger.Logger) *mongo.Client {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://admin:admin@localhost:27017"))
	if err != nil {
		fmt.Println(err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Connection successful")

	createIndexForMongo(client)

	return client
}

func createIndexForMongo(client *mongo.Client) {

	//index so that same username cannot be repeated
	us := client.Database("msg_db").Collection("users")
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "username", Value: 1},
		},
		Options: options.Index().SetUnique(true).SetName("username_index"),
	}
	name1, err := us.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		log.Printf("Error in creating index: %v", err)
	}
	fmt.Println("indexes created for: ", name1)

	//index so that same users cannot send friend requests again
	fr := client.Database("msg_db").Collection("friends")
	indexModel1 := mongo.IndexModel{
		Keys: bson.D{
			{Key: "user1", Value: 1},
			{Key: "user2", Value: 1},
		},
		Options: options.Index().SetUnique(true).SetName("username_friend_index"),
	}
	name2, err := fr.Indexes().CreateOne(context.Background(), indexModel1)
	if err != nil {
		log.Printf("Error in creating index: %v", err)
	}
	fmt.Println("indexes created for: ", name2)
}
