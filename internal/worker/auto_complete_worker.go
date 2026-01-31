package worker

import (
	"log"
	"time"

	"github.com/CashInvoice-Golang-Assignment/internal/repository"
)

type AutoCompleteWorker struct {
	repo  repository.TaskRepository
	queue chan string
	delay time.Duration
}

// Constructor
func NewAutoCompleteWorker(
	repo repository.TaskRepository,
	queue chan string,
	delay time.Duration,
) *AutoCompleteWorker {
	return &AutoCompleteWorker{
		repo:  repo,
		queue: queue,
		delay: delay,
	}
}

func (w *AutoCompleteWorker) Start(numWorkers int) {
	for i := 0; i < numWorkers; i++ {
		go w.workerLoop(i)
	}
}

// Each worker runs forever
func (w *AutoCompleteWorker) workerLoop(id int) {
	log.Printf("Auto-complete worker %d started\n", id)

	for taskID := range w.queue {
		// Process each task in its own goroutine
		go func(tid string) {
			log.Printf("Worker %d received task %s\n", id, tid)

			time.Sleep(w.delay)

			err := w.repo.AutoCompleteIfPending(tid)
			if err != nil {
				log.Printf("Worker %d failed to auto-complete task %s: %v\n", id, tid, err)
			}
			log.Printf("Worker %d Completed task %s\n", id, tid)
		}(taskID)
	}
}
