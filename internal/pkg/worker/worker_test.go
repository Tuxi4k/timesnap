package worker_test

import (
	"testing"
	"time"

	"github.com/Tuxi4k/timesnap/internal/modules/deadline"
	"github.com/Tuxi4k/timesnap/internal/pkg/worker"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}

	err = db.AutoMigrate(&deadline.Deadline{})
	if err != nil {
		t.Fatalf("failed to migrate db: %v", err)
	}

	return db
}

func createDeadline(t *testing.T, db *gorm.DB, d deadline.Deadline) deadline.Deadline {
	t.Helper()

	err := db.Create(&d).Error
	if err != nil {
		t.Fatalf("failed to create deadline: %v", err)
	}

	return d
}

func TestWorker_CloseOverdue_UpdatesOnlyActiveExpired(t *testing.T) {
	db := setupDB(t)
	w := worker.New(db)

	now := time.Now()

	expiredActive := createDeadline(t, db, deadline.Deadline{
		Title:    "expired-active",
		Status:   deadline.StatusActive,
		Priority: deadline.PriorityMedium,
		DueDate:  now.Add(-2 * time.Hour),
	})

	notExpiredActive := createDeadline(t, db, deadline.Deadline{
		Title:    "not-expired-active",
		Status:   deadline.StatusActive,
		Priority: deadline.PriorityMedium,
		DueDate:  now.Add(2 * time.Hour),
	})

	expiredCompleted := createDeadline(t, db, deadline.Deadline{
		Title:    "expired-completed",
		Status:   deadline.StatusCompleted,
		Priority: deadline.PriorityMedium,
		DueDate:  now.Add(-2 * time.Hour),
	})

	rows, err := w.CloseOverdue(now)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), rows)

	var gotExpiredActive deadline.Deadline
	assert.NoError(t, db.First(&gotExpiredActive, expiredActive.ID).Error)
	assert.Equal(t, deadline.StatusCompleted, gotExpiredActive.Status)

	var gotNotExpiredActive deadline.Deadline
	assert.NoError(t, db.First(&gotNotExpiredActive, notExpiredActive.ID).Error)
	assert.Equal(t, deadline.StatusActive, gotNotExpiredActive.Status)

	var gotExpiredCompleted deadline.Deadline
	assert.NoError(t, db.First(&gotExpiredCompleted, expiredCompleted.ID).Error)
	assert.Equal(t, deadline.StatusCompleted, gotExpiredCompleted.Status)
}

func TestWorker_CloseOverdue_NoMatches(t *testing.T) {
	db := setupDB(t)
	w := worker.New(db)

	now := time.Now()

	createDeadline(t, db, deadline.Deadline{
		Title:    "future-active",
		Status:   deadline.StatusActive,
		Priority: deadline.PriorityLow,
		DueDate:  now.Add(24 * time.Hour),
	})

	rows, err := w.CloseOverdue(now)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), rows)
}

