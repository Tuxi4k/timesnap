package deadline

import "time"

const (
	StatusActive    Status = "active"
	StatusOverdue   Status = "overdue"
	StatusCompleted Status = "completed"
)
const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
)

type Status string
type Priority string

type Deadline struct {
	ID        uint      `json:"id"`
	Title     string    `gorm:"not null" json:"title"`
	Priority  string    `gorm:"type:varchar(10);not null;default:medium" json:"priority"`
	Status    string    `gorm:"type:varchar(10);not null;default:active" json:"status"`
	DueDate   time.Time `json:"due_date"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
