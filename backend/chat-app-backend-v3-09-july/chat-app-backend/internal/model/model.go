package model

import "github.com/golang-jwt/jwt/v5"

type User struct {
	Name     string `bson:"name"`
	Username string `bson:"username"`
	Password string `bson:"password"`
}

type Friend struct {
	User1 string `bson:"user1"`
	User2 string `bson:"user2"`
}

type SendMessage struct {
	Message  string `json:"message" bson:"message"`
	Receiver string `json:"receiver" bson:"receiver"`
	Time     string `json:"time" bson:"time"`
}

type ReceiveMessage struct {
	Message string `json:"message" bson:"message"`
	Sender  string `json:"sender" bson:"sender"`
	Time    string `json:"time" bson:"time"`
}

type JWTCustomClaims struct {
	Username string `bson:"username"`
	jwt.RegisteredClaims
}

type UtilsResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type Message struct {
	Message string `json:"message"`
}
