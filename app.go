package main

import (
	"event-sourcing-demo/controller"
	"event-sourcing-demo/handler"
	"event-sourcing-demo/repository"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/streadway/amqp"
	"log"
	"net/http"
)

func main() {
	repository.InitialiseDB()
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("could not connect to rabbitmq: %s", err.Error())
	}
	defer conn.Close()
	amqpCh, err := conn.Channel()
	if err != nil {
		log.Fatalf("could not rabbitmq channel: %s", err.Error())
	}

	server := echo.New()
	server.Use(middleware.Logger())

	server.GET("/ping", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, "pong")
	})
	server.POST("/login", controller.Login)

	mockHandler := handler.MockHandler{}
	payController := controller.PayController{NotificationHandler: &mockHandler, EmailHandler: &mockHandler}
	payWithQueueController, err := controller.NewPayWithQueueController(amqpCh)
	if err != nil {
		log.Fatalf("could not create pay with queue controller: %s", err.Error())
	}
	paymentRoute := server.Group("/pay", middleware.JWT([]byte("secret")))
	paymentRoute.POST("", payController.Pay)
	paymentRoute.POST("/with-queue", payWithQueueController.Pay)

	server.Logger.Fatal(server.Start(":1212"))
}
