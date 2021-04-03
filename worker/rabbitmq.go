package worker

import (
	"github.com/streadway/amqp"
	"log"
)

type RabbitMqWorker struct {
	channel *amqp.Channel
	queue   amqp.Queue
	handler MessageHandler
}

func NewRabbitMqWorker(conn *amqp.Connection, exchange string, handler MessageHandler) *RabbitMqWorker {
	channel, err := conn.Channel()
	if err != nil {
		log.Fatalf("could not rabbitmq channel: %s", err.Error())
	}
	if err := channel.ExchangeDeclare(exchange, "fanout", true, false, false, false, nil); err != nil {
		log.Fatalf("could not declare exchange: %s", err.Error())
	}

	q, err := channel.QueueDeclare("", false, false, true, false, nil)
	if err != nil {
		log.Fatalf("could not declare queue: %s", err.Error())
	}

	if err := channel.QueueBind(q.Name, "", exchange, false, nil); err != nil {
		log.Fatalf("could not bind queue: %s", err.Error())
	}

	return &RabbitMqWorker{
		channel: channel,
		queue:   q,
		handler: handler,
	}

}

func (r *RabbitMqWorker) Consume() {
	defer r.channel.Close()
	messages, err := r.channel.Consume(r.queue.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("could not consume message: %s", err.Error())
	}

	for msg := range messages {
		if err := r.handler.Handle(msg.Body); err != nil {
			// TODO: move to dlq
			log.Printf("could not handle message: %s", err.Error())
		}
	}
}
