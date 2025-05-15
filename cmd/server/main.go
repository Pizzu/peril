package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	amqp "github.com/rabbitmq/amqp091-go"
)

const rabbitMqConnStr string = "amqp://guest:guest@localhost:15672/"

func main() {
	fmt.Println("Starting Peril server...")
	conn, err := amqp.Dial(rabbitMqConnStr)

	if err != nil {
		log.Fatal(err.Error())
	}

	defer conn.Close()
	fmt.Println("Connection successfully made...")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Program gracefully stopped")

}
