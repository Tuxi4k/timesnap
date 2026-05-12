package deadline_test

import (
	"errors"
	"log"
	"sort"
	"testing"
	"time"

	"github.com/Tuxi4k/timesnap/internal/modules/deadline"
	"github.com/Tuxi4k/timesnap/pkg/utils/ptr"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/stretchr/testify/assert"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func pick[T any](currentField, targetField int, newValue, oldValue T) T {
	if currentField == targetField {
		return newValue
	}
	return oldValue
}

func setupService() *deadline.Service {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed initialize database: %v", err)
	}

	err = db.AutoMigrate(&deadline.Deadline{})
	if err != nil {
		log.Fatalf("Failed migrate database: %v", err)
	}

	return deadline.NewService(deadline.NewRepository(db))
}

func TestService_Create_Success(t *testing.T) {
	svc := setupService()

	mockTitle := "Do something..."
	mockStatus := deadline.StatusActive
	mockPriority := deadline.PriorityMedium
	mockTime := time.Now().Add(7 * 24 * time.Hour)

	result, err := svc.Create(deadline.Input{
		Title:    &mockTitle,
		Status:   &mockStatus,
		Priority: &mockPriority,
		DueDate:  &mockTime,
	})

	assert.NoError(t, err)

	assert.Equal(t, mockTitle, result.Title)
	assert.Equal(t, mockStatus, result.Status)
	assert.Equal(t, mockPriority, result.Priority)
	assert.WithinDuration(t, mockTime, result.DueDate, time.Second)

	assert.NotZero(t, result.ID)
	assert.False(t, result.CreatedAt.IsZero())
}

func TestService_Update_Success(t *testing.T) {
	svc := setupService()

	created, err := svc.Create(deadline.Input{
		Title:    ptr.To("Do something..."),
		Status:   ptr.To(deadline.StatusActive),
		Priority: ptr.To(deadline.PriorityMedium),
		DueDate:  ptr.To(time.Now().Add(7 * 24 * time.Hour)),
	})

	assert.NoError(t, err)

	mockTitle := "Closed..."
	mockStatus := deadline.StatusCompleted
	mockPriority := deadline.PriorityLow
	mockTime := time.Now().Add(5 * 24 * time.Hour)

	updated, err := svc.Update(created.ID, deadline.Input{
		Title:    &mockTitle,
		Status:   &mockStatus,
		Priority: &mockPriority,
		DueDate:  &mockTime,
	})

	assert.NoError(t, err)

	assert.Equal(t, mockTitle, updated.Title)
	assert.Equal(t, mockStatus, updated.Status)
	assert.Equal(t, mockPriority, updated.Priority)
	assert.WithinDuration(t, mockTime, updated.DueDate, time.Second)

	assert.NotZero(t, updated.ID)
	assert.False(t, updated.CreatedAt.IsZero())
	assert.False(t, updated.UpdatedAt.IsZero())
}

func TestService_Update_Partial(t *testing.T) {
	svc := setupService()

	mockTitle := "Closed..."
	mockStatus := deadline.StatusCompleted
	mockPriority := deadline.PriorityLow
	mockTime := time.Now().Add(5 * 24 * time.Hour)

	for field := range 4 {
		created, err := svc.Create(deadline.Input{
			Title:    ptr.To("Do something..."),
			Status:   ptr.To(deadline.StatusActive),
			Priority: ptr.To(deadline.PriorityMedium),
			DueDate:  ptr.To(time.Now().Add(7 * 24 * time.Hour)),
		})
		assert.NoError(t, err)

		input := deadline.Input{
			Title:    pick(field, 0, &mockTitle, nil),
			Status:   pick(field, 1, &mockStatus, nil),
			Priority: pick(field, 2, &mockPriority, nil),
			DueDate:  pick(field, 3, &mockTime, nil),
		}

		updated, err := svc.Update(created.ID, input)
		assert.NoError(t, err)

		assert.Equal(t, pick(field, 0, mockTitle, created.Title), updated.Title)
		assert.Equal(t, pick(field, 1, mockStatus, created.Status), updated.Status)
		assert.Equal(t, pick(field, 2, mockPriority, created.Priority), updated.Priority)
		assert.WithinDuration(t, pick(field, 3, mockTime, created.DueDate), updated.DueDate, time.Second)

		assert.NotZero(t, updated.ID)
		assert.False(t, updated.CreatedAt.IsZero())
		assert.False(t, updated.UpdatedAt.IsZero())
	}
}

func TestService_Update_NotFound(t *testing.T) {
	svc := setupService()

	_, err := svc.Update(99, deadline.Input{})

	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
}

