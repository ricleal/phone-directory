package memory

import (
	"phone-directory/internal/repository"
	"sync"
)

var (
	storage Storage
	once    sync.Once
)

type Storage struct {
	users     []repository.User
	phones    []repository.Phone
	addresses []repository.Address
}

// NewStorage returns a new Handler with a database connection.
func NewStorage() (*Storage, error) {
	once.Do(func() {
		storage = Storage{}
	})
	return &storage, nil
}
