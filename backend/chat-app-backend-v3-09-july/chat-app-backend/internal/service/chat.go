package service

import (
	"context"
	"encoding/json"
	"message_app/internal/logger"
	"message_app/internal/model"
	"strings"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type IChatService interface {
	HistoryMsgs(username string, receiver string) ([]model.SendMessage, error)
}

type ChatService struct {
	client *redis.Client
	logger *logger.Logger
}

func (s *ChatService) HistoryMsgs(username string, receiver string) ([]model.SendMessage, error) {
	client := s.client

	username = strings.ToLower(username)
	receiver = strings.ToLower(receiver)
	var key string
	if username < receiver {
		key = username + "_" + receiver
	} else if username > receiver {
		key = receiver + "_" + username
	}

	var history []string

	ans, err := client.Exists(context.Background(), key).Result()
	if err != nil {
		s.logger.L.Error("error occurred", zap.Error(err))
	}
	if ans == 1 {
		s.logger.L.Info("key1 exists")
		history, err = client.LRange(context.Background(), key, 0, -1).Result()
		if err != nil {
			s.logger.L.Info("could not recover old messages")
			return nil, err
		}
	}

	var msgModelArr []model.SendMessage
	for _, h := range history {
		var m model.SendMessage

		hi := []byte(h)
		s.logger.L.Info(string(hi))
		err := json.Unmarshal(hi, &m)
		if err != nil {
			s.logger.L.Info("error in marshalling", zap.Error(err))
		}
		msgModelArr = append(msgModelArr, m)
	}

	return msgModelArr, err
}

func NewChatSvc(client *redis.Client, logger *logger.Logger) IChatService {
	return &ChatService{client: client, logger: logger}
}
