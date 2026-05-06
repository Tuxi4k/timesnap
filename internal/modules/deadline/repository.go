package deadline

import (
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetAll() ([]Deadline, error) {
	var deadlines []Deadline

	err := r.db.Order("created_at DESC").Find(&deadlines).Error
	if err != nil {
		return nil, err
	}

	return deadlines, nil
}

func (r *Repository) GetByID(id uint) (*Deadline, error) {
	var deadline Deadline

	err := r.db.First(&deadline, id).Error
	if err != nil {
		return nil, err
	}

	return &deadline, nil
}

func (r *Repository) Create(deadline *Deadline) error {
	return r.db.Create(deadline).Error
}

func (r *Repository) Update(deadline *Deadline) error {
	return r.db.Select("title", "priority", "status", "due_date").Updates(deadline).Error
}

func (r *Repository) Delete(deadline *Deadline) error {
	return r.db.Delete(deadline).Error
}
