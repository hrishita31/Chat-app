package handler

import (
	"message_app/internal/logger"
	"message_app/internal/service"
	"message_app/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

type IFriendHandler interface {
	AddFriend() echo.HandlerFunc
	FindAllFriends() echo.HandlerFunc
}

type FriendHandler struct {
	svc    service.IFriendService
	logger *logger.Logger
}

func (s *FriendHandler) AddFriend() echo.HandlerFunc {
	return func(c echo.Context) error {
		username1 := c.QueryParam("username1")
		username2 := c.QueryParam("username2")

		rel, err := s.svc.AddFriend(username1, username2)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return utils.ResponseWithError(c, http.StatusUnauthorized, "cannot add as friend")
			}
			s.logger.L.Info("cannot add as friend")
			return utils.ResponseWithError(c, http.StatusUnauthorized, "error occurred")
		}
		return utils.ResponseWithSuccess(c, http.StatusAccepted, rel, "successful")
	}
}

func (s *FriendHandler) FindAllFriends() echo.HandlerFunc {
	return func(c echo.Context) error {
		username := c.Param("username")

		friends, err := s.svc.FindFriends(username)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return utils.ResponseWithError(c, http.StatusUnauthorized, "username not found")
			}
			s.logger.L.Info("username not found")
			return utils.ResponseWithError(c, http.StatusUnauthorized, "error occurred")
		}
		return utils.ResponseWithSuccess(c, http.StatusAccepted, friends, "successful")
		//username, length of friends list
	}
}

func NewFriend(service service.IFriendService, logger *logger.Logger) IFriendHandler {
	return &FriendHandler{svc: service, logger: logger}
}
