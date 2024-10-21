package worker

import (
	"context"
	"log"
	"sync"

	"github.com/srikanthbhandary/todo-server/entity"
)

type WorkerPool struct {
	numOfWorkers int            // The number of workers (goroutines) that will be processing jobs.
	wg           sync.WaitGroup // WaitGroup to keep track of active workers and ensure they complete their tasks before shutdown.
	inputChannel chan Job       // The channel through which jobs are submitted for processing. Workers will consume jobs from this channel.
	WebSocket    *entity.WebSocketConnection
}

// NewWorkerPool creates a new WorkerPool
func NewWorkerPool(numOfWorkers int, inputChannel chan Job) *WorkerPool {
	return &WorkerPool{
		numOfWorkers: numOfWorkers,
		inputChannel: inputChannel,
	}
}

// StartWorker processes jobs in the input channel
func (wp *WorkerPool) StartWorker(ctx context.Context) {
	defer wp.wg.Done()

	for {
		select {
		case <-ctx.Done():
			// Context cancelled, exit worker
			log.Println("worker exiting due to cancellation")
			return
		case job, ok := <-wp.inputChannel:
			if !ok {
				return
			}
			log.Println("JOB Received")
			if err := job.Process(); err != nil {
				log.Printf("error processing job: %v", err)
			}

		}
	}
}

// Init initializes and starts the workers
func (wp *WorkerPool) Init(ctx context.Context) {
	for i := 0; i < wp.numOfWorkers; i++ {
		wp.wg.Add(1)
		log.Println("starting worker", i+1)
		go wp.StartWorker(ctx)
	}
}

// Stop signals the workers to stop gracefully
func (wp *WorkerPool) Stop() {
	close(wp.inputChannel) // Close input channel to stop accepting new jobs
}

// Wait waits for all workers to finish
func (wp *WorkerPool) Wait() {
	wp.wg.Wait()
}

// EnqueueJob
func (wp *WorkerPool) EnqueueJob(job Job) {
	wp.inputChannel <- job
}
