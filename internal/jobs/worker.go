package jobs

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/glueops/autoglue/internal/db/models"
	"gorm.io/gorm"
)

type Worker struct {
	DB       *gorm.DB
	WorkerID string
	PoolSize int
}

func (w *Worker) Run(ctx context.Context) {
	sem := make(chan struct{}, w.PoolSize)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			jobs, err := ClaimBatch(ctx, w.DB, w.WorkerID, w.PoolSize-len(sem))
			if err != nil {
				log.Printf("claim error: %v", err)
				continue
			}
			for i := range jobs {
				sem <- struct{}{}
				go func(j models.Job) {
					defer func() { <-sem }()
					w.handleJob(ctx, j)
				}(jobs[i])
			}
		}
	}
}

func (w *Worker) handleJob(ctx context.Context, j models.Job) {
	var err error
	switch j.Type {
	case "bootstrap_host":
		err = handleBootstrap(ctx, j.Payload)
	case "ansible_playbook":
		err = handleAnsible(ctx, j.Payload)
	default:
		err = fmt.Errorf("unknown job type: %s", j.Type)
	}

	if err == nil {
		_ = FinishSuccess(ctx, w.DB, j.ID)
		return
	}

	if j.Attempts+1 >= j.MaxAttempts {
		_ = FinishFailed(ctx, w.DB, j.ID, err.Error())
	} else {
		_ = FinishRetry(ctx, w.DB, j.ID, j.Attempts, j.MaxAttempts, err.Error())
	}
}
