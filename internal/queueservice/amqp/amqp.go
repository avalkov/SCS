package amqp

import (
	"context"
	"fmt"
	"log"

	"github.com/streadway/amqp"

	qs "github.com/avalkov/SCS/internal/queueservice"
)

type AmqpConfig struct {
	Host      string
	Port      int
	User      string
	Pass      string
	QueueName string
}

type amqpWorker struct {
	config  AmqpConfig
	channel *amqp.Channel
	msgs    <-chan amqp.Delivery
}

func NewAmqpWorker(config AmqpConfig) qs.QueueService {
	return &amqpWorker{
		config:  config,
		channel: nil,
		msgs:    nil,
	}
}

func (aw *amqpWorker) Initialize() error {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d/",
		aw.config.User, aw.config.Pass, aw.config.Host, aw.config.Port))
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %v", err)
	}

	q, err := ch.QueueDeclare(
		aw.config.QueueName, // name
		false,               // durable
		false,               // delete when unused
		false,               // exclusive
		false,               // no-wait
		nil,                 // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare a queue: %v", err)
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return fmt.Errorf("failed to register a consumer: %v", err)
	}

	aw.channel = ch
	aw.msgs = msgs

	return nil
}

func (aw *amqpWorker) Run(ctx context.Context, requests chan<- qs.Message, replies <-chan qs.Message) error {
	if aw.msgs == nil {
		return fmt.Errorf("worker not initialized")
	}
	// TODO: Can run multiple recievers and senders
	go aw.runReceiver(ctx, requests)
	go aw.runSender(ctx, replies)
	return nil
}

func (aw *amqpWorker) runReceiver(ctx context.Context, requests chan<- qs.Message) {
	defer ctx.Done()
	for d := range aw.msgs {
		log.Printf("Received a message: %s", d.Body)

		requests <- qs.Message{
			Body:          string(d.Body),
			ReplyTo:       d.ReplyTo,
			CorrelationId: d.CorrelationId,
		}
	}

	log.Printf("Consumer channel closed")
}

func (aw *amqpWorker) runSender(ctx context.Context, replies <-chan qs.Message) {
	defer ctx.Done()
	for d := range replies {
		if err := aw.channel.Publish(
			"",        // exchange
			d.ReplyTo, // routing key (reply_to)
			false,     // mandatory
			false,     // immediate
			amqp.Publishing{
				ContentType:   "text/plain",
				CorrelationId: d.CorrelationId,
				Body:          []byte(d.Body),
			}); err != nil {
			log.Printf("Failed to publish a message: %v", err)
		}
	}

	log.Printf("Publisher channel closed")
}

func (aw *amqpWorker) Close() error {
	if aw.channel == nil {
		return fmt.Errorf("worker not initialized")
	}
	if err := aw.channel.Close(); err != nil {
		return fmt.Errorf("failed to close channel: %v", err)
	}
	return nil
}
