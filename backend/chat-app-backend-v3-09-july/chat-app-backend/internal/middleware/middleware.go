package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func verifyToken(p []byte) (*jwt.Token, error) {
	token, err := jwt.Parse(string(p), func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, err
	}
	return token, err
}

func VerifyTokenMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		token := c.Request().Header.Get("Authorization")
		pay := strings.Split(token, " ")
		p := pay[1]
		pByte := []byte(p)

		p1, err := verifyToken(pByte)
		if err != nil {
			fmt.Println(err)
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"message": "Invalid token",
			})
		}
		fmt.Printf("payload: %v", p1.Raw)

		claims := jwt.MapClaims{}
		to, err := jwt.ParseWithClaims(p, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil {
			fmt.Printf("error: %v", err)
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"message": "Invalid token",
			})
		}
		fmt.Printf("token: %v", to.Claims)

		username, ok := claims["Username"].(string)
		if !ok {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"message": "invalid token username",
			})
		}
		c.Set("username", username)
		return next(c)
	}
}
