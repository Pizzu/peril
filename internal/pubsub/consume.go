package pubsub

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Acktype int
type SimpleQueueType int

const (
	SimpleQueueDurable SimpleQueueType = iota
	SimpleQueueTransient
)

const (
	Ack Acktype = iota
	NackDiscard
	NackRequeue
)

func DeclareAndBind(conn *amqp.Connection, exchange, queueName, key string, simpleQueueType SimpleQueueType) (*amqp.Channel, amqp.Queue, error) {
	channel, err := conn.Channel()

	if err != nil {
		return nil, amqp.Queue{}, err
	}

	var queue amqp.Queue
	var args amqp.Table = amqp.Table{"x-dead-letter-exchange": routing.ExchangeDeadLetterFanout}

	switch simpleQueueType {
	case SimpleQueueDurable:
		queue, err = channel.QueueDeclare(queueName, true, false, false, false, args)
	case SimpleQueueTransient:
		queue, err = channel.QueueDeclare(queueName, false, true, true, false, args)
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

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange, queueName,
	key string,
	simpleQueueType SimpleQueueType,
	handler func(T) Acktype,
) error {
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
				continue
			}
			switch handler(content) {
			case Ack:
				message.Ack(false)
				fmt.Println("Ack")
			case NackDiscard:
				message.Nack(false, false)
				fmt.Println("NackDiscard")
			case NackRequeue:
				message.Nack(false, true)
				fmt.Println("NackRequeue")
			}
		}
	}()

	return nil
}
