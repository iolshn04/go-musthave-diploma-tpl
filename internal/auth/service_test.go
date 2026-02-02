package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/iolshn04/go-musthave-diploma-tpl/tree/master/internal/models"
)

type mockUserRepo struct {
	users map[string]*models.User
}

func newMockRepo() *mockUserRepo {
	return &mockUserRepo{users: make(map[string]*models.User)}
}

func (m *mockUserRepo) Create(_ context.Context, user *models.User) error {
	if _, ok := m.users[user.Login]; ok {
		return errors.New("exists")
	}
	m.users[user.Login] = user
	return nil
}

func (m *mockUserRepo) GetByLogin(_ context.Context, login string) (*models.User, error) {
	u, ok := m.users[login]
	if !ok {
		return nil, errors.New("not found")
	}
	return u, nil
}

func TestRegisterSuccess(t *testing.T) {
	repo := newMockRepo()
	service := NewService(repo)

	user, err := service.Register(context.Background(), "test", "secret")

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "test", user.Login)
	assert.NotEmpty(t, user.PasswordHash)
}

func TestRegisterDuplicate(t *testing.T) {
	repo := newMockRepo()
	service := NewService(repo)

	_, _ = service.Register(context.Background(), "test", "secret")
	_, err := service.Register(context.Background(), "test", "secret")

	assert.Error(t, err)
}

func TestAuthenticateSuccess(t *testing.T) {
	repo := newMockRepo()
	service := NewService(repo)

	user, _ := service.Register(context.Background(), "test", "secret")

	authUser, err := service.Authenticate(context.Background(), "test", "secret")

	assert.NoError(t, err)
	assert.Equal(t, user.ID, authUser.ID)
}

func TestAuthenticateWrongPassword(t *testing.T) {
	repo := newMockRepo()
	service := NewService(repo)

	_, _ = service.Register(context.Background(), "test", "secret")

	_, err := service.Authenticate(context.Background(), "test", "wrong")

	assert.Error(t, err)
}

func TestAuthenticateUnknownUser(t *testing.T) {
	repo := newMockRepo()
	service := NewService(repo)

	_, err := service.Authenticate(context.Background(), "unknown", "secret")

	assert.Error(t, err)
}
