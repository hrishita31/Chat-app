package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"message_app/internal/logger"
	"message_app/internal/model"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type IUserService interface {
	PostUser(name string, username string, password string) (string, error)
	ValidateUser(username string, password string) (string, error)
	GetUser(username string) (*model.User, error)
	SearchUser(username string) (string, []string, error)
}

type UserService struct {
	mongoClient *mongo.Client
	logger      *logger.Logger
}

func (s *UserService) PostUser(name string, username string, password string) (string, error) {
	client := s.mongoClient
	coll := client.Database("msg_db").Collection("users")

	st := model.User{Name: name, Username: username, Password: password}
	result, err := coll.InsertOne(context.Background(), st)
	if err != nil {
		s.logger.L.Error("Mongo Error", zap.Error(err))
		if mongo.IsDuplicateKeyError(err) {
			return "username already exists", fmt.Errorf("username already exists")
		}
		return "", err

	}
	return fmt.Sprintf("%v", result), err

}

func (s *UserService) ValidateUser(username string, password string) (string, error) {
	client := s.mongoClient
	coll := client.Database("msg_db").Collection("users")

	hashedPassword := sha256.Sum256([]byte(password))
	hashedPasswordStr := hex.EncodeToString(hashedPassword[:])

	userCred := bson.D{{Key: "username", Value: username}, {Key: "password", Value: hashedPasswordStr}}
	var us model.User
	err := coll.FindOne(context.Background(), userCred).Decode(&us)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", fmt.Errorf("username and password do not match")
		}
		s.logger.L.Error("error occurred", zap.Error(err))
		return "", err
	}

	claims := &model.JWTCustomClaims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		s.logger.L.Error("error occurred", zap.Error(err))
		return "", err
	}

	return t, err
}

func (s *UserService) GetUser(username string) (*model.User, error) {
	client := s.mongoClient
	coll := client.Database("msg_db").Collection("users")

	userDet := bson.D{{Key: "username", Value: username}}
	var us model.User
	err := coll.FindOne(context.Background(), userDet).Decode(&us)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("username does not exist")
		}
		s.logger.L.Error("error occurred")
		return nil, err
	}
	return &us, err
}

func (s *UserService) SearchUser(username string) (string, []string, error) {
	client := s.mongoClient
	coll := client.Database("msg_db").Collection("users")

	userDet := bson.D{{Key: "username", Value: username}}
	var us model.User
	err := coll.FindOne(context.Background(), userDet).Decode(&us)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			s.logger.L.Info("username not found, searching for suggestions")

			us := username[:3]
			filter := bson.D{{Key: "username", Value: bson.D{{Key: "$regex", Value: us}}}}
			var suggestions []model.User

			cursor, err := coll.Find(context.Background(), filter)

			if err != nil {
				s.logger.L.Error("error occurred")
				return "", nil, err
			}
			err = cursor.All(context.Background(), &suggestions)
			if err != nil {
				s.logger.L.Error("error occurred")
				return "", nil, err
			}
			var suggestedNames []string
			for _, v := range suggestions {
				suggestedNames = append(suggestedNames, v.Username)
			}
			return "username not found", suggestedNames, err

		}
		s.logger.L.Error("some error has occurred")
	}
	return "username found", []string{us.Username}, err
}

func NewUserSvc(client *mongo.Client, logger *logger.Logger) IUserService {
	return &UserService{logger: logger, mongoClient: client}
}
