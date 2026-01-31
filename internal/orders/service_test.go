package orders

import (
	"context"
	"errors"
	"github.com/iolshn04/go-musthave-diploma-tpl/tree/master/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockRepo struct {
	data    map[string]string
	saveErr error
}

func newMockRepo() *mockRepo {
	return &mockRepo{data: make(map[string]string)}
}

func (m *mockRepo) Save(ctx context.Context, userID, number string) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	m.data[number] = userID
	return nil
}

func (m *mockRepo) Owner(ctx context.Context, number string) (string, error) {
	uid, ok := m.data[number]
	if !ok {
		return "", errors.New("not found")
	}
	return uid, nil
}

func (m *mockRepo) List(ctx context.Context, userID string) ([]models.Order, error) {
	var res []models.Order
	for num, uid := range m.data {
		if uid == userID {
			res = append(res, models.Order{Number: num, Status: "NEW"})
		}
	}
	return res, nil
}

func TestServiceUpload(t *testing.T) {
	repo := newMockRepo()
	s := New(repo)

	user := "user1"

	// пустое тело
	err := s.Upload(context.Background(), user, "")
	assert.Equal(t, ErrEmptyBody, err)

	err = s.Upload(context.Background(), user, "1234")
	assert.Equal(t, ErrInvalidNumber, err)

	err = s.Upload(context.Background(), user, "79927398713")
	assert.NoError(t, err)

	err = s.Upload(context.Background(), user, "79927398713")
	assert.Equal(t, ErrAlreadyMine, err)

	err = s.Upload(context.Background(), "user2", "79927398713")
	assert.Equal(t, ErrAlreadyExists, err)
}

func TestServiceList(t *testing.T) {
	repo := newMockRepo()
	s := New(repo)
	user := "user1"

	s.Upload(context.Background(), user, "79927398713")

	list, err := s.List(context.Background(), user)
	assert.NoError(t, err)
	assert.Len(t, list, 1)
	assert.Equal(t, "79927398713", list[0].Number)
}
