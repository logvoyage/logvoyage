// Consumer service accpets logs and other events from queue, validates,
// and pushes to index and persistance storage.
package main

import (
	"errors"
	"log"
	"regexp"
	"strings"
	"sync"
	"time"

	"bitbucket.org/firstrow/logvoyage/models"
	"github.com/streadway/amqp"
)

const (
	// Default tag name name in case tag does not specified in source message.
	defaultTag = "default"
	// How often in seconds send messages to storge.
	persistTimeout = 2
)

var (
	storage = newInMemStorage()
	// See https://regex101.com/r/sQskdz/1
	msgRegExp = regexp.MustCompile(`^([a-z0-9]{8}-[a-z0-9]{4}-[1-5][a-z0-9]{3}-[a-z0-9]{4}-[a-z0-9]{12})(@([a-zA-Z0-9\_\-\.]{1,20}))?\s*(.*)`)
)

type message struct {
	ProjectUUID string
	Tag         string
	Source      string
	Datetime    time.Time
}

// inMemStorage stores all received valid message in memory.
// Each N seconds messages will be sent to ES via bulk request.
type inMemStorage struct {
	sync.Mutex
	messages []message
}

func newInMemStorage() *inMemStorage {
	s := &inMemStorage{}
	go s.startTimer()
	return s
}

func (s *inMemStorage) Add(msg message) {
	s.Lock()
	s.messages = append(s.messages, msg)
	s.Unlock()
}

func (s *inMemStorage) Persist() {
	s.Lock()
	defer s.Unlock()
	if len(s.messages) == 0 {
		return
	}
	tmp := make([]message, len(s.messages))
	copy(s.messages, tmp)
}

func (s *inMemStorage) startTimer() {
	ticker := time.NewTicker(persistTimeout * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		s.Persist()
	}
}

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
		processDelivery(d)
	}
	log.Fatalln("Deliveries channel closed")
}

func processDelivery(d amqp.Delivery) {
	defer d.Ack(false)

	projectUUID, tag, msgSource, err := processMessage(string(d.Body))

	if err != nil {
		log.Println("Error processing message:", err)
		return
	}

	// TODO: Cache project in mem
	project, err := models.FindProjectByUUID(projectUUID)

	if err != nil {
		log.Println("Project not found")
		return
	}

	storage.Add(message{
		ProjectUUID: project.UUID,
		Tag:         tag,
		Source:      msgSource,
	})
}