func TestService_Update_Validation(t *testing.T) {
	svc := setupService()

	created, err := svc.Create(deadline.Input{
		Title:    ptr.To("Valid"),
		Status:   ptr.To(deadline.StatusActive),
		Priority: ptr.To(deadline.PriorityMedium),
		DueDate:  ptr.To(time.Now().Add(time.Hour)),
	})

	assert.NoError(t, err)

	type step struct {
		input deadline.Input
		field string
	}

	steps := []step{
		{deadline.Input{Title: ptr.To("")}, "title"},
		{deadline.Input{Status: ptr.To(deadline.Status("unknown"))}, "status"},
		{deadline.Input{Priority: ptr.To(deadline.Priority("extreme"))}, "priority"},
		{deadline.Input{DueDate: ptr.To(time.Now().Add(-time.Hour))}, "due_date"},
	}

	for _, tc := range steps {
		_, err := svc.Update(created.ID, tc.input)

		var valErrs validation.Errors
		if assert.Error(t, err) {
			if assert.True(t, errors.As(err, &valErrs)) {
				assert.Contains(t, valErrs, tc.field)
			}
		}
	}
}

func TestService_Create_Validation(t *testing.T) {
	svc := setupService()

	validInput := func() deadline.Input {
		return deadline.Input{
			Title:    ptr.To("Valid Title"),
			Status:   ptr.To(deadline.StatusActive),
			Priority: ptr.To(deadline.PriorityMedium),
			DueDate:  ptr.To(time.Now().Add(time.Hour)),
		}
	}

	type step struct {
		modify func(i *deadline.Input)
		field  string
	}

	steps := []step{
		{func(i *deadline.Input) { i.Title = ptr.To("") }, "title"},
		{func(i *deadline.Input) { i.Status = ptr.To(deadline.Status("bad")) }, "status"},
		{func(i *deadline.Input) { i.Priority = ptr.To(deadline.Priority("none")) }, "priority"},
		{func(i *deadline.Input) { i.DueDate = ptr.To(time.Now().Add(-time.Hour)) }, "due_date"},
		{func(i *deadline.Input) { i.Title = nil }, "title"},
		{func(i *deadline.Input) { i.Status = nil }, "status"},
		{func(i *deadline.Input) { i.Priority = nil }, "priority"},
		{func(i *deadline.Input) { i.DueDate = nil }, "due_date"},
	}

	for _, tc := range steps {
		input := validInput()
		tc.modify(&input)

		_, err := svc.Create(input)

		var valErrs validation.Errors
		if assert.Error(t, err) {
			if assert.True(t, errors.As(err, &valErrs)) {
				assert.Contains(t, valErrs, tc.field)
			}
		}
	}
}

func TestService_Get_Success(t *testing.T) {
	svc := setupService()

	mockData := []deadline.Input{
		{
			Title:    ptr.To("Task 1"),
			Status:   ptr.To(deadline.StatusActive),
			Priority: ptr.To(deadline.PriorityLow),
			DueDate:  ptr.To(time.Now().Add(time.Hour * 2)),
		},
		{
			Title:    ptr.To("Task 2"),
			Status:   ptr.To(deadline.StatusCompleted),
			Priority: ptr.To(deadline.PriorityHigh),
			DueDate:  ptr.To(time.Now().Add(time.Hour * 3)),
		},
	}

	for _, inp := range mockData {
		_, err := svc.Create(inp)
		assert.NoError(t, err)
	}

	results, err := svc.GetAll()
	assert.NoError(t, err)
	assert.Len(t, results, 2)

	sort.Slice(results, func(i, j int) bool {
		return results[i].Title < results[j].Title
	})

	for ind, inp := range mockData {
		assert.Equal(t, *inp.Title, results[ind].Title)
		assert.Equal(t, *inp.Status, results[ind].Status)
		assert.Equal(t, *inp.Priority, results[ind].Priority)
		assert.WithinDuration(t, *inp.DueDate, results[ind].DueDate, time.Second)

		deadline, err := svc.GetByID(results[ind].ID)
		assert.NoError(t, err)

		assert.Equal(t, *inp.Title, deadline.Title)
		assert.Equal(t, *inp.Status, deadline.Status)
		assert.Equal(t, *inp.Priority, deadline.Priority)
		assert.WithinDuration(t, *inp.DueDate, deadline.DueDate, time.Second)
	}
}

func TestService_Get_NotFound(t *testing.T) {
	svc := setupService()

	_, err := svc.GetByID(99)

	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
}

func TestService_GetAll_Empty(t *testing.T) {
	svc := setupService()

	result, err := svc.GetAll()
	assert.NoError(t, err)

	assert.IsType(t, []deadline.Deadline{}, result)
	assert.Empty(t, result)
}

func TestService_Delete_Success(t *testing.T) {
	svc := setupService()

	created, err := svc.Create(deadline.Input{
		Title:    ptr.To("To be deleted"),
		Status:   ptr.To(deadline.StatusActive),
		Priority: ptr.To(deadline.PriorityLow),
		DueDate:  ptr.To(time.Now().Add(time.Hour)),
	})
	assert.NoError(t, err)

	err = svc.Delete(created.ID)
	assert.NoError(t, err)

	_, err = svc.GetByID(created.ID)
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
}

func TestService_Delete_NotFound(t *testing.T) {
	svc := setupService()

	err := svc.Delete(99)

	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
}
