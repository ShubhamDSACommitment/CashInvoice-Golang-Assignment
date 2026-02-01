package worker

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/CashInvoice-Golang-Assignment/internal/repository"
)

type AutoCompleteWorker struct {
	repo  repository.TaskRepository
	queue chan string
	delay time.Duration
	wg    *sync.WaitGroup
}

// Constructor
func NewAutoCompleteWorker(
	repo repository.TaskRepository,
	queue chan string,
	delay time.Duration,
	wg *sync.WaitGroup,
) *AutoCompleteWorker {
	return &AutoCompleteWorker{
		repo:  repo,
		queue: queue,
		delay: delay,
		wg:    wg,
	}
}

func (w *AutoCompleteWorker) Start(ctx context.Context, numWorkers int) {
	for i := 0; i < numWorkers; i++ {
		w.wg.Add(1)
		go w.workerLoop(ctx, i)
	}
}

// Each worker runs forever
func (w *AutoCompleteWorker) workerLoop(ctx context.Context, id int) {
	defer w.wg.Done()
	log.Printf("Auto-complete worker %d started\n", id)
	for {
		select {
		case <-ctx.Done():
			log.Printf("Worker %d shutting down (context cancelled)\n", id)
			return

		case taskID, _ := <-w.queue:

			log.Printf("Worker %d received task %s\n", id, taskID)

			// Wait for delay or shutdown signal
			select {
			case <-time.After(w.delay):
				if err := w.repo.AutoCompleteIfPending(taskID); err != nil {
					log.Printf("Worker %d failed task %s: %v\n", id, taskID, err)
				} else {
					log.Printf("Worker %d completed task %s\n", id, taskID)
				}

			case <-ctx.Done():
				log.Printf("Worker %d cancelled while waiting\n", id)
				return
			}
		}
	}
}
