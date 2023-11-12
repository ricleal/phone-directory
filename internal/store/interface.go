package store

import (
	"context"

	"phone-directory/internal/repository"
)

//go:generate mockgen -source=interface.go -package=mock -destination=mock/interface.go

// Store is the interface that wraps the repositories.
type Store interface {
	Users() repository.UserRepository
	Addresses() repository.AddressRepository
	Phones() repository.PhoneRepository
	ExecTx(ctx context.Context, fn func(Store) error) error
}
