// Consumer service accpets logs and other events from queue, validates,
// and pushes to index and persistance storage.
package main

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"

	"bitbucket.org/firstrow/logvoyage/models"
	"github.com/streadway/amqp"
)

// See https://regex101.com/r/sQskdz/1
var msgRegExp = regexp.MustCompile(`^([a-z0-9]{8}-[a-z0-9]{4}-[1-5][a-z0-9]{3}-[a-z0-9]{4}-[a-z0-9]{12})(@([a-zA-Z0-9\_\-\.]{1,20}))?\s*(.*)`)

// Default tag name name in case tag does not specified in source message.
const defaultTag = "default"

// processMessage extracts apiKey, tag(optional) and source message from log line.
// In case tag not found, "default" tag will be returned.
// Error will be returned if uuid or tag has wrong format.
//
// Examples:
// uuid@tag msg -> uuid, tag, msg, nil
// uuid msg     -> uuid, default, msg, nil
func processMessage(msg string) (string, string, string, error) {
	result := msgRegExp.FindAllStringSubmatch(msg, -1)

	if result == nil {
		return "", "", "", errors.New("Error parsing source message")
	}

	var tag string
	if result[0][3] == "" {
		tag = defaultTag
	} else {
		tag = result[0][3]
	}

	return result[0][1], tag, strings.TrimSpace(result[0][4]), nil
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	amqpConn, err := amqp.Dial("amqp://guest:guest@ubuntu:5672")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer amqpConn.Close()

	channel, err := amqpConn.Channel()
	failOnError(err, "Failed to open a channel")
	defer channel.Close()

	err = channel.ExchangeDeclare(
		"logs",   // name
		"direct", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // noWait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare a exchange")

	queue, err := channel.QueueDeclare(
		"all", // name of the queue
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // noWait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = channel.QueueBind(
		queue.Name, // name of the queue
		"all",      // bindingKey
		"logs",     // sourceExchange
		false,      // noWait
		nil,        // arguments
	)

	deliveries, err := channel.Consume(
		queue.Name, // name
		"tag",      // consumerTag,
		false,      // noAck
		false,      // exclusive
		false,      // noLocal
		false,      // noWait
		nil,        // arguments
	)
	failOnError(err, "Failed to consume")

	// TODO: Handle close method. See example consumer.
	go handle(deliveries)
	select {}
}

func handle(deliveries <-chan amqp.Delivery) {
	for d := range deliveries {
		log.Printf(
			"got %dB delivery: [%v] %q",
			len(d.Body),
			d.DeliveryTag,
			d.Body,
		)

		d.Ack(false)

		projectUUID, _, _, err := processMessage(string(d.Body))

		if err != nil {
			log.Println("Error processing message:", err)
			continue
		}

		// TODO: Cache project in mem
		project, err := models.FindProjectByUUID(projectUUID)

		if err != nil {
			log.Println("Project not found:", err)
			continue
		}

		fmt.Println(project)

		// sendToElastic(project.GetLogsElasticSearchIndexName(), tag, msg)
		// sendToStorage(...)

		// + Extract API key and tag
		// Check if API key exists and store keys in mem
		// Add message to bulk storage. Bulk storage per index?
		// Bulk storage pushes messages to elastic each 2 seconds
		// Send message to persistent storage
		// handle errors
	}
	log.Fatalln("Deliveries channel closed")
}
