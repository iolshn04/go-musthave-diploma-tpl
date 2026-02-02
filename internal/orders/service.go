package orders

import (
	"context"
	"github.com/iolshn04/go-musthave-diploma-tpl/tree/master/internal/models"
	"strings"
)

type Repository interface {
	Save(ctx context.Context, userID, number string) error
	Owner(ctx context.Context, number string) (string, error)
	List(ctx context.Context, userID string) ([]models.Order, error)
}

type ServiceInterface interface {
	Upload(ctx context.Context, userID, number string) error
	List(ctx context.Context, userID string) ([]models.Order, error)
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

func (s *Service) List(ctx context.Context, userID string) ([]models.Order, error) {
	return s.repo.List(ctx, userID)
}
