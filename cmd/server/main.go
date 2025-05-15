package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

const rabbitMqConnStr string = "amqp://guest:guest@localhost:5672/"

func main() {
	fmt.Println("Starting Peril server...")
	conn, err := amqp.Dial(rabbitMqConnStr)

	if err != nil {
		log.Fatal(err.Error())
	}

	defer conn.Close()
	fmt.Println("Connection successfully made...")

	channel, err := conn.Channel()

	if err != nil {
		log.Fatal(err.Error())
	}

	playingState := routing.PlayingState{IsPaused: true}

	err = pubsub.PublishJSON(channel, routing.ExchangePerilDirect, routing.PauseKey, playingState)

	if err != nil {
		fmt.Println(err.Error())
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Program gracefully stopped")

}
