package deadline

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Input struct {
	Title    *string    `json:"title"`
	Status   *string    `json:"status"`
	Priority *string    `json:"priority"`
	DueDate  *time.Time `json:"due_date"`
}

func (i *Input) Validate(required bool) error {
	return validation.ValidateStruct(i,
		validation.Field(
			&i.Title,
			validation.When(required || i.Title != nil,
				validation.Required.Error("Обязательное поле")),
		),
		validation.Field(
			&i.Status,
			validation.When(required || i.Status != nil,
				validation.Required.Error("Обязательное поле")),
			validation.When(required || i.Status != nil,
				validation.In(
					StatusActive,
					StatusOverdue,
					StatusCompleted,
				).Error("Допустимые статусы: active, overdue, completed")),
		),
		validation.Field(
			&i.Priority,
			validation.When(required || i.Status != nil,
				validation.Required.Error("Обязательное поле")),
			validation.When(required || i.Status != nil,
				validation.In(
					PriorityLow,
					PriorityMedium,
					PriorityHigh,
				).Error("Допустимые статусы: high, medium, low")),
		),
		validation.Field(
			&i.DueDate,
			validation.When(required || i.DueDate != nil,
				validation.Required.Error("Обязательное поле"),
				validation.Min(
					time.Now().Add(-time.Minute)).Error("Дата дедлайна должна быть в будущем"),
			),
		),
	)
}
