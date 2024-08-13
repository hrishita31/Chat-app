package handler

import (
	"crypto/sha256"
	"encoding/hex"
	"message_app/internal/logger"
	"message_app/internal/service"
	"message_app/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type IUserHandler interface {
	PostUser() echo.HandlerFunc
	ValidateUser() echo.HandlerFunc
	GetUser() echo.HandlerFunc
	SearchUser() echo.HandlerFunc
}

type UserHandler struct {
	svc    service.IUserService
	logger *logger.Logger
}

func (s *UserHandler) PostUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		name := c.QueryParam("name")
		username := c.QueryParam("username")
		password := c.QueryParam("password")

		ok := 0
		for i := 0; i < len(username); i++ {
			if (string(username[i]) >= "0" && string(username[i]) <= "9") || (string(username[i]) >= "a" && string(username[i]) <= "z") || (string(username[i]) >= "A" && string(username[i]) <= "Z") || string(username[i]) == "." || string(username[i]) == "_" {
				continue
			} else {
				ok = 1
				break
			}
		}
		if ok == 1 {
			return utils.ResponseWithError(c, http.StatusUnauthorized, "username not valid")
		}

		hashedPassword := sha256.Sum256([]byte(password))
		hashedPasswordStr := hex.EncodeToString(hashedPassword[:])

		result, err := s.svc.PostUser(name, username, hashedPasswordStr)
		if err != nil {
			s.logger.L.Info("error occurred", zap.Error(err))
			if mongo.IsDuplicateKeyError(err) {
				return utils.ResponseWithError(c, http.StatusUnauthorized, "username already exists")
			}
			return utils.ResponseWithError(c, http.StatusUnauthorized, "error occurred")
		}
		return utils.ResponseWithSuccess(c, http.StatusAccepted, result, "successful")
	}
}

func (s *UserHandler) ValidateUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		username := c.QueryParam("username")
		password := c.QueryParam("password")

		result, err := s.svc.ValidateUser(username, password)
		if err != nil {
			s.logger.L.Debug("user credentials not validated", zap.Error(err))
			return utils.ResponseWithError(c, http.StatusUnauthorized, "error occurred")
		}
		return utils.ResponseWithSuccess(c, http.StatusAccepted, result, "successful")
	}
}

func (s *UserHandler) GetUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		username := c.Get("username").(string)
		result, err := s.svc.GetUser(username)
		if err != nil {
			s.logger.L.Debug("error occurred", zap.Error(err))
			return utils.ResponseWithError(c, http.StatusUnauthorized, "error occurred")
		}
		return utils.ResponseWithSuccess(c, http.StatusAccepted, result, "successful")
	}
}

func (s *UserHandler) SearchUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		username := c.Param("username")

		result, suggestions, err := s.svc.SearchUser(username)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return utils.ResponseWithError(c, http.StatusUnauthorized, "username not found")
			}
			s.logger.L.Info("error occurred")
			return utils.ResponseWithError(c, http.StatusUnauthorized, "error occurred")
		}
		return utils.ResponseWithSuccess(c, http.StatusAccepted, suggestions, result)
	}
}

func NewUser(service service.IUserService, logger *logger.Logger) IUserHandler {
	return &UserHandler{svc: service, logger: logger}
}
