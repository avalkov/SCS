package queueservice

import "context"

type QueueService interface {
	Initialize() error
	Run(ctx context.Context, requests chan<- Message, replies <-chan Message) error
	Close() error
}

type Message struct {
	Body          string
	ReplyTo       string
	CorrelationId string
}
