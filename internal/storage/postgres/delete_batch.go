package storage

import (
	"context"
	"database/sql"
	"time"
	"github.com/NeozonS/go-shortener-ya.git/internal/storage/postgres"
)

var deleteQueue = make(chan string, 1000)
const batchSize = 100
const batchTimeout = 2 * time.Second

func InitDeleteWorker (ctx context.Context, db *sql.DB){
	go batchDeleteWorker(ctx, db)
}

func EnqueueURLForDeleteion (id string){
	deleteQueue <- id
}

func batchDeleteWorker (ctx context.Context, db *sql.DB) {
	batch := make([]string,0,batchSize)
	timer := time.NewTimer(batchTimeout)
	for {
		select {
		case <- ctx.Done():
			return
		}
		postgres.
		case id := <- deleteQueue:
			batch = append(batch, id)
			if len(batch) >= batchSize {
				post
		}
	}
}