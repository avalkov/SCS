package commandsprocessor

import (
	"sync"
	"testing"

	ds "github.com/avalkov/SCS/internal/datastructures"
	cmdParser "github.com/avalkov/SCS/internal/domain/commands_parser"
	"github.com/avalkov/SCS/internal/queueservice"
)

type mockCommandsParser struct{}

func (m *mockCommandsParser) ParseCommand(command string) (cmdParser.Command, error) {
	cp := cmdParser.NewCommandsParser()
	return cp.ParseCommand(command)
}

func TestCommandsProcessor_Process(t *testing.T) {
	cmdParser := &mockCommandsParser{}
	cp := NewCommandsProcessor(cmdParser, ds.NewOrderedMap())

	tests := []struct {
		name          string
		request       queueservice.Message
		expectedReply queueservice.Message
	}{
		{
			name: "AddItem",
			request: queueservice.Message{
				Body: "addItem('key1', 'value1')",
			},
			expectedReply: queueservice.Message{},
		},
		{
			name: "GetItem Exists",
			request: queueservice.Message{
				Body: "getItem('key1')",
			},
			expectedReply: queueservice.Message{
				Body: "value1",
			},
		},
		{
			name: "DeleteItem",
			request: queueservice.Message{
				Body: "deleteItem('key1')",
			},
			expectedReply: queueservice.Message{},
		},
		{
			name: "GetItem Not Exists",
			request: queueservice.Message{
				Body: "getItem('key1')",
			},
			expectedReply: queueservice.Message{
				Body: "Key not found: key1",
			},
		},
		{
			name: "GetAllItems",
			request: queueservice.Message{
				Body: "getAllItems()",
			},
			expectedReply: queueservice.Message{
				Body: `[]`,
			},
		},
		{
			name: "Invalid Command",
			request: queueservice.Message{
				Body: "invalidCommand('key1')",
			},
			expectedReply: queueservice.Message{
				Body: "Failed to parse command: invalid command format",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requests := make(chan queueservice.Message, 1)
			replies := make(chan queueservice.Message, 1)

			var wg sync.WaitGroup
			wg.Add(1)

			go cp.Process(1, requests, replies, &wg)

			requests <- tt.request
			close(requests)

			wg.Wait()
			close(replies)

			for reply := range replies {
				if reply.Body != tt.expectedReply.Body {
					t.Errorf("expected reply body: %s, got: %s", tt.expectedReply.Body, reply.Body)
				}
			}
		})
	}
}
