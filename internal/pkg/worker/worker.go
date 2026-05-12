package worker

import (
	"context"
	"log"
	"time"

	"github.com/Tuxi4k/timesnap/internal/modules/deadline"
	"gorm.io/gorm"
)

type Worker struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Worker {
	return &Worker{db: db}
}

func (w *Worker) CloseOverdue(now time.Time) (int64, error) {
	res := w.db.Model(&deadline.Deadline{}).
		Where("status = ?", deadline.StatusActive).
		Where("due_date <= ?", now).
		Update("status", deadline.StatusCompleted)

	if res.Error != nil {
		return 0, res.Error
	}

	return res.RowsAffected, nil
}

func (w *Worker) Start(ctx context.Context, interval time.Duration) {
	if interval <= 0 {
		interval = time.Minute
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for tickTime := range ticker.C {
		if ctx.Err() != nil {
			return
		}

		updated, err := w.CloseOverdue(tickTime)
		if err != nil {
			log.Printf("worker: failed to close overdue deadlines: %v", err)
			continue
		}

		if updated > 0 {
			log.Printf("worker: closed overdue deadlines: %d", updated)
		}
	}
}
