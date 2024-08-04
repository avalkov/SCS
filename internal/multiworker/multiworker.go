package multiworker

import (
	"context"
	"log"
	"strconv"
	"sync"

	ds "github.com/avalkov/SCS/internal/datastructures"
	"github.com/avalkov/SCS/internal/queueservice"
)

type MultiWorker struct {
	workersCount      int
	commandsParser    CommandsParser
	commandsProcessor CommandsProcessor
}

func NewMultiWorker(workersCount int, commandsParser CommandsParser, commandsProcessor CommandsProcessor) *MultiWorker {
	return &MultiWorker{
		workersCount:      workersCount,
		commandsParser:    commandsParser,
		commandsProcessor: commandsProcessor,
	}
}

func (mw *MultiWorker) Run(ctx context.Context, requests <-chan queueservice.Message, replies chan<- queueservice.Message) {
	defer ctx.Done()

	consistentHash := ds.NewConsistentHash(mw.workersCount, nil)

	for i := 0; i < mw.workersCount; i++ {
		consistentHash.Add(strconv.Itoa(i))
	}

	workerChans := make([]chan queueservice.Message, mw.workersCount)
	for i := range workerChans {
		workerChans[i] = make(chan queueservice.Message)
	}

	var wg sync.WaitGroup

	// Start worker goroutines
	for i := 0; i < mw.workersCount; i++ {
		wg.Add(1)
		go mw.commandsProcessor.Process(i, workerChans[i], replies, &wg)
	}

	// TODO: Can run multiple routers
	// Route requests to the appropriate worker based on consistent hashing
	go func() {
		for req := range requests {
			commandID, err := mw.commandsParser.GetCommandID(req.Body)
			if err != nil {
				log.Printf("Failed to get command ID: %v", err)
				replies <- queueservice.Message{
					Body:          "Failed to get command ID",
					ReplyTo:       req.ReplyTo,
					CorrelationId: req.CorrelationId,
				}
				continue
			}

			workerID := consistentHash.Get(commandID)
			workerIndex, err := strconv.Atoi(workerID)
			if err != nil {
				log.Printf("Failed to convert worker ID to index: %v", err)
				continue
			}
			workerChans[workerIndex] <- req
		}
	}()

	wg.Wait()

	// Close worker channels
	for i := range workerChans {
		close(workerChans[i])
	}
}

type CommandsProcessor interface {
	Process(processorID int, requests <-chan queueservice.Message, replies chan<- queueservice.Message, wg *sync.WaitGroup)
}

type CommandsParser interface {
	GetCommandID(command string) (string, error)
}
