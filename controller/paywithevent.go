package controller

import (
	"context"
	"encoding/json"
	"errors"
	"event-sourcing-demo/model"
	"event-sourcing-demo/repository"
	"event-sourcing-demo/util"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/labstack/echo"
	"log"
	"net/http"
	"time"
)

type PayWithEventController struct {
	Producer *kafka.Producer
}

type PayEvent struct {
	Sender          model.User `json:"sender"`
	Receiver        model.User `json:"receiver"`
	Amount          uint64     `json:"amount"`
	TransactionTime time.Time  `json:"transaction_time"`
}

var (
	PayEventTopic = "pay-events"
)

func (p *PayWithEventController) Pay(ctx echo.Context) error {
	var request PayRequest
	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusBadRequest, InvalidCredentialsResponse())
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
	if err := repository.ExecuteTransaction(ctx.Request().Context(), transaction); err != nil {
		if errors.Is(err, repository.InsufficientBalance{}) {
			log.Print("insufficient balance")
			return ctx.JSON(http.StatusUnprocessableEntity, GenericResponse("insufficient balance", err.Error()))
		} else {
			log.Printf("could not execute transaction: %s", err.Error())
			return ctx.JSON(http.StatusInternalServerError, InternalErrorResponse())
		}
	}

	payEvent, err := p.createPayEvent(ctx.Request().Context(), transaction)
	if err != nil {
		log.Printf("could not create pay event: %s", err.Error())
		return ctx.JSON(http.StatusInternalServerError, InternalErrorResponse())
	}
	rawPayEvent, err := json.Marshal(payEvent)
	if err != nil {
		log.Printf("could not marshall pay event: %s", err.Error())
		return ctx.JSON(http.StatusInternalServerError, InternalErrorResponse())
	}

	if err := p.Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &PayEventTopic, Partition: kafka.PartitionAny},
		Value:          rawPayEvent,
	}, nil); err != nil {
		log.Printf("could not publish pay event: %s", err.Error())
		return ctx.JSON(http.StatusInternalServerError, InternalErrorResponse())
	}

	return ctx.JSON(http.StatusOK, SuccessResonse(nil))
}

func (p *PayWithEventController) createPayEvent(ctx context.Context, transaction model.Transaction) (PayEvent, error) {
	sender, err := repository.FetchUserByID(ctx, transaction.SenderID)
	if err != nil {
		return PayEvent{}, err
	}
	receiver, err := repository.FetchUserByID(ctx, transaction.ReceiverID)
	if err != nil {
		return PayEvent{}, err
	}
	return PayEvent{
		Sender:          sender,
		Receiver:        receiver,
		Amount:          transaction.Amount,
		TransactionTime: time.Now(),
	}, nil
}
