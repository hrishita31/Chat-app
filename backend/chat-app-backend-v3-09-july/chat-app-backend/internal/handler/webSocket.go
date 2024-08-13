package handler

import (
	"context"
	"encoding/json"
	"message_app/internal/logger"
	"message_app/internal/model"
	"message_app/internal/service"
	"message_app/internal/utils"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type IWebSocketHandler interface {
	WS() echo.HandlerFunc
	HistoryMsgs() echo.HandlerFunc
}

type WebSocketHandler struct {
	logger *logger.Logger
	svc    service.IChatService
	Cr     *redis.Client
	C      *mongo.Client
}

func NewWebSocket(logger *logger.Logger, service service.IChatService, client *redis.Client, mongo *mongo.Client) IWebSocketHandler {
	return &WebSocketHandler{logger: logger, svc: service, Cr: client, C: mongo}
}

func (w *WebSocketHandler) WS() echo.HandlerFunc {
	return func(c echo.Context) error {
		username := c.QueryParam("username")
		upgrader := websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins, adjust for production
			},
		}

		ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			w.logger.L.Error("WebSocket Upgrade failed", zap.Error(err))
			return err
		}
		defer ws.Close()

		ctx, cancel := context.WithCancel(context.Background())
		chMsg := make(chan string)

		go func() {

			select {
			case <-ctx.Done():
				return
			default:

				for {
					_, msg, err := ws.ReadMessage()
					if err != nil {
						w.logger.L.Error("ws read error", zap.Error(err))
						cancel()
						return
					}

					var m model.SendMessage

					err = json.Unmarshal(msg, &m)
					if err != nil {
						w.logger.L.Error("error occurred", zap.Error(err))
					}
					w.logger.L.Info("received message", zap.Any("msg", m))

					err = w.Cr.Publish(context.Background(), m.Receiver, string(msg)).Err()
					if err != nil {
						w.logger.L.Error("publish error", zap.Error(err))
						cancel()
						return
					}
					var m1 model.ReceiveMessage
					m1.Sender = username

					err = json.Unmarshal(msg, &m1)
					if err != nil {
						w.logger.L.Info("error occurred", zap.Error(err))
					}
					w.logger.L.Info("sent message", zap.Any("sent from", m1.Sender))
					client := w.Cr

					m1.Sender = strings.ToLower(m1.Sender)
					m.Receiver = strings.ToLower(m.Receiver)
					var key string
					if m1.Sender > m.Receiver {
						key = m.Receiver + "_" + m1.Sender
					} else if m1.Sender < m.Receiver {
						key = m1.Sender + "_" + m.Receiver
					}

					ans, err := client.Exists(context.Background(), key).Result()
					if err != nil {
						w.logger.L.Error("error occurred", zap.Error(err))
					}
					if ans == 1 {
						w.logger.L.Info("key exists")
						msg, err := json.Marshal(m)
						if err != nil {
							w.logger.L.Error("error marshalling", zap.Error(err))
						}

						err = client.LPush(context.Background(), key, string(msg)).Err()
						if err != nil {
							w.logger.L.Error("could not add to database", zap.Error(err))
							cancel()
							return
						}

						c := w.C
						coll := c.Database("msg_db").Collection(key)
						mes := model.Message{Message: string(msg)}
						_, err = coll.InsertOne(context.Background(), mes)
						if err != nil {
							cancel()
							return
						}

					} else {
						err = client.LPush(context.Background(), key, string(msg)).Err()
						if err != nil {
							w.logger.L.Error("could not add to database", zap.Error(err))
							cancel()
							return
						}
						c := w.C
						coll := c.Database("msg_db").Collection(key)
						mes := model.Message{Message: string(msg)}
						_, err = coll.InsertOne(context.Background(), mes)
						if err != nil {
							cancel()
							return
						}
					}
				}
			}
		}()

		go func() {
			pubsub := w.Cr.Subscribe(context.Background(), username)
			defer pubsub.Close()

			for msg := range pubsub.Channel() {
				chMsg <- msg.Payload
			}
		}()

		for msg := range chMsg {
			if err := ws.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
				w.logger.L.Error("WebSocket write error", zap.Error(err))
				cancel()
				return err
			}
		}

		return nil
	}
}

func (w *WebSocketHandler) HistoryMsgs() echo.HandlerFunc {
	return func(c echo.Context) error {
		username := c.QueryParam("username")
		receiver := c.QueryParam("receiver")
		history, err := w.svc.HistoryMsgs(username, receiver)
		if err != nil {
			w.logger.L.Error("could not find messages", zap.Error(err))
			return utils.ResponseWithError(c, http.StatusUnauthorized, "error occurred")
		}
		return utils.ResponseWithSuccess(c, http.StatusAccepted, history, "successful")
	}
}
