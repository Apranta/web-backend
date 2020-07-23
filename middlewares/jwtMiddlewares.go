package middlewares

import (
	"fmt"
	"net/http"
	"strconv"

	"web-backend-patal/config"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func SetClientJWTmiddlewares(g *echo.Group) {
	jwtConfig := config.App.Config.GetStringMap(fmt.Sprintf("jwt"))

	g.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningMethod: "HS256",
		SigningKey:    []byte(jwtConfig["secret"].(string)),
	}))

	g.Use(validateJWTclient)
}

func validateJWTclient(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user")
		token := user.(*jwt.Token)

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			userID, _ := strconv.Atoi(fmt.Sprintf("%v", claims["vid"]))
			if userID != 0 {
				return next(c)
			} else {
				return echo.NewHTTPError(http.StatusUnauthorized, fmt.Sprintf("%s", "Please Sign In \n Woops! Gonna sign in first\n Only a click away and you can continue to enjoy"))
			}
		}

		return echo.NewHTTPError(http.StatusUnauthorized, fmt.Sprintf("%s", "Please Sign In \n Woops! Gonna sign in first\n Only a click away and you can continue to enjoy"))
	}
}
