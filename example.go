package main

import (
	"fmt"
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[ERROR] %s", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	defer ch.Close()

	// Function to dynamically declare and consume from a new queue
	consumeFromQueue := func(queueName string) {
		q, err := ch.QueueDeclare("", false, false, true, false, nil)
		if err != nil {
			fmt.Fprintln(os.Stderr, "[ERROR] %s", err)
		}

		msgs, err := ch.Consume(q.Name, "consumer", true, false, false, false, nil)
		if err != nil {
			fmt.Fprintln(os.Stderr, "[ERROR] %s", err)
		}

		go func() {
			for d := range msgs { // storing extra memory here with this "permanent" goroutine
				log.Printf("Received a message: %s", d.Body)
			}
		}()

		fmt.Println("Subscribed to:", queueName)
	}

	// Subscribe to the initial queue
	consumeFromQueue("queue1")

	// Dynamically subscribe (in future) to more queues as needed
	go func() {
		consumeFromQueue("queue2")
		consumeFromQueue("queue3")
	}()

	// Block forever (for this example)
	forever := make(chan bool)
	<-forever
}
