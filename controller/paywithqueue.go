package controller

import (
	"encoding/json"
	"errors"
	"event-sourcing-demo/model"
	"event-sourcing-demo/repository"
	"event-sourcing-demo/util"
	"github.com/labstack/echo"
	"github.com/streadway/amqp"
	"log"
	"net/http"
)

type PayWithQueueController struct {
	amqpChannel *amqp.Channel
}

var (
	PushNotificiationExchange = "PUSH_NOTIFICATION_EXCHANGE"
	EmailExchange             = "EMAIL_EXCHANGE"
)

func NewPayWithQueueController(conn *amqp.Connection) (*PayWithQueueController, error) {
	amqpChannel, err := conn.Channel()
	if err != nil {
		log.Fatalf("could not rabbitmq channel: %s", err.Error())
	}
	if err := amqpChannel.ExchangeDeclare(PushNotificiationExchange, "fanout", true, false, false, false, nil); err != nil {
		return nil, err
	}
	if err := amqpChannel.ExchangeDeclare(EmailExchange, "fanout", true, false, false, false, nil); err != nil {
		return nil, err
	}
	return &PayWithQueueController{
		amqpChannel: amqpChannel,
	}, nil
}

func (p *PayWithQueueController) Pay(ctx echo.Context) error {
	var request PayRequest
	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusBadRequest, InvalidRequestResponse(err.Error()))
	}

	senderID, err := util.GetUserIDFromEchoContext(ctx)
	if err != nil {
		log.Printf("could not get user id from context: %s", err.Error())
		return ctx.JSON(http.StatusInternalServerError, InternalErrorResponse())
	}

	transaction := model.Transaction{
		SenderID:   senderID,
		ReceiverID: request.To,
		Amount:     request.Amount,
	}
	if _, err := repository.ExecuteTransaction(ctx.Request().Context(), transaction); err != nil {
		if errors.Is(err, repository.InsufficientBalance{}) {
			log.Print("insufficient balance")
			return ctx.JSON(http.StatusUnprocessableEntity, GenericResponse("insufficient balance", err.Error()))
		} else {
			log.Printf("could not execute transaction: %s", err.Error())
			return ctx.JSON(http.StatusInternalServerError, InternalErrorResponse())
		}
	}

	// Publish to Notification Queue
	if err := p.PublishNotificationQueue(transaction.SenderID, "Money sent!"); err != nil {
		log.Printf("could publish notification to sender: %s", err.Error())
		return ctx.JSON(http.StatusInternalServerError, InternalErrorResponse())
	}
	if err := p.PublishNotificationQueue(transaction.ReceiverID, "Money received!"); err != nil {
		log.Printf("could publish notification to receiver: %s", err.Error())
		return ctx.JSON(http.StatusInternalServerError, InternalErrorResponse())
	}

	// Publish to Email Queue
	sender, err := repository.FetchUserByID(ctx.Request().Context(), transaction.SenderID)
	if err != nil {
		log.Printf("could not fetch sender data: %s", err.Error())
		return ctx.JSON(http.StatusInternalServerError, InternalErrorResponse())
	}
	receiver, err := repository.FetchUserByID(ctx.Request().Context(), transaction.ReceiverID)
	if err != nil {
		log.Printf("could not fetch receiver data: %s", err.Error())
		return ctx.JSON(http.StatusInternalServerError, InternalErrorResponse())
	}
	if sender.ID == 0 || receiver.ID == 0 {
		log.Print("users data not found")
		return ctx.JSON(http.StatusInternalServerError, InternalErrorResponse())
	}
	if err := p.PublishEmailQueue(sender.Email, "Money sent!"); err != nil {
		log.Printf("could publish email to sender: %s", err.Error())
		return ctx.JSON(http.StatusInternalServerError, InternalErrorResponse())
	}
	if err := p.PublishEmailQueue(receiver.Email, "Money received!"); err != nil {
		log.Printf("could publish email to receiver: %s", err.Error())
		return ctx.JSON(http.StatusInternalServerError, InternalErrorResponse())
	}

	return ctx.JSON(http.StatusOK, SuccessResonse(nil))
}

func (p *PayWithQueueController) PublishNotificationQueue(userID uint, message string) error {
	notificationPayload, err := json.Marshal(map[string]interface{}{
		"user_id": userID,
		"message": message,
	})
	if err != nil {
		return err
	}
	if err := p.amqpChannel.Publish(PushNotificiationExchange, "", false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        notificationPayload,
	}); err != nil {
		return err
	}
	return nil
}

func (p *PayWithQueueController) PublishEmailQueue(email, message string) error {
	emailPayload, err := json.Marshal(map[string]interface{}{
		"email":   email,
		"message": message,
	})
	if err != nil {
		return err
	}
	if err := p.amqpChannel.Publish(EmailExchange, "", false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        emailPayload,
	}); err != nil {
		return err
	}
	return nil
}
