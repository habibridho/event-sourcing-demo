package handler

import (
	"encoding/json"
	"errors"
	"event-sourcing-demo/controller"
	"fmt"
	"github.com/go-resty/resty/v2"
	"log"
)

var (
	TestEmail = "neriroge@getnada.com"
)

type KafkaEmailHandler struct{}

func (k *KafkaEmailHandler) Handle(msg []byte) error {
	var payEvent controller.PayEvent
	if err := json.Unmarshal(msg, &payEvent); err != nil {
		log.Printf("could not unmarshal message: %s", err.Error())
		return err
	}
	mailData := map[string]interface{}{
		"sender":   payEvent.Sender.Name,
		"amount":   payEvent.Amount,
		"receiver": payEvent.Receiver.Name,
	}
	if err := SendEmail(payEvent.Sender.Email, "sender-email", mailData); err != nil {
		log.Printf("could not send email: %s", err.Error())
		return err
	}

	if err := SendEmail(payEvent.Receiver.Email, "receiver-email", mailData); err != nil {
		log.Printf("could not send email: %s", err.Error())
		return err
	}
	return nil
}

type RabbitMqEmailHandler struct{}

func (r *RabbitMqEmailHandler) Handle(msg []byte) error {
	var data map[string]interface{}
	if err := json.Unmarshal(msg, &data); err != nil {
		log.Printf("could not unmarshal message: %s", err.Error())
		return err
	}

	mailData := map[string]interface{}{
		"sender":   fmt.Sprintf("%s", data["sender_name"]),
		"amount":   fmt.Sprintf("%v", data["amount"]),
		"receiver": fmt.Sprintf("%s", data["receiver_name"]),
	}
	if err := SendEmail(fmt.Sprintf("%s", data["email"]), fmt.Sprintf("%s", data["template"]), mailData); err != nil {
		log.Printf("could not send email: %s", err.Error())
		return err
	}
	return nil
}

type SynchronousEmailHandler struct{}

func (s *SynchronousEmailHandler) SendEmail(to, template string, data map[string]interface{}) error {
	mailData := map[string]interface{}{
		"sender":   fmt.Sprintf("%s", data["sender_name"]),
		"amount":   fmt.Sprintf("%v", data["amount"]),
		"receiver": fmt.Sprintf("%s", data["receiver_name"]),
	}
	if err := SendEmail(to, template, mailData); err != nil {
		log.Printf("could not send email: %s", err.Error())
		return err
	}
	return nil
}

func SendEmail(destination, template string, data map[string]interface{}) error {
	if destination == "habib@email.com" {
		destination = TestEmail
	} else {
		return nil
	}

	raw, err := json.Marshal(data)
	if err != nil {
		log.Printf("could not marshal data: %s", err.Error())
		return err
	}

	formData := map[string]string{
		"from":                  "Mailgun Sandbox <postmaster@sandboxca3a59cd67284504be30a04d0c648b3a.mailgun.org>",
		"to":                    destination,
		"subject":               "Demo Email",
		"template":              template,
		"h:X-Mailgun-Variables": string(raw),
	}

	url := fmt.Sprintf("https://api.mailgun.net/v3/%s/messages", "sandboxca3a59cd67284504be30a04d0c648b3a.mailgun.org")
	response, err := resty.New().R().
		SetBasicAuth("api", "key-bf1e96060950d45a9df0c6a444dad16e").
		SetFormData(formData).
		Post(url)
	if err != nil {
		log.Printf("could not call mailgun api: %s", err.Error())
		return err
	}

	if response.IsError() {
		responseBody := response.Body()
		err = errors.New(fmt.Sprintf("failed sending mail: %s", string(responseBody)))
		return err
	}

	return nil
}
