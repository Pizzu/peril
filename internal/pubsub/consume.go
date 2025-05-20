package pubsub

import (
	"encoding/json"
	"errors"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type SimpleQueueType int

const (
	SimpleQueueDurable SimpleQueueType = iota
	SimpleQueueTransient
)

func DeclareAndBind(conn *amqp.Connection, exchange, queueName, key string, simpleQueueType SimpleQueueType) (*amqp.Channel, amqp.Queue, error) {
	channel, err := conn.Channel()

	if err != nil {
		return nil, amqp.Queue{}, err
	}

	var queue amqp.Queue

	switch simpleQueueType {
	case SimpleQueueDurable:
		queue, err = channel.QueueDeclare(queueName, true, false, false, false, nil)
	case SimpleQueueTransient:
		queue, err = channel.QueueDeclare(queueName, false, true, true, false, nil)
	default:
		return nil, amqp.Queue{}, errors.New("wrong queue type")
	}

	if err != nil {
		return nil, amqp.Queue{}, err
	}

	err = channel.QueueBind(queue.Name, key, exchange, false, nil)

	if err != nil {
		return nil, amqp.Queue{}, err
	}

	return channel, queue, nil
}

func SubscribeJSON[T any](conn *amqp.Connection, exchange, queueName, key string, simpleQueueType SimpleQueueType, handler func(T)) error {
	channel, queue, err := DeclareAndBind(conn, exchange, queueName, key, simpleQueueType)

	if err != nil {
		return err
	}

	deliveryChan, err := channel.Consume(queue.Name, "", false, false, false, false, nil)

	if err != nil {
		return err
	}

	go func() {
		for message := range deliveryChan {
			var content T
			err := json.Unmarshal(message.Body, &content)
			if err != nil {
				log.Printf("Failed to unmarshal message: %v", err)
				message.Nack(false, false)
				continue
			}
			handler(content)
			message.Ack(false)
		}
	}()

	return nil
}
