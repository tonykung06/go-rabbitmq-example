package main

import (
	"log"

	"github.com/go-rabbitmq-example/common"
	"github.com/streadway/amqp"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	common.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	common.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	common.FailOnError(err, "Failed to declare a queue")
	body := "hello"
	go sendMsg(conn, q, body)
	<-make(chan bool)
}

func sendMsg(conn *amqp.Connection, q amqp.Queue, body string) {
	ch, err := conn.Channel()
	common.FailOnError(err, "Failed to open a channel")
	defer ch.Close()
	for {
		err := ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			},
		)
		log.Printf(" [x] Sent %s", body)
		common.FailOnError(err, "Failed to publish a message")
	}
}
