package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"phone-directory/internal/repository"
	"phone-directory/internal/repository/persistent"
)

// PersistentStore is a store backed by a SQL database.
type persistentStore struct {
	db repository.DBTx
}

// NewPersistentStore creates a new store with the given database connection.
func NewPersistentStore(db repository.DBTx) *persistentStore {
	return &persistentStore{
		db: db,
	}
}

// Users returns a UserRepository for managing users.
func (s *persistentStore) Users() repository.UserRepository {
	return persistent.NewUserStorage(s.db)
}

// Addresses returns a AddressRepository for managing addresses.
func (s *persistentStore) Addresses() repository.AddressRepository {
	return persistent.NewAddressStorage(s.db)
}

// Phones returns a PhoneRepository for managing phones.
func (s *persistentStore) Phones() repository.PhoneRepository {
	return persistent.NewPhoneStorage(s.db)
}

// ExecTx executes the given function within a database transaction.
func (s *persistentStore) ExecTx(ctx context.Context, fn func(Store) error) error {
	db, ok := s.db.(*sql.DB)
	if !ok {
		return errors.New("ExecTx: db is not a *sql.DB")
	}
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("BeginTx: %w", err)
	}
	newStore := NewPersistentStore(tx)
	err = fn(newStore)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("ExecTx: %s: Rollback: %w", err.Error(), rbErr)
		}
		return fmt.Errorf("ExecTx: %w", err)
	}
	return tx.Commit() //nolint:wrapcheck //no need to wrap here
}
