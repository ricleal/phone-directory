package repository

import (
	"context"
	"database/sql"
)

//go:generate mockgen -source=interfaces.go -package=mock -destination=mock/interfaces.go

// DBTx represents a database transaction or connection interface.
type DBTx interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	Get(ctx context.Context, id uint) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id uint) error
}

type AddressRepository interface {
	Create(ctx context.Context, address *Address) error
	Get(ctx context.Context, id uint) (*Address, error)
	Update(ctx context.Context, address *Address) error
	Delete(ctx context.Context, id uint) error
}

type PhoneRepository interface {
	Create(ctx context.Context, phone *Phone) error
	Get(ctx context.Context, id uint) (*Phone, error)
	Update(ctx context.Context, phone *Phone) error
	Delete(ctx context.Context, id uint) error
}
