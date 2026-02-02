package auth

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	"github.com/google/uuid"
	"github.com/iolshn04/go-musthave-diploma-tpl/tree/master/internal/models"
)

type UserRepo interface {
	Create(ctx context.Context, user *models.User) error
	GetByLogin(ctx context.Context, login string) (*models.User, error)
}

type Service struct {
	repo UserRepo
}

func NewService(repo UserRepo) *Service {
	return &Service{repo: repo}
}

func (s *Service) Register(ctx context.Context, login, password string) (*models.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:           uuid.NewString(),
		Login:        login,
		PasswordHash: string(hash),
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) Authenticate(ctx context.Context, login, password string) (*models.User, error) {
	user, err := s.repo.GetByLogin(ctx, login)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	); err != nil {
		return nil, err
	}

	return user, nil
}
