package memory_test

import (
	"context"
	"phone-directory/internal/repository"
	"phone-directory/internal/repository/memory"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	ctx := context.Background()
	s := memory.NewUserStorage()
	user := &repository.User{
		Name: "John Doe",
	}
	err := s.Create(ctx, user)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), user.ID)
}

func TestCreateUserWithTelephoneAndAddress(t *testing.T) {
	ctx := context.Background()
	s := memory.NewUserStorage()
	user := &repository.User{
		Name: "John Doe",
		Phones: []repository.Phone{
			{
				Number: "123456789",
			},
		},
		Addresses: []repository.Address{
			{
				Address: "123 Main St",
			},
		},
	}
	err := s.Create(ctx, user)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), user.ID)
}

func TestGetUser(t *testing.T) {
	ctx := context.Background()
	s := memory.NewUserStorage()
	user := &repository.User{
		Name: "John Doe",
	}
	err := s.Create(ctx, user)
	assert.NoError(t, err)
	u, err := s.Get(ctx, user.ID)
	assert.NoError(t, err)
	assert.Equal(t, user, u)
}

func TestUpdateUser(t *testing.T) {
	ctx := context.Background()
	s := memory.NewUserStorage()
	user := &repository.User{
		Name: "John Doe",
	}
	err := s.Create(ctx, user)
	assert.NoError(t, err)
	user.Name = "Jane Doe"
	err = s.Update(ctx, user)
	assert.NoError(t, err)
	u, err := s.Get(ctx, user.ID)
	assert.NoError(t, err)
	assert.Equal(t, user, u)
}

func TestDeleteUser(t *testing.T) {
	ctx := context.Background()
	s := memory.NewUserStorage()
	user := &repository.User{
		Name: "John Doe",
	}
	err := s.Create(ctx, user)
	assert.NoError(t, err)
	err = s.Delete(ctx, user.ID)
	assert.NoError(t, err)
	_, err = s.Get(ctx, user.ID)
	assert.Equal(t, repository.ErrNotFound, err)
}
