package utils

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func ResponseWithSuccess(c echo.Context, code int, data interface{}, message string) error {
	return c.JSON(http.StatusAccepted, map[string]any{
		"status":  "accepted",
		"message": message,
		"data":    data,
	})

}

func ResponseWithError(c echo.Context, code int, message string) error {
	return c.JSON(http.StatusUnauthorized, map[string]any{
		"status":  "error",
		"message": message,
		"data":    "no data found",
	})

}
