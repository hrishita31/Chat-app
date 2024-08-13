package router

import (
	"message_app/internal/handler"
	"message_app/internal/middleware"
	"net/http"

	"github.com/labstack/echo/v4"
	mdleware "github.com/labstack/echo/v4/middleware"
)

func NewRouter(userHandler handler.IUserHandler, friendHandler handler.IFriendHandler, websocketHandler handler.IWebSocketHandler) *echo.Echo {
	e := echo.New()

	e.Use(mdleware.CORSWithConfig(mdleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	user := e.Group("/user")
	{
		user.POST("/", userHandler.PostUser())
		user.GET("/loginUser", userHandler.ValidateUser())
		user.GET("/", userHandler.GetUser(), middleware.VerifyTokenMiddleware)
		user.GET("/search/:username", userHandler.SearchUser(), middleware.VerifyTokenMiddleware)
	}

	friend := e.Group("/friend")
	friend.Use(middleware.VerifyTokenMiddleware)
	{
		friend.GET("/addFriend", friendHandler.AddFriend())
		friend.GET("/findAllFriends/:username", friendHandler.FindAllFriends())
	}

	ws := e.Group("/ws")
	{
		ws.GET("/", websocketHandler.WS())
		ws.GET("/histories", websocketHandler.HistoryMsgs())
	}
	return e

}
