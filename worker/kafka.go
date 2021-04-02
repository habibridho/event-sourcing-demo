package worker

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
)

type MessageHandler interface {
	Handle(msg []byte) error
}

type KafkaWorker struct {
	consumer *kafka.Consumer
	handler  MessageHandler
}

func NewKafkaConsumer(topic, groupID string, handler MessageHandler) *KafkaWorker {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id":          groupID,
	})
	if err != nil {
		log.Fatalf("could not create kafka consumer: %s", err.Error())
	}
	if err := consumer.Subscribe(topic, nil); err != nil {
		log.Fatalf("could not subscribe to topic: %s", err.Error())
	}

	return &KafkaWorker{
		consumer: consumer,
		handler:  handler,
	}
}

func (k *KafkaWorker) Consume() {
	defer k.consumer.Close()
	for {
		msg, err := k.consumer.ReadMessage(-1)
		if err == nil {
			log.Printf("Message on %s: %s", msg.TopicPartition, string(msg.Value))
			if err := k.handler.Handle(msg.Value); err != nil {
				// TODO: move to dlq
			}
		} else {
			// The client will automatically try to recover from all errors.
			log.Printf("Consumer error: %v (%v)", err, msg)
		}
	}
}
