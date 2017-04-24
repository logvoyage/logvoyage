// Producer accepts logs from various sources and sends them to queue.
// The main task of this service - to be up and running 100% of time and run fast.
// Validation and other things should be done by background workers.
package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"

	"bitbucket.org/firstrow/logvoyage/shared/config"
	"github.com/firstrow/tcp_server"
	"github.com/streadway/amqp"
)

// RabbitMQ channel to write incoming messages
var channel *amqp.Channel

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

// httpHandler accepts incoming logs data. Each request should contain logs separated by new line.
// Example:
// POST https://data.logvoyage.com?apiKey=user1
// log line 1
// log line 2
func httpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}

	apiKey := r.URL.Query().Get("apiKey")
	if apiKey == "" {
		return
	}

	tag := r.URL.Query().Get("tag")
	if tag == "" {
		tag = "default"
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading request body")
	}

	if len(body) == 0 {
		log.Println("Got empty request")
	}

	log.Println("HTTP Message:", string(body))

	// Send message to rabbitmq.
	// In case rabbit node is down - use Amazon SQS.
	// TODO: Iterate over each line and append apiKey
	for _, line := range strings.Split(string(body), "\n") {
		go sendToQueue([]byte(fmt.Sprintf("%s@%s %s", apiKey, tag, line)))
	}
}

func startHttpHandler() {
	http.HandleFunc("/", httpHandler)
	err := http.ListenAndServe(":27000", nil)

	if err != nil {
		log.Fatalln("Error starting http server:", err)
	}
}

func startTCPHandler() {
	tcp := tcp_server.New("localhost:28000")

	tcp.OnNewClient(func(c *tcp_server.Client) {
		log.Println("New TCP client connected")
	})
	tcp.OnNewMessage(func(c *tcp_server.Client, message string) {
		log.Println("New message received:", message)
		sendToQueue([]byte(message))
	})
	tcp.OnClientConnectionClosed(func(c *tcp_server.Client, err error) {
		log.Println("Connection to client lost")
	})

	tcp.Listen()
}

func startUDPHandler() {
	updAddr, err := net.ResolveUDPAddr("udp", ":29000")
	failOnError(err, "Error resolving UDP address")

	conn, err := net.ListenUDP("udp", updAddr)
	failOnError(err, "Error listening to UDP")
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		log.Println("UDP Message", message)
		if err != nil {
			log.Println("Error reading UDP message:", err)
		}
	}
}

// sendToQueue sends data to queue or if it is down - to fallback queue
// body should be an array of bytes. First N bytes should be project uuid.
func sendToQueue(body []byte) {
	err := channel.Publish(
		"logs", // exchange
		"all",  // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})

	if err != nil {
		log.Println("Error publishing message:", string(body))
	}
}

func main() {
	amqpConn, err := amqp.Dial(config.Get("amqp.url"))
	failOnError(err, "Failed to connect to RabbitMQ")
	defer amqpConn.Close()

	// TODO: Put channel into confirm mode.
	channel, err = amqpConn.Channel()
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
	failOnError(err, "Failed to declare a queue")

	go startHttpHandler()
	go startTCPHandler()
	go startUDPHandler()
	select {}
}
