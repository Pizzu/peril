package pubsub

import (
	"errors"

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

	return channel, amqp.Queue{}, nil
}
