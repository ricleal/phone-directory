package memory

import (
	"context"
	"phone-directory/internal/repository"
)

type UserStorage struct {
}

func NewUserStorage() *UserStorage {
	return &UserStorage{}
}

func (s *UserStorage) Create(ctx context.Context, user *repository.User) error {
	user.ID = uint(len(storage.users) + 1)
	for i := range user.Phones {
		user.Phones[i].ID = uint(i + 1)
		user.Phones[i].UserID = user.ID
	}
	for i := range user.Addresses {
		user.Addresses[i].ID = uint(i + 1)
		user.Addresses[i].UserID = user.ID
	}
	storage.users = append(storage.users, *user)
	return nil
}

func (s *UserStorage) Get(ctx context.Context, id uint) (*repository.User, error) {
	for _, u := range storage.users {
		if u.ID == id {
			return &u, nil
		}
	}
	return nil, repository.ErrNotFound
}

func (s *UserStorage) Update(ctx context.Context, user *repository.User) error {
	for i, u := range storage.users {
		if u.ID == user.ID {
			storage.users[i] = *user
			return nil
		}
	}
	return repository.ErrNotFound
}

func (s *UserStorage) Delete(ctx context.Context, id uint) error {
	for i, u := range storage.users {
		if u.ID == id {
			storage.users = append(storage.users[:i], storage.users[i+1:]...)
			return nil
		}
	}
	return repository.ErrNotFound
}
