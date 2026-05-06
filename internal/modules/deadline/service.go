package deadline

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetAll() ([]Deadline, error) {
	return s.repo.GetAll()
}

func (s *Service) GetById(id uint) (*Deadline, error) {
	return s.repo.GetByID(id)
}

func (s *Service) Create(input Input) (*Deadline, error) {
	err := input.Validate(true)
	if err != nil {
		return nil, err
	}

	deadline := &Deadline{
		Title:    *input.Title,
		Status:   *input.Status,
		Priority: *input.Priority,
		DueDate:  *input.DueDate,
	}

	err = s.repo.Create(deadline)
	if err != nil {
		return nil, err
	}

	return deadline, nil
}

func (s *Service) Update(id uint, input Input) (*Deadline, error) {
	err := input.Validate(false)
	if err != nil {
		return nil, err
	}

	deadline, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if input.Title != nil {
		deadline.Title = *input.Title
	}
	if input.Status != nil {
		deadline.Status = *input.Status
	}
	if input.Priority != nil {
		deadline.Priority = *input.Priority
	}
	if input.DueDate != nil {
		deadline.DueDate = *input.DueDate
	}

	err = s.repo.Update(deadline)
	if err != nil {
		return nil, err
	}

	return deadline, nil
}

func (s *Service) Delete(id uint) error {
	deadline, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	return s.repo.Delete(deadline)
}
