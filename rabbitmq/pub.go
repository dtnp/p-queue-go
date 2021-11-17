package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
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
		exchangeName = flag.String("t", "", "Exchange name")
		routingKey = flag.String("r", "", "Routing key")
	)
	flag.Parse()

	if *exchangeName == "" {
		log.Fatal("missing exchange name; use -t to specify the exchange name")
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
		log.Fatalf("Failed to declare an exchange (exchange)")
	}


	// Publish some messages, synchronously
	var i int64
	for {
		i++
		now := time.Now().Format(time.RFC3339)
		payload := fmt.Sprintf("%08d %s", i, now)
		log.Printf("Send: %s, exchange: %s, routingKey: %s\n", payload, *exchangeName, *routingKey)
		body, _ := json.Marshal(payload)
		err = ch.Publish(
			*exchangeName, // exchange
			*routingKey,   // routing key
			false,        // mandatory
			false,        // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			})
		if err != nil {
			log.Fatalf("Failed to publish a message")
		}

		// Sleep for a random time of up to 1s
		time.Sleep(time.Duration(rand.Int63n(1000)) * time.Millisecond)
	}
}
