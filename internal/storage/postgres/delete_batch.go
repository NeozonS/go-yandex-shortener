package postgres

import (
	"context"
	"log"
	"time"
)

type deleteTask struct {
	UserID string
	Token  string
}

var deleteQueue = make(chan deleteTask, 1000)

const batchSize = 100
const batchTimeout = 2 * time.Second

func (p *PostgresDB) InitDeleteWorker(ctx context.Context) {
	log.Println("Worker WORKED  by context")
	go p.batchDeleteWorker(ctx)
}

func EnqueueURLForDeleteion(userID, token string) {
	log.Printf("Enqueuing deletion: user=%s, token=%s", userID, token)
	deleteQueue <- deleteTask{userID, token}
}

func (p *PostgresDB) batchDeleteWorker(ctx context.Context) {
	batch := make([]deleteTask, 0, batchSize)
	timer := time.NewTimer(batchTimeout)
	defer timer.Stop()
	log.Println("Delete worker started")
	defer log.Println("Delete worker stopped")
	for {
		select {
		case <-ctx.Done():
			log.Println("Worker stopped by context")
			return
		case id := <-deleteQueue:
			log.Printf("Received task: user=%s, token=%s", id.UserID, id.Token)
			batch = append(batch, id)
			if len(batch) >= batchSize {
				log.Printf("Flushing full batch (%d items)", batchSize)
				p.processBatch(ctx, batch)
				batch = batch[:0]
				timer.Reset(batchTimeout)
			}

		case <-timer.C:
			if len(batch) > 0 {
				log.Printf("Flushing batch by timeout (%d items)", len(batch))
				p.processBatch(ctx, batch)
				batch = batch[:0]
			}
			timer.Reset(batchTimeout)
		}
	}
}
func (p *PostgresDB) processBatch(ctx context.Context, task []deleteTask) {
	log.Printf("Processing batch of %d tasks", len(task))
	idsByUser := map[string][]string{}
	for _, t := range task {
		idsByUser[t.UserID] = append(idsByUser[t.UserID], t.Token)
	}
	for userID, tokens := range idsByUser {
		if err := p.BatchDeleteURL(ctx, userID, tokens); err != nil {
			log.Printf("Error deleting URLs for user %s: %v", userID, err)
		}
	}
}
