package orders

import (
	"context"
	"strings"
	"time"
)

type Order struct {
	Number     string    `json:"number" db:"number"`
	Status     string    `json:"status" db:"status"`
	Accrual    *float64  `json:"accrual,omitempty" db:"accrual"`
	UploadedAt time.Time `json:"uploaded_at" db:"uploaded_at"`
}

type Repository interface {
	Save(ctx context.Context, userID, number string) error
	Owner(ctx context.Context, number string) (string, error)
	List(ctx context.Context, userID string) ([]Order, error)
}

type ServiceInterface interface {
	Upload(ctx context.Context, userID, number string) error
	List(ctx context.Context, userID string) ([]Order, error)
}

type Service struct {
	repo Repository
}

func New(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Upload(ctx context.Context, userID, number string) error {
	number = strings.TrimSpace(number)
	if number == "" {
		return ErrEmptyBody
	}

	if !ValidLuhn(number) {
		return ErrInvalidNumber
	}

	owner, err := s.repo.Owner(ctx, number)
	if err == nil {
		if owner == userID {
			return ErrAlreadyMine
		}
		return ErrAlreadyExists
	}

	if err := s.repo.Save(ctx, userID, number); err != nil {
		return err
	}

	return nil
}

func (s *Service) List(ctx context.Context, userID string) ([]Order, error) {
	return s.repo.List(ctx, userID)
}
