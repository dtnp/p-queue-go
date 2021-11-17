package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/streadway/amqp"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	// Maybe set from environment variable?
	rabbitUrl := os.Getenv("RABBITMQ_URL")
	if rabbitUrl == "" {
		rabbitUrl = "127.0.0.1:25672"
	}

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	var (
		exchangeName        = flag.String("t", "", "Topic name")
		routingKey 		 = flag.String("r", "", "Routing keys")
	)
	flag.Parse()

	if *exchangeName == "" {
		log.Fatal("missing topic name; use -t to specify the topic name")
	}
	if *routingKey == "" {
		log.Fatal("missing routing key; use -r to specify the routing key")
	}

	// Create an AMQP client
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s/", "guest", "guest", rabbitUrl))
	if err != nil {
		log.Fatalf("unable to create rabbitmq (amqp) client: %v", err)
	}
	// Setup individual Channels
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel")
	}
	// Cleanup
	defer conn.Close()
	defer ch.Close()

	// Setup Topic Queues
	q, err := ch.QueueDeclare(
		"",  // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue")
	}

	rks := strings.Split(*routingKey, ",")
	for _, val := range rks {
		log.Printf("Binding to exchange %s with routing key %s", *exchangeName, val)
		err := ch.QueueBind(
			q.Name,  // queue name
			val, // routing key
			*exchangeName,  // exchange
			false,
			nil,
		)
		if err != nil {
			log.Fatalf("Failed to bind a queue")
		}
	}

	// Setup Exchange
	err = ch.ExchangeDeclare(
		*exchangeName, // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare an exchange (topic)")
	}

	msgs, err := ch.Consume(
		"", // queue
		"",        // consumer
		true,      // auto ack
		false,     // exclusive
		false,     // no local
		false,     // no wait
		nil,       // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer")
	}

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			log.Printf(" [x] %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
	if err != nil {
		log.Fatalf("Receive function failed with %v", err)
	}
}
