package main

import (
	"event-sourcing-demo/controller"
	"event-sourcing-demo/handler"
	"event-sourcing-demo/repository"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/streadway/amqp"
	"log"
	"net/http"
)

func main() {
	// Connecting to database
	repository.InitialiseDB()

	// Connecting to rabbitmq
	conn := connectToRabbitMQ()
	defer conn.Close()
	amqpCh := createRabbitMQChannel(conn)
	defer amqpCh.Close()

	// Connecting to kafka
	producer := createKafkaProducer()
	defer producer.Close()
	go startKafkaDeliveryReport(producer)

	startHttpServer(amqpCh, producer)

	// Wait for 5s for all message to be delivered
	producer.Flush(5000)
}

func createRabbitMQChannel(conn *amqp.Connection) *amqp.Channel {
	amqpCh, err := conn.Channel()
	if err != nil {
		log.Fatalf("could not rabbitmq channel: %s", err.Error())
	}
	return amqpCh
}

func connectToRabbitMQ() *amqp.Connection {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("could not connect to rabbitmq: %s", err.Error())
	}
	return conn
}

func createKafkaProducer() *kafka.Producer {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
	if err != nil {
		log.Fatalf("could not connect to kafka: %s", err.Error())
	}
	return producer
}

func startHttpServer(amqpCh *amqp.Channel, producer *kafka.Producer) {
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
	payWithEventController := controller.PayWithEventController{Producer: producer}
	paymentRoute := server.Group("/pay", middleware.JWT([]byte("secret")))
	paymentRoute.POST("", payController.Pay)
	paymentRoute.POST("/with-queue", payWithQueueController.Pay)
	paymentRoute.POST("/with-event", payWithEventController.Pay)

	server.Logger.Fatal(server.Start(":1212"))
}

func startKafkaDeliveryReport(p *kafka.Producer) {
	for e := range p.Events() {
		switch ev := e.(type) {
		case *kafka.Message:
			if ev.TopicPartition.Error != nil {
				log.Printf("Delivery failed: %v", ev.TopicPartition)
			} else {
				log.Printf("Delivered message to %v", ev.TopicPartition)
			}
		}
	}
}
