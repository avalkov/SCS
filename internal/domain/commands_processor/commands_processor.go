package commandsprocessor

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	ds "github.com/avalkov/SCS/internal/datastructures"
	cmd_parser "github.com/avalkov/SCS/internal/domain/commands_parser"
	"github.com/avalkov/SCS/internal/queueservice"
)

type commandsProcessor struct {
	dataStore KeyValueStorage
	cmdParser CommandsParser
	mu        sync.RWMutex
}

func NewCommandsProcessor(cmdParser CommandsParser, keyValueStorage KeyValueStorage) *commandsProcessor {
	return &commandsProcessor{
		dataStore: keyValueStorage,
		cmdParser: cmdParser,
	}
}

func (cp *commandsProcessor) Process(processorID int, requests <-chan queueservice.Message, replies chan<- queueservice.Message, wg *sync.WaitGroup) {
	defer wg.Done()

	for req := range requests {
		cmd, err := cp.cmdParser.ParseCommand(req.Body)
		if err != nil {
			log.Printf("Worker %d failed to parse command: %v", processorID, err)
			replies <- queueservice.Message{
				Body:          fmt.Sprintf("Failed to parse command: %v", err),
				ReplyTo:       req.ReplyTo,
				CorrelationId: req.CorrelationId,
			}
			continue
		}

		log.Printf("Worker %d processing command: %+v", processorID, cmd)

		switch cmd.Type {
		case cmd_parser.AddItem:
			cp.mu.Lock()
			cp.dataStore.Add(cmd.Key, cmd.Value)
			cp.mu.Unlock()
		case cmd_parser.DeleteItem:
			cp.mu.Lock()
			cp.dataStore.Remove(cmd.Key)
			cp.mu.Unlock()
		case cmd_parser.GetItem:
			cp.mu.RLock()
			value, exists := cp.dataStore.Get(cmd.Key)
			cp.mu.RUnlock()
			if exists {
				replies <- queueservice.Message{
					Body:          value.(string),
					ReplyTo:       req.ReplyTo,
					CorrelationId: req.CorrelationId,
				}
			} else {
				replies <- queueservice.Message{
					Body:          fmt.Sprintf("Key not found: %s", cmd.Key),
					ReplyTo:       req.ReplyTo,
					CorrelationId: req.CorrelationId,
				}
			}
		case cmd_parser.GetAllItems:
			cp.mu.RLock()
			items := cp.dataStore.GetAll()
			cp.mu.RUnlock()
			itemsJSON, err := json.Marshal(items)
			if err != nil {
				log.Printf("Worker %d failed to encode items to JSON: %v", processorID, err)
				replies <- queueservice.Message{
					Body:          fmt.Sprintf("Failed to encode items to JSON: %v", err),
					ReplyTo:       req.ReplyTo,
					CorrelationId: req.CorrelationId,
				}
				continue
			}
			replies <- queueservice.Message{
				Body:          string(itemsJSON),
				ReplyTo:       req.ReplyTo,
				CorrelationId: req.CorrelationId,
			}
		default:
			log.Printf("Worker %d received an unknown command type", processorID)
			replies <- queueservice.Message{
				Body:          fmt.Sprintf("Unknown command: %s", req.Body),
				ReplyTo:       req.ReplyTo,
				CorrelationId: req.CorrelationId,
			}
		}
	}
}

type CommandsParser interface {
	ParseCommand(command string) (cmd_parser.Command, error)
}

type KeyValueStorage interface {
	Add(key string, value interface{})
	Remove(key string)
	Get(key string) (interface{}, bool)
	GetAll() []ds.KeyValue
}
