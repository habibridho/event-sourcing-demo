package controller

import (
	"errors"
	"event-sourcing-demo/model"
	"event-sourcing-demo/repository"
	"event-sourcing-demo/util"
	"github.com/labstack/echo"
	"log"
	"net/http"
)

type PayRequest struct {
	To     uint   `json:"to"`
	Amount uint64 `json:"amount"`
}

type NotificationHandler interface {
	SendNotification(userID uint, message string) error
}

type EmailHandler interface {
	SendEmail(userID uint, message string) error
}

type PayController struct {
	NotificationHandler
	EmailHandler
}

func (p *PayController) Pay(ctx echo.Context) error {
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

	// Send Notification to Users
	if err := p.sendNotification(transaction.SenderID, transaction.ReceiverID); err != nil {
		log.Printf("could not send notification: %s", err.Error())
		return ctx.JSON(http.StatusInternalServerError, InternalErrorResponse())
	}

	// Send Email to Users
	if err := p.sendEmail(transaction.SenderID, transaction.ReceiverID); err != nil {
		log.Printf("could not send email: %s", err.Error())
		return ctx.JSON(http.StatusInternalServerError, InternalErrorResponse())
	}

	return ctx.JSON(http.StatusOK, SuccessResonse(nil))
}

func (p *PayController) sendNotification(senderID, receiverID uint) error {
	if err := p.SendNotification(senderID, "money sent!"); err != nil {
		return err
	}
	if err := p.SendNotification(receiverID, "money received!"); err != nil {
		return err
	}
	return nil
}

func (p *PayController) sendEmail(senderID, receiverID uint) error {
	if err := p.SendEmail(senderID, "money sent!"); err != nil {
		return err
	}
	if err := p.SendEmail(receiverID, "money received!"); err != nil {
		return err
	}
	return nil
}
