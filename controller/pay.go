package controller

import (
	"context"
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
	SendEmail(to, template string, data map[string]interface{}) error
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
	if err := p.sendEmail(ctx.Request().Context(), transaction.SenderID, transaction.ReceiverID, transaction.Amount); err != nil {
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

func (p *PayController) sendEmail(ctx context.Context, senderID, receiverID uint, amount uint64) error {
	sender, err := repository.FetchUserByID(ctx, senderID)
	if err != nil {
		log.Printf("could not fetch sender data: %s", err.Error())
		return err
	}
	receiver, err := repository.FetchUserByID(ctx, receiverID)
	if err != nil {
		log.Printf("could not fetch receiver data: %s", err.Error())
		return err
	}
	if sender.ID == 0 || receiver.ID == 0 {
		err := errors.New("users data not found")
		log.Print(err.Error())
		return err
	}

	emailData := map[string]interface{}{
		"sender_name":   sender.Name,
		"receiver_name": receiver.Name,
		"amount":        amount,
	}

	if err := p.SendEmail(sender.Email, "sender-email", emailData); err != nil {
		return err
	}
	if err := p.SendEmail(receiver.Email, "receiver-email", emailData); err != nil {
		return err
	}
	return nil
}
