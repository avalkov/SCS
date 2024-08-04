package multiworker

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/avalkov/SCS/internal/queueservice"
)

func TestMultiWorker_Run(t *testing.T) {
	commandsParser := &mockCommandsParser{}
	commandsProcessor := &mockCommandProcessor{}
	multiWorker := NewMultiWorker(3, commandsParser, commandsProcessor)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	requests := make(chan queueservice.Message, 10)
	replies := make(chan queueservice.Message, 10)

	go multiWorker.Run(ctx, requests, replies)

	tests := []struct {
		name          string
		request       queueservice.Message
		expectedReply queueservice.Message
	}{
		{
			name: "Valid Command 1",
			request: queueservice.Message{
				Body: "command1",
			},
			expectedReply: queueservice.Message{
				Body: "processed: command1",
			},
		},
		{
			name: "Valid Command 2",
			request: queueservice.Message{
				Body: "command2",
			},
			expectedReply: queueservice.Message{
				Body: "processed: command2",
			},
		},
		{
			name: "Invalid Command",
			request: queueservice.Message{
				Body: "invalid",
			},
			expectedReply: queueservice.Message{
				Body: "Failed to get command ID",
			},
		},
		{
			name: "Process Error",
			request: queueservice.Message{
				Body: "processError",
			},
			expectedReply: queueservice.Message{
				Body: "error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requests <- tt.request

			reply := <-replies
			if reply.Body != tt.expectedReply.Body {
				t.Errorf("expected reply body: %s, got: %s", tt.expectedReply.Body, reply.Body)
			}
		})
	}

	close(requests)
}

type mockCommandsParser struct{}

func (m *mockCommandsParser) GetCommandID(command string) (string, error) {
	if command == "invalid" {
		return "", fmt.Errorf("invalid command")
	}
	return "key", nil
}

type mockCommandProcessor struct{}

func (m *mockCommandProcessor) Process(processorID int, requests <-chan queueservice.Message, replies chan<- queueservice.Message, wg *sync.WaitGroup) {
	defer wg.Done()
	for req := range requests {
		if req.Body == "processError" {
			replies <- queueservice.Message{
				Body: "error",
			}
			continue
		}
		replies <- queueservice.Message{
			Body: "processed: " + req.Body,
		}
	}
}
