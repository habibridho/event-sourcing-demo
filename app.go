package main

import (
	"event-sourcing-demo/controller"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net/http"
)

func main() {
	server := echo.New()
	server.Use(middleware.Logger())

	server.GET("/ping", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, "pong")
	})
	server.POST("/login", controller.Login)

	server.Logger.Fatal(server.Start(":1212"))
}
