package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/avalkov/SCS/internal/configuration"
	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
	"github.com/streadway/amqp"
)

type Message struct {
	Body          string
	ReplyTo       string
	CorrelationId string
}

const (
	replyQueue = "reply_queue"
)

var jsonArrayRegex = regexp.MustCompile(`^\s*\[\s*.*\s*\]\s*$`)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func connectToRabbitMQ(user, password, host string, port int) (*amqp.Connection, *amqp.Channel) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d/", user, password, host, port))
	failOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	return conn, ch
}

func declareQueue(ch *amqp.Channel, name string) amqp.Queue {
	q, err := ch.QueueDeclare(
		name,  // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, fmt.Sprintf("Failed to declare queue: %s", name))
	return q
}

func loadCommandsFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var commands []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		commands = append(commands, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return commands, nil
}

func sendCommand(ch *amqp.Channel, command string, requestQueue, replyQueue string) {
	corrId := fmt.Sprintf("%d", time.Now().UnixNano())

	err := ch.Publish(
		"",           // exchange
		requestQueue, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrId,
			ReplyTo:       replyQueue,
			Body:          []byte(command),
		})
	failOnError(err, "Failed to publish a message")
}

func receiveReplies(ch *amqp.Channel, replyQueue string, clientID int) {
	msgs, err := ch.Consume(
		replyQueue, // queue
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	failOnError(err, "Failed to register a consumer")

	for d := range msgs {
		fmt.Printf("Client %d received reply: %s\n", clientID, d.Body)
		if jsonArrayRegex.Match(d.Body) {
			saveToFile(fmt.Sprintf("getAllItemsResponse_client_%d.json", clientID), d.Body)
		}
	}
}

func saveToFile(filename string, data []byte) {
	file, err := os.Create(filename)
	failOnError(err, fmt.Sprintf("Failed to create file: %s", filename))
	defer file.Close()

	_, err = file.Write(data)
	failOnError(err, fmt.Sprintf("Failed to write to file: %s", filename))
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <commands_file_1> <commands_file_2> ... <commands_file_n>", os.Args[0])
	}

	godotenv.Load(".env")

	var config configuration.Config
	if err := envconfig.Process(context.Background(), &config); err != nil {
		log.Fatalf("Failed to process env var: %v", err)
	}

	conn, ch := connectToRabbitMQ(config.AMQP.User, config.AMQP.Pass, config.AMQP.Host, config.AMQP.Port)
	defer conn.Close()

	replyQueue := declareQueue(ch, replyQueue)

	for i, filename := range os.Args[1:] {
		log.Printf("Processing commands from file: %s", filename)
		commands, err := loadCommandsFromFile(filename)
		failOnError(err, fmt.Sprintf("Failed to load commands from file %s", filename))

		go func(clientID int, commands []string) {
			go receiveReplies(ch, replyQueue.Name, clientID)
			for _, command := range commands {
				sendCommand(ch, command, config.AMQP.QueueName, replyQueue.Name)
			}
		}(i, commands)
	}

	ctx := context.Background()
	<-ctx.Done()
}
