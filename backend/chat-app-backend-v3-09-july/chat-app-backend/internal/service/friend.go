package service

import (
	"context"
	"fmt"
	"message_app/internal/logger"
	"message_app/internal/model"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type IFriendService interface {
	AddFriend(username1 string, username2 string) (string, error)
	FindFriends(username string) ([]string, error)
}

type FriendService struct {
	mongoClient *mongo.Client
	logger      *logger.Logger
	redis       *redis.Client
}

func (s *FriendService) AddFriend(username1 string, username2 string) (string, error) {
	client := s.mongoClient
	coll := client.Database("msg_db").Collection("users")

	//check if username 1 exists
	userDet := bson.D{{Key: "username", Value: username1}}
	var us model.User
	err := coll.FindOne(context.Background(), userDet).Decode(&us)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			s.logger.L.Info("username1 not found")

			return "", err
		}
		s.logger.L.Error("error occurred")
		return "", err
	}

	//check if username 2 exists
	userDet2 := bson.D{{Key: "username", Value: username2}}
	var us2 model.User
	err = coll.FindOne(context.Background(), userDet2).Decode(&us2)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			s.logger.L.Info("username2 not found")

		}
		s.logger.L.Error("some error has occurred")
		return "", err

	}

	coll2 := client.Database("msg_db").Collection("friends")

	//check if friendship already exists
	filter1 := bson.D{
		{Key: "$or", Value: bson.A{
			bson.D{{Key: "user1", Value: username1}, {Key: "user2", Value: username2}},
			bson.D{{Key: "user1", Value: username2}, {Key: "user2", Value: username1}},
		}},
	}

	var existingFriend model.Friend
	err = coll2.FindOne(context.Background(), filter1).Decode(&existingFriend)
	if err == nil {
		s.logger.L.Info("friendship already exists")
		return "friendship already exists", nil
	} else if err != mongo.ErrNoDocuments {
		s.logger.L.Error("error occurred")
		return "", err
	}

	fr := model.Friend{User1: username1, User2: username2}
	_, err = coll2.InsertOne(context.Background(), fr)
	if err != nil {
		s.logger.L.Info("friend not added")
		return "", err
	}

	relation := username1 + "_" + username2
	return relation, err
}

func (s *FriendService) FindFriends(username string) ([]string, error) {
	client := s.mongoClient
	coll2 := client.Database("msg_db").Collection("friends")

	filter := bson.D{
		{Key: "$or", Value: bson.A{
			bson.D{{Key: "user1", Value: username}},
			bson.D{{Key: "user2", Value: username}},
		}},
	}

	var friends []struct {
		User1 string `bson:"user1"`
		User2 string `bson:"user2"`
	}
	fmt.Println(friends)
	cursor, err := coll2.Find(context.Background(), filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			s.logger.L.Info("username not found")
			return nil, err
		}
		s.logger.L.Error("error occurred")
		return nil, err
	}

	err = cursor.All(context.Background(), &friends)
	if err != nil {
		s.logger.L.Error("error occurred")
		return nil, err
	}

	var friendList []string
	for _, v := range friends {
		if v.User1 == username {
			friendList = append(friendList, v.User2)
		} else {
			friendList = append(friendList, v.User1)
		}
	}
	return friendList, nil

}

func NewFriendSvc(client *mongo.Client, logger *logger.Logger, redis *redis.Client) IFriendService {
	return &FriendService{logger: logger, mongoClient: client, redis: redis}
}
