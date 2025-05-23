package deleter

import (
	"context"
	"github.com/NeozonS/go-shortener-ya.git/internal/storage"
	"log"
	"time"
)

type deleteTask struct {
	UserID string
	Token  string
}

type Worker struct {
	deleteQueue  chan deleteTask
	storage      storage.Repository
	batchSize    int
	batchTimeout time.Duration
}

func NewWorker(repo storage.Repository, batchSize int, batchTimeout time.Duration) *Worker {
	return &Worker{
		deleteQueue:  make(chan deleteTask, 1000),
		storage:      repo,
		batchSize:    batchSize,
		batchTimeout: batchTimeout,
	}
}

func (w *Worker) InitDeleteWorker(ctx context.Context) {
	log.Println("Worker WORKED  by context")
	go w.batchDeleteWorker(ctx)
}

func (w *Worker) EnqueueURLForDeleteion(userID, token string) {
	log.Printf("Enqueuing deletion: user=%s, token=%s", userID, token)
	w.deleteQueue <- deleteTask{userID, token}
}

func (w *Worker) batchDeleteWorker(ctx context.Context) {
	batch := make([]deleteTask, 0, w.batchSize)
	timer := time.NewTimer(w.batchTimeout)
	defer timer.Stop()
	log.Println("Delete worker started")
	defer log.Println("Delete worker stopped")
	for {
		select {
		case <-ctx.Done():
			log.Println("Worker stopped by context")
			return
		case id := <-w.deleteQueue:
			log.Printf("Received task: user=%s, token=%s", id.UserID, id.Token)
			batch = append(batch, id)
			if len(batch) >= w.batchSize {
				log.Printf("Flushing full batch (%d items)", w.batchSize)
				w.processBatch(ctx, batch)
				batch = batch[:0]
				timer.Reset(w.batchTimeout)
			}

		case <-timer.C:
			if len(batch) > 0 {
				log.Printf("Flushing batch by timeout (%d items)", len(batch))
				w.processBatch(ctx, batch)
				batch = batch[:0]
			}
			timer.Reset(w.batchTimeout)
		}
	}
}
func (w *Worker) processBatch(ctx context.Context, task []deleteTask) {
	log.Printf("Processing batch of %d tasks", len(task))
	idsByUser := map[string][]string{}
	for _, t := range task {
		idsByUser[t.UserID] = append(idsByUser[t.UserID], t.Token)
	}
	for userID, tokens := range idsByUser {
		if err := w.storage.BatchDeleteURL(ctx, userID, tokens); err != nil {
			log.Printf("Error deleting URLs for user %s: %v", userID, err)
		}
	}
}
