package main

import (
	"os"

	"web-backend-patal/config"
	_ "web-backend-patal/docs" // docs is generated by Swag CLI, we have to import it.
	"web-backend-patal/handlers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title API Documentation for Palembang Digital
// @version 0.0.1
// @description Dokumentasi API untuk website palembangdigital.org

// @contact.name Palembang Digital
// @contact.url https://palembangdigital.org
// @contact.email support@palembangdigital.org

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

func main() {
	defer config.App.Close()

	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"*"},
		AllowHeaders: []string{"*", "X-Accept-Charset", "X-Accept", "Content-Type", "Authorization", "Accept", "Origin", "Access-Control-Request-Method", "Access-Control-Request-Headers"},
	}))

	e.GET("/docs/*", echoSwagger.WrapHandler)

	api := e.Group("/api")
	{
		api.GET("/serviceinfo", handlers.ServiceInfo)

		// Just for example purpose
		accounts := api.Group("/accounts")
		accounts.POST("/login", handlers.Login)
	}

	e.Logger.Fatal(e.Start(":" + config.App.Port))
	os.Exit(0)
}
