package handlers

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/iolshn04/go-musthave-diploma-tpl/tree/master/internal/auth"
	"github.com/iolshn04/go-musthave-diploma-tpl/tree/master/internal/models"
)

type mockRepo struct {
	users map[string]*models.User
}

func newRepo() *mockRepo {
	return &mockRepo{users: make(map[string]*models.User)}
}

func (m *mockRepo) Create(_ context.Context, u *models.User) error {
	if _, ok := m.users[u.Login]; ok {
		return errors.New("exists")
	}
	m.users[u.Login] = u
	return nil
}

func (m *mockRepo) GetByLogin(_ context.Context, login string) (*models.User, error) {
	u, ok := m.users[login]
	if !ok {
		return nil, errors.New("not found")
	}
	return u, nil
}

func TestRegisterHandlerOK(t *testing.T) {
	repo := newRepo()
	service := auth.NewService(repo)
	handler := NewAuthHandler(service, "secret")

	body := []byte(`{"login":"test","password":"123"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/user/register", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.Register(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, w.Result().Cookies())
}

func TestRegisterInvalidBody(t *testing.T) {
	repo := newRepo()
	service := auth.NewService(repo)
	handler := NewAuthHandler(service, "secret")

	req := httptest.NewRequest(http.MethodPost, "/api/user/register", bytes.NewBuffer([]byte(`{`)))
	w := httptest.NewRecorder()

	handler.Register(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLoginUnauthorized(t *testing.T) {
	repo := newRepo()
	service := auth.NewService(repo)
	handler := NewAuthHandler(service, "secret")

	body := []byte(`{"login":"test","password":"bad"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/user/login", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.Login(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
