package main

import (
	"event-sourcing-demo/controller"
	"event-sourcing-demo/repository"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net/http"
)

func main() {
	repository.InitialiseDB()

	server := echo.New()
	server.Use(middleware.Logger())

	server.GET("/ping", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, "pong")
	})
	server.POST("/login", controller.Login)

	paymentRoute := server.Group("/pay", middleware.JWT([]byte("secret")))
	paymentRoute.POST("", controller.Pay)

	server.Logger.Fatal(server.Start(":1212"))
}
