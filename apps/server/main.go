package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/avalkov/SCS/internal/configuration"
	ds "github.com/avalkov/SCS/internal/datastructures"
	commandsParser "github.com/avalkov/SCS/internal/domain/commands_parser"
	commandsProcessor "github.com/avalkov/SCS/internal/domain/commands_processor"
	"github.com/avalkov/SCS/internal/multiworker"
	"github.com/avalkov/SCS/internal/queueservice"
	"github.com/avalkov/SCS/internal/queueservice/amqp"
	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func run() error {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	godotenv.Load(".env")

	var config configuration.Config
	if err := envconfig.Process(ctx, &config); err != nil {
		return err
	}

	amqpWorker := amqp.NewAmqpWorker(amqp.AmqpConfig{
		Host:      config.AMQP.Host,
		Port:      config.AMQP.Port,
		User:      config.AMQP.User,
		Pass:      config.AMQP.Pass,
		QueueName: config.AMQP.QueueName,
	})
	if err := amqpWorker.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize AMQP worker: %v", err)
	}

	defer amqpWorker.Close()

	requestsChans := make(chan queueservice.Message)
	repliesChans := make(chan queueservice.Message)

	defer close(requestsChans)
	defer close(repliesChans)

	if err := amqpWorker.Run(ctx, requestsChans, repliesChans); err != nil {
		return fmt.Errorf("failed to run AMQP worker: %v", err)
	}

	parser := commandsParser.NewCommandsParser()
	processor := commandsProcessor.NewCommandsProcessor(parser, ds.NewOrderedMap())

	multiWorker := multiworker.NewMultiWorker(
		config.PROCESSING_WORKERS_COUNT,
		parser,
		processor,
	)

	go multiWorker.Run(ctx, requestsChans, repliesChans)

	log.Printf("Server started")

	<-ctx.Done()

	return nil
}
